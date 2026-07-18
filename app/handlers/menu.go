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

// MenuHandler 菜单 HTTP 处理器。
type MenuHandler struct {
	handler.Handler
	svc *services.MenuService
}

// NewMenuHandler 创建菜单 Handler。
func NewMenuHandler(svc *services.MenuService) *MenuHandler {
	return &MenuHandler{svc: svc}
}

// List GET /system/menu/list
func (h *MenuHandler) List(ctx *gin.Context) {
	rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, rows)
}

// Get GET /system/menu/:menuId
func (h *MenuHandler) Get(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("menuId"), 10, 64)
	data, err := h.svc.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "菜单不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// TreeSelect GET /system/menu/treeselect
func (h *MenuHandler) TreeSelect(ctx *gin.Context) {
	data, err := h.svc.TreeSelect()
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// RoleMenuTreeSelect GET /system/menu/roleMenuTreeselect/:roleId
func (h *MenuHandler) RoleMenuTreeSelect(ctx *gin.Context) {
	roleID, _ := strconv.ParseUint(ctx.Param("roleId"), 10, 64)
	menus, checked, err := h.svc.RoleMenuTreeSelect(uint(roleID))
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.OK(ctx, gin.H{
		"menus":       menus,
		"checkedKeys": checked,
	})
}

type menuBody struct {
	MenuID    uint   `json:"menuId"`
	MenuName  string `json:"menuName"`
	ParentID  uint   `json:"parentId"`
	OrderNum  int    `json:"orderNum"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Query     string `json:"query"`
	RouteName string `json:"routeName"`
	IsFrame   any    `json:"isFrame"`
	IsCache   any    `json:"isCache"`
	MenuType  string `json:"menuType"`
	Visible   string `json:"visible"`
	Status    string `json:"status"`
	Perms     string `json:"perms"`
	Icon      string `json:"icon"`
	Remark    string `json:"remark"`
}

func (h *MenuHandler) bindMenu(req menuBody) *models.Menu {
	menu := &models.Menu{
		MenuName:  req.MenuName,
		ParentID:  req.ParentID,
		OrderNum:  req.OrderNum,
		Path:      req.Path,
		Component: req.Component,
		Query:     req.Query,
		RouteName: req.RouteName,
		IsFrame:   services.FlexInt(req.IsFrame, 1),
		IsCache:   services.FlexInt(req.IsCache, 0),
		MenuType:  req.MenuType,
		Visible:   req.Visible,
		Status:    req.Status,
		Perms:     req.Perms,
		Icon:      req.Icon,
		Remark:    req.Remark,
	}
	menu.ID = req.MenuID
	if menu.Visible == "" {
		menu.Visible = "0"
	}
	if menu.Status == "" {
		menu.Status = "0"
	}
	return menu
}

// Add POST /system/menu
func (h *MenuHandler) Add(ctx *gin.Context) {
	var req menuBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Add(h.bindMenu(req), services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "新增成功")
}

// Edit PUT /system/menu
func (h *MenuHandler) Edit(ctx *gin.Context) {
	var req menuBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Edit(h.bindMenu(req), services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

type menuSortBody struct {
	MenuIDs   string `json:"menuIds"`
	OrderNums string `json:"orderNums"`
}

// UpdateSort PUT /system/menu/updateSort
func (h *MenuHandler) UpdateSort(ctx *gin.Context) {
	var req menuSortBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.UpdateSort(req.MenuIDs, req.OrderNums); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// Delete DELETE /system/menu/:menuId
func (h *MenuHandler) Delete(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("menuId"), 10, 64)
	if err := h.svc.Delete(uint(id)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "删除成功")
}
