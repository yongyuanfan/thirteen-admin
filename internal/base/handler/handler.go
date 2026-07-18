package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 若依前端业务码（数字）。
const (
	CodeSuccess      = 200
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeWarn         = 601
	CodeError        = 500
)

// Handler 通用 HTTP 处理器基类。
type Handler struct{}

// AjaxResult 统一业务响应信封。
type AjaxResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// Success 成功响应（data 可选）。
func (h *Handler) Success(ctx *gin.Context, data any) {
	h.JSON(ctx, http.StatusOK, AjaxResult{
		Code: CodeSuccess,
		Msg:  "操作成功",
		Data: data,
	})
}

// SuccessMsg 成功响应（自定义 msg，无 data）。
func (h *Handler) SuccessMsg(ctx *gin.Context, msg string) {
	h.JSON(ctx, http.StatusOK, AjaxResult{
		Code: CodeSuccess,
		Msg:  msg,
	})
}

// OK 成功且将字段平铺到根级（登录 token、getInfo 等）。
func (h *Handler) OK(ctx *gin.Context, fields gin.H) {
	if fields == nil {
		fields = gin.H{}
	}
	fields["code"] = CodeSuccess
	if _, ok := fields["msg"]; !ok {
		fields["msg"] = "操作成功"
	}
	h.JSON(ctx, http.StatusOK, fields)
}

// Table 分页列表响应。
func (h *Handler) Table(ctx *gin.Context, total int64, rows any) {
	if rows == nil {
		rows = []any{}
	}
	h.JSON(ctx, http.StatusOK, gin.H{
		"code":  CodeSuccess,
		"msg":   "查询成功",
		"total": total,
		"rows":  rows,
	})
}

// Fail 业务失败（HTTP 仍 200）。
func (h *Handler) Fail(ctx *gin.Context, code int, msg string) {
	h.JSON(ctx, http.StatusOK, AjaxResult{
		Code: code,
		Msg:  msg,
	})
}

// JSON 通用 JSON 响应。
func (h *Handler) JSON(ctx *gin.Context, httpCode int, data any) {
	ctx.JSON(httpCode, data)
}

// BadRequest 参数错误。
func (h *Handler) BadRequest(ctx *gin.Context, message string) {
	h.Fail(ctx, CodeBadRequest, message)
}

// Unauthorized 未认证。
func (h *Handler) Unauthorized(ctx *gin.Context, message string) {
	h.Fail(ctx, CodeUnauthorized, message)
}

// Warn 业务警告。
func (h *Handler) Warn(ctx *gin.Context, message string) {
	h.Fail(ctx, CodeWarn, message)
}

// NotFound 资源不存在。
func (h *Handler) NotFound(ctx *gin.Context, message string) {
	h.Fail(ctx, CodeError, message)
}

// Conflict 资源冲突。
func (h *Handler) Conflict(ctx *gin.Context, message string) {
	h.Fail(ctx, CodeError, message)
}

// InternalError 服务内部错误。
func (h *Handler) InternalError(ctx *gin.Context, err error) {
	if err == nil {
		h.InternalErrorMsg(ctx, "内部服务错误")
		return
	}
	h.InternalErrorMsg(ctx, err.Error())
}

// InternalErrorMsg 服务内部错误。
func (h *Handler) InternalErrorMsg(ctx *gin.Context, message string) {
	h.Fail(ctx, CodeError, message)
}
