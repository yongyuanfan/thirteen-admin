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

// DeptHandler 部门 HTTP 处理器。
type DeptHandler struct {
	handler.Handler
	svc *services.DeptService
}

// NewDeptHandler 创建部门 Handler。
func NewDeptHandler(svc *services.DeptService) *DeptHandler {
	return &DeptHandler{svc: svc}
}

// List GET /system/dept/list
func (h *DeptHandler) List(ctx *gin.Context) {
	rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, rows)
}

// ListExclude GET /system/dept/list/exclude/:deptId
func (h *DeptHandler) ListExclude(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("deptId"), 10, 64)
	rows, err := h.svc.ListExclude(uint(id))
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, rows)
}

// Get GET /system/dept/:deptId
func (h *DeptHandler) Get(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("deptId"), 10, 64)
	data, err := h.svc.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "部门不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
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
func (h *DeptHandler) Add(ctx *gin.Context) {
	var req deptBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
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
	if err := h.svc.Add(dept, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "新增成功")
}

// Edit PUT /system/dept
func (h *DeptHandler) Edit(ctx *gin.Context) {
	var req deptBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
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
	if err := h.svc.Edit(dept, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

type sortBody struct {
	DeptIDs   string `json:"deptIds"`
	OrderNums string `json:"orderNums"`
}

// UpdateSort PUT /system/dept/updateSort
func (h *DeptHandler) UpdateSort(ctx *gin.Context) {
	var req sortBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.UpdateSort(req.DeptIDs, req.OrderNums); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// Delete DELETE /system/dept/:deptId
func (h *DeptHandler) Delete(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("deptId"), 10, 64)
	if err := h.svc.Delete(uint(id)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "删除成功")
}
