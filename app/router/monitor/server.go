package monitor

import (
	hmonitor "thirteen-admin/app/handlers/monitor"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerServerRouter)
}

func registerServerRouter(g *gin.RouterGroup) {
	h := hmonitor.ServerHandler{}
	g.GET("/monitor/server", h.GetInfo)
}
