package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/product"
	"github.com/yycy134679/school-secondhand-trading-system/backend/middleware"
)

// SetupProductRoutes 设置商品相关路由
func SetupProductRoutes(engine *gin.Engine, productController *product.ProductController, imageController *product.ImageController) {
	// API路由组
	api := engine.Group("/api/v1")

	// 公开接口（无需认证，但可选登录以记录浏览）
	public := api.Group("/")
	{
		// 获取商品详情 - 使用可选认证中间件以便记录浏览
		public.GET("/products/:id", middleware.OptionalAuthMiddleware(), productController.GetProductDetail)
		// 搜索商品
		public.GET("/products/search", productController.SearchProducts)
		// 获取分类商品
		public.GET("/products/category/:categoryId", productController.GetProductsByCategory)
	}

	// 需要认证的接口
	auth := api.Group("/")
	auth.Use(middleware.AuthMiddleware()) // 使用认证中间件
	{
		// 创建商品
		auth.POST("/products", productController.CreateProduct)
		// 更新商品
		auth.PUT("/products/:id", productController.UpdateProduct)
		// 变更商品状态
		auth.POST("/products/:id/status", productController.ChangeProductStatus)
		// 撤销状态变更
		auth.POST("/products/:id/status/undo", productController.UndoLastStatusChange)
		// 获取我的商品列表
		auth.GET("/products/my", productController.ListMyProducts)

		// 图片管理接口
		auth.POST("/products/:id/images", imageController.UploadProductImage)
		auth.PUT("/products/:id/images/:imageId/primary", imageController.SetPrimaryImage)
		auth.PATCH("/products/:id/images/:imageId", imageController.UpdateImageSortOrder)
		auth.DELETE("/products/:id/images/:imageId", imageController.DeleteProductImage)
	}
}
