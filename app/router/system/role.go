package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerRoleRouter)
}

func registerRoleRouter(g *gin.RouterGroup) {
	h := hsystem.RoleHandler{}
	r := g.Group("/system/role")
	{
		r.GET("/list", h.List)
		r.GET("/authUser/allocatedList", h.AllocatedList)
		r.GET("/authUser/unallocatedList", h.UnallocatedList)
		r.GET("/deptTree/:roleId", h.DeptTree)
		r.GET("/:roleId", h.Get)
		r.POST("", h.Add)
		r.POST("/export", h.Export)
		r.PUT("", h.Edit)
		r.PUT("/dataScope", h.DataScope)
		r.PUT("/changeStatus", h.ChangeStatus)
		r.PUT("/authUser/cancel", h.CancelAuthUser)
		r.PUT("/authUser/cancelAll", h.CancelAuthUserAll)
		r.PUT("/authUser/selectAll", h.SelectAuthUserAll)
		r.DELETE("/:roleId", h.Delete)
	}
}
