package middleware

import (
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里是简化的认证逻辑，实际项目中需要实现完整的JWT验证等
		// 暂时模拟一个已登录的用户
		c.Set("user_id", "1")
		c.Set("role", "user")
		c.Next()
	}
}