package handlers

import (
	"errors"

	"thirteen-admin/app/middlewares"
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

func (e *AuthHandler) authService(c *gin.Context) (*services.AuthService, error) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		return nil, err
	}
	return services.NewAuthService(e.Orm, runtime.Redis(), runtime.JWT(), runtime.Captcha(), runtime.Online()), nil
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
	token, err := s.Login(c, req.Username, req.Password, req.Code, req.UUID, c.ClientIP(), c.Request.UserAgent())
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
	e.OK(c, gin.H{"token": token})
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

// Logout POST /logout
func (e AuthHandler) Logout(c *gin.Context) {
	s, err := e.authService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	token, _ := c.Get(middlewares.ContextTokenKey)
	tokenStr, _ := token.(string)
	_ = s.Logout(c, tokenStr)
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
