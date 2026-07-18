package models

import "thirteen-admin/internal/base/model"

// DictType 字典类型表。
type DictType struct {
	model.BaseModel
	DictName  string `gorm:"size:100;default:''" json:"dictName"`
	DictType  string `gorm:"size:100;default:'';uniqueIndex" json:"dictType"`
	Status    string `gorm:"size:1;default:0" json:"status"` // 0正常 1停用
	CreatedBy string `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy string `gorm:"size:64;default:''" json:"updatedBy"`
	Remark    string `gorm:"size:500" json:"remark"`
}

func (DictType) TableName() string { return "sys_dict_type" }
