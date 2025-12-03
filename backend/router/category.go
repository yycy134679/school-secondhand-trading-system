package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/category"
	"github.com/yycy134679/school-secondhand-trading-system/backend/middleware"
)

// SetupCategoryRoutes 设置分类相关路由
func SetupCategoryRoutes(engine *gin.Engine, categoryController *category.CategoryController) {
	// API路由组
	api := engine.Group("/api/v1")

	// 前台公开接口（无需认证）
	public := api.Group("/")
	{
		// 获取所有分类
		public.GET("/categories", categoryController.ListCategories)
	}

	// 管理端接口（需要认证）
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware())  // 使用认证中间件
	admin.Use(middleware.AdminMiddleware()) // 使用管理员权限中间件
	{
		// 创建分类
		admin.POST("/categories", categoryController.CreateCategory)
		// 更新分类
		admin.PUT("/categories/:id", categoryController.UpdateCategory)
		// 删除分类
		admin.DELETE("/categories/:id", categoryController.DeleteCategory)
	}
}
