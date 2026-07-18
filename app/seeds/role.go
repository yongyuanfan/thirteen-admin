package seeds

import (
	"fmt"

	"clash-admin/app/models"

	"gorm.io/gorm"
)

// AutoMigrateRoles 迁移角色及相关联表。
func AutoMigrateRoles(db *gorm.DB) error {
	return db.AutoMigrate(&models.Role{}, &models.RoleMenu{}, &models.RoleButton{})
}

// SeedRoles 幂等写入基础角色。
func SeedRoles(db *gorm.DB) error {
	items := []models.Role{
		{RoleName: "超级管理员", RoleCode: "R_SUPER", RoleDesc: "拥有系统全部权限", Home: "home", Status: 1, CreateBy: "system", UpdateBy: "system"},
		{RoleName: "管理员", RoleCode: "R_ADMIN", RoleDesc: "系统管理权限", Home: "home", Status: 1, CreateBy: "system", UpdateBy: "system"},
		{RoleName: "普通用户", RoleCode: "R_USER", RoleDesc: "普通业务权限", Home: "home", Status: 1, CreateBy: "system", UpdateBy: "system"},
	}
	for _, item := range items {
		var count int64
		if err := db.Model(&models.Role{}).Where("role_code = ?", item.RoleCode).Count(&count).Error; err != nil {
			return fmt.Errorf("count role %s: %w", item.RoleCode, err)
		}
		if count > 0 {
			continue
		}
		if err := db.Create(&item).Error; err != nil {
			return fmt.Errorf("create role %s: %w", item.RoleCode, err)
		}
	}
	return nil
}
