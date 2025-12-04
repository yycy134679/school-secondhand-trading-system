package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/recommend"
	"github.com/yycy134679/school-secondhand-trading-system/backend/middleware"
)

// SetupRecommendRoutes 设置推荐模块路由
//
// 参数：
//   - r: Gin引擎实例
//   - recommendController: 推荐控制器实例
//
// 功能：
//  1. 注册首页数据接口（公开）
//  2. 注册最近浏览接口（需要登录）
//  3. 注册记录浏览接口（需要登录）
func SetupRecommendRoutes(r *gin.Engine, recommendController *recommend.RecommendController) {
	api := r.Group("/api/v1")
	{
		// 公开接口 - 首页数据（可选登录）
		api.GET("/home", middleware.OptionalAuthMiddleware(), recommendController.GetHomeData)

		// 需要登录的接口
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			// 获取最近浏览记录
			users.GET("/recent-views", recommendController.GetRecentViews)
		}

		// 记录商品浏览（需要登录）
		products := api.Group("/products")
		products.Use(middleware.AuthMiddleware())
		{
			// 记录浏览
			products.POST("/:id/view", recommendController.RecordProductView)
		}
	}
}
