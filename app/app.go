package app

import (
	"fmt"
	"path/filepath"

	"thirteen-admin/app/middlewares"
	"thirteen-admin/app/router"
	"thirteen-admin/app/runtime"
	"thirteen-admin/app/seeds"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/config"
	"thirteen-admin/internal/storage"
	"thirteen-admin/internal/utils/pathx"

	"github.com/gin-gonic/gin"
)

// Bootstrap 初始化应用并注册路由。
func Bootstrap(engine *gin.Engine) error {
	appConfig := config.GetApplication()

	db, err := storage.RequireMySQL()
	if err != nil {
		return fmt.Errorf("auth 依赖 mysql: %w", err)
	}
	rdb, err := storage.RequireRedis()
	if err != nil {
		return fmt.Errorf("auth 依赖 redis: %w", err)
	}

	if err := seeds.AutoMigrate(db); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	if err := seeds.SeedIfEmpty(db); err != nil {
		return fmt.Errorf("seed: %w", err)
	}
	if err := seeds.CleanupObsolete(db); err != nil {
		return fmt.Errorf("cleanup: %w", err)
	}

	runtime.SetDB(db)
	runtime.SetRedis(rdb)
	runtime.SetJWT(services.NewJWTManager(appConfig.JWT))
	runtime.SetOnline(services.NewOnlineService(rdb, 0))
	runtime.SetCaptcha(services.NewCaptchaService(rdb, true))
	runtime.SetFile(services.NewFileService())

	profileRoot := filepath.Join(pathx.ResolveProjectRoot(), "profile")
	engine.Static("/profile", profileRoot)
	engine.Use(middlewares.WithDB(), middlewares.OperLog())

	router.InitRouter(engine)
	return nil
}
