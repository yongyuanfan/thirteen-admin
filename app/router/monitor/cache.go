package monitor

import (
	hmonitor "thirteen-admin/app/handlers/monitor"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerCacheRouter)
}

func registerCacheRouter(g *gin.RouterGroup) {
	h := hmonitor.CacheHandler{}
	r := g.Group("/monitor/cache")
	{
		r.GET("", h.GetInfo)
		r.GET("/getNames", h.GetNames)
		r.DELETE("/clearCacheAll", h.ClearCacheAll)
		r.DELETE("/clearCacheName/:cacheName", h.ClearCacheName)
		r.DELETE("/clearCacheKey/:cacheKey", h.ClearCacheKey)
		r.GET("/getKeys/:cacheName", h.GetKeys)
		r.GET("/getValue/:cacheName/*cacheKey", h.GetValue)
	}
}
