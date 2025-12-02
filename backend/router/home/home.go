// Package home 处理首页相关的路由注册
package home

import (
	"github.com/gin-gonic/gin"
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/home"
)

// RegisterHomeRoutes 注册首页相关路由
// GET /api/v1/home 获取首页数据
// GET /api/v1/recommendations 获取个性化推荐
func RegisterHomeRoutes(router *gin.RouterGroup, controller *home.HomeController, optionalAuth gin.HandlerFunc) {
	// 首页路由 - 可选认证（用于获取个性化推荐）
	router.GET("/home", optionalAuth, controller.GetHomeData)
	// 个性化推荐路由 - 需要登录
	router.GET("/recommendations", optionalAuth, controller.GetRecommendations)
}
