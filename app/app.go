package app

import (
	"fmt"
	"path/filepath"

	"thirteen-admin/app/handlers"
	hmonitor "thirteen-admin/app/handlers/monitor"
	hsystem "thirteen-admin/app/handlers/system"
	htool "thirteen-admin/app/handlers/tool"
	"thirteen-admin/app/middlewares"
	"thirteen-admin/app/router"
	rmonitor "thirteen-admin/app/router/monitor"
	rsystem "thirteen-admin/app/router/system"
	rtool "thirteen-admin/app/router/tool"
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
	deptHandler := hsystem.NewDeptHandler(deptService)
	menuHandler := hsystem.NewMenuHandler(menuService)
	postHandler := hsystem.NewPostHandler(postService)
	configHandler := hsystem.NewConfigHandler(configService)
	dictHandler := hsystem.NewDictHandler(dictService)
	noticeHandler := hsystem.NewNoticeHandler(noticeService)
	userHandler := hsystem.NewUserHandler(userService, fileService)
	roleHandler := hsystem.NewRoleHandler(roleService)
	operLogHandler := hmonitor.NewOperLogHandler(operLogService)
	logininforHandler := hmonitor.NewLogininforHandler(logininforService)
	onlineHandler := hmonitor.NewOnlineHandler(onlineSvc)
	cacheHandler := hmonitor.NewCacheHandler(cacheService)
	serverHandler := hmonitor.NewServerHandler(serverService)
	jobHandler := hmonitor.NewJobHandler(jobService)
	jobLogHandler := hmonitor.NewJobLogHandler(jobLogService)
	genHandler := htool.NewGenHandler(genService)

	router.RegisterAuthRoutes(authHandler)
	router.RegisterCommonRoutes(commonHandler)
	rsystem.RegisterDeptRoutes(deptHandler)
	rsystem.RegisterMenuRoutes(menuHandler)
	rsystem.RegisterPostRoutes(postHandler)
	rsystem.RegisterConfigRoutes(configHandler)
	rsystem.RegisterDictRoutes(dictHandler)
	rsystem.RegisterNoticeRoutes(noticeHandler)
	rsystem.RegisterUserRoutes(userHandler)
	rsystem.RegisterRoleRoutes(roleHandler)
	rmonitor.RegisterOperLogRoutes(operLogHandler)
	rmonitor.RegisterLogininforRoutes(logininforHandler)
	rmonitor.RegisterOnlineRoutes(onlineHandler)
	rmonitor.RegisterCacheRoutes(cacheHandler)
	rmonitor.RegisterServerRoutes(serverHandler)
	rmonitor.RegisterJobRoutes(jobHandler)
	rmonitor.RegisterJobLogRoutes(jobLogHandler)
	rtool.RegisterGenRoutes(genHandler)

	router.InitRouter(engine, middlewares.Auth(jwtManager))
	return nil
}
