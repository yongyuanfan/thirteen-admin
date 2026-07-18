package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterDeptRoutes 注册部门路由。
func RegisterDeptRoutes(h *handlers.DeptHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		sys := g.Group("/system/dept", auth)
		sys.GET("/list", h.List)
		sys.GET("/list/exclude/:deptId", h.ListExclude)
		sys.GET("/:deptId", h.Get)
		sys.POST("", h.Add)
		sys.PUT("", h.Edit)
		sys.PUT("/updateSort", h.UpdateSort)
		sys.DELETE("/:deptId", h.Delete)
	})
}
