package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerConfigRouter)
}

func registerConfigRouter(g *gin.RouterGroup) {
	h := hsystem.ConfigHandler{}
	r := g.Group("/system/config")
	{
		r.GET("/list", h.List)
		r.GET("/configKey/:configKey", h.GetByKey)
		r.DELETE("/refreshCache", h.RefreshCache)
		r.GET("/:configId", h.Get)
		r.POST("", h.Add)
		r.POST("/export", h.Export)
		r.PUT("", h.Edit)
		r.DELETE("/:configId", h.Delete)
	}
}
