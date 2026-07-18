package models

import "thirteen-admin/internal/base/model"

// Job 定时任务调度表。
type Job struct {
	model.BaseModel
	JobName        string `gorm:"size:64;default:'';uniqueIndex:uk_job_name_group" json:"jobName"`
	JobGroup       string `gorm:"size:64;default:DEFAULT;uniqueIndex:uk_job_name_group" json:"jobGroup"`
	InvokeTarget   string `gorm:"size:500;not null" json:"invokeTarget"`
	CronExpression string `gorm:"size:255;default:''" json:"cronExpression"`
	MisfirePolicy  string `gorm:"size:20;default:3" json:"misfirePolicy"` // 1立即 2执行一次 3放弃
	Concurrent     string `gorm:"size:1;default:1" json:"concurrent"`     // 0允许 1禁止
	Status         string `gorm:"size:1;default:0" json:"status"`         // 0正常 1暂停
	CreatedBy      string `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy      string `gorm:"size:64;default:''" json:"updatedBy"`
	Remark         string `gorm:"size:500;default:''" json:"remark"`
}

func (Job) TableName() string { return "sys_job" }
