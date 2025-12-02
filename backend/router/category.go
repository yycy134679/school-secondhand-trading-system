package router

import (
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/category"

	"github.com/gin-gonic/gin"
)

// RegisterCategoryRoutes 注册分类相关路由
func RegisterCategoryRoutes(router *gin.RouterGroup, categoryController *category.CategoryController, adminMiddleware gin.HandlerFunc) {
	// 公开接口路由组
	publicGroup := router.Group("/")
	// 注册分类公开接口路由
	categoryController.RegisterRoutes(publicGroup)

	// 为管理员接口应用管理员中间件
	adminGroup := router.Group("/admin/categories")
	adminGroup.Use(adminMiddleware)
	// 重新注册管理员相关的分类路由
	adminGroup.POST("", categoryController.CreateCategory)
	adminGroup.PUT("/:id", categoryController.UpdateCategory)
	adminGroup.DELETE("/:id", categoryController.DeleteCategory)
}
