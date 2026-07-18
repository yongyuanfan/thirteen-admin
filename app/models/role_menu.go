package models

// RoleMenu 角色-菜单关联。
type RoleMenu struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	RoleID uint `gorm:"not null;uniqueIndex:uk_role_menu" json:"roleId"`
	MenuID uint `gorm:"not null;uniqueIndex:uk_role_menu;index" json:"menuId"`
}

func (RoleMenu) TableName() string {
	return "role_menus"
}
