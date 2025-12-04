package admin

import (
	"strconv"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/category"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类管理控制器
type CategoryController struct {
	categoryService category.CategoryService
}

// NewCategoryController 创建分类管理控制器
func NewCategoryController(categoryService category.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

// ListCategories 获取分类列表
// GET /api/v1/admin/categories
func (cc *CategoryController) ListCategories(c *gin.Context) {
	// 调用服务层获取分类列表
	categories, err := cc.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		resp.Error(c, 500, "获取分类列表失败: "+err.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, categories)
}

// CreateCategory 创建分类
// POST /api/v1/admin/categories
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	// 绑定请求体
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if parseErr := c.ShouldBindJSON(&req); parseErr != nil {
		resp.Error(c, 400, "请求参数无效: "+parseErr.Error())
		return
	}

	// 创建分类模型
	category := &model.Category{
		Name: req.Name,
	}

	// 调用服务层创建分类
	if err := cc.categoryService.CreateCategory(c.Request.Context(), category); err != nil {
		resp.Error(c, 500, "创建分类失败: "+err.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, category)
}

// UpdateCategory 更新分类
// PUT /api/v1/admin/categories/:id
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	// 获取分类ID
	idStr := c.Param("id")
	categoryID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || categoryID <= 0 {
		resp.Error(c, 400, "无效的分类ID")
		return
	}

	// 绑定请求体
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if parseErr := c.ShouldBindJSON(&req); parseErr != nil {
		resp.Error(c, 400, "请求参数无效: "+parseErr.Error())
		return
	}

	// 创建分类模型
	category := &model.Category{
		ID:   categoryID,
		Name: req.Name,
	}

	// 调用服务层更新分类
	if err := cc.categoryService.UpdateCategory(c.Request.Context(), category); err != nil {
		if err.Error() == "分类不存在" {
			resp.Error(c, 404, "分类不存在")
			return
		}
		resp.Error(c, 500, "更新分类失败: "+err.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, category)
}

// DeleteCategory 删除分类
// DELETE /api/v1/admin/categories/:id
func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	// 获取分类ID
	idStr := c.Param("id")
	categoryID, parseErr := strconv.ParseInt(idStr, 10, 64)
	if parseErr != nil || categoryID <= 0 {
		resp.Error(c, 400, "无效的分类ID")
		return
	}

	// 调用服务层删除分类
	deleteErr := cc.categoryService.DeleteCategory(c.Request.Context(), categoryID)
	if deleteErr != nil {
		// 根据错误信息返回不同的错误码
		errMsg := deleteErr.Error()
		switch errMsg {
		case "分类不存在":
			resp.Error(c, 404, "分类不存在")
		case "该分类下存在商品，无法删除":
			// 当有商品引用该分类时，返回4001错误码
			resp.Error(c, 4001, errMsg)
		default:
			resp.Error(c, 500, "删除分类失败: "+errMsg)
		}
		return
	}

	// 返回成功响应
	resp.Success(c, gin.H{"message": "分类删除成功"})
}
