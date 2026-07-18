package monitor

import (
	"net/url"
	"strings"

	"thirteen-admin/app/runtime"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// CacheHandler 缓存监控 HTTP 处理器。
type CacheHandler struct {
	handler.Api
}

func (e CacheHandler) cacheSvc() *services.CacheService {
	return services.NewCacheService(runtime.Cache())
}

// GetInfo GET /monitor/cache
func (e CacheHandler) GetInfo(c *gin.Context) {
	data, err := e.cacheSvc().Overview(c.Request.Context())
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

// GetNames GET /monitor/cache/getNames
func (e CacheHandler) GetNames(c *gin.Context) {
	e.Success(c, e.cacheSvc().GetNames())
}

// GetKeys GET /monitor/cache/getKeys/:cacheName
func (e CacheHandler) GetKeys(c *gin.Context) {
	keys, err := e.cacheSvc().GetKeys(c.Request.Context(), c.Param("cacheName"))
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Success(c, keys)
}

// GetValue GET /monitor/cache/getValue/:cacheName/:cacheKey
func (e CacheHandler) GetValue(c *gin.Context) {
	cacheKey := strings.TrimPrefix(c.Param("cacheKey"), "/")
	data, err := e.cacheSvc().GetValue(c.Request.Context(), c.Param("cacheName"), cacheKey)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

// ClearCacheName DELETE /monitor/cache/clearCacheName/:cacheName
func (e CacheHandler) ClearCacheName(c *gin.Context) {
	if err := e.cacheSvc().ClearCacheName(c.Request.Context(), c.Param("cacheName")); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "清理成功")
}

// ClearCacheKey DELETE /monitor/cache/clearCacheKey/:cacheKey
func (e CacheHandler) ClearCacheKey(c *gin.Context) {
	key, _ := url.PathUnescape(c.Param("cacheKey"))
	if err := e.cacheSvc().ClearCacheKey(c.Request.Context(), key); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "清理成功")
}

// ClearCacheAll DELETE /monitor/cache/clearCacheAll
func (e CacheHandler) ClearCacheAll(c *gin.Context) {
	if err := e.cacheSvc().ClearCacheAll(c.Request.Context()); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "清理成功")
}
