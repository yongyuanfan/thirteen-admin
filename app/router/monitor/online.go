package monitor

import (
	hmonitor "thirteen-admin/app/handlers/monitor"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerOnlineRouter)
}

func registerOnlineRouter(g *gin.RouterGroup) {
	h := hmonitor.OnlineHandler{}
	r := g.Group("/monitor/online")
	{
		r.GET("/list", h.List)
		r.DELETE("/:tokenId", h.ForceLogout)
	}
}
