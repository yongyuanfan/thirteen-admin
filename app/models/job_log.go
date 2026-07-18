package models

import (
	"time"

	"thirteen-admin/internal/base/model"
)

// JobLog 定时任务调度日志表。
type JobLog struct {
	model.BaseModel
	JobName       string     `gorm:"size:64;not null" json:"jobName"`
	JobGroup      string     `gorm:"size:64;not null" json:"jobGroup"`
	InvokeTarget  string     `gorm:"size:500;not null" json:"invokeTarget"`
	JobMessage    string     `gorm:"size:500" json:"jobMessage"`
	Status        string     `gorm:"size:1;default:0" json:"status"` // 0正常 1失败
	ExceptionInfo string     `gorm:"size:2000;default:''" json:"exceptionInfo"`
	StartTime     *time.Time `gorm:"column:start_time" json:"startTime"`
	EndTime       *time.Time `gorm:"column:end_time" json:"endTime"`
}

func (JobLog) TableName() string { return "sys_job_log" }
