package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	freeRouterGroups = make([]func(*gin.RouterGroup), 0)
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

func registerBaseRouter(engine *gin.Engine) {
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "pong",
		})
	})
}

func registerBusinessRouter(engine *gin.Engine, authMiddleware gin.HandlerFunc) {
	root := engine.Group("")
	for _, f := range freeRouterGroups {
		f(root)
	}
	for _, f := range authRouterGroups {
		f(root, authMiddleware)
	}
}
