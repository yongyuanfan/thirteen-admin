package models

import (
	"time"

	"thirteen-admin/internal/base/model"
)

// Logininfor 系统访问记录。
type Logininfor struct {
	model.BaseModel
	UserName      string     `gorm:"size:50;default:''" json:"userName"`
	IPAddr        string     `gorm:"column:ipaddr;size:128;default:''" json:"ipaddr"`
	LoginLocation string     `gorm:"size:255;default:''" json:"loginLocation"`
	Browser       string     `gorm:"size:50;default:''" json:"browser"`
	OS            string     `gorm:"column:os;size:50;default:''" json:"os"`
	Status        string     `gorm:"size:1;default:0;index:idx_sys_logininfor_s" json:"status"` // 0成功 1失败
	Msg           string     `gorm:"size:255;default:''" json:"msg"`
	LoginTime     *time.Time `gorm:"column:login_time;index:idx_sys_logininfor_lt" json:"loginTime"`
}

func (Logininfor) TableName() string { return "sys_logininfor" }
