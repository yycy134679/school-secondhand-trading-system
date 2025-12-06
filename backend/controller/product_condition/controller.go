package productcondition

import (
	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/product_condition"
)

// Controller 新旧程度控制器
type Controller struct {
	service productcondition.Service
}

// NewController 创建控制器实例
func NewController(service productcondition.Service) *Controller {
	return &Controller{service: service}
}

// ListProductConditions 获取所有新旧程度
// GET /api/v1/product-conditions
func (pc *Controller) ListProductConditions(c *gin.Context) {
	conditions, err := pc.service.ListProductConditions(c.Request.Context())
	if err != nil {
		resp.Error(c, 500, "获取新旧程度失败: "+err.Error())
		return
	}
	resp.Success(c, conditions)
}
