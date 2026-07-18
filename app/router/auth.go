package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes 注册鉴权相关路由。
func RegisterAuthRoutes(h *handlers.AuthHandler) {
	RegisterFreeRouter(func(g *gin.RouterGroup) {
		g.POST("/login", h.Login)
		g.POST("/register", h.Register)
		g.GET("/captchaImage", h.CaptchaImage)
	})
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		g.GET("/getInfo", auth, h.GetInfo)
		g.POST("/logout", auth, h.Logout)
		g.POST("/unlockscreen", auth, h.UnlockScreen)
		g.GET("/getRouters", auth, h.GetRouters)
	})
}
