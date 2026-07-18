package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterServerRoutes 注册服务监控路由。
func RegisterServerRoutes(h *handlers.ServerHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		g.GET("/monitor/server", auth, h.GetInfo)
	})
}
