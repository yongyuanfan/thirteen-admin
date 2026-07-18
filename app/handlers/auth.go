package handlers

import (
	"errors"

	"thirteen-admin/app/middlewares"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// AuthHandler 鉴权 HTTP 处理器。
type AuthHandler struct {
	handler.Handler
	svc *services.AuthService
}

// NewAuthHandler 创建鉴权 Handler。
func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
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

// Login POST /login
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req loginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	token, err := h.svc.Login(ctx, req.Username, req.Password, req.Code, req.UUID, ctx.ClientIP(), ctx.Request.UserAgent())
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCaptchaInvalid):
			h.Fail(ctx, handler.CodeError, err.Error())
		case errors.Is(err, services.ErrUserDisabled):
			h.Fail(ctx, handler.CodeError, err.Error())
		case errors.Is(err, services.ErrInvalidCredentials):
			h.Fail(ctx, handler.CodeError, err.Error())
		default:
			h.InternalError(ctx, err)
		}
		return
	}
	h.OK(ctx, gin.H{"token": token})
}

// Register POST /register
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req registerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Register(ctx, req.Username, req.Password, req.Code, req.UUID); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "注册成功")
}

// CaptchaImage GET /captchaImage
func (h *AuthHandler) CaptchaImage(ctx *gin.Context) {
	uuid, img, enabled, err := h.svc.CaptchaImage(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.OK(ctx, gin.H{
		"uuid":           uuid,
		"img":            img,
		"captchaEnabled": enabled,
	})
}

// GetInfo GET /getInfo
func (h *AuthHandler) GetInfo(ctx *gin.Context) {
	uid := services.ContextUserID(ctx)
	user, roles, perms, defPwd, expired, err := h.svc.GetInfo(uid)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.OK(ctx, gin.H{
		"user":               user,
		"roles":              roles,
		"permissions":        perms,
		"isDefaultModifyPwd": defPwd,
		"isPasswordExpired":  expired,
	})
}

// Logout POST /logout
func (h *AuthHandler) Logout(ctx *gin.Context) {
	token, _ := ctx.Get(middlewares.ContextTokenKey)
	tokenStr, _ := token.(string)
	_ = h.svc.Logout(ctx, tokenStr)
	h.SuccessMsg(ctx, "退出成功")
}

// UnlockScreen POST /unlockscreen
func (h *AuthHandler) UnlockScreen(ctx *gin.Context) {
	var req unlockReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.UnlockScreen(services.ContextUserID(ctx), req.Password); err != nil {
		h.Fail(ctx, handler.CodeError, "密码错误")
		return
	}
	h.SuccessMsg(ctx, "解锁成功")
}

// GetRouters GET /getRouters
func (h *AuthHandler) GetRouters(ctx *gin.Context) {
	data, err := h.svc.GetRouters(services.ContextUserID(ctx))
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}
