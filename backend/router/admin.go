// Package router 负责HTTP路由的初始化和配置
// 统一管理所有API端点和中间件
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/admin"
	"github.com/yycy134679/school-secondhand-trading-system/backend/middleware"
)

// RegisterAdminRoutes 注册管理后台相关路由
//
// 功能说明：
//   - 配置管理后台的路由组和中间件
//   - 注册所有管理员接口
//   - 确保接口安全访问
//
// 参数：
//   - api: API v1版本的路由组
//   - dashboardController: 仪表盘控制器实例
//   - userController: 用户管理控制器实例
//   - productController: 商品管理控制器实例
//   - categoryController: 分类管理控制器实例
//   - tagController: 标签管理控制器实例
//   - adminMiddleware: 管理员权限验证中间件
func RegisterAdminRoutes(api *gin.RouterGroup,
	dashboardController *admin.DashboardController,
	userController *admin.UserController,
	productController *admin.ProductController,
	categoryController *admin.CategoryController,
	tagController *admin.TagController,
	adminMiddleware gin.HandlerFunc) {
	// 创建管理员路由组
	// 路径前缀：/api/v1/admin
	adminGroup := api.Group("/admin")

	// 应用认证中间件（必须先登录）
	adminGroup.Use(middleware.AuthMiddleware())

	// 应用管理员权限验证中间件（必须是管理员）
	adminGroup.Use(adminMiddleware)

	// 注册仪表盘相关接口
	// GET /api/v1/admin/dashboard - 获取仪表盘统计数据
	adminGroup.GET("/dashboard", dashboardController.GetDashboard)

	// 注册用户管理相关接口
	// GET /api/v1/admin/users - 获取用户列表
	adminGroup.GET("/users", userController.ListUsers)

	// 注册商品管理相关接口
	// GET /api/v1/admin/products - 获取商品列表
	adminGroup.GET("/products", productController.ListProducts)
	// PUT /api/v1/admin/products/:id - 更新商品信息
	adminGroup.PUT("/products/:id", productController.UpdateProduct)

	// 注册分类管理相关接口
	// GET /api/v1/admin/categories - 获取分类列表
	adminGroup.GET("/categories", categoryController.ListCategories)
	// POST /api/v1/admin/categories - 创建分类
	adminGroup.POST("/categories", categoryController.CreateCategory)
	// PUT /api/v1/admin/categories/:id - 更新分类
	adminGroup.PUT("/categories/:id", categoryController.UpdateCategory)
	// DELETE /api/v1/admin/categories/:id - 删除分类
	adminGroup.DELETE("/categories/:id", categoryController.DeleteCategory)

	// 注册标签管理相关接口
	// GET /api/v1/admin/tags - 获取标签列表
	adminGroup.GET("/tags", tagController.ListTags)
	// POST /api/v1/admin/tags - 创建标签
	adminGroup.POST("/tags", tagController.CreateTag)
	// PUT /api/v1/admin/tags/:id - 更新标签
	adminGroup.PUT("/tags/:id", tagController.UpdateTag)
	// DELETE /api/v1/admin/tags/:id - 删除标签
	adminGroup.DELETE("/tags/:id", tagController.DeleteTag)
}
