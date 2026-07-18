package monitor

import (
	hmonitor "thirteen-admin/app/handlers/monitor"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

// RegisterServerRoutes 注册服务监控路由。
func RegisterServerRoutes(h *hmonitor.ServerHandler) {
	router.RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		g.GET("/monitor/server", auth, h.GetInfo)
	})
}
