package models

// RoleButton 角色-按钮权限关联。
type RoleButton struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	RoleID     uint   `gorm:"not null;uniqueIndex:uk_role_button" json:"roleId"`
	ButtonCode string `gorm:"size:128;not null;uniqueIndex:uk_role_button" json:"buttonCode"`
}

func (RoleButton) TableName() string {
	return "role_buttons"
}
