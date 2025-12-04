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

// OptionalAuthMiddleware 可选的认证中间件
// 如果有token则验证并设置用户信息，没有token则继续执行但不设置用户信息
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从Header获取token
		token := c.GetHeader("Authorization")

		// 如果有token，则验证并设置用户信息
		// 这里是简化的逻辑，实际项目中需要实现完整的JWT验证
		if token != "" {
			// 模拟验证成功，设置用户ID
			c.Set("user_id", "1")
			c.Set("role", "user")
		}

		// 无论是否有token都继续执行
		c.Next()
	}
}
