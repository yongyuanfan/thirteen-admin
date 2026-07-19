package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("用户不存在/密码错误")
	ErrUserDisabled       = errors.New("用户已停用，请联系管理员")
	ErrCaptchaInvalid     = errors.New("验证码错误")
	ErrRegisterDisabled   = errors.New("当前系统没有开启注册功能")
	ErrSessionRevoked     = errors.New("登录状态已过期")
)

// AuthService 鉴权业务。
type AuthService struct {
	service.Service
	jwt     *JWTManager
	captcha *CaptchaService
	online  *OnlineService
}

// NewAuthService 创建鉴权服务。
func NewAuthService(db *gorm.DB, jwt *JWTManager, captcha *CaptchaService, online *OnlineService) *AuthService {
	return &AuthService{
		Service: service.New(db),
		jwt:     jwt,
		captcha: captcha,
		online:  online,
	}
}

// CaptchaImage 生成验证码。
func (s *AuthService) CaptchaImage(ctx context.Context) (uuid, img string, enabled bool, err error) {
	enabled = s.captcha.Enabled()
	if !enabled {
		return "", "", false, nil
	}
	uuid, img, err = s.captcha.Generate(ctx)
	return uuid, img, enabled, err
}

// Login 用户名密码登录。
func (s *AuthService) Login(ctx context.Context, username, password, code, uuid, ip, ua string) (*TokenPair, error) {
	if !s.captcha.Verify(ctx, uuid, code) {
		s.recordLogin(username, ip, "1", "验证码错误")
		return nil, ErrCaptchaInvalid
	}
	user, err := s.findByUserName(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.recordLogin(username, ip, "1", "用户不存在")
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if user.Status == "1" || user.DelFlag == "2" {
		s.recordLogin(username, ip, "1", "用户已停用")
		return nil, ErrUserDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.recordLogin(username, ip, "1", "密码错误")
		return nil, ErrInvalidCredentials
	}
	pair, err := s.jwt.IssuePair(user.ID, user.UserName, "", time.Time{})
	if err != nil {
		return nil, err
	}
	now := time.Now()
	_ = s.DB.Model(user).Updates(map[string]any{
		"login_ip":   ip,
		"login_date": now,
	}).Error
	s.recordLogin(username, ip, "0", "登录成功")
	if s.online != nil {
		_ = s.online.Put(ctx, OnlineUser{
			TokenID:   pair.Sid,
			UserName:  user.UserName,
			IPAddr:    ip,
			Browser:   ua,
			LoginTime: now.UnixMilli(),
		}, time.Duration(pair.RefreshExpiresIn)*time.Second)
	}
	return pair, nil
}

// Refresh 使用 refresh token 换发新的双 token。
func (s *AuthService) Refresh(ctx context.Context, refreshToken, ip, ua string) (*TokenPair, error) {
	if s.jwt == nil {
		return nil, ErrInvalidToken
	}
	claims, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, err
	}
	if s.jwt.AutoRevokeToken() && s.online != nil {
		ok, err := s.online.TokenExists(ctx, claims.Sid)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrSessionRevoked
		}
	}
	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		return nil, ErrInvalidToken
	}
	var user models.User
	if err := s.DB.First(&user, uint(userID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if user.Status == "1" || user.DelFlag == "2" {
		return nil, ErrUserDisabled
	}

	var refreshExp time.Time
	if s.jwt.AutoRefreshToken() {
		refreshExp = time.Time{} // IssuePair 内按 now+refreshTTL
	} else if claims.ExpiresAt != nil {
		refreshExp = claims.ExpiresAt.Time
	} else {
		return nil, ErrExpiredToken
	}

	pair, err := s.jwt.IssuePair(user.ID, user.UserName, claims.Sid, refreshExp)
	if err != nil {
		return nil, err
	}
	if s.online != nil {
		loginTime := time.Now().UnixMilli()
		if existing, err := s.online.Get(ctx, claims.Sid); err == nil && existing.LoginTime > 0 {
			loginTime = existing.LoginTime
		}
		_ = s.online.Put(ctx, OnlineUser{
			TokenID:   pair.Sid,
			UserName:  user.UserName,
			IPAddr:    ip,
			Browser:   ua,
			LoginTime: loginTime,
		}, time.Duration(pair.RefreshExpiresIn)*time.Second)
	}
	return pair, nil
}

// Register 注册。
func (s *AuthService) Register(ctx context.Context, username, password, code, uuid string) error {
	if !s.isRegisterEnabled() {
		return ErrRegisterDisabled
	}
	if !s.captcha.Verify(ctx, uuid, code) {
		return ErrCaptchaInvalid
	}
	var count int64
	if err := s.DB.Model(&models.User{}).Where("user_name = ?", username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("保存用户'" + username + "'失败，注册账号已存在")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &models.User{
		UserName: username,
		NickName: username,
		Password: string(hash),
		Status:   "0",
		DelFlag:  "0",
	}
	return s.DB.Create(user).Error
}

// Logout 退出：按 sid 删除在线会话（auto_revoke_token 开启时后续鉴权立即失败）。
func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	if s.online == nil {
		return nil
	}
	sid := ""
	if accessToken != "" && s.jwt != nil {
		if claims, err := s.jwt.ParseIgnoreExpiry(accessToken); err == nil && claims.Sid != "" {
			sid = claims.Sid
		}
	}
	if sid == "" && refreshToken != "" && s.jwt != nil {
		if claims, err := s.jwt.ParseIgnoreExpiry(refreshToken); err == nil && claims.Sid != "" {
			sid = claims.Sid
		}
	}
	if sid == "" {
		return nil
	}
	return s.online.Remove(ctx, sid)
}

// UnlockScreen 解锁屏幕。
func (s *AuthService) UnlockScreen(userID uint, password string) error {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

// GetInfo 当前用户信息。
func (s *AuthService) GetInfo(userID uint) (user map[string]any, roles, perms []string, isDefaultModifyPwd, isPasswordExpired bool, err error) {
	var u models.User
	if err = s.DB.First(&u, userID).Error; err != nil {
		return nil, nil, nil, false, false, err
	}
	roles, err = s.RoleKeys(userID)
	if err != nil {
		return nil, nil, nil, false, false, err
	}
	perms, err = s.Permissions(userID, roles)
	if err != nil {
		return nil, nil, nil, false, false, err
	}
	user = UserToMap(&u)
	return user, roles, perms, false, false, nil
}

// GetRouters 动态路由。
func (s *AuthService) GetRouters(userID uint) ([]map[string]any, error) {
	roles, err := s.RoleKeys(userID)
	if err != nil {
		return nil, err
	}
	var menus []models.Menu
	q := s.DB.Model(&models.Menu{}).Where("status = ? AND menu_type IN ?", "0", []string{"M", "C"}).Order("parent_id asc, order_num asc")
	if !isAdminRole(roles) {
		var menuIDs []uint
		err = s.DB.Table("sys_role_menu rm").
			Joins("JOIN sys_user_role ur ON ur.role_id = rm.role_id").
			Where("ur.user_id = ?", userID).
			Pluck("rm.menu_id", &menuIDs).Error
		if err != nil {
			return nil, err
		}
		if len(menuIDs) == 0 {
			return []map[string]any{}, nil
		}
		q = q.Where("id IN ?", menuIDs)
	}
	if err := q.Find(&menus).Error; err != nil {
		return nil, err
	}
	return buildRouterTree(menus, 0), nil
}

// RoleKeys 用户角色权限字符。
func (s *AuthService) RoleKeys(userID uint) ([]string, error) {
	var keys []string
	err := s.DB.Table("sys_role r").
		Joins("JOIN sys_user_role ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.status = ? AND r.del_flag = ?", userID, "0", "0").
		Pluck("r.role_key", &keys).Error
	if err != nil {
		return nil, err
	}
	if keys == nil {
		keys = []string{}
	}
	return keys, nil
}

// Permissions 权限标识。
func (s *AuthService) Permissions(userID uint, roles []string) ([]string, error) {
	if isAdminRole(roles) {
		return []string{"*:*:*"}, nil
	}
	var perms []string
	err := s.DB.Table("sys_menu m").
		Joins("JOIN sys_role_menu rm ON rm.menu_id = m.id").
		Joins("JOIN sys_user_role ur ON ur.role_id = rm.role_id").
		Where("ur.user_id = ? AND m.status = ? AND m.perms <> '' AND m.perms IS NOT NULL", userID, "0").
		Distinct().
		Pluck("m.perms", &perms).Error
	if err != nil {
		return nil, err
	}
	if perms == nil {
		perms = []string{}
	}
	return perms, nil
}

func (s *AuthService) findByUserName(username string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("user_name = ? AND del_flag = ?", username, "0").First(&user).Error
	return &user, err
}

func (s *AuthService) isRegisterEnabled() bool {
	var cfg models.Config
	if err := s.DB.Where("config_key = ?", "sys.account.registerUser").First(&cfg).Error; err != nil {
		return false
	}
	return strings.EqualFold(cfg.ConfigValue, "true")
}

func (s *AuthService) recordLogin(username, ip, status, msg string) {
	now := time.Now()
	_ = s.DB.Create(&models.Logininfor{
		UserName:  username,
		IPAddr:    ip,
		Status:    status,
		Msg:       msg,
		LoginTime: &now,
	}).Error
}

func isAdminRole(roles []string) bool {
	for _, r := range roles {
		if r == "admin" {
			return true
		}
	}
	return false
}

// UserToMap 用户 VO。
func UserToMap(u *models.User) map[string]any {
	if u == nil {
		return nil
	}
	m := map[string]any{
		"userId":      u.ID,
		"deptId":      u.DeptID,
		"userName":    u.UserName,
		"nickName":    u.NickName,
		"email":       u.Email,
		"phonenumber": u.Phonenumber,
		"sex":         u.Sex,
		"avatar":      u.Avatar,
		"status":      u.Status,
		"loginIp":     u.LoginIP,
		"loginDate":   FormatTime(u.LoginDate),
		"createdBy":   u.CreatedBy,
		"createdAt":   FormatTimeVal(u.CreatedAt),
		"updatedBy":   u.UpdatedBy,
		"updatedAt":   FormatTimeVal(u.UpdatedAt),
		"remark":      u.Remark,
	}
	return m
}

func buildRouterTree(menus []models.Menu, parentID uint) []map[string]any {
	out := make([]map[string]any, 0)
	for _, m := range menus {
		if m.ParentID != parentID {
			continue
		}
		children := buildRouterTree(menus, m.ID)
		node := map[string]any{
			"name":      routeName(m),
			"path":      routePath(m),
			"hidden":    m.Visible == "1",
			"component": routeComponent(m),
			"meta": map[string]any{
				"title":   m.MenuName,
				"icon":    m.Icon,
				"noCache": m.IsCache == 1,
				"link":    nil,
			},
		}
		if m.Query != "" {
			node["query"] = m.Query
		}
		if len(children) > 0 {
			if m.MenuType == "M" {
				node["alwaysShow"] = true
				node["redirect"] = "noRedirect"
			}
			node["children"] = children
		}
		out = append(out, node)
	}
	return out
}

func routeName(m models.Menu) string {
	if m.RouteName != "" {
		return m.RouteName
	}
	if m.Path == "" {
		return fmt.Sprintf("Menu%d", m.ID)
	}
	parts := strings.Split(m.Path, "/")
	last := parts[len(parts)-1]
	if last == "" {
		return fmt.Sprintf("Menu%d", m.ID)
	}
	return strings.ToUpper(last[:1]) + last[1:]
}

func isExternalPath(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "mailto:") ||
		strings.HasPrefix(path, "tel:")
}

func routePath(m models.Menu) string {
	if isExternalPath(m.Path) {
		return m.Path
	}
	if m.ParentID == 0 && !strings.HasPrefix(m.Path, "/") {
		return "/" + m.Path
	}
	return m.Path
}

func routeComponent(m models.Menu) string {
	if m.ParentID == 0 && m.MenuType == "M" {
		return "Layout"
	}
	if m.MenuType == "M" {
		return "ParentView"
	}
	if m.Component != "" {
		return m.Component
	}
	return "ParentView"
}
