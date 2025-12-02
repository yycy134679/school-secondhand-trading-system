package category

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"school-secondhand-trading-system/model"
	"school-secondhand-trading-system/service/category"
	"school-secondhand-trading-system/util/response"
)

// CategoryController 分类控制器
type CategoryController struct {
	categoryService category.CategoryService
}

// NewCategoryController 创建分类控制器实例
func NewCategoryController(categoryService category.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

// ListCategories 获取所有分类（前台公开接口）
// GET /api/v1/categories
func (cc *CategoryController) ListCategories(c *gin.Context) {
	categories, err := cc.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "获取分类列表失败: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, categories)
}

// CreateCategory 创建分类（管理端接口）
// POST /api/v1/admin/categories
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	// 解析请求体
	type CreateRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	var req CreateRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "请求参数格式错误: "+err.Error())
		return
	}

	// 创建分类模型
	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	// 调用服务层创建分类
	err = cc.categoryService.CreateCategory(c.Request.Context(), category)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "创建分类失败: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, category)
}

// UpdateCategory 更新分类（管理端接口）
// PUT /api/v1/admin/categories/:id
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	// 获取分类ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的分类ID")
		return
	}

	// 解析请求体
	type UpdateRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	var req UpdateRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "请求参数格式错误: "+err.Error())
		return
	}

	// 创建分类模型
	category := &model.Category{
		ID:          categoryID,
		Name:        req.Name,
		Description: req.Description,
	}

	// 调用服务层更新分类
	err = cc.categoryService.UpdateCategory(c.Request.Context(), category)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "更新分类失败: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, category)
}

// DeleteCategory 删除分类（管理端接口）
// DELETE /api/v1/admin/categories/:id
func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	// 获取分类ID
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的分类ID")
		return
	}

	// 调用服务层删除分类
	err = cc.categoryService.DeleteCategory(c.Request.Context(), categoryID)
	if err != nil {
		// 处理特定错误码
		if strings.Contains(err.Error(), "category has products") {
			response.Error(c, http.StatusBadRequest, category.ErrCodeCategoryHasProducts, err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, 500, "删除分类失败: "+err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "分类删除成功"})
}
