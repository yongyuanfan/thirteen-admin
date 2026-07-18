package handlers

import (
	"clash-admin/app/models"
	"clash-admin/app/services"

	"github.com/gin-gonic/gin"
)

type menuUpsertReq struct {
	ID              uint                `json:"id"`
	ParentID        uint                `json:"parentId"`
	MenuType        string              `json:"menuType"`
	MenuName        string              `json:"menuName"`
	RouteName       string              `json:"routeName"`
	RoutePath       string              `json:"routePath"`
	Component       string              `json:"component"`
	Icon            string              `json:"icon"`
	IconType        string              `json:"iconType"`
	Buttons         []models.MenuButton `json:"buttons"`
	Query           []models.MenuQuery  `json:"query"`
	I18nKey         *string             `json:"i18nKey"`
	KeepAlive       *bool               `json:"keepAlive"`
	Constant        *bool               `json:"constant"`
	Order           *int                `json:"order"`
	Href            *string             `json:"href"`
	HideInMenu      *bool               `json:"hideInMenu"`
	ActiveMenu      *string             `json:"activeMenu"`
	MultiTab        *bool               `json:"multiTab"`
	FixedIndexInTab *int                `json:"fixedIndexInTab"`
	Status          string              `json:"status"`
}

func (r menuUpsertReq) toInput(operator string) services.MenuUpsertInput {
	return services.MenuUpsertInput{
		ID:              r.ID,
		ParentID:        r.ParentID,
		MenuType:        r.MenuType,
		MenuName:        r.MenuName,
		RouteName:       r.RouteName,
		RoutePath:       r.RoutePath,
		Component:       r.Component,
		Icon:            r.Icon,
		IconType:        r.IconType,
		Buttons:         r.Buttons,
		Query:           r.Query,
		I18nKey:         r.I18nKey,
		KeepAlive:       r.KeepAlive,
		Constant:        r.Constant,
		Order:           r.Order,
		Href:            r.Href,
		HideInMenu:      r.HideInMenu,
		ActiveMenu:      r.ActiveMenu,
		MultiTab:        r.MultiTab,
		FixedIndexInTab: r.FixedIndexInTab,
		Status:          r.Status,
		Operator:        operator,
	}
}

// GetMenuList GET /systemManage/getMenuList/v2
func (h *SystemManageHandler) GetMenuList(ctx *gin.Context) {
	data, err := h.menus.GetMenuList(ctx.Request.Context())
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// GetMenuTree GET /systemManage/getMenuTree
func (h *SystemManageHandler) GetMenuTree(ctx *gin.Context) {
	data, err := h.menus.GetMenuTree(ctx.Request.Context())
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// GetAllPages GET /systemManage/getAllPages
func (h *SystemManageHandler) GetAllPages(ctx *gin.Context) {
	data, err := h.menus.GetAllPages(ctx.Request.Context())
	if err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// AddMenu POST /systemManage/addMenu
func (h *SystemManageHandler) AddMenu(ctx *gin.Context) {
	var req menuUpsertReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	if err := h.menus.AddMenu(ctx.Request.Context(), req.toInput(h.operator(ctx))); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// UpdateMenu POST /systemManage/updateMenu
func (h *SystemManageHandler) UpdateMenu(ctx *gin.Context) {
	var req menuUpsertReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	if err := h.menus.UpdateMenu(ctx.Request.Context(), req.toInput(h.operator(ctx))); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// DeleteMenu DELETE /systemManage/deleteMenu
func (h *SystemManageHandler) DeleteMenu(ctx *gin.Context) {
	var req idReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	if err := h.menus.DeleteMenus(ctx.Request.Context(), []uint{req.ID}); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}

// BatchDeleteMenu DELETE /systemManage/batchDeleteMenu
func (h *SystemManageHandler) BatchDeleteMenu(ctx *gin.Context) {
	var req idsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		h.BadRequest(ctx, "请求参数错误")
		return
	}
	if err := h.menus.DeleteMenus(ctx.Request.Context(), req.IDs); err != nil {
		h.failManage(ctx, err)
		return
	}
	h.Success(ctx, nil)
}
