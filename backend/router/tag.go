package router

import (
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/tag"

	"github.com/gin-gonic/gin"
)

// RegisterTagRoutes 注册标签相关路由
func RegisterTagRoutes(router *gin.RouterGroup, tagController *tag.TagController, adminMiddleware gin.HandlerFunc) {
	// 公开接口路由组
	publicGroup := router.Group("/")
	// 注册标签公开接口路由
	tagController.RegisterRoutes(publicGroup)

	// 为管理员接口应用管理员中间件
	adminGroup := router.Group("/admin/tags")
	adminGroup.Use(adminMiddleware)
	// 重新注册管理员相关的标签路由
	adminGroup.POST("", tagController.CreateTag)
	adminGroup.PUT("/:id", tagController.UpdateTag)
	adminGroup.DELETE("/:id", tagController.DeleteTag)
}
