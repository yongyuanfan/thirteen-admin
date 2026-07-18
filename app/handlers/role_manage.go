package handlers

import (
	"strconv"

	"clash-admin/app/services"

	"github.com/gin-gonic/gin"
)

type roleUpsertReq struct {
	ID       uint   `json:"id"`
	RoleName string `json:"roleName"`
	RoleCode string `json:"roleCode"`
	RoleDesc string `json:"roleDesc"`
	Status   string `json:"status"`
}

type roleMenusReq struct {
	RoleID  uint   `json:"roleId"`
	MenuIDs []uint `json:"menuIds"`
}

type roleHomeReq struct {
	RoleID uint   `json:"roleId"`
	Home   string `json:"home"`
}

type roleButtonsReq struct {
	RoleID    uint  `json:"roleId"`
	ButtonIDs []int `json:"buttonIds"`
}

// GetRoleList GET /systemManage/getRoleList
func (h *SystemManageHandler) GetRoleList(ctx *gin.Context) {
	current, _ := strconv.Atoi(ctx.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	data, err := h.roles.GetRoleList(ctx.Request.Context(), services.RoleListQuery{
		Current:  current,
		Size:     size,
		RoleName: ctx.Query("roleName"),
		RoleCode: ctx.Query("roleCode"),
		Status:   ctx.Query("status"),
	})
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// GetAllRoles GET /systemManage/getAllRoles
func (h *SystemManageHandler) GetAllRoles(ctx *gin.Context) {
	data, err := h.roles.GetAllRoles(ctx.Request.Context())
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// AddRole POST /systemManage/addRole
func (h *SystemManageHandler) AddRole(ctx *gin.Context) {
	var req roleUpsertReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	if err := h.roles.AddRole(ctx.Request.Context(), services.RoleUpsertInput{
		RoleName: req.RoleName,
		RoleCode: req.RoleCode,
		RoleDesc: req.RoleDesc,
		Status:   req.Status,
		Operator: h.operator(ctx),
	}); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// UpdateRole POST /systemManage/updateRole
func (h *SystemManageHandler) UpdateRole(ctx *gin.Context) {
	var req roleUpsertReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	if err := h.roles.UpdateRole(ctx.Request.Context(), services.RoleUpsertInput{
		ID:       req.ID,
		RoleName: req.RoleName,
		RoleCode: req.RoleCode,
		RoleDesc: req.RoleDesc,
		Status:   req.Status,
		Operator: h.operator(ctx),
	}); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// DeleteRole DELETE /systemManage/deleteRole
func (h *SystemManageHandler) DeleteRole(ctx *gin.Context) {
	var req idReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	if err := h.roles.DeleteRoles(ctx.Request.Context(), []uint{req.ID}); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// BatchDeleteRole DELETE /systemManage/batchDeleteRole
func (h *SystemManageHandler) BatchDeleteRole(ctx *gin.Context) {
	var req idsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	if err := h.roles.DeleteRoles(ctx.Request.Context(), req.IDs); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// GetRoleMenus GET /systemManage/getRoleMenus
func (h *SystemManageHandler) GetRoleMenus(ctx *gin.Context) {
	roleID, _ := strconv.ParseUint(ctx.Query("roleId"), 10, 64)
	data, err := h.roles.GetRoleMenus(ctx.Request.Context(), uint(roleID))
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// UpdateRoleMenus POST /systemManage/updateRoleMenus
func (h *SystemManageHandler) UpdateRoleMenus(ctx *gin.Context) {
	var req roleMenusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	if err := h.roles.UpdateRoleMenus(ctx.Request.Context(), req.RoleID, req.MenuIDs); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// GetRoleHome GET /systemManage/getRoleHome
func (h *SystemManageHandler) GetRoleHome(ctx *gin.Context) {
	roleID, _ := strconv.ParseUint(ctx.Query("roleId"), 10, 64)
	home, err := h.roles.GetRoleHome(ctx.Request.Context(), uint(roleID))
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, home)
}

// UpdateRoleHome POST /systemManage/updateRoleHome
func (h *SystemManageHandler) UpdateRoleHome(ctx *gin.Context) {
	var req roleHomeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	if err := h.roles.UpdateRoleHome(ctx.Request.Context(), req.RoleID, req.Home); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// GetAllButtons GET /systemManage/getAllButtons
func (h *SystemManageHandler) GetAllButtons(ctx *gin.Context) {
	data, err := h.roles.ListAllButtons(ctx.Request.Context())
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// GetRoleButtons GET /systemManage/getRoleButtons
func (h *SystemManageHandler) GetRoleButtons(ctx *gin.Context) {
	roleID, _ := strconv.ParseUint(ctx.Query("roleId"), 10, 64)
	data, err := h.roles.GetRoleButtons(ctx.Request.Context(), uint(roleID))
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// UpdateRoleButtons POST /systemManage/updateRoleButtons
func (h *SystemManageHandler) UpdateRoleButtons(ctx *gin.Context) {
	var req roleButtonsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	if err := h.roles.UpdateRoleButtons(ctx.Request.Context(), req.RoleID, req.ButtonIDs); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}
