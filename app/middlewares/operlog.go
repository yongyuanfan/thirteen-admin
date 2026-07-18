package middlewares

import (
	"bytes"
	"io"
	"strings"
	"time"

	"thirteen-admin/app/models"
	"thirteen-admin/app/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OperLog 记录写操作日志。
func OperLog(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		if method == "GET" || method == "OPTIONS" || method == "HEAD" {
			ctx.Next()
			return
		}
		path := ctx.Request.URL.Path
		if strings.HasPrefix(path, "/common/upload") || strings.HasPrefix(path, "/common/download") {
			ctx.Next()
			return
		}
		var body []byte
		if ctx.Request.Body != nil {
			body, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		start := time.Now()
		ctx.Next()
		cost := time.Since(start).Milliseconds()

		status := 0
		if ctx.Writer.Status() >= 400 {
			status = 1
		}
		param := string(body)
		if len(param) > 2000 {
			param = param[:2000]
		}
		now := time.Now()
		_ = db.Create(&models.OperLog{
			Title:         path,
			BusinessType:  mapBusinessType(method),
			Method:        path,
			RequestMethod: method,
			OperatorType:  1,
			OperName:      services.ContextUserName(ctx),
			OperURL:       path,
			OperIP:        ctx.ClientIP(),
			OperParam:     param,
			Status:        status,
			OperTime:      &now,
			CostTime:      cost,
		}).Error
	}
}

func mapBusinessType(method string) int {
	switch method {
	case "POST":
		return 1
	case "PUT":
		return 2
	case "DELETE":
		return 3
	default:
		return 0
	}
}
