package handlers

import (
	"strconv"

	"clash-admin/app/services"

	"github.com/gin-gonic/gin"
)

type userUpsertReq struct {
	ID         uint     `json:"id"`
	UserName   string   `json:"userName"`
	UserGender *string  `json:"userGender"`
	NickName   string   `json:"nickName"`
	UserPhone  string   `json:"userPhone"`
	UserEmail  string   `json:"userEmail"`
	UserRoles  []string `json:"userRoles"`
	Status     string   `json:"status"`
}

type idReq struct {
	ID uint `json:"id"`
}

type idsReq struct {
	IDs []uint `json:"ids"`
}

// GetUserList GET /systemManage/getUserList
func (h *SystemManageHandler) GetUserList(ctx *gin.Context) {
	current, _ := strconv.Atoi(ctx.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	q := services.UserListQuery{
		Current:    current,
		Size:       size,
		UserName:   ctx.Query("userName"),
		UserGender: ctx.Query("userGender"),
		NickName:   ctx.Query("nickName"),
		UserPhone:  ctx.Query("userPhone"),
		UserEmail:  ctx.Query("userEmail"),
		Status:     ctx.Query("status"),
	}
	data, err := h.users.GetUserList(ctx.Request.Context(), q)
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// AddUser POST /systemManage/addUser
func (h *SystemManageHandler) AddUser(ctx *gin.Context) {
	var req userUpsertReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	err := h.users.AddUser(ctx.Request.Context(), services.UserUpsertInput{
		UserName:   req.UserName,
		UserGender: req.UserGender,
		NickName:   req.NickName,
		UserPhone:  req.UserPhone,
		UserEmail:  req.UserEmail,
		UserRoles:  req.UserRoles,
		Status:     req.Status,
		Operator:   h.operator(ctx),
	})
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// UpdateUser POST /systemManage/updateUser
func (h *SystemManageHandler) UpdateUser(ctx *gin.Context) {
	var req userUpsertReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	err := h.users.UpdateUser(ctx.Request.Context(), services.UserUpsertInput{
		ID:         req.ID,
		UserName:   req.UserName,
		UserGender: req.UserGender,
		NickName:   req.NickName,
		UserPhone:  req.UserPhone,
		UserEmail:  req.UserEmail,
		UserRoles:  req.UserRoles,
		Status:     req.Status,
		Operator:   h.operator(ctx),
	})
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// DeleteUser DELETE /systemManage/deleteUser
func (h *SystemManageHandler) DeleteUser(ctx *gin.Context) {
	var req idReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	uid, ok := h.currentUserID(ctx)
	if !ok {
		h.Fail(ctx, "8888", "未登录")
		return
	}
	if err := h.users.DeleteUsers(ctx.Request.Context(), []uint{req.ID}, uid); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// BatchDeleteUser DELETE /systemManage/batchDeleteUser
func (h *SystemManageHandler) BatchDeleteUser(ctx *gin.Context) {
	var req idsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	uid, ok := h.currentUserID(ctx)
	if !ok {
		h.Fail(ctx, "8888", "未登录")
		return
	}
	if err := h.users.DeleteUsers(ctx.Request.Context(), req.IDs, uid); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}
