package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"thirteen-admin/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

var (
	ErrInvalidToken = errors.New("无效的令牌")
	ErrExpiredToken = errors.New("令牌已过期")
)

// Claims JWT 载荷。
type Claims struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Typ      string `json:"typ"`
	Sid      string `json:"sid"`
	jwt.RegisteredClaims
}

// TokenPair 双 token 签发结果。
type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	Sid              string
	ExpiresIn        int // access 有效秒数
	RefreshExpiresIn int // refresh 剩余有效秒数
}

// JWTManager JWT 签发与解析。
type JWTManager struct {
	secret           []byte
	accessTTL        time.Duration
	refreshTTL       time.Duration
	autoRefreshToken bool
	autoRevokeToken  bool
}

// NewJWTManager 从配置创建 JWT 管理器。
func NewJWTManager(cfg config.JWT) *JWTManager {
	accessTTL := time.Duration(cfg.ExpiresIn) * time.Second
	if accessTTL <= 0 {
		accessTTL = 30 * time.Minute
	}
	refreshTTL := time.Duration(cfg.RefreshExpiresIn) * time.Second
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	autoRefresh := true
	if cfg.AutoRefreshToken != nil {
		autoRefresh = *cfg.AutoRefreshToken
	}
	autoRevoke := true
	if cfg.AutoRevokeToken != nil {
		autoRevoke = *cfg.AutoRevokeToken
	}
	return &JWTManager{
		secret:           []byte(cfg.Secret),
		accessTTL:        accessTTL,
		refreshTTL:       refreshTTL,
		autoRefreshToken: autoRefresh,
		autoRevokeToken:  autoRevoke,
	}
}

// AccessTTL 返回 access token 有效期。
func (m *JWTManager) AccessTTL() time.Duration { return m.accessTTL }

// RefreshTTL 返回 refresh token 窗口。
func (m *JWTManager) RefreshTTL() time.Duration { return m.refreshTTL }

// AutoRefreshToken 是否开启滑动续期。
func (m *JWTManager) AutoRefreshToken() bool { return m.autoRefreshToken }

// AutoRevokeToken 是否开启主动吊销。
func (m *JWTManager) AutoRevokeToken() bool { return m.autoRevokeToken }

// IssuePair 签发 access + refresh。sid 为空时新建会话；refreshExp 为零时使用 now+refreshTTL。
func (m *JWTManager) IssuePair(userID uint, userName, sid string, refreshExp time.Time) (*TokenPair, error) {
	uid := strconv.FormatUint(uint64(userID), 10)
	if sid == "" {
		var err error
		sid, err = newJTI()
		if err != nil {
			return nil, err
		}
	}
	now := time.Now()
	if refreshExp.IsZero() {
		refreshExp = now.Add(m.refreshTTL)
	}
	if !refreshExp.After(now) {
		return nil, ErrExpiredToken
	}

	accessJTI, err := newJTI()
	if err != nil {
		return nil, err
	}
	refreshJTI, err := newJTI()
	if err != nil {
		return nil, err
	}

	access, err := m.sign(uid, userName, tokenTypeAccess, sid, accessJTI, now.Add(m.accessTTL))
	if err != nil {
		return nil, err
	}
	refresh, err := m.sign(uid, userName, tokenTypeRefresh, sid, refreshJTI, refreshExp)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:      access,
		RefreshToken:     refresh,
		Sid:              sid,
		ExpiresIn:        int(m.accessTTL.Seconds()),
		RefreshExpiresIn: int(refreshExp.Sub(now).Seconds()),
	}, nil
}

func newJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate jti: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (m *JWTManager) sign(userID, userName, typ, sid, jti string, exp time.Time) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		Typ:      typ,
		Sid:      sid,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   userID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) parse(tokenStr string, validateExp bool) (*Claims, error) {
	opts := []jwt.ParserOption{jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})}
	if !validateExp {
		opts = append(opts, jwt.WithoutClaimsValidation())
	}
	parser := jwt.NewParser(opts...)
	token, err := parser.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		if validateExp && errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || (validateExp && !token.Valid) {
		return nil, ErrInvalidToken
	}
	if claims.Sid == "" {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// ParseAccess 解析并校验 access token。
func (m *JWTManager) ParseAccess(tokenStr string) (*Claims, error) {
	claims, err := m.parse(tokenStr, true)
	if err != nil {
		return nil, err
	}
	if claims.Typ != tokenTypeAccess {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// ParseRefresh 解析并校验 refresh token。
func (m *JWTManager) ParseRefresh(tokenStr string) (*Claims, error) {
	claims, err := m.parse(tokenStr, true)
	if err != nil {
		return nil, err
	}
	if claims.Typ != tokenTypeRefresh {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// ParseIgnoreExpiry 解析 token（忽略过期），用于退出时取 sid。
func (m *JWTManager) ParseIgnoreExpiry(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, false)
}
