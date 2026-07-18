package monitor

import (
	"net/url"
	"strings"

	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// CacheHandler 缓存监控 HTTP 处理器。
type CacheHandler struct {
	handler.Handler
	svc *services.CacheService
}

// NewCacheHandler 创建缓存 Handler。
func NewCacheHandler(svc *services.CacheService) *CacheHandler {
	return &CacheHandler{svc: svc}
}

// GetInfo GET /monitor/cache
func (h *CacheHandler) GetInfo(ctx *gin.Context) {
	data, err := h.svc.Overview(ctx.Request.Context())
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// GetNames GET /monitor/cache/getNames
func (h *CacheHandler) GetNames(ctx *gin.Context) {
	h.Success(ctx, h.svc.GetNames())
}

// GetKeys GET /monitor/cache/getKeys/:cacheName
func (h *CacheHandler) GetKeys(ctx *gin.Context) {
	keys, err := h.svc.GetKeys(ctx.Request.Context(), ctx.Param("cacheName"))
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, keys)
}

// GetValue GET /monitor/cache/getValue/:cacheName/:cacheKey
func (h *CacheHandler) GetValue(ctx *gin.Context) {
	cacheKey := strings.TrimPrefix(ctx.Param("cacheKey"), "/")
	data, err := h.svc.GetValue(ctx.Request.Context(), ctx.Param("cacheName"), cacheKey)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// ClearCacheName DELETE /monitor/cache/clearCacheName/:cacheName
func (h *CacheHandler) ClearCacheName(ctx *gin.Context) {
	if err := h.svc.ClearCacheName(ctx.Request.Context(), ctx.Param("cacheName")); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "清理成功")
}

// ClearCacheKey DELETE /monitor/cache/clearCacheKey/:cacheKey
func (h *CacheHandler) ClearCacheKey(ctx *gin.Context) {
	key, _ := url.PathUnescape(ctx.Param("cacheKey"))
	if err := h.svc.ClearCacheKey(ctx.Request.Context(), key); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "清理成功")
}

// ClearCacheAll DELETE /monitor/cache/clearCacheAll
func (h *CacheHandler) ClearCacheAll(ctx *gin.Context) {
	if err := h.svc.ClearCacheAll(ctx.Request.Context()); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "清理成功")
}
