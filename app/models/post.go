package models

import "thirteen-admin/internal/base/model"

// Post 岗位信息表。
type Post struct {
	model.BaseModel
	PostCode  string `gorm:"size:64;not null" json:"postCode"`
	PostName  string `gorm:"size:50;not null" json:"postName"`
	PostSort  int    `gorm:"not null" json:"postSort"`
	Status    string `gorm:"size:1;not null" json:"status"` // 0正常 1停用
	CreatedBy string `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy string `gorm:"size:64;default:''" json:"updatedBy"`
	Remark    string `gorm:"size:500" json:"remark"`
}

func (Post) TableName() string { return "sys_post" }
