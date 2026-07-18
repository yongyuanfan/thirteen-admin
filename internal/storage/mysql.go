package storage

import (
	"errors"
	"fmt"

	"thirteen-admin/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlDB *gorm.DB

var ErrMySQLNotInitialized = errors.New("mysql 未初始化")

// InitMySQL 初始化 MySQL 连接。
func InitMySQL(mysqlConfig config.MySQL) error {
	charset := mysqlConfig.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig.User,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database,
		charset,
	)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect mysql: %w", err)
	}

	mysqlDB = conn
	db = conn
	dbDriver = "mysql"
	return nil
}

// GetMySQL 返回全局 MySQL 连接实例。
func GetMySQL() *gorm.DB {
	return mysqlDB
}

// RequireMySQL 返回全局 MySQL 连接实例，未初始化时返回错误。
func RequireMySQL() (*gorm.DB, error) {
	if mysqlDB == nil {
		return nil, ErrMySQLNotInitialized
	}
	return mysqlDB, nil
}
