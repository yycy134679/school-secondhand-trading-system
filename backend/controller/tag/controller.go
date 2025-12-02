package tag

import (
	"net/http"
	"strconv"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	tagservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/tag"

	"github.com/gin-gonic/gin"
)

// TagController 标签控制器
type TagController struct {
	tagService tagservice.TagService
}

// NewTagController 创建标签控制器实例
func NewTagController(tagService tagservice.TagService) *TagController {
	return &TagController{
		tagService: tagService,
	}
}

// RegisterRoutes 注册标签相关路由
func (c *TagController) RegisterRoutes(router *gin.RouterGroup) {
	// 公开接口，不需要登录
	router.GET("/tags", c.ListTags)
	router.GET("/tags/:id", c.GetTagByID)

	// 管理员接口，需要管理员权限
	adminGroup := router.Group("/admin/tags")
	// 这里会在路由注册时添加管理员中间件
	adminGroup.POST("", c.CreateTag)
	adminGroup.PUT("/:id", c.UpdateTag)
	adminGroup.DELETE("/:id", c.DeleteTag)
}

// ListTags 获取所有标签列表
// @Summary 获取所有标签列表
// @Description 获取系统中所有标签的列表，供前台使用
// @Tags 标签管理
// @Accept json
// @Produce json
// @Success 200 {array} model.TagDTO
// @Router /api/tags [get]
func (c *TagController) ListTags(ctx *gin.Context) {
	tags, err := c.tagService.ListTags(ctx)
	if err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "获取标签列表失败")
		return
	}

	resp.Success(ctx, tags)
}

// GetTagByID 根据ID获取标签详情
// @Summary 根据ID获取标签详情
// @Description 根据标签ID获取标签的详细信息
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} model.TagDTO
// @Router /api/tags/{id} [get]
func (c *TagController) GetTagByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, "无效的标签ID")
		return
	}

	tag, err := c.tagService.GetTagByID(ctx, id)
	if err != nil {
		resp.Error(ctx, http.StatusNotFound, "获取标签失败")
		return
	}

	resp.Success(ctx, tag)
}

// CreateTag 创建标签
// @Summary 创建标签
// @Description 创建新的商品标签，需要管理员权限
// @Tags 标签管理-管理员
// @Accept json
// @Produce json
// @Param tag body CreateTagRequest true "标签信息"
// @Success 200 {object} model.TagDTO
// @Router /api/admin/tags [post]
func (c *TagController) CreateTag(ctx *gin.Context) {
	var req CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Error(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	tag, err := c.tagService.CreateTag(ctx, req.Name)
	if err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "创建标签失败")
		return
	}

	resp.Success(ctx, tag)
}

// UpdateTag 更新标签
// @Summary 更新标签
// @Description 更新标签信息，需要管理员权限
// @Tags 标签管理-管理员
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param tag body UpdateTagRequest true "标签信息"
// @Success 200 {object} model.TagDTO
// @Router /api/admin/tags/{id} [put]
func (c *TagController) UpdateTag(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, "无效的标签ID")
		return
	}

	var req UpdateTagRequest
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		resp.Error(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	tag, err := c.tagService.UpdateTag(ctx, id, req.Name)
	if err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "更新标签失败")
		return
	}

	resp.Success(ctx, tag)
}

// DeleteTag 删除标签
// @Summary 删除标签
// @Description 删除标签，如果标签下有商品则无法删除，需要管理员权限
// @Tags 标签管理-管理员
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} resp.Response
// @Router /api/admin/tags/{id} [delete]
func (c *TagController) DeleteTag(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, "无效的标签ID")
		return
	}

	err = c.tagService.DeleteTag(ctx, id)
	if err != nil {
		// 如果是因为有商品引用而无法删除，返回特定的错误码
		if err.Error() == "该标签下存在商品，无法删除" {
			resp.Error(ctx, http.StatusBadRequest, "删除失败")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "删除标签失败")
		return
	}

	resp.Success(ctx, gin.H{"message": "删除成功"})
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}
