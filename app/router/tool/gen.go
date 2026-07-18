package tool

import (
	htool "thirteen-admin/app/handlers/tool"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerGenRouter)
}

func registerGenRouter(g *gin.RouterGroup) {
	h := htool.GenHandler{}
	r := g.Group("/tool/gen")
	{
		r.GET("/list", h.List)
		r.GET("/db/list", h.DbList)
		r.GET("/preview/:tableId", h.Preview)
		r.GET("/genCode/:tableName", h.GenCode)
		r.GET("/synchDb/:tableName", h.SynchDb)
		r.GET("/batchGenCode", h.BatchGenCode)
		r.POST("/importTable", h.ImportTable)
		r.POST("/createTable", h.CreateTable)
		r.PUT("", h.Update)
		r.DELETE("/:tableId", h.Delete)
		r.GET("/:tableId", h.Get)
	}
}
