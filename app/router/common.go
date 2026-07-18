package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

func init() {
	RegisterCheckRole(registerCommonRouter)
}

func registerCommonRouter(g *gin.RouterGroup) {
	h := handlers.CommonHandler{}
	g.POST("/common/upload", h.Upload)
	g.GET("/common/download", h.Download)
	g.GET("/common/download/resource", h.DownloadResource)
}
