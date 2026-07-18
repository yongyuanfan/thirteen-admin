package handlers

import (
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// ServerHandler 服务监控 HTTP 处理器。
type ServerHandler struct {
	handler.Handler
	svc *services.ServerService
}

// NewServerHandler 创建服务监控 Handler。
func NewServerHandler(svc *services.ServerService) *ServerHandler {
	return &ServerHandler{svc: svc}
}

// GetInfo GET /monitor/server
func (h *ServerHandler) GetInfo(ctx *gin.Context) {
	h.Success(ctx, h.svc.Info())
}
