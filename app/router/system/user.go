package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerUserRouter)
}

func registerUserRouter(g *gin.RouterGroup) {
	h := hsystem.UserHandler{}
	r := g.Group("/system/user")
	{
		r.GET("/list", h.List)
		r.GET("/profile", h.Profile)
		r.GET("/deptTree", h.DeptTree)
		r.GET("/authRole/:userId", h.AuthRole)
		r.GET("", h.Get)
		r.GET("/:userId", h.Get)
		r.POST("", h.Add)
		r.POST("/export", h.Export)
		r.POST("/importData", h.ImportData)
		r.POST("/importTemplate", h.ImportTemplate)
		r.POST("/profile/avatar", h.Avatar)
		r.PUT("", h.Edit)
		r.PUT("/resetPwd", h.ResetPwd)
		r.PUT("/changeStatus", h.ChangeStatus)
		r.PUT("/profile", h.UpdateProfile)
		r.PUT("/profile/updatePwd", h.UpdatePwd)
		r.PUT("/authRole", h.SaveAuthRole)
		r.DELETE("/:userId", h.Delete)
	}
}
