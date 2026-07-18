package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerPostRouter)
}

func registerPostRouter(g *gin.RouterGroup) {
	h := hsystem.PostHandler{}
	r := g.Group("/system/post")
	{
		r.GET("/list", h.List)
		r.GET("/:postId", h.Get)
		r.POST("", h.Add)
		r.POST("/export", h.Export)
		r.PUT("", h.Edit)
		r.DELETE("/:postId", h.Delete)
	}
}
