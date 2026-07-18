package monitor

import (
	hmonitor "thirteen-admin/app/handlers/monitor"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerOperLogRouter)
}

func registerOperLogRouter(g *gin.RouterGroup) {
	h := hmonitor.OperLogHandler{}
	r := g.Group("/monitor/operlog")
	{
		r.GET("/list", h.List)
		r.DELETE("/clean", h.Clean)
		r.POST("/export", h.Export)
		r.DELETE("/:operId", h.Delete)
	}
}
