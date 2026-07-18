package monitor

import (
	"thirteen-admin/app/runtime"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// OnlineHandler 在线用户 HTTP 处理器。
type OnlineHandler struct {
	handler.Api
}

// List GET /monitor/online/list
func (e OnlineHandler) List(c *gin.Context) {
	svc := runtime.Online()
	users, err := svc.List(c.Request.Context(), c.Query("ipaddr"), c.Query("userName"))
	if err != nil {
		e.InternalError(c, err)
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
	e.Table(c, int64(len(rows)), rows)
}

// ForceLogout DELETE /monitor/online/:tokenId
func (e OnlineHandler) ForceLogout(c *gin.Context) {
	tokenID := c.Param("tokenId")
	if tokenID == "" {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := runtime.Online().Remove(c.Request.Context(), tokenID); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "强退成功")
}
