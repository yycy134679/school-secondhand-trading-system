package category

import (
	"net/http"
	"strconv"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	categoryservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/category"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类控制器
type CategoryController struct {
	categoryService categoryservice.CategoryService
}

// NewCategoryController 创建分类控制器实例
func NewCategoryController(categoryService categoryservice.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

// RegisterRoutes 注册分类相关路由
func (c *CategoryController) RegisterRoutes(router *gin.RouterGroup) {
	// 公开接口，不需要登录
	router.GET("/categories", c.ListCategories)
	router.GET("/categories/:id", c.GetCategoryByID)

	// 管理员接口，需要管理员权限
	adminGroup := router.Group("/admin/categories")
	// 这里会在路由注册时添加管理员中间件
	adminGroup.POST("", c.CreateCategory)
	adminGroup.PUT("/:id", c.UpdateCategory)
	adminGroup.DELETE("/:id", c.DeleteCategory)
}

// ListCategories 获取所有分类列表
// @Summary 获取所有分类列表
// @Description 获取系统中所有分类的列表，供前台使用
// @Tags 分类管理
// @Accept json
// @Produce json
// @Success 200 {array} model.CategoryDTO
// @Router /api/categories [get]
func (c *CategoryController) ListCategories(ctx *gin.Context) {
	categories, err := c.categoryService.ListCategories(ctx)
	if err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "获取分类列表失败")
		return
	}

	resp.Success(ctx, categories)
}

// GetCategoryByID 根据ID获取分类详情
// @Summary 根据ID获取分类详情
// @Description 根据分类ID获取分类的详细信息
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} model.CategoryDTO
// @Router /api/categories/{id} [get]
func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, "无效的分类ID")
		return
	}

	category, err := c.categoryService.GetCategoryByID(ctx, id)
	if err != nil {
		resp.Error(ctx, http.StatusNotFound, "获取分类失败")
		return
	}

	resp.Success(ctx, category)
}

// CreateCategory 创建分类
// @Summary 创建分类
// @Description 创建新的商品分类，需要管理员权限
// @Tags 分类管理-管理员
// @Accept json
// @Produce json
// @Param category body CreateCategoryRequest true "分类信息"
// @Success 200 {object} model.CategoryDTO
// @Router /api/admin/categories [post]
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var req CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Error(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	category, err := c.categoryService.CreateCategory(ctx, req.Name)
	if err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "创建分类失败")
		return
	}

	resp.Success(ctx, category)
}

// UpdateCategory 更新分类
// @Summary 更新分类
// @Description 更新分类信息，需要管理员权限
// @Tags 分类管理-管理员
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param category body UpdateCategoryRequest true "分类信息"
// @Success 200 {object} model.CategoryDTO
// @Router /api/admin/categories/{id} [put]
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, "无效的分类ID")
		return
	}

	var req UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Error(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	category, err := c.categoryService.UpdateCategory(ctx, id, req.Name)
	if err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "更新分类失败")
		return
	}

	resp.Success(ctx, category)
}

// DeleteCategory 删除分类
// @Summary 删除分类
// @Description 删除分类，如果分类下有商品则无法删除，需要管理员权限
// @Tags 分类管理-管理员
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} resp.Response
// @Router /api/admin/categories/{id} [delete]
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, "无效的分类ID")
		return
	}

	err = c.categoryService.DeleteCategory(ctx, id)
	if err != nil {
		// 如果是因为有商品引用而无法删除，返回特定的错误码
		if err.Error() == "该分类下存在商品，无法删除" {
			resp.Error(ctx, http.StatusBadRequest, "删除失败")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "删除分类失败")
		return
	}

	resp.Success(ctx, gin.H{"message": "删除成功"})
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}
