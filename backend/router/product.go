// Package router 提供路由注册功能
package router

import (
	"github.com/gin-gonic/gin"
	productController "github.com/yycy134679/school-secondhand-trading-system/backend/controller/product"
	"github.com/yycy134679/school-secondhand-trading-system/backend/middleware"
	productService "github.com/yycy134679/school-secondhand-trading-system/backend/service/product"
)

// RegisterProductRoutes 注册商品模块的路由
func RegisterProductRoutes(rg *gin.RouterGroup, svc *productService.ProductService) {
	// 创建商品控制器实例
	controller := productController.NewProductController(svc)
	// 设置控制器处理器
	productController.SetProductController(controller)

	// 创建商品路由组
	products := rg.Group("/products")
	{
		// ============ 公开接口（无需登录）============
		// GET /api/v1/products/:id - 获取商品详情
		products.GET("/:id", productController.GetProductDetailHandler)
		// GET /api/v1/products/search - 搜索商品
		products.GET("/search", productController.SearchProductsHandler)
		// GET /api/v1/products/category/:categoryId - 获取分类商品列表
		products.GET("/category/:categoryId", productController.GetProductsByCategoryHandler)

		// ============ 需要登录的接口 ============
		authorized := products.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 商品管理
			authorized.POST("", productController.CreateProductHandler)                    // 发布商品
			authorized.PUT("/:id", productController.UpdateProductHandler)                 // 编辑商品
			authorized.GET("/my", productController.GetMyProductsHandler)                  // 我的发布
			authorized.POST("/:id/status", productController.ChangeProductStatusHandler)   // 修改状态
			authorized.POST("/:id/status/undo", productController.UndoStatusChangeHandler) // 撤销状态

			// ============ 图片管理接口 ============
			// 注：图片管理接口暂时未实现，待common/util/file.go实现后再启用
			// authorized.POST("/:id/images", product.UploadProductImageHandler)                // 上传图片
			// authorized.PUT("/:id/images/:imageId/primary", product.SetPrimaryImageHandler)   // 设置主图
			// authorized.PATCH("/:id/images/:imageId", product.UpdateImageSortHandler)         // 更新排序
			// authorized.DELETE("/:id/images/:imageId", product.DeleteProductImageHandler)     // 删除图片
		}
	}
}
