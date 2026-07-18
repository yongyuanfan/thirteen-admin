package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 业务错误码，与前端 VITE_SERVICE_*_CODES 对齐。
const (
	CodeSuccess      = "0000"
	CodeBadRequest   = "4001"
	CodeAuthFailed   = "4010"
	CodeUserNotFound = "4011"
	CodeUserExists   = "4012"
	CodeUnauthorized = "8888"
	CodeTokenExpired = "9999"
)

// Handler 通用 HTTP 处理器基类。
type Handler struct{}

// Response 统一业务响应信封。
type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Success 成功响应，HTTP 始终 200。
func (h *Handler) Success(ctx *gin.Context, data any) {
	h.JSON(ctx, http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "ok",
		Data: data,
	})
}

// Fail 业务失败响应，HTTP 始终 200，由 code 区分语义。
func (h *Handler) Fail(ctx *gin.Context, code, msg string) {
	h.JSON(ctx, http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// JSON 通用 JSON 响应
func (h *Handler) JSON(ctx *gin.Context, code int, data any) {
	ctx.JSON(code, data)
}

// BadRequest 参数错误响应（业务信封）
func (h *Handler) BadRequest(ctx *gin.Context, message string) {
	h.Fail(ctx, CodeBadRequest, message)
}

// NotFound 资源不存在响应（业务信封）
func (h *Handler) NotFound(ctx *gin.Context, message string) {
	h.Fail(ctx, CodeUserNotFound, message)
}

// Conflict 资源冲突响应（业务信封）
func (h *Handler) Conflict(ctx *gin.Context, message string) {
	h.Fail(ctx, CodeUserExists, message)
}

// InternalError 服务内部错误响应（业务信封）
func (h *Handler) InternalError(ctx *gin.Context, err error) {
	if err == nil {
		h.InternalErrorMsg(ctx, "内部服务错误")
		return
	}
	h.InternalErrorMsg(ctx, err.Error())
}

// InternalErrorMsg 服务内部错误响应（业务信封）
func (h *Handler) InternalErrorMsg(ctx *gin.Context, message string) {
	h.Fail(ctx, "5000", message)
}

// ErrorMessage 返回标准错误体（非业务信封，供特殊场景使用）
func (h *Handler) ErrorMessage(ctx *gin.Context, code int, message string) {
	h.JSON(ctx, code, gin.H{"error": message})
}
