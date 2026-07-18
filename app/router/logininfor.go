package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterLogininforRoutes 注册登录日志路由。
func RegisterLogininforRoutes(h *handlers.LogininforHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		m := g.Group("/monitor/logininfor", auth)
		m.GET("/list", h.List)
		m.GET("/unlock/:userName", h.Unlock)
		m.DELETE("/clean", h.Clean)
		m.POST("/export", h.Export)
		m.DELETE("/:infoId", h.Delete)
	})
}
