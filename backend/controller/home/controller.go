package home

import (
	"github.com/gin-gonic/gin"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/errors"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/recommend"
)

// HomeController 处理首页与推荐相关接口
type HomeController struct {
	recSvc *recommend.RecommendService
}

func NewHomeController(recSvc *recommend.RecommendService) *HomeController {
	return &HomeController{recSvc: recSvc}
}

// GetHomeData GET /api/v1/home
// 可选读取登录用户ID（中间件可在上下文设置 userID, isAuthenticated）
func (hc *HomeController) GetHomeData(c *gin.Context) {
	var userIDPtr *int64
	if val, exists := c.Get("isAuthenticated"); exists {
		if ok, _ := val.(bool); ok {
			if uidVal, ok2 := c.Get("userID"); ok2 {
				if uid, ok3 := uidVal.(int64); ok3 {
					userIDPtr = &uid
				}
			}
		}
	}

	data, err := hc.recSvc.GetHomeData(c.Request.Context(), userIDPtr, 1, 10)
	if err != nil {
		resp.Error(c, errors.CodeInternal, "获取首页数据失败")
		return
	}
	resp.Success(c, data)
}

// GetRecommendations GET /api/v1/recommendations
// 返回推荐的商品ID列表（可用于前端二次查询详情）
func (hc *HomeController) GetRecommendations(c *gin.Context) {
	var userID int64
	if uidVal, ok := c.Get("userID"); !ok {
		resp.Error(c, errors.CodeUnauthenticated, "需要用户登录")
		return
	} else {
		if uid, ok2 := uidVal.(int64); ok2 {
			userID = uid
		} else {
			resp.Error(c, errors.CodeInvalidParams, "用户ID无效")
			return
		}
	}

	ids, err := hc.recSvc.GetRecommendations(c.Request.Context(), userID, 10)
	if err != nil {
		resp.Error(c, errors.CodeInternal, "获取推荐失败")
		return
	}
	resp.Success(c, gin.H{"ids": ids})
}
