// Package router 负责HTTP路由的初始化和配置
// 统一管理所有API端点和中间件
package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/cache"
	"github.com/yycy134679/school-secondhand-trading-system/backend/config"
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/admin"
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/category"
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/product"
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/recommend"
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/tag"
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/user"
	"github.com/yycy134679/school-secondhand-trading-system/backend/middleware"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
	adminservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/admin"
	categoryservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/category"
	productservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/product"
	recommendservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/recommend"
	tagservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/tag"
	userservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/user"
)

// SetupRouter 初始化并配置HTTP路由引擎
//
// 功能说明：
//   - 创建Gin引擎实例并注册默认中间件（Logger和Recovery）
//   - 注册健康检查端点
//   - 注册API v1版本的所有业务路由
//   - 将数据库和内存缓存注入到各个模块（通过依赖注入）
//
// 参数：
//   - db: GORM数据库连接实例，用于数据持久化操作
//   - memCache: 内存缓存服务实例，用于缓存和状态管理
//   - cfg: 应用配置对象，包含JWT密钥、文件存储路径等
//
// 返回值：
//   - *gin.Engine: 配置好的Gin引擎实例，可直接用于启动HTTP服务器
//
// 路由结构：
//
//	/health              - 健康检查端点（用于负载均衡器和监控）
//	/api/v1/users/*      - 用户相关接口（注册、登录、个人信息等）
//	/api/v1/products/*   - 商品相关接口（发布、搜索、详情等）
//	/api/v1/categories/* - 分类管理接口（待实现）
//	/api/v1/tags/*       - 标签管理接口（待实现）
//	/api/v1/admin/*      - 后台管理接口（待实现）
func SetupRouter(db *gorm.DB, memCache *cache.MemoryCache, cfg *config.Config) *gin.Engine {
	// 创建Gin引擎实例
	// gin.Default() 会自动附加两个中间件：
	// 1. Logger() - 记录每个HTTP请求的日志（方法、路径、状态码、耗时等）
	// 2. Recovery() - 捕获panic并返回500错误，防止服务器崩溃
	r := gin.Default()

	// 注册健康检查端点
	// 用途：
	// - 负载均衡器（如Nginx、K8s）用于检测服务是否存活
	// - 监控系统用于健康状态检查
	// - 开发调试用于快速验证服务是否启动
	// 返回格式：{"status":"ok"}
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 创建API v1路由组
	// 路由组的优势：
	// - 统一的URL前缀（/api/v1）
	// - 可以在组级别应用中间件（如鉴权、CORS等）
	// - 便于API版本管理（将来可以添加/api/v2）
	api := r.Group("/api/v1")
	{
		// 初始化用户相关组件
		// 创建用户仓库实例
		userRepo := repository.NewUserRepository(db)
		// 创建用户服务实例
		userService := userservice.NewUserService(userRepo)

		// 注册用户模块路由
		// 包含的接口：
		// POST /api/v1/users/register  - 用户注册
		// POST /api/v1/users/login     - 用户登录
		// GET  /api/v1/users/profile   - 获取个人信息
		// PUT  /api/v1/users/profile   - 更新个人信息
		// PUT  /api/v1/users/password  - 修改密码
		user.RegisterRoutes(api, userService)

		// 注册商品模块路由
		// 包含的接口（示例）：
		// POST /api/v1/products         - 发布商品
		// GET  /api/v1/products/:id     - 获取商品详情
		// PUT  /api/v1/products/:id     - 编辑商品
		// GET  /api/v1/products/search  - 搜索商品
		// GET  /api/v1/products/my      - 我的发布
		// 创建商品相关组件
		productService := productservice.NewProductService()
		productController := product.NewProductController(productService)
		imageController := product.NewImageController(productService)
		SetupProductRoutes(r, productController, imageController)

		// 初始化推荐服务和浏览记录相关组件
		viewRecordRepo := repository.NewViewRecordRepository(db)
		productRepo := repository.NewProductRepository(db)
		recommendService := recommendservice.NewRecommendService(viewRecordRepo, productRepo, db, nil) // Redis设为nil,可选
		recommendController := recommend.NewRecommendController(recommendService)
		SetupRecommendRoutes(r, recommendController)

		// 初始化分类和标签相关组件
		// 创建仓库层实例
		categoryRepo := repository.NewCategoryRepository(db)
		tagRepo := repository.NewTagRepository(db)

		// 创建服务层实例
		categoryService := categoryservice.NewCategoryService(categoryRepo)
		tagService := tagservice.NewTagService(tagRepo)

		// 创建控制器实例
		categoryController := category.NewCategoryController(categoryService)
		tagController := tag.NewTagController(tagService)

		// 创建管理员中间件
		adminMiddleware := middleware.AdminMiddleware()

		// 注册分类模块路由
		SetupCategoryRoutes(r, categoryController)

		// 注册标签模块路由
		SetupTagRoutes(r, tagController)

		// 初始化管理后台相关组件
		// 创建服务层实例
		adminService := adminservice.NewAdminService(db)

		// 创建其他管理后台控制器实例
		dashboardController := admin.NewDashboardController(adminService)
		userController := admin.NewUserController(adminService)
		adminProductController := admin.NewProductController(adminService)

		// 注册管理后台路由（不包括分类和标签，因为已经在上面注册了）
		RegisterAdminRoutes(api, dashboardController, userController, adminProductController, adminMiddleware)
	}

	// 返回配置好的Gin引擎实例
	// 调用方可以直接使用 engine.Run(":8080") 启动服务器
	return r
}
