package runtime

import (
	"context"
	"os"
	"os/signal"
	"thirteen-admin/internal/config"
	"thirteen-admin/internal/logging"
	"thirteen-admin/internal/storage"
	"time"
)

// InitializeRuntime 加载配置并初始化服务运行所需基础组件。
func InitializeRuntime(configPath string) error {
	if err := config.Load(configPath); err != nil {
		return err
	}
	if err := logging.Init(config.GetApplication().Logging); err != nil {
		return err
	}
	if err := validateRuntimeDependencies(); err != nil {
		return err
	}
	if err := storage.Init(); err != nil {
		return err
	}
	return nil
}

func validateRuntimeDependencies() error {
	// 基础框架运行时依赖验证
	// 当前无强制依赖，后续业务模块可按需添加验证逻辑
	return nil
}

// RunServer 完成服务初始化、启动、等待退出和优雅关闭流程。
func RunServer(configPath string) error {
	if err := InitializeRuntime(configPath); err != nil {
		return err
	}
	defer func() {
		_ = logging.Close()
	}()

	runtime, err := NewServiceRuntime()
	if err != nil {
		return err
	}
	errCh := runtime.Start()
	runtime.PrintEndpoints()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	select {
	case <-quit:
	case err := <-errCh:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logCoref("shutdown server ...")

	if err := runtime.Shutdown(ctx); err != nil {
		return err
	}
	logCoref("server exiting")
	return nil
}
