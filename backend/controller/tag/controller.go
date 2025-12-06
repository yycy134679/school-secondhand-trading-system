package tag

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/tag"
)

// TagController 标签控制器
type TagController struct {
	tagService tag.TagService
}

// NewTagController 创建标签控制器实例
func NewTagController(tagService tag.TagService) *TagController {
	return &TagController{
		tagService: tagService,
	}
}

// ListTags 获取所有标签列表（前台公开接口）
// @Summary 获取所有标签
// @Description 获取系统中所有可用的标签列表，无需登录
// @Tags 前台-标签
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]model.Tag}
// @Router /api/v1/tags [get]
func (tc *TagController) ListTags(c *gin.Context) {
	// 调用服务层获取标签列表
	tags, err := tc.tagService.ListTags(c)
	if err != nil {
		resp.Error(c, 500, "获取标签列表失败: "+err.Error())
		return
	}

	resp.Success(c, tags)
}

// CreateTag 创建新标签（管理端接口）
// @Summary 创建标签
// @Description 管理员创建新标签
// @Tags 管理端-标签
// @Accept json
// @Produce json
// @Param tag body TagCreateRequest true "标签创建请求"
// @Success 200 {object} response.Response{data=model.Tag}
// @Router /api/v1/admin/tags [post]
func (tc *TagController) CreateTag(c *gin.Context) {
	var req TagCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, 400, "请求参数无效: "+err.Error())
		return
	}

	// 创建标签模型
	tag := &model.Tag{
		Name:       req.Name,
		CategoryID: req.CategoryID,
	}

	// 调用服务层创建标签
	err := tc.tagService.CreateTag(c.Request.Context(), tag)
	if err != nil {
		resp.Error(c, 500, "创建标签失败: "+err.Error())
		return
	}

	resp.Success(c, tag)
}

// UpdateTag 更新标签（管理端接口）
// @Summary 更新标签
// @Description 管理员更新标签信息
// @Tags 管理端-标签
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param tag body TagUpdateRequest true "标签更新请求"
// @Success 200 {object} response.Response{data=model.Tag}
// @Router /api/v1/admin/tags/{id} [put]
func (tc *TagController) UpdateTag(c *gin.Context) {
	// 获取标签ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的标签ID")
		return
	}

	// 绑定请求参数
	var req TagUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, 400, "请求参数无效: "+err.Error())
		return
	}

	// 创建标签模型
	tag := &model.Tag{
		ID:         id,
		Name:       req.Name,
		CategoryID: req.CategoryID,
	}

	// 调用服务层更新标签
	err = tc.tagService.UpdateTag(c.Request.Context(), tag)
	if err != nil {
		resp.Error(c, 500, "更新标签失败: "+err.Error())
		return
	}

	resp.Success(c, tag)
}

// DeleteTag 删除标签（管理端接口）
// @Summary 删除标签
// @Description 管理员删除标签，若标签下有关联商品则返回错误码4002
// @Tags 管理端-标签
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/tags/{id} [delete]
func (tc *TagController) DeleteTag(c *gin.Context) {
	// 获取标签ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的标签ID")
		return
	}

	// 调用服务层删除标签
	err = tc.tagService.DeleteTag(c, id)
	if err != nil {
		// 判断是否为标签下有商品的错误（错误码4002）
		if err == tag.ErrTagHasProducts {
			resp.Error(c, tag.ErrCodeTagHasProducts, err.Error())
			return
		}
		resp.Error(c, 500, "删除标签失败: "+err.Error())
		return
	}

	resp.Success(c, nil)
}

// TagCreateRequest 标签创建请求
type TagCreateRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=50"`
	CategoryID int64  `json:"categoryId" binding:"required,gt=0"`
}

// TagUpdateRequest 标签更新请求
type TagUpdateRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=50"`
	CategoryID int64  `json:"categoryId" binding:"required,gt=0"`
}
