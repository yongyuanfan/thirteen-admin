package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerMenuRouter)
}

func registerMenuRouter(g *gin.RouterGroup) {
	h := hsystem.MenuHandler{}
	r := g.Group("/system/menu")
	{
		r.GET("/list", h.List)
		r.GET("/treeselect", h.TreeSelect)
		r.GET("/roleMenuTreeselect/:roleId", h.RoleMenuTreeSelect)
		r.GET("/:menuId", h.Get)
		r.POST("", h.Add)
		r.PUT("", h.Edit)
		r.PUT("/updateSort", h.UpdateSort)
		r.DELETE("/:menuId", h.Delete)
	}
}
