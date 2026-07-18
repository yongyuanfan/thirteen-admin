package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

// RegisterMenuRoutes 注册菜单路由。
func RegisterMenuRoutes(h *hsystem.MenuHandler) {
	router.RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		sys := g.Group("/system/menu", auth)
		sys.GET("/list", h.List)
		sys.GET("/treeselect", h.TreeSelect)
		sys.GET("/roleMenuTreeselect/:roleId", h.RoleMenuTreeSelect)
		sys.GET("/:menuId", h.Get)
		sys.POST("", h.Add)
		sys.PUT("", h.Edit)
		sys.PUT("/updateSort", h.UpdateSort)
		sys.DELETE("/:menuId", h.Delete)
	})
}
