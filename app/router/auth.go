package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

func init() {
	RegisterNoCheckRole(registerAuthNoCheckRole)
	RegisterCheckRole(registerAuthCheckRole)
}

func registerAuthNoCheckRole(g *gin.RouterGroup) {
	h := handlers.AuthHandler{}
	g.POST("/login", h.Login)
	g.POST("/register", h.Register)
	g.GET("/captchaImage", h.CaptchaImage)
	g.POST("/refresh", h.Refresh)
	g.POST("/logout", h.Logout)
}

func registerAuthCheckRole(g *gin.RouterGroup) {
	h := handlers.AuthHandler{}
	g.GET("/getInfo", h.GetInfo)
	g.POST("/unlockscreen", h.UnlockScreen)
	g.GET("/getRouters", h.GetRouters)
}
