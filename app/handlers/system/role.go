package system

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
	handler.Api
}

// List GET /system/role/list
func (e RoleHandler) List(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	total, rows, err := s.List(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Table(c, total, rows)
}

// Get GET /system/role/:roleId
func (e RoleHandler) Get(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	id, _ := strconv.ParseUint(c.Param("roleId"), 10, 64)
	data, err := s.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "角色不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
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

func (e RoleHandler) bindRole(req roleBody) *models.Role {
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
func (e RoleHandler) Add(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	var req roleBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.Add(e.bindRole(req), req.MenuIDs, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "新增成功")
}

// Edit PUT /system/role
func (e RoleHandler) Edit(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	var req roleBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.Edit(e.bindRole(req), req.MenuIDs, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// DataScope PUT /system/role/dataScope
func (e RoleHandler) DataScope(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	var req roleBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.DataScope(e.bindRole(req), req.DeptIDs, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "修改成功")
}

type roleStatusBody struct {
	RoleID uint   `json:"roleId"`
	Status string `json:"status"`
}

// ChangeStatus PUT /system/role/changeStatus
func (e RoleHandler) ChangeStatus(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	var req roleStatusBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.ChangeStatus(req.RoleID, req.Status, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// Delete DELETE /system/role/:roleId
func (e RoleHandler) Delete(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	ids := services.ParseUintIDs(c.Param("roleId"))
	if len(ids) == 0 {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.Delete(ids); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "删除成功")
}

// AllocatedList GET /system/role/authUser/allocatedList
func (e RoleHandler) AllocatedList(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	total, rows, err := s.AllocatedList(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Table(c, total, rows)
}

// UnallocatedList GET /system/role/authUser/unallocatedList
func (e RoleHandler) UnallocatedList(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	total, rows, err := s.UnallocatedList(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Table(c, total, rows)
}

type authUserBody struct {
	UserID uint `json:"userId"`
	RoleID uint `json:"roleId"`
}

// CancelAuthUser PUT /system/role/authUser/cancel
func (e RoleHandler) CancelAuthUser(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	var req authUserBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.CancelAuthUser(req.UserID, req.RoleID); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "取消授权成功")
}

// CancelAuthUserAll PUT /system/role/authUser/cancelAll
func (e RoleHandler) CancelAuthUserAll(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	roleID, _ := strconv.ParseUint(c.Query("roleId"), 10, 64)
	userIDs := services.ParseUintIDs(c.Query("userIds"))
	if roleID == 0 || len(userIDs) == 0 {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.CancelAuthUserAll(uint(roleID), userIDs); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "取消授权成功")
}

// SelectAuthUserAll PUT /system/role/authUser/selectAll
func (e RoleHandler) SelectAuthUserAll(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	roleID, _ := strconv.ParseUint(c.Query("roleId"), 10, 64)
	userIDs := services.ParseUintIDs(c.Query("userIds"))
	if roleID == 0 || len(userIDs) == 0 {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.SelectAuthUserAll(uint(roleID), userIDs); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "授权成功")
}

// DeptTree GET /system/role/deptTree/:roleId
func (e RoleHandler) DeptTree(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	roleID, _ := strconv.ParseUint(c.Param("roleId"), 10, 64)
	depts, checked, err := s.DeptTree(uint(roleID))
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.OK(c, gin.H{
		"depts":       depts,
		"checkedKeys": checked,
	})
}

// Export POST /system/role/export
func (e RoleHandler) Export(c *gin.Context) {
	s := services.RoleService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	s.WithDeps()
	rows, err := s.Export(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	services.WriteCSV(c, "role.csv", []string{"角色编号", "角色名称", "权限字符", "显示顺序", "状态", "创建时间"}, rows)
}
