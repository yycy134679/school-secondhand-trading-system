// Package product 提供商品模块的HTTP控制器
// 负责处理商品相关的HTTP请求（发布、搜索、详情、状态管理等）
package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yycy134679/school-secondhand-trading-system/backend/middleware"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/product"
)

// SearchRequest represents search parameters for product search
type SearchRequest struct {
	Keyword            string  `form:"keyword"`
	CategoryID         *int64  `form:"categoryId"`
	TagIDs             []int64 `form:"tagIds"`
	ConditionIDs       []int64 `form:"conditionIds"`
	MinPrice           *int64  `form:"minPrice"`
	MaxPrice           *int64  `form:"maxPrice"`
	PublishedTimeRange string  `form:"publishedTimeRange"`
	Sort               string  `form:"sort"`
	Page               int     `form:"page" binding:"min=1"`
	PageSize           int     `form:"pageSize" binding:"min=1,max=100"`
}

// toSearchParams converts SearchRequest to model.SearchParams
func (r *SearchRequest) toSearchParams() model.SearchParams {
	params := model.SearchParams{
		Keyword:  r.Keyword,
		Page:     r.Page,
		PageSize: r.PageSize,
	}

	// 处理指针类型转换
	if r.CategoryID != nil {
		params.CategoryID = *r.CategoryID
	}
	if r.MinPrice != nil {
		params.MinPrice = *r.MinPrice
	}
	if r.MaxPrice != nil {
		params.MaxPrice = *r.MaxPrice
	}

	// 只使用第一个标签ID和条件ID（因为模型只支持单个）
	if len(r.TagIDs) > 0 {
		params.TagID = r.TagIDs[0]
	}
	if len(r.ConditionIDs) > 0 {
		params.ConditionID = r.ConditionIDs[0]
	}

	// 处理排序参数
	if r.Sort != "" {
		// 假设r.Sort格式为"created_at:desc"或"price:asc"
		params.SortBy = "created_at"
		params.SortOrder = "desc"
	}

	return params
}

// ProductController 处理商品相关请求

type ProductController struct {
	productService *product.ProductService
}

// NewProductController 创建商品控制器实例
func NewProductController(productService *product.ProductService) *ProductController {
	if productService == nil {
		panic("product service cannot be nil")
	}
	return &ProductController{
		productService: productService,
	}
}

// RegisterRoutes 注册商品模块的所有路由
//
// 参数：
//   - rg: 父路由组，通常是 /api/v1
func RegisterRoutes(rg *gin.RouterGroup) {
	// 注意：在实际应用中，应该从依赖注入容器获取ProductService实例
	// 这里简化处理，仅注册路由结构
	// 后续应通过依赖注入框架提供正确的服务实例

	// 创建商品路由组，前缀为 /products
	// 最终路径为：/api/v1/products/*
	p := rg.Group("/products")
	{
		// ============ 公开接口（无需登录）============

		// GET /api/v1/products/:id - 获取商品详情
		p.GET("/:id", GetProductDetailHandler)

		// GET /api/v1/products/search - 搜索商品
		p.GET("/search", SearchProductsHandler)

		// GET /api/v1/products/category/:categoryId - 获取分类商品列表
		p.GET("/category/:categoryId", GetProductsByCategoryHandler)

		// ============ 需要登录的接口 ============
		// 添加AuthMiddleware中间件
		authorized := p.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			authorized.POST("", CreateProductHandler)                    // 发布商品
			authorized.PUT("/:id", UpdateProductHandler)                 // 编辑商品
			authorized.GET("/my", GetMyProductsHandler)                  // 我的发布
			authorized.POST("/:id/status", ChangeProductStatusHandler)   // 修改状态
			authorized.POST("/:id/status/undo", UndoStatusChangeHandler) // 撤销状态

			// ============ 图片管理接口 ============
			// POST /api/v1/products/:id/images - 上传图片
			authorized.POST("/:id/images", UploadProductImageHandler)
			// PUT /api/v1/products/:id/images/:imageId/primary - 设置主图
			authorized.PUT("/:id/images/:imageId/primary", SetPrimaryImageHandler)
			// PATCH /api/v1/products/:id/images/:imageId - 更新排序
			authorized.PATCH("/:id/images/:imageId", UpdateImageSortHandler)
			// DELETE /api/v1/products/:id/images/:imageId - 删除图片
			authorized.DELETE("/:id/images/:imageId", DeleteProductImageHandler)
		}
	}
}

// 以下是临时的处理器函数，实际应用中应通过依赖注入框架设置
// 这些函数将在应用启动时替换为实际的控制器实例方法
var (
	GetProductDetailHandler      = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
	SearchProductsHandler        = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
	GetProductsByCategoryHandler = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
	CreateProductHandler         = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
	UpdateProductHandler         = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
	GetMyProductsHandler         = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
	ChangeProductStatusHandler   = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
	UndoStatusChangeHandler      = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
	// 图片管理相关处理器
	UploadProductImageHandler = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
	SetPrimaryImageHandler    = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
	UpdateImageSortHandler    = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
	DeleteProductImageHandler = func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
)

// SetProductController 设置商品控制器实例，用于依赖注入
func SetProductController(controller *ProductController) {
	GetProductDetailHandler = controller.GetProductDetail
	SearchProductsHandler = controller.SearchProducts
	GetProductsByCategoryHandler = controller.GetProductsByCategory
	CreateProductHandler = controller.CreateProduct
	UpdateProductHandler = controller.UpdateProduct
	GetMyProductsHandler = controller.GetMyProducts
	ChangeProductStatusHandler = controller.ChangeProductStatus
	UndoStatusChangeHandler = controller.UndoStatusChange
	// 图片管理相关处理器暂时未实现，待实现common/util/file.go后再设置
}

// GetProductDetail 处理GET /api/v1/products/:id请求（公开接口）
// 根据是否登录传入viewerID
func (pc *ProductController) GetProductDetail(c *gin.Context) {
	// 从路径参数获取商品ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品ID格式错误"})
		return
	}

	// 检查用户是否登录，获取viewerID（如果已登录）
	var viewerID *int64 // 默认未登录为nil
	userID, exists := c.Get("userID")
	if exists {
		uid := userID.(int64)
		viewerID = &uid
	}

	// 调用服务层获取商品详情
	productDetail, err := pc.productService.GetProductDetail(c.Request.Context(), productID, viewerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if productDetail == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    productDetail,
	})
}

// GetProductsByCategory 处理GET /api/v1/products/category/:categoryId请求（公开接口）
// 根据分类ID获取商品列表
func (pc *ProductController) GetProductsByCategory(c *gin.Context) {
	// 从路径参数获取分类ID
	categoryIDStr := c.Param("categoryId")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil || categoryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类ID格式错误"})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 参数校验
	if page < 1 || pageSize < 1 || pageSize > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：page必须大于等于1，pageSize必须在1-100之间"})
		return
	}

	// 调用服务层获取分类商品
	products, total, err := pc.productService.ListByCategory(c.Request.Context(), categoryID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items":    products,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// SearchProducts 处理GET /api/v1/products/search请求（公开接口）
// 解析query参数，调用Search服务
func (pc *ProductController) SearchProducts(c *gin.Context) {
	// 使用当前包定义的SearchRequest结构体
	var req SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 参数校验
	if req.Page < 1 || req.PageSize < 1 || req.PageSize > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：page必须大于等于1，pageSize必须在1-100之间"})
		return
	}

	// 转换为repository参数并调用Search方法
	searchParams := req.toSearchParams()
	products, total, err := pc.productService.Search(c.Request.Context(), searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items":    products,
			"total":    total,
			"page":     req.Page,
			"pageSize": req.PageSize,
		},
	})
}

// CreateProduct 处理POST /api/v1/products请求
// 解析 multipart 表单，调用 CreateProduct 服务
func (pc *ProductController) CreateProduct(c *gin.Context) {
	// 从上下文获取用户ID（假设通过中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	sellerID := userID.(int64)

	// 解析表单数据
	title := c.PostForm("title")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	categoryIDStr := c.PostForm("categoryId")
	conditionIDStr := c.PostForm("conditionId")
	tagIDsStr := c.PostFormArray("tagIds")

	// 转换数据类型
	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "价格格式错误"})
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类ID格式错误"})
		return
	}

	conditionID, err := strconv.ParseInt(conditionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品状态ID格式错误"})
		return
	}

	// 解析标签ID
	var tagIDs []int64
	for _, tagIDStr := range tagIDsStr {
		tagID, parseErr := strconv.ParseInt(tagIDStr, 10, 64)
		if parseErr != nil {
			continue // 跳过无效的标签ID
		}
		tagIDs = append(tagIDs, tagID)
	}

	// 获取上传的图片文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析表单失败"})
		return
	}

	files := form.File["images"]
	var imageURLs []string

	// 处理图片上传（实际项目中应该保存文件到存储服务并返回URL）
	for _, file := range files {
		// 这里简化处理，实际应该保存文件并返回URL
		imageURLs = append(imageURLs, "/uploads/"+file.Filename)
	}

	// 构建请求参数
	req := product.CreateProductRequest{
		Title:       title,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
		ConditionID: conditionID,
		TagIDs:      tagIDs,
		Images:      imageURLs,
	}

	// 调用服务层创建商品
	product, err := pc.productService.CreateProduct(c.Request.Context(), sellerID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})
}

// UpdateProduct 处理PUT /api/v1/products/:id请求
// 解析 JSON，调用 UpdateProduct 服务
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	// 从上下文获取用户ID（假设通过中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	updaterID := userID.(int64)

	// 从路径参数获取商品ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品ID格式错误"})
		return
	}

	// 解析JSON请求体
	var req product.UpdateProductRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误"})
		return
	}

	// 调用服务层更新商品
	updatedProduct, err := pc.productService.UpdateProduct(c.Request.Context(), updaterID, productID, req, false) // 添加isAdmin参数
	if err != nil {
		// 处理错误码
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedProduct,
	})
}

// ChangeProductStatus 处理POST /api/v1/products/:id/status请求
// 解析action参数，调用ChangeStatus服务
func (pc *ProductController) ChangeProductStatus(c *gin.Context) {
	// 从上下文获取用户ID（假设通过中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	operatorID := userID.(int64)

	// 从路径参数获取商品ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品ID格式错误"})
		return
	}

	// 解析请求体中的action参数
	var req struct {
		Action string `json:"action" binding:"required"`
	}
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误，缺少action参数"})
		return
	}

	// 调用服务层改变商品状态
	err = pc.productService.ChangeStatus(c.Request.Context(), operatorID, productID, req.Action)
	if err != nil {
		// 处理错误码
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 重新获取更新后的商品信息
	updatedProduct, _ := pc.productService.GetProductDetail(c.Request.Context(), productID, &operatorID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedProduct,
	})
}

// UndoStatusChange 处理POST /api/v1/products/:id/status/undo请求
// 调用UndoLastStatusChange服务撤销最近的状态变更
func (pc *ProductController) UndoStatusChange(c *gin.Context) {
	// 从上下文获取用户ID（假设通过中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	operatorID := userID.(int64)

	// 从路径参数获取商品ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品ID格式错误"})
		return
	}

	// 调用服务层撤销状态变更
	err = pc.productService.UndoLastStatusChange(c.Request.Context(), operatorID, productID)
	if err != nil {
		// 处理错误码
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 重新获取更新后的商品信息
	updatedProduct, _ := pc.productService.GetProductDetail(c.Request.Context(), productID, &operatorID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedProduct,
	})
}

// GetMyProducts 处理GET /api/v1/products/my请求（登录接口）
// 调用ListMyProducts服务获取用户发布的商品列表
func (pc *ProductController) GetMyProducts(c *gin.Context) {
	// 从上下文获取用户ID（假设通过中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	sellerID := userID.(int64)

	// 解析分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 调用服务层获取我的商品列表
	products, total, err := pc.productService.ListMyProducts(c.Request.Context(), sellerID, "", page, pageSize) // 添加空状态过滤参数
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items":    products,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}
