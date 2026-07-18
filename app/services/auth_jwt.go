package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"clash-admin/internal/config"

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
	jwt.RegisteredClaims
}

// TokenPair access + refresh 双令牌。
type TokenPair struct {
	Token        string
	RefreshToken string
	RefreshJTI   string
	UserID       string
}

// JWTManager JWT 签发与解析。
type JWTManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewJWTManager 从配置创建 JWT 管理器。
func NewJWTManager(cfg config.JWT) *JWTManager {
	accessTTL := time.Duration(cfg.ExpiresIn) * time.Second
	if accessTTL <= 0 {
		accessTTL = 2 * time.Hour
	}
	refreshTTL := time.Duration(cfg.RefreshExpiresIn) * time.Second
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	return &JWTManager{
		secret:     []byte(cfg.Secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// AccessTTL 返回 access token 有效期。
func (m *JWTManager) AccessTTL() time.Duration { return m.accessTTL }

// RefreshTTL 返回 refresh token 有效期。
func (m *JWTManager) RefreshTTL() time.Duration { return m.refreshTTL }

// IssueTokenPair 为用户签发 access + refresh。
func (m *JWTManager) IssueTokenPair(userID uint, userName string) (*TokenPair, error) {
	uid := strconv.FormatUint(uint64(userID), 10)
	accessJTI, err := newJTI()
	if err != nil {
		return nil, err
	}
	access, err := m.sign(uid, userName, tokenTypeAccess, m.accessTTL, accessJTI)
	if err != nil {
		return nil, fmt.Errorf("签发 access token: %w", err)
	}
	refreshJTI, err := newJTI()
	if err != nil {
		return nil, err
	}
	refresh, err := m.sign(uid, userName, tokenTypeRefresh, m.refreshTTL, refreshJTI)
	if err != nil {
		return nil, fmt.Errorf("签发 refresh token: %w", err)
	}
	return &TokenPair{
		Token:        access,
		RefreshToken: refresh,
		RefreshJTI:   refreshJTI,
		UserID:       uid,
	}, nil
}

func newJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate jti: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (m *JWTManager) sign(userID, userName, typ string, ttl time.Duration, jti string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		Typ:      typ,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Subject:   userID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseAccess 解析并校验 access token。
func (m *JWTManager) ParseAccess(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, tokenTypeAccess)
}

// ParseRefresh 解析并校验 refresh token。
func (m *JWTManager) ParseRefresh(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, tokenTypeRefresh)
}

func (m *JWTManager) parse(tokenStr, expectTyp string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.Typ != expectTyp {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
