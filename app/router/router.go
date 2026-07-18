package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// 无需认证的路由
	freeRouterGroups = make([]func(*gin.RouterGroup), 0)

	// 需要认证的路由
	authRouterGroups = make([]func(*gin.RouterGroup, gin.HandlerFunc), 0)
)

// RegisterFreeRouter 注册无需认证的路由组。
func RegisterFreeRouter(fn func(*gin.RouterGroup)) {
	freeRouterGroups = append(freeRouterGroups, fn)
}

// RegisterAuthRouter 注册需要认证的路由组。
func RegisterAuthRouter(fn func(*gin.RouterGroup, gin.HandlerFunc)) {
	authRouterGroups = append(authRouterGroups, fn)
}

// InitRouter 初始化路由。
func InitRouter(engine *gin.Engine, authMiddleware gin.HandlerFunc) {
	registerBaseRouter(engine)
	registerBusinessRouter(engine, authMiddleware)
}

// registerBaseRouter 注册基础路由
func registerBaseRouter(engine *gin.Engine) {
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "pong",
		})
	})
}

// registerBusinessRouter 注册业务路由
func registerBusinessRouter(engine *gin.Engine, authMiddleware gin.HandlerFunc) {
	v1 := engine.Group("/api/v1")

	for _, f := range freeRouterGroups {
		f(v1)
	}

	for _, f := range authRouterGroups {
		f(v1, authMiddleware)
	}
}
