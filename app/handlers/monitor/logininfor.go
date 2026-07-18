package monitor

import (
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// LogininforHandler 登录日志 HTTP 处理器。
type LogininforHandler struct {
	handler.Handler
	svc *services.LogininforService
}

// NewLogininforHandler 创建登录日志 Handler。
func NewLogininforHandler(svc *services.LogininforService) *LogininforHandler {
	return &LogininforHandler{svc: svc}
}

// List GET /monitor/logininfor/list
func (h *LogininforHandler) List(ctx *gin.Context) {
	total, rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// Delete DELETE /monitor/logininfor/:infoId
func (h *LogininforHandler) Delete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("infoId"))
	if len(ids) == 0 {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Delete(ids); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "删除成功")
}

// Unlock GET /monitor/logininfor/unlock/:userName
func (h *LogininforHandler) Unlock(ctx *gin.Context) {
	userName := ctx.Param("userName")
	if err := h.svc.Unlock(ctx.Request.Context(), userName); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "解锁成功")
}

// Clean DELETE /monitor/logininfor/clean
func (h *LogininforHandler) Clean(ctx *gin.Context) {
	if err := h.svc.Clean(); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "清空成功")
}

// Export POST /monitor/logininfor/export
func (h *LogininforHandler) Export(ctx *gin.Context) {
	rows, err := h.svc.Export(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	services.WriteCSV(ctx, "logininfor.csv", []string{"访问编号", "用户名称", "登录地址", "登录状态", "描述", "访问时间"}, rows)
}
