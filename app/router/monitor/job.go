package monitor

import (
	hmonitor "thirteen-admin/app/handlers/monitor"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerJobRouter)
	router.RegisterCheckRole(registerJobLogRouter)
}

func registerJobRouter(g *gin.RouterGroup) {
	h := hmonitor.JobHandler{}
	r := g.Group("/monitor/job")
	{
		r.GET("/list", h.List)
		r.PUT("/changeStatus", h.ChangeStatus)
		r.PUT("/run", h.Run)
		r.POST("/export", h.Export)
		r.GET("/:jobId", h.Get)
		r.POST("", h.Add)
		r.PUT("", h.Edit)
		r.DELETE("/:jobId", h.Delete)
	}
}

func registerJobLogRouter(g *gin.RouterGroup) {
	h := hmonitor.JobLogHandler{}
	r := g.Group("/monitor/jobLog")
	{
		r.GET("/list", h.List)
		r.DELETE("/clean", h.Clean)
		r.POST("/export", h.Export)
		r.DELETE("/:jobLogId", h.Delete)
	}
}
