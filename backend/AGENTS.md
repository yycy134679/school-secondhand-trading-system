# AGENTS.md

这是一个用于学校二手交易系统的 Go 后端服务。后端使用以下技术：

- **Gin** 作为 HTTP Web 框架
## 项目概述
- **GORM** 进行 PostgreSQL 数据库操作
- **JWT** (golang-jwt/jwt/v5) 用于认证
- **Viper** 进行配置管理
- **内存缓存** 用于推荐和状态管理（无 Redis 依赖）

## 常用命令

### 开发
```bash
# 运行应用程序
go run ./cmd

# 构建应用程序
go build -o backend ./cmd

# 安装/更新依赖
go mod tidy

# 格式化代码
gofmt -w .
```

注意：当前代码库中不存在测试文件。添加测试时：
```bash
# 运行所有测试
go test ./...
# 运行特定包的测试

go test ./service/user
# 运行单个测试

go test ./service/user -run TestUserService_Register
```
### 测试

## 架构

### 分层结构

代码库遵循 **3 层架构** 模式：

1. **控制器层** (`controller/`) - HTTP 处理器、请求验证、响应格式化
2. **服务层** (`service/`) - 业务逻辑实现
3. **仓库层** (`repository/`) - 数据库访问和查询

数据流向：`HTTP 请求 → 控制器 → 服务 → 仓库 → 数据库`

### 关键目录

- `cmd/` - 应用程序入口点（main.go）
- `config/` - 配置加载和数据库初始化
- `router/` - HTTP 路由注册和中间件设置
- `model/` - 数据模型（GORM 实体）
- `middleware/` - HTTP 中间件（认证、跨域资源共享、管理员权限）
- `common/` - 共享工具：
  - `auth/` - JWT 令牌生成与解析
  - `cache/` - 内存缓存实现
  - `errors/` - 标准错误码
  - `resp/` - HTTP 响应辅助函数
  - `util/` - 工具函数（文件处理、验证）

### 模块组织

代码库按功能模块组织：
- `user/` - 用户注册、登录、个人资料管理
- `product/` - 产品 CRUD、搜索、状态管理
- `category/` - 分类管理（仅管理员）
- `tag/` - 标签管理（仅管理员）
- `recommend/` - 推荐系统及浏览记录
- `admin/` - 管理员仪表盘、用户/产品管理
- `upload/` - 文件上传处理

每个模块通常包含：
- 控制器（`controller/<模块>/`）
- 服务（`service/<模块>/`）
- 代码仓库（如需，放在 `repository/` 目录下）

## 配置

配置按以下顺序从 `.env` 文件（或环境变量）加载：
1. 环境变量（优先级最高）
2. `.env` 文件
3. 代码中的默认值（优先级最低）

### 必需的 .env 变量：
```bash
APP_ENV=development
HTTP_PORT=8080
DB_DSN="host=... user=... password=... dbname=... port=5432 sslmode=disable"
JWT_SECRET=your-secret-key
FILE_STORAGE_DIR=./uploads
```

## 认证与授权

### JWT 认证
- 令牌在 `common/auth/jwt.go` 中生成
- 默认过期时间：1 小时
- 令牌格式：`Authorization: Bearer <token>` 或 `Authorization: <token>`

### 中间件
- `AuthMiddleware()` - 需要有效的 JWT，将 `user_id` 和 `role` 注入上下文
- `OptionalAuthMiddleware()` - 如果存在 JWT 则解析，允许匿名访问
- `AdminMiddleware()` - 需要 JWT 和 `is_admin=true`

### 在控制器中访问用户信息
```go
userIDStr := c.GetString("user_id")
userID, _ := strconv.ParseInt(userIDStr, 10, 64)
```

## 数据库

### 连接
数据库连接在 `config/db.go` 中初始化。应用程序可以在没有数据库的情况下启动（此时 DB 将为 nil）。

### 模型
所有 GORM 模型都在 `model/` 目录下：
- `User` - 用户账户（映射到 `users` 表）
- `Product` - 产品列表（映射到 `products` 表）
- `ProductImage` - 产品图片
- `Category`, `Tag`, `ViewRecord` 等

### 仓储模式
`repository/` 目录下的仓储处理所有数据库查询。始终通过构造函数注入 `*gorm.DB`：
```go
userRepo := repository.NewUserRepository(数据库)
```

## 重要实现细节

### 文件上传
- 上传目录：通过 `.env` 中的 `FILE_STORAGE_DIR` 配置
- 静态文件服务：`/uploads` 路由映射到 `FILE_STORAGE_DIR`
- 文件命名：`<timestamp>.<ext>`（例如：`1764933128914.png`）
- 上传控制器：`controller/upload/upload_controller.go`

### 内存缓存
系统使用内存缓存（非 Redis）用于：
- 推荐内容缓存
- 产品状态撤销（状态变更的 3 秒恢复窗口）

缓存在 `cmd/main.go` 中初始化，并通过 `router.SetupRouter()` 注入。

### 错误处理
标准错误码在 `common/errors/codes.go` 中定义：
- 使用 `resp.Error(c, errors.CodeXXX, "message")` 返回错误响应
- 使用 `resp.OK(c, data)` 返回成功响应

示例：
```go
resp.Error(c, errors.CodeUnauthenticated, "请先登录")
resp.OK(c, gin.H{"user": user})
```

### 跨域资源共享（CORS）
CORS 中间件在 `router.SetupRouter()` 中全局注册。当前允许所有源（`*`）——生产环境中应进行限制。

## 代码规范

### 添加新的 API 端点

1. **在路由文件中定义路由**（`router/<module>_routes.go`）：
```go
api.POST("/products", middleware.AuthMiddleware(), productController.CreateProduct)
```

2. **实现控制器方法**（`controller/<module>/`）：
```go
func (pc *ProductController) CreateProduct(c *gin.Context) {
    var req CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        resp.Error(c, errors.CodeBadRequest, "invalid request")
        return
    }

    userID := c.GetString("user_id")
    product, err := pc.service.CreateProduct(userID, &req)
if err != nil {
        resp.Error(c, errors.CodeInternalError, err.Error())
        return
    }

    resp.OK(c, product)
}
```

3. **实现服务逻辑** (`service/<module>/`):
```go
func (s *ProductService) CreateProduct(userID string, req *CreateProductRequest) (*model.Product, error) {
    // 这里编写业务逻辑
    return product, nil
}
```

4. **如果需要添加仓库方法** (`repository/`):
```go
func (r *ProductRepository) Create(product *model.Product) error {
    return r.db.Create(product).Error
}
```

### 依赖注入模式
所有依赖项均通过构造函数注入：
```go
// 在 router/router.go 中
userRepo := repository.NewUserRepository(db)
userService := userservice.NewUserService(userRepo)
userController := user.NewUserController(userService)
```

## 注意事项

- 项目使用 Go 1.23
- 模块路径：`github.com/yycy134679/school-secondhand-trading-system/backend`
- 当前无测试覆盖率 - 应添加关键流程的测试
- 管理员功能要求用户记录中包含 `is_admin=true`
- 静态文件（上传文件）通过 `/uploads/*` 路由提供服务
