package monitor

import (
	hmonitor "thirteen-admin/app/handlers/monitor"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

// RegisterJobRoutes 注册定时任务路由。
func RegisterJobRoutes(h *hmonitor.JobHandler) {
	router.RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		m := g.Group("/monitor/job", auth)
		m.GET("/list", h.List)
		m.PUT("/changeStatus", h.ChangeStatus)
		m.PUT("/run", h.Run)
		m.POST("/export", h.Export)
		m.GET("/:jobId", h.Get)
		m.POST("", h.Add)
		m.PUT("", h.Edit)
		m.DELETE("/:jobId", h.Delete)
	})
}

// RegisterJobLogRoutes 注册任务日志路由。
func RegisterJobLogRoutes(h *hmonitor.JobLogHandler) {
	router.RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		m := g.Group("/monitor/jobLog", auth)
		m.GET("/list", h.List)
		m.DELETE("/clean", h.Clean)
		m.POST("/export", h.Export)
		m.DELETE("/:jobLogId", h.Delete)
	})
}
