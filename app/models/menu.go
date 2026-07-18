package models

import "clash-admin/internal/base/model"

// MenuButton 菜单按钮权限项。
type MenuButton struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

// MenuQuery 菜单路由 query 项。
type MenuQuery struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Menu 系统菜单。
type Menu struct {
	model.BaseModel
	ParentID        uint                        `gorm:"not null;default:0;index" json:"parentId"`
	MenuType        string                      `gorm:"size:8;not null" json:"menuType"` // 1=目录 2=菜单
	MenuName        string                      `gorm:"size:64;not null" json:"menuName"`
	RouteName       string                      `gorm:"size:128;not null;index" json:"routeName"`
	RoutePath       string                      `gorm:"size:255;not null;default:''" json:"routePath"`
	Component       string                      `gorm:"size:255;not null;default:''" json:"component"`
	Icon            string                      `gorm:"size:128;not null;default:''" json:"icon"`
	IconType        string                      `gorm:"size:8;not null;default:1" json:"iconType"` // 1=iconify 2=local
	Buttons         model.JSONSlice[MenuButton] `gorm:"type:json" json:"buttons"`
	Query           model.JSONSlice[MenuQuery]  `gorm:"type:json" json:"query"`
	I18nKey         *string                     `gorm:"size:128" json:"i18nKey"`
	KeepAlive       *bool                       `json:"keepAlive"`
	Constant        *bool                       `json:"constant"`
	Order           *int                        `gorm:"column:sort_order" json:"order"`
	Href            *string                     `gorm:"size:512" json:"href"`
	HideInMenu      *bool                       `json:"hideInMenu"`
	ActiveMenu      *string                     `gorm:"size:128" json:"activeMenu"`
	MultiTab        *bool                       `json:"multiTab"`
	FixedIndexInTab *int                        `json:"fixedIndexInTab"`
	Status          int8                        `gorm:"not null;default:1" json:"status"` // 1=启用 2=禁用
	CreateBy        string                      `gorm:"size:64;not null;default:''" json:"createBy"`
	UpdateBy        string                      `gorm:"size:64;not null;default:''" json:"updateBy"`
}

func (Menu) TableName() string {
	return "menus"
}
