package storage

import (
	"fmt"
	"strings"

	"thirteen-admin/internal/config"
)

// Init 按配置初始化数据库与缓存。
func Init() error {
	app := config.GetApplication()

	dbDriverName := strings.ToLower(strings.TrimSpace(app.Database.Driver))
	switch dbDriverName {
	case "mysql", "":
		if err := InitMySQL(app.Database.MySQL); err != nil {
			return fmt.Errorf("init mysql: %w", err)
		}
	case "postgres":
		if err := InitPostgres(app.Database.Postgres); err != nil {
			return fmt.Errorf("init postgres: %w", err)
		}
	default:
		return fmt.Errorf("unsupported database driver: %s", app.Database.Driver)
	}

	cacheDriverName := strings.ToLower(strings.TrimSpace(app.Cache.Driver))
	switch cacheDriverName {
	case "file", "":
		if err := InitFileCache(app.Cache.File); err != nil {
			return fmt.Errorf("init file cache: %w", err)
		}
	case "redis":
		if err := InitRedis(app.Cache.Redis); err != nil {
			return fmt.Errorf("init redis: %w", err)
		}
	default:
		return fmt.Errorf("unsupported cache driver: %s", app.Cache.Driver)
	}

	return nil
}
