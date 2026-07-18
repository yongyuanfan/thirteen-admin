package handlers

import (
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// OperLogHandler 操作日志 HTTP 处理器。
type OperLogHandler struct {
	handler.Handler
	svc *services.OperLogService
}

// NewOperLogHandler 创建操作日志 Handler。
func NewOperLogHandler(svc *services.OperLogService) *OperLogHandler {
	return &OperLogHandler{svc: svc}
}

// List GET /monitor/operlog/list
func (h *OperLogHandler) List(ctx *gin.Context) {
	total, rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// Delete DELETE /monitor/operlog/:operId
func (h *OperLogHandler) Delete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("operId"))
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

// Clean DELETE /monitor/operlog/clean
func (h *OperLogHandler) Clean(ctx *gin.Context) {
	if err := h.svc.Clean(); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "清空成功")
}

// Export POST /monitor/operlog/export
func (h *OperLogHandler) Export(ctx *gin.Context) {
	rows, err := h.svc.Export(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	services.WriteCSV(ctx, "operlog.csv", []string{"日志编号", "系统模块", "业务类型", "操作人员", "主机", "操作状态", "操作时间"}, rows)
}
