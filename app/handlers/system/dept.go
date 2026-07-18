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

// DeptHandler 部门 HTTP 处理器。
type DeptHandler struct {
	handler.Api
}

// List GET /system/dept/list
func (e DeptHandler) List(c *gin.Context) {
	s := services.DeptService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	rows, err := s.List(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Success(c, rows)
}

// ListExclude GET /system/dept/list/exclude/:deptId
func (e DeptHandler) ListExclude(c *gin.Context) {
	s := services.DeptService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("deptId"), 10, 64)
	rows, err := s.ListExclude(uint(id))
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Success(c, rows)
}

// Get GET /system/dept/:deptId
func (e DeptHandler) Get(c *gin.Context) {
	s := services.DeptService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("deptId"), 10, 64)
	data, err := s.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "部门不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

type deptBody struct {
	DeptID   uint   `json:"deptId"`
	ParentID uint   `json:"parentId"`
	DeptName string `json:"deptName"`
	OrderNum int    `json:"orderNum"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

// Add POST /system/dept
func (e DeptHandler) Add(c *gin.Context) {
	s := services.DeptService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req deptBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	dept := &models.Dept{
		ParentID: req.ParentID,
		DeptName: req.DeptName,
		OrderNum: req.OrderNum,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   req.Status,
	}
	if dept.Status == "" {
		dept.Status = "0"
	}
	if err := s.Add(dept, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "新增成功")
}

// Edit PUT /system/dept
func (e DeptHandler) Edit(c *gin.Context) {
	s := services.DeptService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req deptBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	dept := &models.Dept{
		ParentID: req.ParentID,
		DeptName: req.DeptName,
		OrderNum: req.OrderNum,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   req.Status,
	}
	dept.ID = req.DeptID
	if err := s.Edit(dept, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "修改成功")
}

type sortBody struct {
	DeptIDs   string `json:"deptIds"`
	OrderNums string `json:"orderNums"`
}

// UpdateSort PUT /system/dept/updateSort
func (e DeptHandler) UpdateSort(c *gin.Context) {
	s := services.DeptService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req sortBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.UpdateSort(req.DeptIDs, req.OrderNums); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// Delete DELETE /system/dept/:deptId
func (e DeptHandler) Delete(c *gin.Context) {
	s := services.DeptService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("deptId"), 10, 64)
	if err := s.Delete(uint(id)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "删除成功")
}
