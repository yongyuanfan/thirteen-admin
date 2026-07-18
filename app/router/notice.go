package router

import (
	"thirteen-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterNoticeRoutes 注册公告路由。
func RegisterNoticeRoutes(h *handlers.NoticeHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, auth gin.HandlerFunc) {
		sys := g.Group("/system/notice", auth)
		sys.GET("/list", h.List)
		sys.GET("/listTop", h.ListTop)
		sys.GET("/readUsers/list", h.ReadUsersList)
		sys.GET("/:noticeId", h.Get)
		sys.POST("", h.Add)
		sys.POST("/markRead", h.MarkRead)
		sys.POST("/markReadAll", h.MarkReadAll)
		sys.PUT("", h.Edit)
		sys.DELETE("/:noticeId", h.Delete)
	})
}
