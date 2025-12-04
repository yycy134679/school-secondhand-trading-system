package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware 管理员权限中间件
// 检查用户是否为管理员角色，否则拒绝访问
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或登录已过期"})
			c.Abort()
			return
		}

		// 检查角色是否为管理员
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足，需要管理员权限"})
			c.Abort()
			return
		}

		// 权限通过，继续处理请求
		c.Next()
	}
}
