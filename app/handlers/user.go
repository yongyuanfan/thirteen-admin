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

// UserHandler 用户 HTTP 处理器。
type UserHandler struct {
	handler.Handler
	svc  *services.UserService
	file *services.FileService
}

// NewUserHandler 创建用户 Handler。
func NewUserHandler(svc *services.UserService, file *services.FileService) *UserHandler {
	return &UserHandler{svc: svc, file: file}
}

// List GET /system/user/list
func (h *UserHandler) List(ctx *gin.Context) {
	total, rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// Get GET /system/user/ 或 /system/user/:userId
func (h *UserHandler) Get(ctx *gin.Context) {
	idStr := ctx.Param("userId")
	var userID uint
	if idStr != "" && idStr != "/" {
		n, _ := strconv.ParseUint(idStr, 10, 64)
		userID = uint(n)
	}
	user, posts, roles, postIDs, roleIDs, err := h.svc.FormData(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "用户不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.OK(ctx, gin.H{
		"data":    user,
		"posts":   posts,
		"roles":   roles,
		"postIds": postIDs,
		"roleIds": roleIDs,
	})
}

type userBody struct {
	UserID      uint   `json:"userId"`
	DeptID      *uint  `json:"deptId"`
	UserName    string `json:"userName"`
	NickName    string `json:"nickName"`
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
	Sex         string `json:"sex"`
	Status      string `json:"status"`
	Password    string `json:"password"`
	Remark      string `json:"remark"`
	RoleIDs     []uint `json:"roleIds"`
	PostIDs     []uint `json:"postIds"`
}

func (h *UserHandler) bindUser(req userBody) *models.User {
	u := &models.User{
		DeptID:      req.DeptID,
		UserName:    req.UserName,
		NickName:    req.NickName,
		Email:       req.Email,
		Phonenumber: req.Phonenumber,
		Sex:         req.Sex,
		Status:      req.Status,
		Password:    req.Password,
		Remark:      req.Remark,
	}
	u.ID = req.UserID
	if u.Sex == "" {
		u.Sex = "0"
	}
	if u.Status == "" {
		u.Status = "0"
	}
	return u
}

// Add POST /system/user
func (h *UserHandler) Add(ctx *gin.Context) {
	var req userBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Add(h.bindUser(req), req.RoleIDs, req.PostIDs, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "新增成功")
}

// Edit PUT /system/user
func (h *UserHandler) Edit(ctx *gin.Context) {
	var req userBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Edit(h.bindUser(req), req.RoleIDs, req.PostIDs, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// Delete DELETE /system/user/:userId
func (h *UserHandler) Delete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("userId"))
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

type resetPwdBody struct {
	UserID   uint   `json:"userId"`
	Password string `json:"password"`
}

// ResetPwd PUT /system/user/resetPwd
func (h *UserHandler) ResetPwd(ctx *gin.Context) {
	var req resetPwdBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.ResetPwd(req.UserID, req.Password, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

type changeStatusBody struct {
	UserID uint   `json:"userId"`
	Status string `json:"status"`
}

// ChangeStatus PUT /system/user/changeStatus
func (h *UserHandler) ChangeStatus(ctx *gin.Context) {
	var req changeStatusBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.ChangeStatus(req.UserID, req.Status, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// Profile GET /system/user/profile
func (h *UserHandler) Profile(ctx *gin.Context) {
	user, roleGroup, postGroup, err := h.svc.Profile(services.ContextUserID(ctx))
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.OK(ctx, gin.H{
		"data":      user,
		"roleGroup": roleGroup,
		"postGroup": postGroup,
	})
}

// UpdateProfile PUT /system/user/profile
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	var req userBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.UpdateProfile(services.ContextUserID(ctx), req.NickName, req.Email, req.Phonenumber, req.Sex); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

type updatePwdBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// UpdatePwd PUT /system/user/profile/updatePwd
func (h *UserHandler) UpdatePwd(ctx *gin.Context) {
	var req updatePwdBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.UpdatePwd(services.ContextUserID(ctx), req.OldPassword, req.NewPassword); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// Avatar POST /system/user/profile/avatar
func (h *UserHandler) Avatar(ctx *gin.Context) {
	file, err := ctx.FormFile("avatarfile")
	if err != nil {
		h.BadRequest(ctx, "请选择头像文件")
		return
	}
	f, err := file.Open()
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	defer f.Close()
	_, _, urlPath, err := h.file.SaveUpload(file.Filename, f)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	if err := h.svc.UpdateAvatar(services.ContextUserID(ctx), urlPath); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.OK(ctx, gin.H{"imgUrl": urlPath})
}

// AuthRole GET /system/user/authRole/:userId
func (h *UserHandler) AuthRole(ctx *gin.Context) {
	userID, _ := strconv.ParseUint(ctx.Param("userId"), 10, 64)
	user, roles, err := h.svc.AuthRole(uint(userID))
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.OK(ctx, gin.H{
		"user":  user,
		"roles": roles,
	})
}

// SaveAuthRole PUT /system/user/authRole
func (h *UserHandler) SaveAuthRole(ctx *gin.Context) {
	userID, _ := strconv.ParseUint(ctx.Query("userId"), 10, 64)
	roleIDs := services.ParseUintIDs(ctx.Query("roleIds"))
	if userID == 0 {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.SaveAuthRole(uint(userID), roleIDs); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "授权成功")
}

// DeptTree GET /system/user/deptTree
func (h *UserHandler) DeptTree(ctx *gin.Context) {
	data, err := h.svc.DeptTree()
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// Export POST /system/user/export
func (h *UserHandler) Export(ctx *gin.Context) {
	rows, err := h.svc.Export(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	services.WriteCSV(ctx, "user.csv", []string{"用户编号", "登录名称", "用户名称", "部门", "手机号码", "状态", "创建时间"}, rows)
}

// ImportTemplate POST /system/user/importTemplate
func (h *UserHandler) ImportTemplate(ctx *gin.Context) {
	rows := h.svc.ImportTemplate()
	services.WriteCSV(ctx, "user_template.csv", rows[0], rows[1:])
}

// ImportData POST /system/user/importData
func (h *UserHandler) ImportData(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		h.BadRequest(ctx, "请选择导入文件")
		return
	}
	f, err := file.Open()
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	defer f.Close()
	updateSupport := ctx.DefaultQuery("updateSupport", "0") == "1"
	msg, err := h.svc.ImportData(f, updateSupport, services.ContextUserName(ctx))
	if err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, msg)
}
