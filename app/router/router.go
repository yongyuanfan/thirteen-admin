package router

import (
	"net/http"

	"thirteen-admin/app/middlewares"

	"github.com/gin-gonic/gin"
)

var (
	routerNoCheckRole = make([]func(*gin.RouterGroup), 0)
	routerCheckRole   = make([]func(*gin.RouterGroup), 0)
)

// RegisterNoCheckRole 注册无需认证的路由。
func RegisterNoCheckRole(fn func(*gin.RouterGroup)) {
	routerNoCheckRole = append(routerNoCheckRole, fn)
}

// RegisterCheckRole 注册需要认证的路由。
func RegisterCheckRole(fn func(*gin.RouterGroup)) {
	routerCheckRole = append(routerCheckRole, fn)
}

// InitRouter 初始化路由。
func InitRouter(engine *gin.Engine) {
	registerBaseRouter(engine)

	free := engine.Group("")
	for _, f := range routerNoCheckRole {
		f(free)
	}

	authed := engine.Group("", middlewares.Auth())
	for _, f := range routerCheckRole {
		f(authed)
	}
}

func registerBaseRouter(engine *gin.Engine) {
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "pong",
		})
	})
}
