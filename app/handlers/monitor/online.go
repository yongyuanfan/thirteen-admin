package monitor

import (
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// OnlineHandler 在线用户 HTTP 处理器。
type OnlineHandler struct {
	handler.Handler
	svc *services.OnlineService
}

// NewOnlineHandler 创建在线用户 Handler。
func NewOnlineHandler(svc *services.OnlineService) *OnlineHandler {
	return &OnlineHandler{svc: svc}
}

// List GET /monitor/online/list
func (h *OnlineHandler) List(ctx *gin.Context) {
	users, err := h.svc.List(ctx.Request.Context(), ctx.Query("ipaddr"), ctx.Query("userName"))
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	rows := make([]map[string]any, 0, len(users))
	for _, u := range users {
		rows = append(rows, map[string]any{
			"tokenId":       u.TokenID,
			"deptName":      u.DeptName,
			"userName":      u.UserName,
			"ipaddr":        u.IPAddr,
			"loginLocation": u.LoginLocation,
			"browser":       u.Browser,
			"os":            u.OS,
			"loginTime":     u.LoginTime,
		})
	}
	h.Table(ctx, int64(len(rows)), rows)
}

// ForceLogout DELETE /monitor/online/:tokenId
func (h *OnlineHandler) ForceLogout(ctx *gin.Context) {
	tokenID := ctx.Param("tokenId")
	if tokenID == "" {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Remove(ctx.Request.Context(), tokenID); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "强退成功")
}
