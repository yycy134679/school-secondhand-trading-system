// Package middleware 提供HTTP中间件函数
// 中间件在请求到达业务逻辑之前或之后执行，用于处理横切关注点
package middleware

import (
	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件（待完善）
//
// 功能说明：
//   - 从请求头中提取JWT token
//   - 验证token的合法性（签名、过期时间等）
//   - 将用户信息（userID、isAdmin等）存入Gin上下文
//   - 如果验证失败，返回401错误并中止请求
//
// 使用场景：
//   - 需要登录才能访问的接口（如发布商品、修改个人信息）
//   - 可以在路由组级别应用，也可以单独应用到某个路由
//
// 工作流程：
//  1. 从请求头 "Authorization" 中获取token
//     格式：Authorization: Bearer <token>
//  2. 解析token字符串，提取JWT claims
//  3. 验证token签名是否正确（使用配置中的JWT_SECRET）
//  4. 检查token是否过期
//  5. 如果验证通过：
//     - 将userID存入上下文：c.Set("userID", userID)
//     - 将isAdmin存入上下文：c.Set("isAdmin", isAdmin)
//     - 调用c.Next()继续执行后续处理器
//  6. 如果验证失败：
//     - 返回错误响应：resp.Error(c, errors.CodeUnauthenticated, "未登录或token无效")
//     - 调用c.Abort()中止请求
//
// 使用示例：
//
//	// 在路由中应用认证中间件
//	authorized := r.Group("/api/v1")
//	authorized.Use(middleware.AuthMiddleware())
//	{
//	    authorized.POST("/products", productController.Create)  // 需要登录
//	    authorized.GET("/users/profile", userController.Profile) // 需要登录
//	}
//
// TODO: 完整实现步骤
//  1. 实现 common/auth/jwt.go 中的 ParseToken 函数
//  2. 在这里调用 ParseToken 解析和验证token
//  3. 将解析出的用户信息存入上下文
//  4. 处理各种错误情况（token缺失、格式错误、过期等）
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现JWT认证逻辑
		// 伪代码示例：
		//
		// // 1. 从请求头获取token
		// authHeader := c.GetHeader("Authorization")
		// if authHeader == "" {
		//     resp.Error(c, errors.CodeUnauthenticated, "请先登录")
		//     c.Abort()
		//     return
		// }
		//
		// // 2. 提取token（去掉"Bearer "前缀）
		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		//
		// // 3. 解析和验证token
		// claims, err := auth.ParseToken(tokenString)
		// if err != nil {
		//     resp.Error(c, errors.CodeUnauthenticated, "token无效或已过期")
		//     c.Abort()
		//     return
		// }
		//
		// // 4. 将用户信息存入上下文
		// c.Set("userID", claims.UserID)
		// c.Set("isAdmin", claims.IsAdmin)
		//
		// // 5. 继续执行后续处理器
		c.Next()
	}
}
