package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/auth"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/errors"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
)

// AuthMiddleware 验证请求中的 JWT，失败则返回未登录错误
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			resp.Error(c, errors.CodeUnauthenticated, "请先登录")
			c.Abort()
			return
		}

		token := authHeader
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token = strings.TrimSpace(authHeader[7:])
		}
		if token == "" {
			resp.Error(c, errors.CodeUnauthenticated, "请先登录")
			c.Abort()
			return
		}

		userID, err := auth.ParseToken(token)
		if err != nil {
			resp.Error(c, errors.CodeUnauthenticated, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		c.Set("user_id", strconv.FormatInt(userID, 10))
		c.Set("role", "user")
		c.Next()
	}
}

// OptionalAuthMiddleware 允许匿名访问，有 token 则解析并注入用户信息
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.Next()
			return
		}

		token := authHeader
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token = strings.TrimSpace(authHeader[7:])
		}
		if token == "" {
			c.Next()
			return
		}

		if userID, err := auth.ParseToken(token); err == nil {
			c.Set("user_id", strconv.FormatInt(userID, 10))
			c.Set("role", "user")
		}

		c.Next()
	}
}
