package models

// RoleMenu 角色和菜单关联表。
type RoleMenu struct {
	RoleID uint `gorm:"column:role_id;primaryKey" json:"roleId"`
	MenuID uint `gorm:"column:menu_id;primaryKey" json:"menuId"`
}

func (RoleMenu) TableName() string { return "sys_role_menu" }
