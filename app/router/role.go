package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterRoleRoutes 注册角色路由。
func RegisterRoleRoutes(h *handlers.RoleHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		sys := g.Group("/system/role", auth)
		sys.GET("/list", h.List)
		sys.GET("/authUser/allocatedList", h.AllocatedList)
		sys.GET("/authUser/unallocatedList", h.UnallocatedList)
		sys.GET("/deptTree/:roleId", h.DeptTree)
		sys.GET("/:roleId", h.Get)
		sys.POST("", h.Add)
		sys.POST("/export", h.Export)
		sys.PUT("", h.Edit)
		sys.PUT("/dataScope", h.DataScope)
		sys.PUT("/changeStatus", h.ChangeStatus)
		sys.PUT("/authUser/cancel", h.CancelAuthUser)
		sys.PUT("/authUser/cancelAll", h.CancelAuthUserAll)
		sys.PUT("/authUser/selectAll", h.SelectAuthUserAll)
		sys.DELETE("/:roleId", h.Delete)
	})
}
