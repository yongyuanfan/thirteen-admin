package handlers

import (
	"errors"
	"strconv"
	"strings"

	"clash-admin/app/middlewares"
	"clash-admin/app/services"
	"clash-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// SystemManageHandler 系统管理 HTTP 处理器。
type SystemManageHandler struct {
	handler.Handler
	users *services.UserManageService
	roles *services.RoleManageService
	menus *services.MenuManageService
}

// NewSystemManageHandler 创建系统管理 Handler。
func NewSystemManageHandler(
	users *services.UserManageService,
	roles *services.RoleManageService,
	menus *services.MenuManageService,
) *SystemManageHandler {
	return &SystemManageHandler{users: users, roles: roles, menus: menus}
}

func (h *SystemManageHandler) operator(ctx *gin.Context) string {
	if v, ok := ctx.Get(middlewares.ContextUserNameKey); ok {
		if name, ok := v.(string); ok && name != "" {
			return name
		}
	}
	return "unknown"
}

func (h *SystemManageHandler) currentUserID(ctx *gin.Context) (uint, bool) {
	v, ok := ctx.Get(middlewares.ContextUserIDKey)
	if !ok {
		return 0, false
	}
	switch id := v.(type) {
	case string:
		n, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return 0, false
		}
		return uint(n), true
	case uint:
		return id, true
	case uint64:
		return uint(id), true
	default:
		return 0, false
	}
}

func (h *SystemManageHandler) failManage(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrManageInvalid):
		h.BadRequest(ctx, err.Error())
	case errors.Is(err, services.ErrManageNotFound):
		h.NotFound(ctx, err.Error())
	case errors.Is(err, services.ErrManageConflict):
		h.Conflict(ctx, err.Error())
	case errors.Is(err, services.ErrManageForbidden), errors.Is(err, services.ErrManageRoleInUse):
		h.Fail(ctx, handler.CodeBadRequest, err.Error())
	default:
		msg := err.Error()
		if strings.Contains(msg, "不能为空") || strings.Contains(msg, "仅支持") || strings.Contains(msg, "无效") || strings.Contains(msg, "已存在") {
			h.BadRequest(ctx, msg)
			return
		}
		h.InternalError(ctx, err)
	}
}
