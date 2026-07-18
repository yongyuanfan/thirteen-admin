package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterGenRoutes 注册代码生成路由。
func RegisterGenRoutes(h *handlers.GenHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		m := g.Group("/tool/gen", auth)
		m.GET("/list", h.List)
		m.GET("/db/list", h.DbList)
		m.GET("/preview/:tableId", h.Preview)
		m.GET("/genCode/:tableName", h.GenCode)
		m.GET("/synchDb/:tableName", h.SynchDb)
		m.GET("/batchGenCode", h.BatchGenCode)
		m.POST("/importTable", h.ImportTable)
		m.POST("/createTable", h.CreateTable)
		m.PUT("", h.Update)
		m.DELETE("/:tableId", h.Delete)
		m.GET("/:tableId", h.Get)
	})
}
