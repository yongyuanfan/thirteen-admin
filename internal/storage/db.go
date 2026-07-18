package storage

import (
	"errors"

	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	dbDriver string
)

var ErrDBNotInitialized = errors.New("database 未初始化")

// GetDB 返回当前驱动对应的全局 DB 连接。
func GetDB() *gorm.DB {
	return db
}

// RequireDB 返回当前驱动对应的全局 DB 连接，未初始化时返回错误。
func RequireDB() (*gorm.DB, error) {
	if db == nil {
		return nil, ErrDBNotInitialized
	}
	return db, nil
}

// DBDriver 返回当前数据库驱动名称。
func DBDriver() string {
	return dbDriver
}
