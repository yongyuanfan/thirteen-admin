package models

import "thirteen-admin/internal/base/model"

// Menu 菜单权限表。
type Menu struct {
	model.BaseModel
	MenuName  string `gorm:"size:50;not null" json:"menuName"`
	ParentID  uint   `gorm:"column:parent_id;default:0" json:"parentId"`
	OrderNum  int    `gorm:"column:order_num;default:0" json:"orderNum"`
	Path      string `gorm:"size:200;default:''" json:"path"`
	Component string `gorm:"size:255" json:"component"`
	Query     string `gorm:"size:255" json:"query"`
	RouteName string `gorm:"size:50;default:''" json:"routeName"`
	IsFrame   int    `gorm:"column:is_frame;default:1" json:"isFrame"` // 0是 1否
	IsCache   int    `gorm:"column:is_cache;default:0" json:"isCache"` // 0缓存 1不缓存
	MenuType  string `gorm:"size:1;default:''" json:"menuType"`        // M目录 C菜单 F按钮
	Visible   string `gorm:"size:1;default:0" json:"visible"`          // 0显示 1隐藏
	Status    string `gorm:"size:1;default:0" json:"status"`           // 0正常 1停用
	Perms     string `gorm:"size:100" json:"perms"`
	Icon      string `gorm:"size:100;default:#" json:"icon"`
	CreatedBy string `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy string `gorm:"size:64;default:''" json:"updatedBy"`
	Remark    string `gorm:"size:500;default:''" json:"remark"`
}

func (Menu) TableName() string { return "sys_menu" }
