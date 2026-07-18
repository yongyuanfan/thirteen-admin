package router

import (
	"clash-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes 注册鉴权相关路由。
func RegisterAuthRoutes(h *handlers.AuthHandler) {
	RegisterFreeRouter(func(g *gin.RouterGroup) {
		group := g.Group("/auth")
		{
			group.POST("/login", h.Login)
			group.POST("/refreshToken", h.RefreshToken)
			group.POST("/sendCaptcha", h.SendCaptcha)
			group.POST("/verifyCaptcha", h.VerifyCaptcha)
			group.POST("/codeLogin", h.CodeLogin)
			group.POST("/register", h.Register)
			group.POST("/resetPassword", h.ResetPassword)
		}
	})

	RegisterAuthRouter(func(g *gin.RouterGroup, authMW gin.HandlerFunc) {
		group := g.Group("/auth")
		if authMW != nil {
			group.Use(authMW)
		}
		{
			group.GET("/getUserInfo", h.GetUserInfo)
		}
	})
}
