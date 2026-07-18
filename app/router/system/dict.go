package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerDictRouter)
}

func registerDictRouter(g *gin.RouterGroup) {
	h := hsystem.DictHandler{}
	typ := g.Group("/system/dict/type")
	{
		typ.GET("/list", h.TypeList)
		typ.GET("/optionselect", h.TypeOptionSelect)
		typ.DELETE("/refreshCache", h.TypeRefreshCache)
		typ.GET("/:dictId", h.TypeGet)
		typ.POST("", h.TypeAdd)
		typ.POST("/export", h.TypeExport)
		typ.PUT("", h.TypeEdit)
		typ.DELETE("/:dictId", h.TypeDelete)
	}

	data := g.Group("/system/dict/data")
	{
		data.GET("/list", h.DataList)
		data.GET("/type/:dictType", h.DataByType)
		data.GET("/:dictCode", h.DataGet)
		data.POST("", h.DataAdd)
		data.POST("/export", h.DataExport)
		data.PUT("", h.DataEdit)
		data.DELETE("/:dictCode", h.DataDelete)
	}
}
