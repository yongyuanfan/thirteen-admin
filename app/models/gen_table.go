package models

import "thirteen-admin/internal/base/model"

// GenTable 代码生成业务表。
type GenTable struct {
	model.BaseModel
	Name           string `gorm:"column:table_name;size:200;default:''" json:"tableName"`
	TableComment   string `gorm:"size:500;default:''" json:"tableComment"`
	SubTableName   string `gorm:"size:64" json:"subTableName"`
	SubTableFkName string `gorm:"column:sub_table_fk_name;size:64" json:"subTableFkName"`
	ClassName      string `gorm:"size:100;default:''" json:"className"`
	TplCategory    string `gorm:"size:200;default:crud" json:"tplCategory"` // crud / tree
	TplWebType     string `gorm:"size:30;default:''" json:"tplWebType"`
	PackageName    string `gorm:"size:100" json:"packageName"`
	ModuleName     string `gorm:"size:30" json:"moduleName"`
	BusinessName   string `gorm:"size:30" json:"businessName"`
	FunctionName   string `gorm:"size:50" json:"functionName"`
	FunctionAuthor string `gorm:"size:50" json:"functionAuthor"`
	FormColNum     int    `gorm:"column:form_col_num;default:1" json:"formColNum"`
	GenType        string `gorm:"size:1;default:0" json:"genType"` // 0zip 1自定义路径
	GenPath        string `gorm:"size:200;default:'/'" json:"genPath"`
	Options        string `gorm:"size:1000" json:"options"`
	CreatedBy      string `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy      string `gorm:"size:64;default:''" json:"updatedBy"`
	Remark         string `gorm:"size:500" json:"remark"`
}

func (GenTable) TableName() string { return "gen_table" }
