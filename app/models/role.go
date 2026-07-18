package models

import "thirteen-admin/internal/base/model"

// Role 角色信息表。
type Role struct {
	model.BaseModel
	RoleName          string `gorm:"size:30;not null" json:"roleName"`
	RoleKey           string `gorm:"size:100;not null" json:"roleKey"`
	RoleSort          int    `gorm:"not null" json:"roleSort"`
	DataScope         string `gorm:"size:1;default:1" json:"dataScope"` // 1全部 2自定 3本部门 4本部门及以下
	MenuCheckStrictly bool   `gorm:"default:true" json:"menuCheckStrictly"`
	DeptCheckStrictly bool   `gorm:"default:true" json:"deptCheckStrictly"`
	Status            string `gorm:"size:1;not null" json:"status"`   // 0正常 1停用
	DelFlag           string `gorm:"size:1;default:0" json:"delFlag"` // 0存在 2删除
	CreatedBy         string `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy         string `gorm:"size:64;default:''" json:"updatedBy"`
	Remark            string `gorm:"size:500" json:"remark"`
}

func (Role) TableName() string { return "sys_role" }
