package system

import (
	hsystem "thirteen-admin/app/handlers/system"
	"thirteen-admin/app/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.RegisterCheckRole(registerNoticeRouter)
}

func registerNoticeRouter(g *gin.RouterGroup) {
	h := hsystem.NoticeHandler{}
	r := g.Group("/system/notice")
	{
		r.GET("/list", h.List)
		r.GET("/listTop", h.ListTop)
		r.GET("/readUsers/list", h.ReadUsersList)
		r.GET("/:noticeId", h.Get)
		r.POST("", h.Add)
		r.POST("/markRead", h.MarkRead)
		r.POST("/markReadAll", h.MarkReadAll)
		r.PUT("", h.Edit)
		r.DELETE("/:noticeId", h.Delete)
	}
}
