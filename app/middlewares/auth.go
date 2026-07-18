package middlewares

import (
	"errors"
	"strings"

	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// Gin Context 中注入的鉴权字段 key。
const (
	ContextUserIDKey   = "authUserID"
	ContextUserNameKey = "authUserName"
)

// Auth JWT 鉴权中间件。
func Auth(jwtManager *services.JWTManager) gin.HandlerFunc {
	h := &handler.Handler{}
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			h.Fail(ctx, handler.CodeUnauthorized, "未登录")
			ctx.Abort()
			return
		}
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			h.Fail(ctx, handler.CodeUnauthorized, "未登录")
			ctx.Abort()
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(header, prefix))
		if tokenStr == "" {
			h.Fail(ctx, handler.CodeUnauthorized, "未登录")
			ctx.Abort()
			return
		}

		claims, err := jwtManager.ParseAccess(tokenStr)
		if err != nil {
			if errors.Is(err, services.ErrExpiredToken) {
				h.Fail(ctx, handler.CodeTokenExpired, "令牌已过期")
			} else {
				h.Fail(ctx, handler.CodeUnauthorized, "未登录")
			}
			ctx.Abort()
			return
		}

		ctx.Set(ContextUserIDKey, claims.UserID)
		ctx.Set(ContextUserNameKey, claims.UserName)
		ctx.Next()
	}
}
