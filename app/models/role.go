package models

import "clash-admin/internal/base/model"

// Role 系统角色。
type Role struct {
	model.BaseModel
	RoleName string `gorm:"size:64;not null" json:"roleName"`
	RoleCode string `gorm:"size:64;uniqueIndex;not null" json:"roleCode"`
	RoleDesc string `gorm:"size:255;not null;default:''" json:"roleDesc"`
	Home     string `gorm:"size:128;not null;default:home" json:"home"`
	Status   int8   `gorm:"not null;default:1" json:"status"` // 1=启用 2=禁用
	CreateBy string `gorm:"size:64;not null;default:''" json:"createBy"`
	UpdateBy string `gorm:"size:64;not null;default:''" json:"updateBy"`
}

func (Role) TableName() string {
	return "roles"
}
