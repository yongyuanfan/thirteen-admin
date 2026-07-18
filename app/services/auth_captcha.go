package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"clash-admin/internal/logging"

	"github.com/redis/go-redis/v9"
)

const (
	captchaTTL      = 5 * time.Minute
	captchaRateTTL  = 60 * time.Second
	captchaKeyPref  = "auth:captcha:"
	captchaRatePref = "auth:captcha:rate:"
)

var (
	ErrCaptchaTooFrequent = errors.New("验证码发送过于频繁，请稍后再试")
	ErrCaptchaInvalid     = errors.New("验证码错误或已过期")
)

// CaptchaStore Redis 验证码存储。
type CaptchaStore struct {
	rdb  *redis.Client
	mode string // application.mode，非 prod 时打日志
}

// NewCaptchaStore 创建验证码存储。
func NewCaptchaStore(rdb *redis.Client, mode string) *CaptchaStore {
	return &CaptchaStore{rdb: rdb, mode: mode}
}

// Send 生成并保存验证码，受 60s 频率限制。
func (s *CaptchaStore) Send(ctx context.Context, phone string) (string, error) {
	rateKey := captchaRatePref + phone
	ok, err := s.rdb.SetNX(ctx, rateKey, "1", captchaRateTTL).Result()
	if err != nil {
		return "", fmt.Errorf("captcha rate limit: %w", err)
	}
	if !ok {
		return "", ErrCaptchaTooFrequent
	}

	code, err := generateCode(6)
	if err != nil {
		return "", err
	}

	if err := s.rdb.Set(ctx, captchaKeyPref+phone, code, captchaTTL).Err(); err != nil {
		return "", fmt.Errorf("store captcha: %w", err)
	}

	if s.mode != "prod" {
		logging.Info("send captcha (dev)", "phone", phone, "code", code)
	}
	return code, nil
}

// Verify 校验验证码，成功后删除（一次性）。
func (s *CaptchaStore) Verify(ctx context.Context, phone, code string) error {
	key := captchaKeyPref + phone
	stored, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCaptchaInvalid
		}
		return fmt.Errorf("get captcha: %w", err)
	}
	if stored != code {
		return ErrCaptchaInvalid
	}
	_ = s.rdb.Del(ctx, key).Err()
	return nil
}

func generateCode(digits int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("generate captcha: %w", err)
	}
	format := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(format, n.Int64()), nil
}
