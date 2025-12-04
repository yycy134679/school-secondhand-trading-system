package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/category"
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
}
