package models

// UserRole 用户和角色关联表。
type UserRole struct {
	UserID uint `gorm:"column:user_id;primaryKey" json:"userId"`
	RoleID uint `gorm:"column:role_id;primaryKey" json:"roleId"`
}

func (UserRole) TableName() string { return "sys_user_role" }
