package app

import (
	"fmt"

	"thirteen-admin/app/handlers"
	"thirteen-admin/app/middlewares"
	"thirteen-admin/app/router"
	"thirteen-admin/app/seeds"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/config"
	"thirteen-admin/internal/storage"

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

	if err := seeds.AutoMigrateUsers(db); err != nil {
		return fmt.Errorf("migrate users: %w", err)
	}
	if err := seeds.AutoMigrateRoles(db); err != nil {
		return fmt.Errorf("migrate roles: %w", err)
	}
	if err := seeds.AutoMigrateMenus(db); err != nil {
		return fmt.Errorf("migrate menus: %w", err)
	}
	if err := seeds.SeedRoles(db); err != nil {
		return fmt.Errorf("seed roles: %w", err)
	}
	if err := seeds.SeedUsers(db); err != nil {
		return fmt.Errorf("seed users: %w", err)
	}
	if err := seeds.SeedMenus(db); err != nil {
		return fmt.Errorf("seed menus: %w", err)
	}

	jwtManager := services.NewJWTManager(appConfig.JWT)
	captchaStore := services.NewCaptchaStore(rdb, appConfig.Mode)
	authService := services.NewAuthService(db, rdb, jwtManager, captchaStore)
	authHandler := handlers.NewAuthHandler(authService)
	router.RegisterAuthRoutes(authHandler)

	router.InitRouter(engine, middlewares.Auth(jwtManager))
	return nil
}
