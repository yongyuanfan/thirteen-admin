package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterOperLogRoutes 注册操作日志路由。
func RegisterOperLogRoutes(h *handlers.OperLogHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		m := g.Group("/monitor/operlog", auth)
		m.GET("/list", h.List)
		m.DELETE("/clean", h.Clean)
		m.POST("/export", h.Export)
		m.DELETE("/:operId", h.Delete)
	})
}
