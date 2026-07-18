package handlers

import (
	"errors"
	"strings"

	"clash-admin/app/middlewares"
	"clash-admin/app/services"
	"clash-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type refreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

type phoneReq struct {
	Phone string `json:"phone"`
}

type captchaVerifyReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type phoneCodeLoginReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type registerReq struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

type resetPasswordReq struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

// Login POST /auth/login
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req loginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数无效")
		return
	}
	req.UserName = strings.TrimSpace(req.UserName)
	if req.UserName == "" || req.Password == "" {
		h.BadRequest(ctx, "用户名和密码不能为空")
		return
	}
	token, err := h.svc.Login(ctx.Request.Context(), req.UserName, req.Password)
	if err != nil {
		h.mapError(ctx, err)
		return
	}
	h.Success(ctx, token)
}

// GetUserInfo GET /auth/getUserInfo
func (h *AuthHandler) GetUserInfo(ctx *gin.Context) {
	userID, ok := ctx.Get(middlewares.ContextUserIDKey)
	if !ok {
		h.Fail(ctx, handler.CodeUnauthorized, "未登录")
		return
	}
	info, err := h.svc.GetUserInfo(ctx.Request.Context(), userID.(string))
	if err != nil {
		h.mapError(ctx, err)
		return
	}
	h.Success(ctx, info)
}

// RefreshToken POST /auth/refreshToken
func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var req refreshReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数无效")
		return
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		h.BadRequest(ctx, "refreshToken 不能为空")
		return
	}
	token, err := h.svc.RefreshToken(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		h.mapError(ctx, err)
		return
	}
	h.Success(ctx, token)
}

// SendCaptcha POST /auth/sendCaptcha
func (h *AuthHandler) SendCaptcha(ctx *gin.Context) {
	var req phoneReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数无效")
		return
	}
	phone, err := services.NormalizePhone(req.Phone)
	if err != nil {
		h.BadRequest(ctx, err.Error())
		return
	}
	if err := h.svc.SendCaptcha(ctx.Request.Context(), phone); err != nil {
		h.mapError(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// VerifyCaptcha POST /auth/verifyCaptcha
func (h *AuthHandler) VerifyCaptcha(ctx *gin.Context) {
	var req captchaVerifyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数无效")
		return
	}
	phone, err := services.NormalizePhone(req.Phone)
	if err != nil {
		h.BadRequest(ctx, err.Error())
		return
	}
	if strings.TrimSpace(req.Code) == "" {
		h.BadRequest(ctx, "验证码不能为空")
		return
	}
	if err := h.svc.VerifyCaptcha(ctx.Request.Context(), phone, req.Code); err != nil {
		h.mapError(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// CodeLogin POST /auth/codeLogin
func (h *AuthHandler) CodeLogin(ctx *gin.Context) {
	var req phoneCodeLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数无效")
		return
	}
	phone, err := services.NormalizePhone(req.Phone)
	if err != nil {
		h.BadRequest(ctx, err.Error())
		return
	}
	if strings.TrimSpace(req.Code) == "" {
		h.BadRequest(ctx, "验证码不能为空")
		return
	}
	token, err := h.svc.CodeLogin(ctx.Request.Context(), phone, req.Code)
	if err != nil {
		h.mapError(ctx, err)
		return
	}
	h.Success(ctx, token)
}

// Register POST /auth/register
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req registerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数无效")
		return
	}
	phone, err := services.NormalizePhone(req.Phone)
	if err != nil {
		h.BadRequest(ctx, err.Error())
		return
	}
	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Password) == "" {
		h.BadRequest(ctx, "验证码和密码不能为空")
		return
	}
	if len(req.Password) < 6 {
		h.BadRequest(ctx, "密码长度至少 6 位")
		return
	}
	token, err := h.svc.Register(ctx.Request.Context(), phone, req.Code, req.Password)
	if err != nil {
		h.mapError(ctx, err)
		return
	}
	h.Success(ctx, token)
}

// ResetPassword POST /auth/resetPassword
func (h *AuthHandler) ResetPassword(ctx *gin.Context) {
	var req resetPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数无效")
		return
	}
	phone, err := services.NormalizePhone(req.Phone)
	if err != nil {
		h.BadRequest(ctx, err.Error())
		return
	}
	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Password) == "" {
		h.BadRequest(ctx, "验证码和密码不能为空")
		return
	}
	if len(req.Password) < 6 {
		h.BadRequest(ctx, "密码长度至少 6 位")
		return
	}
	if err := h.svc.ResetPassword(ctx.Request.Context(), phone, req.Code, req.Password); err != nil {
		h.mapError(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

func (h *AuthHandler) mapError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidCredentials), errors.Is(err, services.ErrCaptchaInvalid):
		h.Fail(ctx, handler.CodeAuthFailed, err.Error())
	case errors.Is(err, services.ErrCaptchaTooFrequent):
		h.Fail(ctx, handler.CodeBadRequest, err.Error())
	case errors.Is(err, gorm.ErrRecordNotFound):
		h.Fail(ctx, handler.CodeUserNotFound, "用户不存在或手机号未注册")
	case errors.Is(err, services.ErrPhoneRegistered):
		h.Fail(ctx, handler.CodeUserExists, err.Error())
	case errors.Is(err, services.ErrRefreshInvalid), errors.Is(err, services.ErrInvalidToken):
		h.Fail(ctx, handler.CodeUnauthorized, "未登录或令牌已失效")
	case errors.Is(err, services.ErrUserDisabled):
		h.Fail(ctx, handler.CodeUnauthorized, err.Error())
	case errors.Is(err, services.ErrExpiredToken):
		h.Fail(ctx, handler.CodeTokenExpired, "令牌已过期")
	default:
		h.InternalError(ctx, err)
	}
}
