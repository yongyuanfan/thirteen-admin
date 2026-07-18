package models

// UserPost 用户与岗位关联表。
type UserPost struct {
	UserID uint `gorm:"column:user_id;primaryKey" json:"userId"`
	PostID uint `gorm:"column:post_id;primaryKey" json:"postId"`
}

func (UserPost) TableName() string { return "sys_user_post" }
