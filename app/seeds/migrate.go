package seeds

import (
	"thirteen-admin/app/models"

	"gorm.io/gorm"
)

// AutoMigrate 迁移全部业务表。
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Dept{},
		&models.User{},
		&models.Post{},
		&models.Role{},
		&models.Menu{},
		&models.UserRole{},
		&models.RoleMenu{},
		&models.RoleDept{},
		&models.UserPost{},
		&models.OperLog{},
		&models.DictType{},
		&models.DictData{},
		&models.Config{},
		&models.Logininfor{},
		&models.Job{},
		&models.JobLog{},
		&models.Notice{},
		&models.NoticeRead{},
		&models.GenTable{},
		&models.GenTableColumn{},
	)
}
