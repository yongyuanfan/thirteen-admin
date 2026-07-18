package models

import "clash-admin/internal/base/model"

// User 系统用户。
type User struct {
	model.BaseModel
	UserName     string            `gorm:"size:64;uniqueIndex;not null" json:"userName"`
	PasswordHash string            `gorm:"size:255;not null" json:"-"`
	NickName     string            `gorm:"size:64;not null;default:''" json:"nickName"`
	UserGender   *string           `gorm:"size:8" json:"userGender"` // 1=男 2=女
	Phone        *string           `gorm:"size:20;uniqueIndex" json:"phone"`
	UserEmail    string            `gorm:"size:128;not null;default:''" json:"userEmail"`
	Roles        model.StringSlice `gorm:"type:json;not null" json:"roles"`
	Buttons      model.StringSlice `gorm:"type:json;not null" json:"buttons"`
	Status       int8              `gorm:"not null;default:1" json:"status"` // 1=启用 2=禁用
	CreateBy     string            `gorm:"size:64;not null;default:''" json:"createBy"`
	UpdateBy     string            `gorm:"size:64;not null;default:''" json:"updateBy"`
}

func (User) TableName() string {
	return "users"
}
