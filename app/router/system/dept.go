package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

// RegisterDeptRoutes 注册部门路由。
func RegisterDeptRoutes(h *hsystem.DeptHandler) {
	router.RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
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
