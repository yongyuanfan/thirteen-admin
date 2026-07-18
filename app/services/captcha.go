package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const captchaKeyPref = "captcha:"

// CaptchaService 图形验证码（简易数字图）。
type CaptchaService struct {
	rdb     *redis.Client
	enabled bool
}

// NewCaptchaService 创建验证码服务。
func NewCaptchaService(rdb *redis.Client, enabled bool) *CaptchaService {
	return &CaptchaService{rdb: rdb, enabled: enabled}
}

// Enabled 是否启用验证码。
func (s *CaptchaService) Enabled() bool { return s.enabled }

// Generate 生成验证码。
func (s *CaptchaService) Generate(ctx context.Context) (id, b64s string, err error) {
	if !s.enabled {
		return "", "", nil
	}
	code, err := randomDigits(4)
	if err != nil {
		return "", "", err
	}
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return "", "", err
	}
	id = fmt.Sprintf("%x", idBytes)
	img, err := renderDigitPNG(code)
	if err != nil {
		return "", "", err
	}
	b64s = "data:image/png;base64," + base64.StdEncoding.EncodeToString(img)
	if s.rdb != nil {
		if err := s.rdb.Set(ctx, captchaKeyPref+id, code, 2*time.Minute).Err(); err != nil {
			return "", "", err
		}
	}
	return id, b64s, nil
}

// Verify 校验验证码。
func (s *CaptchaService) Verify(ctx context.Context, id, code string) bool {
	if !s.enabled {
		return true
	}
	if id == "" || code == "" || s.rdb == nil {
		return false
	}
	key := captchaKeyPref + id
	val, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	_ = s.rdb.Del(ctx, key).Err()
	return strings.EqualFold(val, code)
}

func randomDigits(n int) (string, error) {
	var b strings.Builder
	for i := 0; i < n; i++ {
		v, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		b.WriteByte(byte('0' + v.Int64()))
	}
	return b.String(), nil
}

func renderDigitPNG(code string) ([]byte, error) {
	const w, h = 120, 40
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	bg := color.RGBA{R: 240, G: 240, B: 245, A: 255}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, bg)
		}
	}
	fg := color.RGBA{R: 30, G: 30, B: 30, A: 255}
	for i, ch := range code {
		drawDigit(img, 12+i*26, 8, ch, fg)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 极简 5x7 点阵数字。
func drawDigit(img *image.RGBA, ox, oy int, ch rune, c color.Color) {
	glyphs := map[rune][7]string{
		'0': {"111", "101", "101", "101", "101", "101", "111"},
		'1': {"010", "110", "010", "010", "010", "010", "111"},
		'2': {"111", "001", "001", "111", "100", "100", "111"},
		'3': {"111", "001", "001", "111", "001", "001", "111"},
		'4': {"101", "101", "101", "111", "001", "001", "001"},
		'5': {"111", "100", "100", "111", "001", "001", "111"},
		'6': {"111", "100", "100", "111", "101", "101", "111"},
		'7': {"111", "001", "001", "001", "001", "001", "001"},
		'8': {"111", "101", "101", "111", "101", "101", "111"},
		'9': {"111", "101", "101", "111", "001", "001", "111"},
	}
	g, ok := glyphs[ch]
	if !ok {
		return
	}
	scale := 3
	for row, line := range g {
		for col, p := range line {
			if p != '1' {
				continue
			}
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.Set(ox+col*scale+dx, oy+row*scale+dy, c)
				}
			}
		}
	}
}
