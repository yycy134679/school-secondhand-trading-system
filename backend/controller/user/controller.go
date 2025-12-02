// Package user 提供用户模块的HTTP控制器
// 负责处理用户相关的HTTP请求（注册、登录、个人信息管理等）
package user

import (
	"github.com/gin-gonic/gin"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/errors"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
)

// RegisterRoutes 注册用户模块的所有路由
//
// 路由列表（待完善）：
//
//	POST   /users/register         - 用户注册
//	POST   /users/login            - 用户登录
//	GET    /users/profile          - 获取个人信息（需要登录）
//	PUT    /users/profile          - 更新个人信息（需要登录）
//	PUT    /users/password         - 修改密码（需要登录）
//	GET    /users/recent-views     - 最近浏览记录（需要登录）
//
// 参数：
//   - rg: 父路由组，通常是 /api/v1
//
// 设计说明：
//   - 公开接口（register、login）不需要认证中间件
//   - 个人信息相关接口需要AuthMiddleware验证登录状态
//   - 控制器层负责：
//     1. 接收和验证HTTP请求参数
//     2. 调用Service层处理业务逻辑
//     3. 构造统一格式的响应
//     4. 处理和转换错误
//
// TODO: 完整实现步骤
//  1. 创建UserController结构体，注入UserService依赖
//  2. 实现各个Handler方法（Register、Login、GetProfile等）
//  3. 在需要登录的路由上应用AuthMiddleware
//  4. 使用validator验证请求参数
//  5. 使用resp.Success/resp.Error返回统一响应
func RegisterRoutes(rg *gin.RouterGroup) {
	// 创建用户路由组，前缀为 /users
	// 最终路径为：/api/v1/users/*
	usr := rg.Group("/users")
	{
		// ============ 公开接口（无需登录）============

		// POST /api/v1/users/register - 用户注册
		// 请求体示例：
		// {
		//   "account": "zhang123",
		//   "nickname": "张三",
		//   "password": "password123",
		//   "confirmPassword": "password123",
		//   "wechatId": "wx_zhang123"  // 可选
		// }
		// 响应：
		// {
		//   "code": 0,
		//   "message": "ok",
		//   "data": {
		//     "user": {...},
		//     "token": "eyJhbGciOiJIUzI1NiIs..."
		//   }
		// }
		usr.POST("/register", func(c *gin.Context) {
			c.JSON(200, gin.H{"msg": "register stub - 待实现"})
			// TODO: 实现注册逻辑
		})

		// POST /api/v1/users/login - 用户登录
		// 请求体示例：
		// {
		//   "account": "zhang123",
		//   "password": "password123",
		//   "rememberMe": true  // 是否记住登录（token有效期更长）
		// }
		// 响应：
		// {
		//   "code": 0,
		//   "message": "ok",
		//   "data": {
		//     "user": {...},
		//     "token": "eyJhbGciOiJIUzI1NiIs..."
		//   }
		// }
		usr.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{"msg": "login stub - 待实现"})
			// TODO: 实现登录逻辑
		})

		// ============ 需要登录的接口 ============
		// TODO: 添加AuthMiddleware中间件
		// authorized := usr.Group("")
		// authorized.Use(middleware.AuthMiddleware())
		// {
		//     authorized.GET("/profile", handleGetProfile)
		//     authorized.PUT("/profile", handleUpdateProfile)
		//     authorized.PUT("/password", handleChangePassword)
		//     authorized.GET("/recent-views", handleRecentViews)
		// }
	}
}

// RegisterRecentViewRoutes 注册最近浏览接口（需要登录）
// GET /api/v1/users/recent-views
// 从浏览记录与商品表联查，返回最近浏览商品卡片
func RegisterRecentViewRoutes(rg *gin.RouterGroup, vr repository.ViewRecordRepository, pr repository.ProductRepository, auth gin.HandlerFunc) {
	users := rg.Group("/users")
	authorized := users.Group("")
	authorized.Use(auth)
	authorized.GET("/recent-views", func(c *gin.Context) {
		uidVal, ok := c.Get("userID")
		if !ok {
			resp.Error(c, errors.CodeUnauthenticated, "未登录")
			return
		}
		userID, ok := uidVal.(int64)
		if !ok {
			resp.Error(c, errors.CodeInvalidParams, "用户ID无效")
			return
		}

		views, err := vr.ListRecentViews(c.Request.Context(), userID, 20)
		if err != nil {
			resp.Error(c, errors.CodeInternal, "查询浏览记录失败")
			return
		}
		cards := make([]model.ProductCardDTO, 0, len(views))
		for _, v := range views {
			p, _, _, err := pr.GetByID(c.Request.Context(), v.ProductID)
			if err == nil && p != nil {
				cards = append(cards, model.ProductCardDTO{
					ID:           p.ID,
					Title:        p.Title,
					Price:        p.Price,
					MainImageUrl: p.MainImageURL,
					Status:       p.Status,
					CreatedAt:    p.CreatedAt,
				})
			}
		}
		resp.Success(c, gin.H{"items": cards})
	})
}
