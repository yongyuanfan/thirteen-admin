package models

import "thirteen-admin/internal/base/model"

// Config 参数配置表。
type Config struct {
	model.BaseModel
	ConfigName  string `gorm:"size:100;default:''" json:"configName"`
	ConfigKey   string `gorm:"size:100;default:''" json:"configKey"`
	ConfigValue string `gorm:"size:500;default:''" json:"configValue"`
	ConfigType  string `gorm:"size:1;default:N" json:"configType"` // Y是 N否
	CreatedBy   string `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy   string `gorm:"size:64;default:''" json:"updatedBy"`
	Remark      string `gorm:"size:500" json:"remark"`
}

func (Config) TableName() string { return "sys_config" }
