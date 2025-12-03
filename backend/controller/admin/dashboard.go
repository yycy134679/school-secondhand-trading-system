package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/admin"
)

// DashboardController 仪表盘控制器
type DashboardController struct {
	adminService *admin.AdminService
}

// NewDashboardController 创建仪表盘控制器实例
func NewDashboardController(adminService *admin.AdminService) *DashboardController {
	return &DashboardController{
		adminService: adminService,
	}
}

// GetDashboard 获取仪表盘统计数据
//
// 功能说明：
//   - 提供管理后台仪表盘的统计数据接口
//   - 包含用户总数、商品总数、在售数、已售数等统计信息
//
// API路径：GET /api/v1/admin/dashboard
//
// 访问权限：
//   - 需要管理员权限（通过AdminMiddleware验证）
//
// 返回数据结构：
//   {
//     "code": 0,
//     "data": {
//       "userCount": 100,
//       "productCount": 200,
//       "forSaleCount": 150,
//       "soldCount": 50
//     },
//     "message": "success"
//   }
func (ctrl *DashboardController) GetDashboard(c *gin.Context) {
	// 调用服务层获取统计数据
	stats, err := ctrl.adminService.GetDashboardStats(c.Request.Context())
	if err != nil {
		resp.Error(c, 500, "获取统计数据失败: "+err.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, stats)
}
