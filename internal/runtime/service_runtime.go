package runtime

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"thirteen-admin/app"
	"thirteen-admin/internal/config"

	"github.com/gin-gonic/gin"
)

// ServiceRuntime 统一调度已启用服务的初始化、启动和停止逻辑。
type ServiceRuntime struct {
	services []service
}

type service interface {
	Name() string
	Start(chan<- error)
	Shutdown(context.Context) error
	Endpoints() []string
	Status() ServiceStatus
}

// ServiceStatus 表示单个服务的运行时状态快照。
type ServiceStatus struct {
	Name      string   `json:"name"`
	Enabled   bool     `json:"enabled"`
	Listening []string `json:"listening,omitempty"`
}

// NewServiceRuntime 根据当前配置构建服务运行时。
func NewServiceRuntime() (*ServiceRuntime, error) {
	appConfig := config.GetApplication()
	if !appConfig.HTTP.Enabled {
		return nil, errors.New("invalid service config: application.http.enabled must be true")
	}

	var engine *gin.Engine
	if appConfig.HTTP.Enabled {
		engine = NewRouter()
	}
	if err := app.Bootstrap(engine); err != nil {
		return nil, err
	}

	services := make([]service, 0, 1)
	if appConfig.HTTP.Enabled {
		services = append(services, NewHTTPService(appConfig.HTTP, engine))
	}

	return &ServiceRuntime{services: services}, nil
}

// Start 启动已启用的服务。
func (r *ServiceRuntime) Start() <-chan error {
	errCh := make(chan error, len(r.services))
	logCoref("starting %d enabled service(s)", len(r.services))
	for _, service := range r.services {
		service.Start(errCh)
	}
	return errCh
}

// PrintEndpoints 输出当前已启用服务的监听地址。
func (r *ServiceRuntime) PrintEndpoints() {
	fmt.Println("Server run at:")
	for _, service := range r.services {
		for _, endpoint := range service.Endpoints() {
			fmt.Printf("-  %s\r\n", endpoint)
		}
	}
	fmt.Println("Enter Control + C Shutdown Server")
}

// Shutdown 优雅停止所有已启用服务。
func (r *ServiceRuntime) Shutdown(ctx context.Context) error {
	var shutdownErr error
	logCoref("shutting down %d enabled service(s)", len(r.services))
	for _, service := range r.services {
		if err := service.Shutdown(ctx); err != nil {
			shutdownErr = errors.Join(shutdownErr, err)
		}
	}
	return shutdownErr
}

// ServiceStatuses 返回所有已启用服务的状态快照。
func (r *ServiceRuntime) ServiceStatuses() []ServiceStatus {
	statuses := make([]ServiceStatus, 0, len(r.services))
	for _, service := range r.services {
		statuses = append(statuses, service.Status())
	}
	return statuses
}

func joinEndpoint(label string, value string) string {
	return fmt.Sprintf("%-8s %s", strings.TrimSpace(label)+":", value)
}
