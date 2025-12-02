package router

import (
	"github.com/gin-gonic/gin"

	"school-secondhand-trading-system/controller/tag"
	"school-secondhand-trading-system/middleware"
)

// SetupTagRoutes 设置标签相关路由
func SetupTagRoutes(engine *gin.Engine, tagController *tag.TagController) {
	// API路由组
	api := engine.Group("/api/v1")

	// 前台公开接口（无需认证）
	public := api.Group("/")
	{
		// 获取所有标签
		public.GET("/tags", tagController.ListTags)
	}

	// 管理端接口（需要认证）
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware()) // 使用认证中间件
	admin.Use(middleware.AdminMiddleware()) // 使用管理员权限中间件
	{
		// 创建标签
		admin.POST("/tags", tagController.CreateTag)
		// 更新标签
		admin.PUT("/tags/:id", tagController.UpdateTag)
		// 删除标签
		admin.DELETE("/tags/:id", tagController.DeleteTag)
	}
}
