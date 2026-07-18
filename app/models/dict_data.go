package models

import "thirteen-admin/internal/base/model"

// DictData 字典数据表。
type DictData struct {
	model.BaseModel
	DictSort  int    `gorm:"column:dict_sort;default:0" json:"dictSort"`
	DictLabel string `gorm:"size:100;default:''" json:"dictLabel"`
	DictValue string `gorm:"size:100;default:''" json:"dictValue"`
	DictType  string `gorm:"size:100;default:''" json:"dictType"`
	CSSClass  string `gorm:"column:css_class;size:100" json:"cssClass"`
	ListClass string `gorm:"size:100" json:"listClass"`
	IsDefault string `gorm:"size:1;default:N" json:"isDefault"` // Y是 N否
	Status    string `gorm:"size:1;default:0" json:"status"`    // 0正常 1停用
	CreatedBy string `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy string `gorm:"size:64;default:''" json:"updatedBy"`
	Remark    string `gorm:"size:500" json:"remark"`
}

func (DictData) TableName() string { return "sys_dict_data" }
