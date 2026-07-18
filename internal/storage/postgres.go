package storage

import (
	"errors"
	"fmt"

	"thirteen-admin/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDB *gorm.DB

var ErrPostgresNotInitialized = errors.New("postgres 未初始化")

// InitPostgres 初始化 PostgreSQL 连接。
func InitPostgres(postgresConfig config.Postgres) error {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.User,
		postgresConfig.Password,
		postgresConfig.Database,
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}

	postgresDB = conn
	db = conn
	dbDriver = "postgres"
	return nil
}

// GetPostgres 返回全局 PostgreSQL 连接实例。
func GetPostgres() *gorm.DB {
	return postgresDB
}

// RequirePostgres 返回全局 PostgreSQL 连接实例，未初始化时返回错误。
func RequirePostgres() (*gorm.DB, error) {
	if postgresDB == nil {
		return nil, ErrPostgresNotInitialized
	}
	return postgresDB, nil
}
