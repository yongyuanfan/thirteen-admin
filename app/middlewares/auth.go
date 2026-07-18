package middlewares

import (
	"errors"
	"strings"

	"thirteen-admin/app/runtime"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey   = "authUserID"
	ContextUserNameKey = "authUserName"
	ContextTokenKey    = "authToken"
)

// Auth JWT 鉴权中间件（JWT 从 runtime 读取）。
func Auth() gin.HandlerFunc {
	h := &handler.Handler{}
	return func(ctx *gin.Context) {
		jwtManager := runtime.JWT()
		if jwtManager == nil {
			h.Unauthorized(ctx, "认证失败，无法访问系统资源")
			ctx.Abort()
			return
		}

		header := ctx.GetHeader("Authorization")
		if header == "" {
			h.Unauthorized(ctx, "认证失败，无法访问系统资源")
			ctx.Abort()
			return
		}
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			h.Unauthorized(ctx, "认证失败，无法访问系统资源")
			ctx.Abort()
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(header, prefix))
		if tokenStr == "" {
			h.Unauthorized(ctx, "认证失败，无法访问系统资源")
			ctx.Abort()
			return
		}

		claims, err := jwtManager.ParseAccess(tokenStr)
		if err != nil {
			if errors.Is(err, services.ErrExpiredToken) {
				h.Unauthorized(ctx, "登录状态已过期")
			} else {
				h.Unauthorized(ctx, "认证失败，无法访问系统资源")
			}
			ctx.Abort()
			return
		}

		ctx.Set(ContextUserIDKey, claims.UserID)
		ctx.Set(ContextUserNameKey, claims.UserName)
		ctx.Set(ContextTokenKey, tokenStr)
		ctx.Next()
	}
}
