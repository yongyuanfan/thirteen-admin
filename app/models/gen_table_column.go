package models

import "thirteen-admin/internal/base/model"

// GenTableColumn 代码生成业务表字段。
type GenTableColumn struct {
	model.BaseModel
	TableID       *uint  `gorm:"column:table_id" json:"tableId"`
	ColumnName    string `gorm:"size:200" json:"columnName"`
	ColumnComment string `gorm:"size:500" json:"columnComment"`
	ColumnType    string `gorm:"size:100" json:"columnType"`
	JavaType      string `gorm:"size:500" json:"javaType"`
	JavaField     string `gorm:"size:200" json:"javaField"`
	IsPk          string `gorm:"size:1" json:"isPk"`
	IsIncrement   string `gorm:"size:1" json:"isIncrement"`
	IsRequired    string `gorm:"size:1" json:"isRequired"`
	IsInsert      string `gorm:"size:1" json:"isInsert"`
	IsEdit        string `gorm:"size:1" json:"isEdit"`
	IsList        string `gorm:"size:1" json:"isList"`
	IsQuery       string `gorm:"size:1" json:"isQuery"`
	QueryType     string `gorm:"size:200;default:EQ" json:"queryType"`
	HTMLType      string `gorm:"column:html_type;size:200" json:"htmlType"`
	DictType      string `gorm:"size:200;default:''" json:"dictType"`
	Sort          int    `json:"sort"`
	CreatedBy     string `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy     string `gorm:"size:64;default:''" json:"updatedBy"`
}

func (GenTableColumn) TableName() string { return "gen_table_column" }
