package storage

import (
	"fmt"

	"clash-admin/internal/config"
)

// Init 按配置初始化所有存储组件。
func Init() error {
	storageConfig := config.GetApplication().Storage
	if storageConfig.MySQL.Enabled {
		if err := InitMySQL(storageConfig.MySQL); err != nil {
			return fmt.Errorf("init mysql: %w", err)
		}
	}

	if storageConfig.Postgres.Enabled {
		if err := InitPostgres(storageConfig.Postgres); err != nil {
			return fmt.Errorf("init postgres: %w", err)
		}
	}

	if storageConfig.Redis.Enabled {
		if err := InitRedis(storageConfig.Redis); err != nil {
			return fmt.Errorf("init redis: %w", err)
		}
	}

	return nil
}
