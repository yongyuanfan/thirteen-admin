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

// MenuHandler 菜单 HTTP 处理器。
type MenuHandler struct {
	handler.Api
}

// List GET /system/menu/list
func (e MenuHandler) List(c *gin.Context) {
	s := services.MenuService{}
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

// Get GET /system/menu/:menuId
func (e MenuHandler) Get(c *gin.Context) {
	s := services.MenuService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("menuId"), 10, 64)
	data, err := s.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "菜单不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

// TreeSelect GET /system/menu/treeselect
func (e MenuHandler) TreeSelect(c *gin.Context) {
	s := services.MenuService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	data, err := s.TreeSelect()
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

// RoleMenuTreeSelect GET /system/menu/roleMenuTreeselect/:roleId
func (e MenuHandler) RoleMenuTreeSelect(c *gin.Context) {
	s := services.MenuService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	roleID, _ := strconv.ParseUint(c.Param("roleId"), 10, 64)
	menus, checked, err := s.RoleMenuTreeSelect(uint(roleID))
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.OK(c, gin.H{
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

func (e MenuHandler) bindMenu(req menuBody) *models.Menu {
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
func (e MenuHandler) Add(c *gin.Context) {
	s := services.MenuService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req menuBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.Add(e.bindMenu(req), services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "新增成功")
}

// Edit PUT /system/menu
func (e MenuHandler) Edit(c *gin.Context) {
	s := services.MenuService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req menuBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.Edit(e.bindMenu(req), services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "修改成功")
}

type menuSortBody struct {
	MenuIDs   string `json:"menuIds"`
	OrderNums string `json:"orderNums"`
}

// UpdateSort PUT /system/menu/updateSort
func (e MenuHandler) UpdateSort(c *gin.Context) {
	s := services.MenuService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req menuSortBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.UpdateSort(req.MenuIDs, req.OrderNums); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// Delete DELETE /system/menu/:menuId
func (e MenuHandler) Delete(c *gin.Context) {
	s := services.MenuService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("menuId"), 10, 64)
	if err := s.Delete(uint(id)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "删除成功")
}
