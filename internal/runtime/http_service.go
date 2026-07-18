package runtime

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"thirteen-admin/internal/config"
)

// HTTPService 封装 HTTP 服务的构建、启动和停止逻辑。
type HTTPService struct {
	httpConfig config.HTTP
	server     *http.Server
}

// Name 返回服务名称。
func (s *HTTPService) Name() string {
	return "http"
}

// NewHTTPService 创建 HTTP 服务实例。
func NewHTTPService(httpConfig config.HTTP, handler http.Handler) *HTTPService {
	return &HTTPService{
		httpConfig: httpConfig,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", httpConfig.Host, httpConfig.Port),
			Handler: handler,
		},
	}
}

// Start 启动 HTTP 服务监听。
func (s *HTTPService) Start(errCh chan<- error) {
	logServicef(s.Name(), "starting server on %s", s.server.Addr)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http server listen failed on %s: %w", s.server.Addr, err)
		}
	}()
}

// Shutdown 优雅停止 HTTP 服务。
func (s *HTTPService) Shutdown(ctx context.Context) error {
	logServicef(s.Name(), "stopping server on %s", s.server.Addr)
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http shutdown failed on %s: %w", s.server.Addr, err)
	}
	logServicef(s.Name(), "server stopped on %s", s.server.Addr)
	return nil
}

// Endpoints 返回 HTTP 服务地址信息。
func (s *HTTPService) Endpoints() []string {
	return []string{
		joinEndpoint("Local", fmt.Sprintf("http://localhost:%d/", s.httpConfig.Port)),
		joinEndpoint("Network", fmt.Sprintf("http://%s:%d/", s.httpConfig.Host, s.httpConfig.Port)),
	}
}

// Status 返回 HTTP 服务状态快照。
func (s *HTTPService) Status() ServiceStatus {
	return ServiceStatus{
		Name:      s.Name(),
		Enabled:   true,
		Listening: []string{s.server.Addr},
	}
}
