package runtime

import (
	"thirteen-admin/internal/config"

	"github.com/gin-gonic/gin"
)

// NewRouter 创建 HTTP 服务路由引擎。
func NewRouter() *gin.Engine {
	appConfig := config.GetApplication()
	if appConfig.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	return router
}
