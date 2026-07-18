package seeds

import (
	"fmt"

	"clash-admin/app/models"
	"clash-admin/internal/base/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AutoMigrateUsers 迁移用户表。
func AutoMigrateUsers(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	// 将历史禁用状态 0 规范为 2
	return db.Model(&models.User{}).Where("status = ?", 0).Update("status", 2).Error
}

// SeedUsers 幂等写入演示账号（Super / Admin / User）。
func SeedUsers(db *gorm.DB) error {
	items := []struct {
		userName string
		password string
		nickName string
		phone    string
		email    string
		roles    model.StringSlice
		buttons  model.StringSlice
	}{
		{"Super", "123456", "超级管理员", "13800000001", "super@example.com", model.StringSlice{"R_SUPER"}, model.StringSlice{"*"}},
		{"Admin", "123456", "管理员", "13800000002", "admin@example.com", model.StringSlice{"R_ADMIN"}, model.StringSlice{}},
		{"User", "123456", "普通用户", "13800000003", "user@example.com", model.StringSlice{"R_USER"}, model.StringSlice{}},
	}

	for _, s := range items {
		var count int64
		if err := db.Model(&models.User{}).Where("user_name = ?", s.userName).Count(&count).Error; err != nil {
			return fmt.Errorf("count user %s: %w", s.userName, err)
		}
		if count > 0 {
			continue
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(s.password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password for %s: %w", s.userName, err)
		}
		phone := s.phone
		user := models.User{
			UserName:     s.userName,
			PasswordHash: string(hash),
			NickName:     s.nickName,
			Phone:        &phone,
			UserEmail:    s.email,
			Roles:        s.roles,
			Buttons:      s.buttons,
			Status:       1,
			CreateBy:     "system",
			UpdateBy:     "system",
		}
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("create user %s: %w", s.userName, err)
		}
	}
	return nil
}
