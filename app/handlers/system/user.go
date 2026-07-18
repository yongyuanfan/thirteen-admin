package system

import (
	"errors"
	"strconv"

	"thirteen-admin/app/models"
	"thirteen-admin/app/runtime"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler 用户 HTTP 处理器。
type UserHandler struct {
	handler.Api
}

func (e *UserHandler) userService(c *gin.Context) (*services.UserService, error) {
	s := &services.UserService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		return nil, err
	}
	return s.WithDeps(), nil
}

// List GET /system/user/list
func (e UserHandler) List(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	total, rows, err := svc.List(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Table(c, total, rows)
}

// Get GET /system/user/ 或 /system/user/:userId
func (e UserHandler) Get(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	idStr := c.Param("userId")
	var userID uint
	if idStr != "" && idStr != "/" {
		n, _ := strconv.ParseUint(idStr, 10, 64)
		userID = uint(n)
	}
	user, posts, roles, postIDs, roleIDs, err := svc.FormData(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "用户不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.OK(c, gin.H{
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

func (e UserHandler) bindUser(req userBody) *models.User {
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
func (e UserHandler) Add(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	var req userBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := svc.Add(e.bindUser(req), req.RoleIDs, req.PostIDs, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "新增成功")
}

// Edit PUT /system/user
func (e UserHandler) Edit(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	var req userBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := svc.Edit(e.bindUser(req), req.RoleIDs, req.PostIDs, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// Delete DELETE /system/user/:userId
func (e UserHandler) Delete(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Param("userId"))
	if len(ids) == 0 {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := svc.Delete(ids); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "删除成功")
}

type resetPwdBody struct {
	UserID   uint   `json:"userId"`
	Password string `json:"password"`
}

// ResetPwd PUT /system/user/resetPwd
func (e UserHandler) ResetPwd(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	var req resetPwdBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := svc.ResetPwd(req.UserID, req.Password, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "修改成功")
}

type changeStatusBody struct {
	UserID uint   `json:"userId"`
	Status string `json:"status"`
}

// ChangeStatus PUT /system/user/changeStatus
func (e UserHandler) ChangeStatus(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	var req changeStatusBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := svc.ChangeStatus(req.UserID, req.Status, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// Profile GET /system/user/profile
func (e UserHandler) Profile(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	user, roleGroup, postGroup, err := svc.Profile(services.ContextUserID(c))
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.OK(c, gin.H{
		"data":      user,
		"roleGroup": roleGroup,
		"postGroup": postGroup,
	})
}

// UpdateProfile PUT /system/user/profile
func (e UserHandler) UpdateProfile(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	var req userBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := svc.UpdateProfile(services.ContextUserID(c), req.NickName, req.Email, req.Phonenumber, req.Sex); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "修改成功")
}

type updatePwdBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// UpdatePwd PUT /system/user/profile/updatePwd
func (e UserHandler) UpdatePwd(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	var req updatePwdBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := svc.UpdatePwd(services.ContextUserID(c), req.OldPassword, req.NewPassword); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// Avatar POST /system/user/profile/avatar
func (e UserHandler) Avatar(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	file, err := c.FormFile("avatarfile")
	if err != nil {
		e.BadRequest(c, "请选择头像文件")
		return
	}
	f, err := file.Open()
	if err != nil {
		e.InternalError(c, err)
		return
	}
	defer f.Close()
	_, _, urlPath, err := runtime.File().SaveUpload(file.Filename, f)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	if err := svc.UpdateAvatar(services.ContextUserID(c), urlPath); err != nil {
		e.InternalError(c, err)
		return
	}
	e.OK(c, gin.H{"imgUrl": urlPath})
}

// AuthRole GET /system/user/authRole/:userId
func (e UserHandler) AuthRole(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	userID, _ := strconv.ParseUint(c.Param("userId"), 10, 64)
	user, roles, err := svc.AuthRole(uint(userID))
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.OK(c, gin.H{
		"user":  user,
		"roles": roles,
	})
}

// SaveAuthRole PUT /system/user/authRole
func (e UserHandler) SaveAuthRole(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	userID, _ := strconv.ParseUint(c.Query("userId"), 10, 64)
	roleIDs := services.ParseUintIDs(c.Query("roleIds"))
	if userID == 0 {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := svc.SaveAuthRole(uint(userID), roleIDs); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "授权成功")
}

// DeptTree GET /system/user/deptTree
func (e UserHandler) DeptTree(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	data, err := svc.DeptTree()
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

// Export POST /system/user/export
func (e UserHandler) Export(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	rows, err := svc.Export(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	services.WriteCSV(c, "user.csv", []string{"用户编号", "登录名称", "用户名称", "部门", "手机号码", "状态", "创建时间"}, rows)
}

// ImportTemplate POST /system/user/importTemplate
func (e UserHandler) ImportTemplate(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	rows := svc.ImportTemplate()
	services.WriteCSV(c, "user_template.csv", rows[0], rows[1:])
}

// ImportData POST /system/user/importData
func (e UserHandler) ImportData(c *gin.Context) {
	svc, err := e.userService(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		e.BadRequest(c, "请选择导入文件")
		return
	}
	f, err := file.Open()
	if err != nil {
		e.InternalError(c, err)
		return
	}
	defer f.Close()
	updateSupport := c.DefaultQuery("updateSupport", "0") == "1"
	msg, err := svc.ImportData(f, updateSupport, services.ContextUserName(c))
	if err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, msg)
}
