package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/tag"
)

// TagController 标签管理控制器
type TagController struct {
	tagService tag.TagService
}

// NewTagController 创建标签管理控制器
func NewTagController(tagService tag.TagService) *TagController {
	return &TagController{
		tagService: tagService,
	}
}

// ListTags 获取标签列表
// GET /api/v1/admin/tags
func (tc *TagController) ListTags(c *gin.Context) {
	// 调用服务层获取标签列表
	tags, err := tc.tagService.ListTags(c.Request.Context())
	if err != nil {
		resp.Error(c, 500, "获取标签列表失败: "+err.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, tags)
}

// CreateTag 创建标签
// POST /api/v1/admin/tags
func (tc *TagController) CreateTag(c *gin.Context) {
	// 绑定请求体
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, 400, "请求参数无效: "+err.Error())
		return
	}

	// 创建标签模型
	tag := &model.Tag{
		Name: req.Name,
	}

	// 调用服务层创建标签
	if err := tc.tagService.CreateTag(c.Request.Context(), tag); err != nil {
		resp.Error(c, 500, "创建标签失败: "+err.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, tag)
}

// UpdateTag 更新标签
// PUT /api/v1/admin/tags/:id
func (tc *TagController) UpdateTag(c *gin.Context) {
	// 获取标签ID
	idStr := c.Param("id")
	tagID, parseErr := strconv.ParseInt(idStr, 10, 64)
	if parseErr != nil || tagID <= 0 {
		resp.Error(c, 400, "无效的标签ID")
		return
	}

	// 绑定请求体
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, 400, "请求参数无效: "+err.Error())
		return
	}

	// 创建标签模型
	tag := &model.Tag{
		ID:   tagID,
		Name: req.Name,
	}

	// 调用服务层更新标签
	if err := tc.tagService.UpdateTag(c.Request.Context(), tag); err != nil {
		if err.Error() == "标签不存在" {
			resp.Error(c, 404, "标签不存在")
			return
		}
		resp.Error(c, 500, "更新标签失败: "+err.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, tag)
}

// DeleteTag 删除标签
// DELETE /api/v1/admin/tags/:id
func (tc *TagController) DeleteTag(c *gin.Context) {
	// 获取标签ID
	idStr := c.Param("id")
	tagID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || tagID <= 0 {
		resp.Error(c, 400, "无效的标签ID")
		return
	}

	// 调用服务层删除标签
	tagDeleteErr := tc.tagService.DeleteTag(c.Request.Context(), tagID)
	if tagDeleteErr != nil {
		// 根据错误信息返回不同的错误码
		errMsg := tagDeleteErr.Error()
		switch errMsg {
		case "标签不存在":
			resp.Error(c, 404, "标签不存在")
		case "该标签下存在商品，无法删除":
			// 当有商品引用该标签时，返回4002错误码
			resp.Error(c, 4002, errMsg)
		default:
			resp.Error(c, 500, "删除标签失败: "+errMsg)
		}
		return
	}

	// 返回成功响应
	resp.Success(c, gin.H{"message": "标签删除成功"})
}
