package middlewares

import (
	"thirteen-admin/app/runtime"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// WithDB 将请求级 DB（WithContext）写入 gin.Context。
func WithDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := runtime.DB()
		if db != nil {
			c.Set(handler.ContextDBKey, db.WithContext(c.Request.Context()))
		}
		c.Next()
	}
}
