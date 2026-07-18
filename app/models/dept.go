package models

import "thirteen-admin/internal/base/model"

// Dept 部门表。
type Dept struct {
	model.BaseModel
	ParentID  uint   `gorm:"column:parent_id;default:0" json:"parentId"`
	Ancestors string `gorm:"size:50;default:''" json:"ancestors"`
	DeptName  string `gorm:"size:30;default:''" json:"deptName"`
	OrderNum  int    `gorm:"column:order_num;default:0" json:"orderNum"`
	Leader    string `gorm:"size:20" json:"leader"`
	Phone     string `gorm:"size:11" json:"phone"`
	Email     string `gorm:"size:50" json:"email"`
	Status    string `gorm:"size:1;default:0" json:"status"`  // 0正常 1停用
	DelFlag   string `gorm:"size:1;default:0" json:"delFlag"` // 0存在 2删除
	CreatedBy string `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy string `gorm:"size:64;default:''" json:"updatedBy"`
}

func (Dept) TableName() string { return "sys_dept" }
