package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

// RegisterPostRoutes 注册岗位路由。
func RegisterPostRoutes(h *hsystem.PostHandler) {
	router.RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		sys := g.Group("/system/post", auth)
		sys.GET("/list", h.List)
		sys.GET("/:postId", h.Get)
		sys.POST("", h.Add)
		sys.POST("/export", h.Export)
		sys.PUT("", h.Edit)
		sys.DELETE("/:postId", h.Delete)
	})
}
