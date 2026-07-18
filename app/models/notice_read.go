package models

import (
	"time"

	"thirteen-admin/internal/base/model"
)

// NoticeRead 公告已读记录表。
type NoticeRead struct {
	model.BaseModel
	NoticeID uint      `gorm:"column:notice_id;not null;uniqueIndex:uk_user_notice" json:"noticeId"`
	UserID   uint      `gorm:"column:user_id;not null;uniqueIndex:uk_user_notice" json:"userId"`
	ReadTime time.Time `gorm:"column:read_time;not null" json:"readTime"`
}

func (NoticeRead) TableName() string { return "sys_notice_read" }
