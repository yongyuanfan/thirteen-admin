package middlewares

import (
	"github.com/gin-gonic/gin"
)

// ResponseStream 流式输出
func ResponseStream() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.Writer.Header()
		header.Set("Content-Type", "text/event-stream")
		header.Set("Cache-Control", "no-cache")
		header.Set("Connection", "keep-alive")
		ctx.Next()
	}
}
