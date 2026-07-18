package models

import (
	"time"

	"thirteen-admin/internal/base/model"
)

// User 用户信息表。
type User struct {
	model.BaseModel
	DeptID        *uint      `gorm:"column:dept_id" json:"deptId"`
	UserName      string     `gorm:"size:30;not null" json:"userName"`
	NickName      string     `gorm:"size:30;not null" json:"nickName"`
	UserType      string     `gorm:"size:2;default:00" json:"userType"` // 00系统用户
	Email         string     `gorm:"size:50;default:''" json:"email"`
	Phonenumber   string     `gorm:"size:11;default:''" json:"phonenumber"`
	Sex           string     `gorm:"size:1;default:0" json:"sex"` // 0男 1女 2未知
	Avatar        string     `gorm:"size:100;default:''" json:"avatar"`
	Password      string     `gorm:"size:100;default:''" json:"-"`
	Status        string     `gorm:"size:1;default:0" json:"status"`  // 0正常 1停用
	DelFlag       string     `gorm:"size:1;default:0" json:"delFlag"` // 0存在 2删除
	LoginIP       string     `gorm:"column:login_ip;size:128;default:''" json:"loginIp"`
	LoginDate     *time.Time `gorm:"column:login_date" json:"loginDate"`
	PwdUpdateDate *time.Time `gorm:"column:pwd_update_date" json:"pwdUpdateDate"`
	CreatedBy     string     `gorm:"size:64;default:''" json:"createdBy"`
	UpdatedBy     string     `gorm:"size:64;default:''" json:"updatedBy"`
	Remark        string     `gorm:"size:500" json:"remark"`
}

func (User) TableName() string { return "sys_user" }
