package monitor

import (
	hmonitor "thirteen-admin/app/handlers/monitor"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

// RegisterOnlineRoutes 注册在线用户路由。
func RegisterOnlineRoutes(h *hmonitor.OnlineHandler) {
	router.RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		m := g.Group("/monitor/online", auth)
		m.GET("/list", h.List)
		m.DELETE("/:tokenId", h.ForceLogout)
	})
}
