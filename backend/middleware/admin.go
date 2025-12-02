package middleware

import (
	"github.com/gin-gonic/gin"
)

// AdminMiddleware 管理员权限验证中间件（待完善）
//
// 功能说明：
//   - 验证当前登录用户是否具有管理员权限
//   - 必须在 AuthMiddleware 之后使用（依赖认证中间件设置的用户信息）
//   - 如果用户不是管理员，返回403错误并中止请求
//
// 使用场景：
//   - 后台管理接口（如用户管理、商品审核、分类/标签管理等）
//   - 需要管理员权限才能执行的操作（如删除用户、编辑已售商品等）
//
// 工作流程：
//  1. 从Gin上下文中获取 isAdmin 标志
//     （这个标志由AuthMiddleware在验证JWT时设置）
//  2. 检查 isAdmin 是否为 true
//  3. 如果是管理员：
//     - 调用c.Next()继续执行后续处理器
//  4. 如果不是管理员：
//     - 返回403错误：resp.Error(c, errors.CodeForbidden, "需要管理员权限")
//     - 调用c.Abort()中止请求
//
// 使用示例：
//
//	// 管理员路由组
//	admin := r.Group("/api/v1/admin")
//	admin.Use(middleware.AuthMiddleware())   // 先验证登录
//	admin.Use(middleware.AdminMiddleware())  // 再验证管理员权限
//	{
//	    admin.GET("/users", adminController.ListUsers)           // 用户列表
//	    admin.DELETE("/users/:id", adminController.DeleteUser)   // 删除用户
//	    admin.PUT("/products/:id", adminController.EditProduct)  // 编辑商品
//	}
//
// 注意事项：
//   - AdminMiddleware 必须在 AuthMiddleware 之后使用
//   - 如果没有先执行 AuthMiddleware，上下文中不会有 isAdmin 信息
//   - 普通用户尝试访问管理员接口时，会收到错误码 1003（CodeForbidden）
//
// TODO: 完整实现
//  1. 从上下文获取 isAdmin 标志
//  2. 验证是否为管理员
//  3. 根据验证结果决定继续或中止请求
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现管理员权限检查
		// 伪代码示例：
		//
		// // 1. 从上下文获取isAdmin标志
		// isAdmin, exists := c.Get("isAdmin")
		// if !exists {
		//     // 如果没有isAdmin信息，说明AuthMiddleware未执行
		//     resp.Error(c, errors.CodeUnauthenticated, "请先登录")
		//     c.Abort()
		//     return
		// }
		//
		// // 2. 检查是否为管理员
		// if !isAdmin.(bool) {
		//     // 不是管理员，返回403错误
		//     resp.Error(c, errors.CodeForbidden, "需要管理员权限")
		//     c.Abort()
		//     return
		// }
		//
		// // 3. 是管理员，继续执行
		c.Next()
	}
}
