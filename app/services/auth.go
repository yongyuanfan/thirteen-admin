package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"clash-admin/app/models"
	"clash-admin/internal/base/model"
	"clash-admin/internal/base/service"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const refreshKeyPref = "auth:refresh:"

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserDisabled       = errors.New("用户已被禁用")
	ErrPhoneRegistered    = errors.New("手机号已注册")
	ErrRefreshInvalid     = errors.New("刷新令牌无效")
)

// LoginToken 登录令牌响应。
type LoginToken struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// UserInfo 用户信息响应。
type UserInfo struct {
	UserID   string   `json:"userId"`
	UserName string   `json:"userName"`
	Roles    []string `json:"roles"`
	Buttons  []string `json:"buttons"`
}

// AuthService 鉴权业务。
type AuthService struct {
	service.Service
	rdb     *redis.Client
	jwt     *JWTManager
	captcha *CaptchaStore
}

// NewAuthService 创建鉴权服务。
func NewAuthService(db *gorm.DB, rdb *redis.Client, jwt *JWTManager, captcha *CaptchaStore) *AuthService {
	return &AuthService{
		Service: service.New(db),
		rdb:     rdb,
		jwt:     jwt,
		captcha: captcha,
	}
}

// Login 用户名密码登录。
func (s *AuthService) Login(ctx context.Context, userName, password string) (*LoginToken, error) {
	user, err := s.findByUserName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return s.issueAndStore(ctx, user)
}

// CodeLogin 手机验证码登录。
func (s *AuthService) CodeLogin(ctx context.Context, phone, code string) (*LoginToken, error) {
	if err := s.captcha.Verify(ctx, phone, code); err != nil {
		return nil, err
	}
	user, err := s.findByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}
	return s.issueAndStore(ctx, user)
}

// Register 手机号注册并直接登录。
func (s *AuthService) Register(ctx context.Context, phone, code, password string) (*LoginToken, error) {
	if err := s.captcha.Verify(ctx, phone, code); err != nil {
		return nil, err
	}
	var count int64
	if err := s.DB.WithContext(ctx).Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrPhoneRegistered
	}
	userName := phone
	var nameCount int64
	if err := s.DB.WithContext(ctx).Model(&models.User{}).Where("user_name = ?", userName).Count(&nameCount).Error; err != nil {
		return nil, err
	}
	if nameCount > 0 {
		userName = phone + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	p := phone
	user := &models.User{
		UserName:     userName,
		PasswordHash: string(hash),
		Phone:        &p,
		Roles:        model.StringSlice{"R_USER"},
		Buttons:      model.StringSlice{},
		Status:       1,
	}
	if err := s.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return s.issueAndStore(ctx, user)
}

// ResetPassword 通过验证码重置密码。
func (s *AuthService) ResetPassword(ctx context.Context, phone, code, password string) error {
	if err := s.captcha.Verify(ctx, phone, code); err != nil {
		return err
	}
	user, err := s.findByPhone(phone)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	return s.DB.WithContext(ctx).Model(user).Update("password_hash", string(hash)).Error
}

// SendCaptcha 发送验证码。
func (s *AuthService) SendCaptcha(ctx context.Context, phone string) error {
	_, err := s.captcha.Send(ctx, phone)
	return err
}

// VerifyCaptcha 独立校验验证码。
func (s *AuthService) VerifyCaptcha(ctx context.Context, phone, code string) error {
	return s.captcha.Verify(ctx, phone, code)
}

// GetUserInfo 按用户 ID 查询信息。
func (s *AuthService) GetUserInfo(ctx context.Context, userID string) (*UserInfo, error) {
	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, ErrInvalidToken
	}
	var user models.User
	if err := s.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	roles := []string(user.Roles)
	if roles == nil {
		roles = []string{}
	}
	buttons, err := ResolveButtonsByRoleCodes(ctx, s.DB, roles)
	if err != nil {
		return nil, err
	}
	// 兼容种子用户上直接配置的 "*" 按钮
	for _, b := range user.Buttons {
		if b == "*" {
			buttons = []string{"*"}
			break
		}
	}
	return &UserInfo{
		UserID:   strconv.FormatUint(uint64(user.ID), 10),
		UserName: user.UserName,
		Roles:    roles,
		Buttons:  buttons,
	}, nil
}

// RefreshToken 轮转刷新令牌。
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*LoginToken, error) {
	claims, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, ErrRefreshInvalid
	}
	key := refreshRedisKey(claims.UserID, claims.ID)
	stored, err := s.rdb.Get(ctx, key).Result()
	if err != nil || stored != refreshToken {
		return nil, ErrRefreshInvalid
	}
	_ = s.rdb.Del(ctx, key).Err()

	id, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		return nil, ErrRefreshInvalid
	}
	var user models.User
	if err := s.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, ErrRefreshInvalid
	}
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}
	return s.issueAndStore(ctx, &user)
}

func (s *AuthService) issueAndStore(ctx context.Context, user *models.User) (*LoginToken, error) {
	pair, err := s.jwt.IssueTokenPair(user.ID, user.UserName)
	if err != nil {
		return nil, err
	}
	key := refreshRedisKey(pair.UserID, pair.RefreshJTI)
	if err := s.rdb.Set(ctx, key, pair.RefreshToken, s.jwt.RefreshTTL()).Err(); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}
	return &LoginToken{
		Token:        pair.Token,
		RefreshToken: pair.RefreshToken,
	}, nil
}

func (s *AuthService) findByUserName(userName string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("user_name = ?", userName).First(&user).Error
	return &user, err
}

func (s *AuthService) findByPhone(phone string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func refreshRedisKey(userID, jti string) string {
	return refreshKeyPref + userID + ":" + jti
}

// NormalizePhone 简单校验手机号。
func NormalizePhone(phone string) (string, error) {
	phone = strings.TrimSpace(phone)
	if len(phone) < 6 || len(phone) > 20 {
		return "", errors.New("手机号格式不正确")
	}
	return phone, nil
}
