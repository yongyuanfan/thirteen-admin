package handlers

import (
	"errors"
	"strconv"

	"thirteen-admin/app/models"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoleHandler 角色 HTTP 处理器。
type RoleHandler struct {
	handler.Handler
	svc *services.RoleService
}

// NewRoleHandler 创建角色 Handler。
func NewRoleHandler(svc *services.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

// List GET /system/role/list
func (h *RoleHandler) List(ctx *gin.Context) {
	total, rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// Get GET /system/role/:roleId
func (h *RoleHandler) Get(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("roleId"), 10, 64)
	data, err := h.svc.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "角色不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

type roleBody struct {
	RoleID            uint   `json:"roleId"`
	RoleName          string `json:"roleName"`
	RoleKey           string `json:"roleKey"`
	RoleSort          int    `json:"roleSort"`
	DataScope         string `json:"dataScope"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly"`
	DeptCheckStrictly bool   `json:"deptCheckStrictly"`
	Status            string `json:"status"`
	Remark            string `json:"remark"`
	MenuIDs           []uint `json:"menuIds"`
	DeptIDs           []uint `json:"deptIds"`
}

func (h *RoleHandler) bindRole(req roleBody) *models.Role {
	role := &models.Role{
		RoleName:          req.RoleName,
		RoleKey:           req.RoleKey,
		RoleSort:          req.RoleSort,
		DataScope:         req.DataScope,
		MenuCheckStrictly: req.MenuCheckStrictly,
		DeptCheckStrictly: req.DeptCheckStrictly,
		Status:            req.Status,
		Remark:            req.Remark,
	}
	role.ID = req.RoleID
	if role.DataScope == "" {
		role.DataScope = "1"
	}
	if role.Status == "" {
		role.Status = "0"
	}
	return role
}

// Add POST /system/role
func (h *RoleHandler) Add(ctx *gin.Context) {
	var req roleBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Add(h.bindRole(req), req.MenuIDs, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "新增成功")
}

// Edit PUT /system/role
func (h *RoleHandler) Edit(ctx *gin.Context) {
	var req roleBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Edit(h.bindRole(req), req.MenuIDs, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// DataScope PUT /system/role/dataScope
func (h *RoleHandler) DataScope(ctx *gin.Context) {
	var req roleBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.DataScope(h.bindRole(req), req.DeptIDs, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

type roleStatusBody struct {
	RoleID uint   `json:"roleId"`
	Status string `json:"status"`
}

// ChangeStatus PUT /system/role/changeStatus
func (h *RoleHandler) ChangeStatus(ctx *gin.Context) {
	var req roleStatusBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.ChangeStatus(req.RoleID, req.Status, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// Delete DELETE /system/role/:roleId
func (h *RoleHandler) Delete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("roleId"))
	if len(ids) == 0 {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Delete(ids); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "删除成功")
}

// AllocatedList GET /system/role/authUser/allocatedList
func (h *RoleHandler) AllocatedList(ctx *gin.Context) {
	total, rows, err := h.svc.AllocatedList(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// UnallocatedList GET /system/role/authUser/unallocatedList
func (h *RoleHandler) UnallocatedList(ctx *gin.Context) {
	total, rows, err := h.svc.UnallocatedList(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

type authUserBody struct {
	UserID uint `json:"userId"`
	RoleID uint `json:"roleId"`
}

// CancelAuthUser PUT /system/role/authUser/cancel
func (h *RoleHandler) CancelAuthUser(ctx *gin.Context) {
	var req authUserBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.CancelAuthUser(req.UserID, req.RoleID); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "取消授权成功")
}

// CancelAuthUserAll PUT /system/role/authUser/cancelAll
func (h *RoleHandler) CancelAuthUserAll(ctx *gin.Context) {
	roleID, _ := strconv.ParseUint(ctx.Query("roleId"), 10, 64)
	userIDs := services.ParseUintIDs(ctx.Query("userIds"))
	if roleID == 0 || len(userIDs) == 0 {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.CancelAuthUserAll(uint(roleID), userIDs); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "取消授权成功")
}

// SelectAuthUserAll PUT /system/role/authUser/selectAll
func (h *RoleHandler) SelectAuthUserAll(ctx *gin.Context) {
	roleID, _ := strconv.ParseUint(ctx.Query("roleId"), 10, 64)
	userIDs := services.ParseUintIDs(ctx.Query("userIds"))
	if roleID == 0 || len(userIDs) == 0 {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.SelectAuthUserAll(uint(roleID), userIDs); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "授权成功")
}

// DeptTree GET /system/role/deptTree/:roleId
func (h *RoleHandler) DeptTree(ctx *gin.Context) {
	roleID, _ := strconv.ParseUint(ctx.Param("roleId"), 10, 64)
	depts, checked, err := h.svc.DeptTree(uint(roleID))
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.OK(ctx, gin.H{
		"depts":       depts,
		"checkedKeys": checked,
	})
}

// Export POST /system/role/export
func (h *RoleHandler) Export(ctx *gin.Context) {
	rows, err := h.svc.Export(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	services.WriteCSV(ctx, "role.csv", []string{"角色编号", "角色名称", "权限字符", "显示顺序", "状态", "创建时间"}, rows)
}
