package monitor

import (
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// ServerHandler 服务监控 HTTP 处理器。
type ServerHandler struct {
	handler.Api
}

// GetInfo GET /monitor/server
func (e ServerHandler) GetInfo(c *gin.Context) {
	s := services.ServerService{}
	e.Success(c, s.Info())
}
