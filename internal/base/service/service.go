package service

import "gorm.io/gorm"

// Service 通用业务服务基类。
type Service struct {
	DB *gorm.DB
}

// New 创建带 DB 的基础 Service。
func New(db *gorm.DB) Service {
	return Service{DB: db}
}
