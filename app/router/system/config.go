package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

// RegisterConfigRoutes 注册参数路由。
func RegisterConfigRoutes(h *hsystem.ConfigHandler) {
	router.RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		sys := g.Group("/system/config", auth)
		sys.GET("/list", h.List)
		sys.GET("/configKey/:configKey", h.GetByKey)
		sys.DELETE("/refreshCache", h.RefreshCache)
		sys.GET("/:configId", h.Get)
		sys.POST("", h.Add)
		sys.POST("/export", h.Export)
		sys.PUT("", h.Edit)
		sys.DELETE("/:configId", h.Delete)
	})
}
