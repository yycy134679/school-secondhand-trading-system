package recommend

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/recommend"
)

// RecommendController 推荐控制器
type RecommendController struct {
	recommendService *recommend.RecommendService
}

// NewRecommendController 创建推荐控制器实例
func NewRecommendController(recommendService *recommend.RecommendService) *RecommendController {
	return &RecommendController{
		recommendService: recommendService,
	}
}

// GetHomeData 获取首页数据
// GET /api/v1/home
func (rc *RecommendController) GetHomeData(c *gin.Context) {
	// 尝试获取登录用户ID（可选）
	var userID *int64
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := strconv.ParseInt(userIDStr.(string), 10, 64); err == nil {
			userID = &uid
		}
	}

	// 获取分页参数
	page := 1
	pageSize := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// 获取首页数据
	homeData, err := rc.recommendService.GetHomeData(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		resp.Error(c, 500, "获取首页数据失败")
		return
	}

	resp.Success(c, homeData)
}

// GetRecentViews 获取用户最近浏览记录
// GET /api/v1/users/recent-views
func (rc *RecommendController) GetRecentViews(c *gin.Context) {
	// 从上下文获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		resp.Error(c, 401, "用户未登录")
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的用户ID")
		return
	}

	// 获取limit参数，默认20
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	// 获取最近浏览记录
	views, err := rc.recommendService.GetRecentViewsWithProducts(c.Request.Context(), userID, limit)
	if err != nil {
		resp.Error(c, 500, "获取浏览记录失败")
		return
	}

	resp.Success(c, gin.H{
		"views": views,
		"total": len(views),
	})
}

// RecordProductView 记录商品浏览
// POST /api/v1/products/:id/view
func (rc *RecommendController) RecordProductView(c *gin.Context) {
	// 从上下文获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		// 未登录用户不记录浏览
		resp.Success(c, gin.H{"recorded": false})
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的用户ID")
		return
	}

	// 获取商品ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的商品ID")
		return
	}

	// 记录浏览
	err = rc.recommendService.RecordView(c.Request.Context(), userID, productID)
	if err != nil {
		// 记录失败不影响主要业务，只记录日志
		resp.Success(c, gin.H{"recorded": false, "error": err.Error()})
		return
	}

	resp.Success(c, gin.H{"recorded": true})
}
