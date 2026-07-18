package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户路由。
func RegisterUserRoutes(h *handlers.UserHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		sys := g.Group("/system/user", auth)
		sys.GET("/list", h.List)
		sys.GET("/profile", h.Profile)
		sys.GET("/deptTree", h.DeptTree)
		sys.GET("/authRole/:userId", h.AuthRole)
		sys.GET("", h.Get)
		sys.GET("/:userId", h.Get)
		sys.POST("", h.Add)
		sys.POST("/export", h.Export)
		sys.POST("/importData", h.ImportData)
		sys.POST("/importTemplate", h.ImportTemplate)
		sys.POST("/profile/avatar", h.Avatar)
		sys.PUT("", h.Edit)
		sys.PUT("/resetPwd", h.ResetPwd)
		sys.PUT("/changeStatus", h.ChangeStatus)
		sys.PUT("/profile", h.UpdateProfile)
		sys.PUT("/profile/updatePwd", h.UpdatePwd)
		sys.PUT("/authRole", h.SaveAuthRole)
		sys.DELETE("/:userId", h.Delete)
	})
}
