package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterCommonRoutes 注册公共上传下载路由。
func RegisterCommonRoutes(h *handlers.CommonHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		g.POST("/common/upload", auth, h.Upload)
		g.GET("/common/download", auth, h.Download)
		g.GET("/common/download/resource", auth, h.DownloadResource)
	})
}
