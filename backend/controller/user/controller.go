// Package user 提供用户模块的HTTP控制器
// 负责处理用户相关的HTTP请求（注册、登录、个人信息管理等）
package user

import (
	"github.com/gin-gonic/gin"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/errors"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/middleware"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/user"
)

// UserController 处理用户相关的HTTP请求
// 依赖于UserService来处理业务逻辑
// 负责参数绑定、验证、错误处理和响应格式封装

// UserController 用户控制器结构体
// 注入UserService依赖

// RegisterRoutes 注册用户模块的所有路由
//
// 路由列表：
//
//	POST   /users/register         - 用户注册
//	POST   /users/login            - 用户登录
//	GET    /users/profile          - 获取个人信息（需要登录）
//	PUT    /users/profile          - 更新个人信息（需要登录）
//	PUT    /users/password         - 修改密码（需要登录）
//
// 参数：
//   - rg: 父路由组，通常是 /api/v1
//   - userService: 用户服务实例，用于处理业务逻辑
func RegisterRoutes(rg *gin.RouterGroup, userService *user.UserService) {
	// 创建用户路由组，前缀为 /users
	// 最终路径为：/api/v1/users/*
	usr := rg.Group("/users")
	{
		// ============ 公开接口（无需登录）============

		// POST /api/v1/users/register - 用户注册
		usr.POST("/register", func(c *gin.Context) {
			var req struct {
				Account         string  `json:"account" binding:"required"`
				Nickname        string  `json:"nickname" binding:"required"`
				Password        string  `json:"password" binding:"required"`
				ConfirmPassword string  `json:"confirmPassword" binding:"required"`
				WechatID        *string `json:"wechatId,omitempty"`
			}

			// 绑定请求体
			if err := c.ShouldBindJSON(&req); err != nil {
				resp.Error(c, errors.CodeInvalidParams, "请求参数错误: "+err.Error())
				return
			}

			// 验证密码一致性
			if req.Password != req.ConfirmPassword {
				resp.Error(c, errors.CodeInvalidParams, "两次输入的密码不一致")
				return
			}

			// 调用服务层注册用户
			authResp, err := userService.Register(c.Request.Context(), req.Account, req.Nickname, req.Password, req.WechatID)
			if err != nil {
				// 根据错误类型返回对应的错误信息
				resp.Error(c, errors.CodeInvalidParams, err.Error())
				return
			}

			// 返回成功响应
			resp.Success(c, authResp)
		})

		// POST /api/v1/users/login - 用户登录
		usr.POST("/login", func(c *gin.Context) {
			var req struct {
				Account    string `json:"account" binding:"required"`
				Password   string `json:"password" binding:"required"`
				RememberMe bool   `json:"rememberMe"`
			}

			// 绑定请求体
			if err := c.ShouldBindJSON(&req); err != nil {
				resp.Error(c, errors.CodeInvalidParams, "请求参数错误: "+err.Error())
				return
			}

			// 调用服务层登录
			authResp, err := userService.Login(c.Request.Context(), req.Account, req.Password, req.RememberMe)
			if err != nil {
				// 根据错误类型返回对应的错误信息
				resp.Error(c, errors.CodeInvalidParams, err.Error())
				return
			}

			// 返回成功响应
			resp.Success(c, authResp)
		})

		// ============ 需要登录的接口 ============
		// 创建需要认证的路由组
		authorized := usr.Group("")
		// 添加认证中间件，确保用户已登录
		authorized.Use(middleware.AuthMiddleware())
		{
			// GET /api/v1/users/profile - 获取个人信息
			authorized.GET("/profile", func(c *gin.Context) {
				// 从上下文获取userID（由AuthMiddleware注入）
				userIDInterface, exists := c.Get("userID")
				if !exists {
					resp.Error(c, errors.CodeUnauthenticated, "用户未登录")
					return
				}
				userID, ok := userIDInterface.(uint)
				if !ok {
					resp.Error(c, errors.CodeInvalidParams, "用户ID格式错误")
					return
				}

				// 调用服务层获取用户信息
				userResp, err := userService.GetProfile(c.Request.Context(), userID)
				if err != nil {
					resp.Error(c, errors.CodeInvalidParams, err.Error())
					return
				}

				// 返回成功响应
				resp.Success(c, userResp)
			})

			// PUT /api/v1/users/profile - 更新个人信息
			authorized.PUT("/profile", func(c *gin.Context) {
				var req struct {
					Nickname  string  `json:"nickname"`
					AvatarURL string  `json:"avatarUrl"`
					WechatID  *string `json:"wechatId,omitempty"`
				}

				// 绑定请求体
				if err := c.ShouldBindJSON(&req); err != nil {
					resp.Error(c, errors.CodeInvalidParams, "请求参数错误: "+err.Error())
					return
				}

				// 从上下文获取userID（由AuthMiddleware注入）
				userIDInterface, exists := c.Get("userID")
				if !exists {
					resp.Error(c, errors.CodeUnauthenticated, "用户未登录")
					return
				}
				userID, ok := userIDInterface.(uint)
				if !ok {
					resp.Error(c, errors.CodeInvalidParams, "用户ID格式错误")
					return
				}

				// 调用服务层更新用户信息
				userResp, err := userService.UpdateProfile(c.Request.Context(), userID, req.Nickname, req.AvatarURL, req.WechatID)
				if err != nil {
					resp.Error(c, errors.CodeInvalidParams, err.Error())
					return
				}

				// 返回成功响应
				resp.Success(c, userResp)
			})

			// PUT /api/v1/users/password - 修改密码
			authorized.PUT("/password", func(c *gin.Context) {
				var req struct {
					OldPassword     string `json:"oldPassword" binding:"required"`
					NewPassword     string `json:"newPassword" binding:"required"`
					ConfirmPassword string `json:"confirmPassword" binding:"required"`
				}

				// 绑定请求体
				if err := c.ShouldBindJSON(&req); err != nil {
					resp.Error(c, errors.CodeInvalidParams, "请求参数错误: "+err.Error())
					return
				}

				// 验证新密码一致性
				if req.NewPassword != req.ConfirmPassword {
					resp.Error(c, errors.CodeInvalidParams, "两次输入的新密码不一致")
					return
				}

				// 从上下文获取userID（由AuthMiddleware注入）
				userIDInterface, exists := c.Get("userID")
				if !exists {
					resp.Error(c, errors.CodeUnauthenticated, "用户未登录")
					return
				}
				userID, ok := userIDInterface.(uint)
				if !ok {
					resp.Error(c, errors.CodeInvalidParams, "用户ID格式错误")
					return
				}

				// 调用服务层修改密码
				userResp, err := userService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword)
				if err != nil {
					resp.Error(c, errors.CodeInvalidParams, err.Error())
					return
				}

				// 返回成功响应
				resp.Success(c, userResp)
			})
		}
	}
}
