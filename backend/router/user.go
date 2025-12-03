// Package router 提供用户模块的路由配置
// 负责注册用户相关的所有路由和中间件
package router

import (
	"github.com/gin-gonic/gin"
	userController "github.com/yycy134679/school-secondhand-trading-system/backend/controller/user"
	userService "github.com/yycy134679/school-secondhand-trading-system/backend/service/user"
)

// SetupUserRoutes 设置用户模块路由
//
// 参数：
//   - rg: 父路由组，通常是 /api/v1
//   - userService: 用户服务实例
//
// 功能：
//  1. 注册用户模块的所有路由
//  2. 配置认证中间件
//  3. 处理路由与控制器的绑定
func SetupUserRoutes(rg *gin.RouterGroup, userService *userService.UserService) {
	// 调用控制器的路由注册函数
	userController.RegisterRoutes(rg, userService)
}
