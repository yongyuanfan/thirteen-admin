package models

import (
	"time"

	"thirteen-admin/internal/base/model"
)

// OperLog 操作日志记录。
type OperLog struct {
	model.BaseModel
	Title         string     `gorm:"size:50;default:''" json:"title"`
	BusinessType  int        `gorm:"column:business_type;default:0;index:idx_sys_oper_log_bt" json:"businessType"` // 0其它 1新增 2修改 3删除
	Method        string     `gorm:"size:200;default:''" json:"method"`
	RequestMethod string     `gorm:"size:10;default:''" json:"requestMethod"`
	OperatorType  int        `gorm:"column:operator_type;default:0" json:"operatorType"` // 0其它 1后台用户 2手机端
	OperName      string     `gorm:"size:50;default:''" json:"operName"`
	DeptName      string     `gorm:"size:50;default:''" json:"deptName"`
	OperURL       string     `gorm:"column:oper_url;size:255;default:''" json:"operUrl"`
	OperIP        string     `gorm:"column:oper_ip;size:128;default:''" json:"operIp"`
	OperLocation  string     `gorm:"size:255;default:''" json:"operLocation"`
	OperParam     string     `gorm:"size:2000;default:''" json:"operParam"`
	JSONResult    string     `gorm:"column:json_result;size:2000;default:''" json:"jsonResult"`
	Status        int        `gorm:"default:0;index:idx_sys_oper_log_s" json:"status"` // 0正常 1异常
	ErrorMsg      string     `gorm:"size:2000;default:''" json:"errorMsg"`
	OperTime      *time.Time `gorm:"column:oper_time;index:idx_sys_oper_log_ot" json:"operTime"`
	CostTime      int64      `gorm:"column:cost_time;default:0" json:"costTime"`
}

func (OperLog) TableName() string { return "sys_oper_log" }
