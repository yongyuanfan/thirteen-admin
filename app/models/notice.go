package models

import "thirteen-admin/internal/base/model"

// Notice 通知公告表。
type Notice struct {
	model.BaseModel
	NoticeTitle   string `gorm:"size:50;not null" json:"noticeTitle"`
	NoticeType    string `gorm:"size:1;not null" json:"noticeType"` // 1通知 2公告
	NoticeContent []byte `gorm:"type:longblob" json:"noticeContent"`
	Status        string `gorm:"size:1;default:0" json:"status"` // 0正常 1关闭
	CreatedBy     string `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy     string `gorm:"size:64;default:''" json:"updatedBy"`
	Remark        string `gorm:"size:255" json:"remark"`
}

func (Notice) TableName() string { return "sys_notice" }
