package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterOnlineRoutes 注册在线用户路由。
func RegisterOnlineRoutes(h *handlers.OnlineHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		m := g.Group("/monitor/online", auth)
		m.GET("/list", h.List)
		m.DELETE("/:tokenId", h.ForceLogout)
	})
}
