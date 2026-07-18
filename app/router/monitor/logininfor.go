package monitor

import (
	hmonitor "thirteen-admin/app/handlers/monitor"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerLogininforRouter)
}

func registerLogininforRouter(g *gin.RouterGroup) {
	h := hmonitor.LogininforHandler{}
	r := g.Group("/monitor/logininfor")
	{
		r.GET("/list", h.List)
		r.GET("/unlock/:userName", h.Unlock)
		r.DELETE("/clean", h.Clean)
		r.POST("/export", h.Export)
		r.DELETE("/:infoId", h.Delete)
	}
}
