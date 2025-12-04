package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/tag"
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
}
