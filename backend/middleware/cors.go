// Package middleware 提供HTTP请求处理的中间件
// cors.go 实现了跨域资源共享(CORS)中间件
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 创建跨域资源共享(CORS)中间件
// 功能说明：
// - 允许所有跨域请求（开发环境）
// - 设置必要的CORS响应头
// - 处理预检请求(OPTIONS)
// 
// 返回值：
// - gin.HandlerFunc: CORS中间件处理函数
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源访问
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		
		// 允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		// 允许暴露的响应头
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		
		// 允许凭证（如cookies）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}
