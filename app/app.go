package app

import (
	"fmt"
	"path/filepath"

	"thirteen-admin/app/handlers"
	"thirteen-admin/app/middlewares"
	"thirteen-admin/app/router"
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

	profileRoot := filepath.Join(pathx.ResolveProjectRoot(), "profile")
	engine.Static("/profile", profileRoot)
	engine.Use(middlewares.OperLog(db))

	jwtManager := services.NewJWTManager(appConfig.JWT)
	captchaSvc := services.NewCaptchaService(rdb, true)
	onlineSvc := services.NewOnlineService(rdb, 0)
	authService := services.NewAuthService(db, rdb, jwtManager, captchaSvc, onlineSvc)
	fileService := services.NewFileService()

	deptService := services.NewDeptService(db)
	postService := services.NewPostService(db)
	menuService := services.NewMenuService(db)
	configService := services.NewConfigService(db)
	dictService := services.NewDictService(db)
	noticeService := services.NewNoticeService(db)
	userService := services.NewUserService(db, deptService, postService)
	roleService := services.NewRoleService(db, deptService)

	operLogService := services.NewOperLogService(db)
	logininforService := services.NewLogininforService(db, rdb)
	cacheService := services.NewCacheService(rdb)
	serverService := services.NewServerService()
	jobService := services.NewJobService(db)
	jobLogService := services.NewJobLogService(db)
	genService := services.NewGenService(db)

	authHandler := handlers.NewAuthHandler(authService)
	commonHandler := handlers.NewCommonHandler(fileService)
	deptHandler := handlers.NewDeptHandler(deptService)
	menuHandler := handlers.NewMenuHandler(menuService)
	postHandler := handlers.NewPostHandler(postService)
	configHandler := handlers.NewConfigHandler(configService)
	dictHandler := handlers.NewDictHandler(dictService)
	noticeHandler := handlers.NewNoticeHandler(noticeService)
	userHandler := handlers.NewUserHandler(userService, fileService)
	roleHandler := handlers.NewRoleHandler(roleService)
	operLogHandler := handlers.NewOperLogHandler(operLogService)
	logininforHandler := handlers.NewLogininforHandler(logininforService)
	onlineHandler := handlers.NewOnlineHandler(onlineSvc)
	cacheHandler := handlers.NewCacheHandler(cacheService)
	serverHandler := handlers.NewServerHandler(serverService)
	jobHandler := handlers.NewJobHandler(jobService)
	jobLogHandler := handlers.NewJobLogHandler(jobLogService)
	genHandler := handlers.NewGenHandler(genService)

	router.RegisterAuthRoutes(authHandler)
	router.RegisterCommonRoutes(commonHandler)
	router.RegisterDeptRoutes(deptHandler)
	router.RegisterMenuRoutes(menuHandler)
	router.RegisterPostRoutes(postHandler)
	router.RegisterConfigRoutes(configHandler)
	router.RegisterDictRoutes(dictHandler)
	router.RegisterNoticeRoutes(noticeHandler)
	router.RegisterUserRoutes(userHandler)
	router.RegisterRoleRoutes(roleHandler)
	router.RegisterOperLogRoutes(operLogHandler)
	router.RegisterLogininforRoutes(logininforHandler)
	router.RegisterOnlineRoutes(onlineHandler)
	router.RegisterCacheRoutes(cacheHandler)
	router.RegisterServerRoutes(serverHandler)
	router.RegisterJobRoutes(jobHandler)
	router.RegisterJobLogRoutes(jobLogHandler)
	router.RegisterGenRoutes(genHandler)

	router.InitRouter(engine, middlewares.Auth(jwtManager))
	return nil
}
