// Package product 提供商品模块的HTTP控制器
// 负责处理商品相关的HTTP请求（发布、搜索、详情、状态管理等）
package product

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册商品模块的所有路由
//
// 路由列表（待完善）：
//   公开接口：
//   GET    /products/:id            - 获取商品详情
//   GET    /products/search         - 搜索商品
//   GET    /products/category/:id   - 按分类浏览商品
//
//   需要登录：
//   POST   /products                - 发布商品
//   PUT    /products/:id            - 编辑商品
//   GET    /products/my             - 我的发布
//   POST   /products/:id/status     - 修改商品状态（上架/下架/已售）
//   POST   /products/:id/status/undo - 撤销最近的状态变更（3秒窗口）
//   GET    /products/:id/contact    - 获取卖家联系方式
//   POST   /products/:id/images     - 上传商品图片
//   DELETE /products/:id/images/:imageId - 删除商品图片
//
// 参数：
//   - rg: 父路由组，通常是 /api/v1
//
// 商品状态机说明：
//   ForSale（在售） --下架--> Delisted（已下架）
//   Delisted（已下架） --上架--> ForSale（在售）
//   ForSale（在售） --标记已售--> Sold（已售出）
//   Sold（已售出）是终态，不可再改变
//
// 撤销功能说明：
//   - 用户执行下架/上架操作后，3秒内可以撤销
//   - 撤销信息存储在Redis中，过期自动删除
//   - 已售状态不支持撤销
//
// TODO: 完整实现步骤
//   1. 创建ProductController结构体，注入ProductService依赖
//   2. 实现各个Handler方法
//   3. 在需要登录的路由上应用AuthMiddleware
//   4. 实现图片上传和文件存储
//   5. 实现状态机和撤销逻辑
func RegisterRoutes(rg *gin.RouterGroup) {
	// 创建商品路由组，前缀为 /products
	// 最终路径为：/api/v1/products/*
	p := rg.Group("/products")
	{
		// ============ 公开接口（无需登录）============

		// GET /api/v1/products/:id - 获取商品详情
		// 参数：
		//   - id: 商品ID（路径参数）
		//   - viewerID: 浏览者ID（可选，用于记录浏览历史）
		// 响应：
		// {
		//   "code": 0,
		//   "message": "ok",
		//   "data": {
		//     "id": 123,
		//     "title": "二手教材 高等数学",
		//     "description": "九成新，无笔记",
		//     "price": 30,
		//     "images": [...],
		//     "seller": {...},
		//     "condition": "九成新",
		//     "tags": ["教材", "几乎全新"],
		//     "status": "ForSale",
		//     "createdAt": "2024-01-01T10:00:00Z"
		//   }
		// }
		p.GET("/:id", func(c *gin.Context) {
			c.JSON(200, gin.H{"msg": "product detail stub - 待实现"})
			// TODO: 实现商品详情逻辑
		})

		// GET /api/v1/products/search - 搜索商品
		// 查询参数：
		//   - keyword: 关键词（标题、描述）
		//   - categoryId: 分类ID
		//   - conditionId: 新旧程度ID
		//   - minPrice: 最低价格
		//   - maxPrice: 最高价格
		//   - tagIds: 标签ID列表（逗号分隔）
		//   - page: 页码（默认1）
		//   - pageSize: 每页条数（默认20）
		// 响应：
		// {
		//   "code": 0,
		//   "message": "ok",
		//   "data": {
		//     "items": [...],
		//     "total": 100,
		//     "page": 1,
		//     "pageSize": 20
		//   }
		// }
		p.GET("/search", func(c *gin.Context) {
			c.JSON(200, gin.H{"msg": "search stub - 待实现"})
			// TODO: 实现商品搜索逻辑
		})

		// ============ 需要登录的接口 ============
		// TODO: 添加AuthMiddleware中间件
		// authorized := p.Group("")
		// authorized.Use(middleware.AuthMiddleware())
		// {
		//     authorized.POST("", handleCreateProduct)          // 发布商品
		//     authorized.PUT("/:id", handleUpdateProduct)       // 编辑商品
		//     authorized.GET("/my", handleMyProducts)           // 我的发布
		//     authorized.POST("/:id/status", handleChangeStatus) // 修改状态
		//     authorized.POST("/:id/status/undo", handleUndoStatus) // 撤销状态
		//     authorized.GET("/:id/contact", handleGetContact)  // 获取联系方式
		//     authorized.POST("/:id/images", handleUploadImage) // 上传图片
		// }
	}
}
