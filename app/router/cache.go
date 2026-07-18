package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterCacheRoutes 注册缓存监控路由。
func RegisterCacheRoutes(h *handlers.CacheHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		m := g.Group("/monitor/cache", auth)
		m.GET("", h.GetInfo)
		m.GET("/getNames", h.GetNames)
		m.DELETE("/clearCacheAll", h.ClearCacheAll)
		m.DELETE("/clearCacheName/:cacheName", h.ClearCacheName)
		m.DELETE("/clearCacheKey/:cacheKey", h.ClearCacheKey)
		m.GET("/getKeys/:cacheName", h.GetKeys)
		m.GET("/getValue/:cacheName/*cacheKey", h.GetValue)
	})
}
