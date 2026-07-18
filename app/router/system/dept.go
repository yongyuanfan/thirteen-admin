package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerDeptRouter)
}

func registerDeptRouter(g *gin.RouterGroup) {
	h := hsystem.DeptHandler{}
	r := g.Group("/system/dept")
	{
		r.GET("/list", h.List)
		r.GET("/list/exclude/:deptId", h.ListExclude)
		r.GET("/:deptId", h.Get)
		r.POST("", h.Add)
		r.PUT("", h.Edit)
		r.PUT("/updateSort", h.UpdateSort)
		r.DELETE("/:deptId", h.Delete)
	}
}
