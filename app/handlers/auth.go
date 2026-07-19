package handlers

import (
	"errors"
	"strings"

	"thirteen-admin/app/runtime"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// AuthHandler 鉴权 HTTP 处理器。
type AuthHandler struct {
	handler.Api
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
	UUID     string `json:"uuid"`
}

type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
	UUID     string `json:"uuid"`
}

type unlockReq struct {
	Password string `json:"password"`
}

type refreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

type logoutReq struct {
	RefreshToken string `json:"refreshToken"`
}

func (e *AuthHandler) authService(c *gin.Context) (*services.AuthService, error) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		return nil, err
	}
	return services.NewAuthService(e.Orm, runtime.JWT(), runtime.Captcha(), runtime.Online()), nil
}

func tokenPairResp(pair *services.TokenPair) gin.H {
	return gin.H{
		"token":            pair.AccessToken,
		"refreshToken":     pair.RefreshToken,
		"expiresIn":        pair.ExpiresIn,
		"refreshExpiresIn": pair.RefreshExpiresIn,
	}
}

// Login POST /login
func (e AuthHandler) Login(c *gin.Context) {
	s, err := e.authService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	pair, err := s.Login(c, req.Username, req.Password, req.Code, req.UUID, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCaptchaInvalid):
			e.Fail(c, handler.CodeError, err.Error())
		case errors.Is(err, services.ErrUserDisabled):
			e.Fail(c, handler.CodeError, err.Error())
		case errors.Is(err, services.ErrInvalidCredentials):
			e.Fail(c, handler.CodeError, err.Error())
		default:
			e.InternalError(c, err)
		}
		return
	}
	e.OK(c, tokenPairResp(pair))
}

// Refresh POST /refresh
func (e AuthHandler) Refresh(c *gin.Context) {
	s, err := e.authService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.RefreshToken) == "" {
		e.BadRequest(c, "参数错误")
		return
	}
	pair, err := s.Refresh(c, req.RefreshToken, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		switch {
		case errors.Is(err, services.ErrExpiredToken),
			errors.Is(err, services.ErrInvalidToken),
			errors.Is(err, services.ErrSessionRevoked):
			e.Unauthorized(c, "登录状态已过期")
		case errors.Is(err, services.ErrUserDisabled):
			e.Fail(c, handler.CodeError, err.Error())
		default:
			e.InternalError(c, err)
		}
		return
	}
	e.OK(c, tokenPairResp(pair))
}

// Register POST /register
func (e AuthHandler) Register(c *gin.Context) {
	s, err := e.authService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.Register(c, req.Username, req.Password, req.Code, req.UUID); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "注册成功")
}

// CaptchaImage GET /captchaImage
func (e AuthHandler) CaptchaImage(c *gin.Context) {
	s, err := e.authService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	uuid, img, enabled, err := s.CaptchaImage(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.OK(c, gin.H{
		"uuid":           uuid,
		"img":            img,
		"captchaEnabled": enabled,
	})
}

// GetInfo GET /getInfo
func (e AuthHandler) GetInfo(c *gin.Context) {
	s, err := e.authService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	uid := services.ContextUserID(c)
	user, roles, perms, defPwd, expired, err := s.GetInfo(uid)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.OK(c, gin.H{
		"user":               user,
		"roles":              roles,
		"permissions":        perms,
		"isDefaultModifyPwd": defPwd,
		"isPasswordExpired":  expired,
	})
}

// Logout POST /logout（可选鉴权：Bearer access 或 body.refreshToken）
func (e AuthHandler) Logout(c *gin.Context) {
	s, err := e.authService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	var req logoutReq
	_ = c.ShouldBindJSON(&req)

	accessToken := ""
	header := c.GetHeader("Authorization")
	const prefix = "Bearer "
	if strings.HasPrefix(header, prefix) {
		accessToken = strings.TrimSpace(strings.TrimPrefix(header, prefix))
	}
	_ = s.Logout(c, accessToken, req.RefreshToken)
	e.SuccessMsg(c, "退出成功")
}

// UnlockScreen POST /unlockscreen
func (e AuthHandler) UnlockScreen(c *gin.Context) {
	s, err := e.authService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	var req unlockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.UnlockScreen(services.ContextUserID(c), req.Password); err != nil {
		e.Fail(c, handler.CodeError, "密码错误")
		return
	}
	e.SuccessMsg(c, "解锁成功")
}

// GetRouters GET /getRouters
func (e AuthHandler) GetRouters(c *gin.Context) {
	s, err := e.authService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	data, err := s.GetRouters(services.ContextUserID(c))
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}
