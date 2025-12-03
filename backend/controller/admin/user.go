package admin

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	adminservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/admin"
)

// UserController 管理后台用户控制器
type UserController struct {
	adminService *adminservice.AdminService
}

// NewUserController 创建用户控制器实例
func NewUserController(adminService *adminservice.AdminService) *UserController {
	return &UserController{
		adminService: adminService,
	}
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表，支持根据账号或昵称模糊搜索
// @Tags 管理后台-用户管理
// @Accept json
// @Produce json
// @Param keyword query string false "搜索关键词（账号或昵称）"
// @Param page query int false "页码，默认1" default(1)
// @Param pageSize query int false "每页数量，默认10" default(10)
// @Success 200 {object} resp.Response{data=adminservice.UserListResponse}
// @Router /api/v1/admin/users [get]
func (uc *UserController) ListUsers(ctx *gin.Context) {
	// 解析查询参数
	keyword := ctx.Query("keyword")
	
	// 解析页码，默认为1
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		resp.Error(ctx, 1001, "无效的页码参数")
		return
	}
	
	// 解析每页数量，默认为10，最大限制为100
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		resp.Error(ctx, 1001, "无效的每页数量参数，范围1-100")
		return
	}
	
	// 调用服务层方法获取用户列表
	result, err := uc.adminService.ListUsers(ctx, keyword, page, pageSize)
	if err != nil {
		resp.Error(ctx, 500, "获取用户列表失败: "+err.Error())
		return
	}
	
	// 返回成功响应
	resp.Success(ctx, result)
}
