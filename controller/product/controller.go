package product

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"school-secondhand-trading-system/service/product"
	"school-secondhand-trading-system/util/response"
)

// ProductController 商品控制器
type ProductController struct {
	productService *product.ProductService
}

// NewProductController 创建商品控制器实例
func NewProductController(productService *product.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// CreateProduct 创建商品
// POST /api/v1/products
func (pc *ProductController) CreateProduct(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 401, "用户未登录")
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的用户ID")
		return
	}

	// 解析表单数据
	title := c.PostForm("title")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	categoryIDStr := c.PostForm("categoryId")
	conditionIDStr := c.PostForm("conditionId")
	tagIDsStr := c.PostForm("tagIds")

	// 基本参数校验
	if title == "" || priceStr == "" || categoryIDStr == "" || conditionIDStr == "" {
		response.Error(c, http.StatusBadRequest, 400, "标题、价格、分类ID和商品状态为必填项")
		return
	}

	// 解析数字参数
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的价格格式")
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的分类ID")
		return
	}

	conditionID, err := strconv.ParseInt(conditionIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的商品状态ID")
		return
	}

	// 解析标签ID
	var tagIDs []int64
	if tagIDsStr != "" {
		for _, idStr := range strings.Split(tagIDsStr, ",") {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				response.Error(c, http.StatusBadRequest, 400, "无效的标签ID格式")
				return
			}
			tagIDs = append(tagIDs, id)
		}
	}

	// 获取上传的文件
	files := c.Request.MultipartForm.File["images"]

	// 构建请求参数
	req := &product.CreateProductRequest{
		Title:       title,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
		ConditionID: conditionID,
		TagIDs:      tagIDs,
		Images:      files,
	}

	// 调用服务层方法
	productDTO, err := pc.productService.CreateProduct(c.Request.Context(), userID, req)
	if err != nil {
		// 根据错误类型返回不同的错误码
		if strings.Contains(err.Error(), "请先完善微信号") {
			response.Error(c, http.StatusBadRequest, 1001, err.Error())
		} else {
			response.Error(c, http.StatusBadRequest, 400, err.Error())
		}
		return
	}

	response.Success(c, http.StatusCreated, productDTO)
}

// UpdateProduct 更新商品
// PUT /api/v1/products/:id
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	// 从上下文中获取用户ID和角色
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 401, "用户未登录")
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的用户ID")
		return
	}

	// 判断是否为管理员
	isAdmin := false
	role, exists := c.Get("role")
	if exists && role == "admin" {
		isAdmin = true
	}

	// 获取商品ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的商品ID")
		return
	}

	// 解析JSON请求体
	var req product.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "请求参数格式错误")
		return
	}

	// 调用服务层方法
	productDTO, err := pc.productService.UpdateProduct(c.Request.Context(), userID, productID, &req, isAdmin)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, http.StatusOK, productDTO)
}

// ChangeProductStatus 变更商品状态
// POST /api/v1/products/:id/status
func (pc *ProductController) ChangeProductStatus(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 401, "用户未登录")
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的用户ID")
		return
	}

	// 获取商品ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的商品ID")
		return
	}

	// 解析动作参数
	type StatusRequest struct {
		Action string `json:"action" binding:"required,oneof=delist relist sold"`
	}
	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的动作参数，支持：delist, relist, sold")
		return
	}

	// 调用服务层方法
	err = pc.productService.ChangeStatus(c.Request.Context(), userID, productID, req.Action)
	if err != nil {
		// 根据错误类型返回不同的错误码
		if strings.Contains(err.Error(), "终态") {
			response.Error(c, http.StatusBadRequest, 3004, err.Error())
		} else {
			response.Error(c, http.StatusBadRequest, 400, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "状态变更成功"})
}

// UndoLastStatusChange 撤销最近的状态变更
// POST /api/v1/products/:id/status/undo
func (pc *ProductController) UndoLastStatusChange(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 401, "用户未登录")
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的用户ID")
		return
	}

	// 获取商品ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的商品ID")
		return
	}

	// 调用服务层方法
	err = pc.productService.UndoLastStatusChange(c.Request.Context(), userID, productID)
	if err != nil {
		// 根据错误类型返回不同的错误码
		if strings.Contains(err.Error(), "终态") {
			response.Error(c, http.StatusBadRequest, 3004, err.Error())
		} else if strings.Contains(err.Error(), "不存在") || strings.Contains(err.Error(), "超时") {
			response.Error(c, http.StatusBadRequest, 3005, err.Error())
		} else {
			response.Error(c, http.StatusBadRequest, 400, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "状态撤销成功"})
}

// GetProductDetail 获取商品详情
// GET /api/v1/products/:id
func (pc *ProductController) GetProductDetail(c *gin.Context) {
	// 获取商品ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的商品ID")
		return
	}

	// 获取viewerID（如果已登录）
	var viewerID *int64
	if userIDStr, exists := c.Get("user_id"); exists {
		id, err := strconv.ParseInt(userIDStr.(string), 10, 64)
		if err == nil {
			viewerID = &id
		}
	}

	// 调用服务层方法
	productDTO, err := pc.productService.GetProductDetail(c.Request.Context(), productID, viewerID)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, "商品不存在")
		return
	}

	response.Success(c, http.StatusOK, productDTO)
}

// ListMyProducts 获取我的商品列表
// GET /api/v1/products/my
func (pc *ProductController) ListMyProducts(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 401, "用户未登录")
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的用户ID")
		return
	}

	// 解析查询参数
	keyword := c.Query("keyword")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	// 调用服务层方法
	products, total, err := pc.productService.ListMyProducts(c.Request.Context(), userID, keyword, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "获取商品列表失败")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"list":  products,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// SearchProducts 搜索商品
// GET /api/v1/products/search
func (pc *ProductController) SearchProducts(c *gin.Context) {
	// 解析查询参数
	keyword := c.Query("keyword")
	categoryIDStr := c.Query("categoryId")
	minPriceStr := c.Query("minPrice")
	maxPriceStr := c.Query("maxPrice")
	conditionIDStr := c.Query("conditionId")
	tagIDStr := c.Query("tagId")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	// 构建搜索参数
	params := &product.SearchParams{
		Keyword:    keyword,
		Page:       page,
		PageSize:   pageSize,
	}

	// 解析可选参数
	if categoryIDStr != "" {
		if categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			params.CategoryID = &categoryID
		}
	}

	if minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			params.MinPrice = &minPrice
		}
	}

	if maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			params.MaxPrice = &maxPrice
		}
	}

	if conditionIDStr != "" {
		if conditionID, err := strconv.ParseInt(conditionIDStr, 10, 64); err == nil {
			params.ConditionID = &conditionID
		}
	}

	if tagIDStr != "" {
		if tagID, err := strconv.ParseInt(tagIDStr, 10, 64); err == nil {
			params.TagID = &tagID
		}
	}

	// 调用服务层方法
	products, total, err := pc.productService.Search(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "搜索商品失败")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"list":  products,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// GetProductsByCategory 获取分类商品
// GET /api/v1/products/category/:categoryId
func (pc *ProductController) GetProductsByCategory(c *gin.Context) {
	// 获取分类ID
	categoryIDStr := c.Param("categoryId")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的分类ID")
		return
	}

	// 解析查询参数
	minPriceStr := c.Query("minPrice")
	maxPriceStr := c.Query("maxPrice")
	conditionIDStr := c.Query("conditionId")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	// 构建查询参数
	params := &product.SearchParams{
		CategoryID: &categoryID,
		Page:       page,
		PageSize:   pageSize,
	}

	// 解析可选参数
	if minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			params.MinPrice = &minPrice
		}
	}

	if maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			params.MaxPrice = &maxPrice
		}
	}

	if conditionIDStr != "" {
		if conditionID, err := strconv.ParseInt(conditionIDStr, 10, 64); err == nil {
			params.ConditionID = &conditionID
		}
	}

	// 调用服务层方法
	products, total, err := pc.productService.ListByCategory(c.Request.Context(), categoryID, params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "获取分类商品失败")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"list":  products,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}