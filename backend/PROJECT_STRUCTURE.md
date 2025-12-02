# 后端项目结构说明

## 项目概述

这是一个基于 Go + Gin + GORM 的校园二手交易平台后端项目，采用经典的分层架构设计。

## 技术栈

- **Web框架**: Gin (高性能HTTP框架)
- **ORM**: GORM (对象关系映射)
- **数据库**: PostgreSQL
- **缓存**: Redis
- **配置管理**: Viper
- **密码加密**: bcrypt
- **JWT认证**: golang-jwt/jwt
- **参数验证**: go-playground/validator

## 目录结构

```
backend/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主函数：初始化配置、DB、Redis并启动服务器
│
├── config/                 # 配置管理
│   ├── config.go          # 配置加载（从.env和环境变量）
│   ├── database.go        # 数据库连接初始化（GORM + PostgreSQL）
│   └── redis.go           # Redis连接初始化
│
├── router/                 # 路由配置
│   └── router.go          # 路由初始化和注册
│
├── middleware/             # HTTP中间件
│   ├── auth.go            # JWT认证中间件
│   ├── admin.go           # 管理员权限验证中间件
│   ├── logger.go          # 请求日志中间件
│   └── recovery.go        # Panic恢复中间件
│
├── controller/             # 控制器层（处理HTTP请求）
│   ├── user/              # 用户模块控制器
│   │   └── controller.go  # 用户注册、登录、个人信息等接口
│   ├── product/           # 商品模块控制器
│   │   └── controller.go  # 商品发布、搜索、详情等接口
│   ├── category/          # 分类模块控制器（待实现）
│   ├── tag/               # 标签模块控制器（待实现）
│   ├── recommend/         # 推荐模块控制器（待实现）
│   └── admin/             # 后台管理控制器（待实现）
│
├── service/                # 服务层（业务逻辑）
│   ├── user/              # 用户业务逻辑
│   │   └── service.go     # 注册、登录、信息更新、密码修改等
│   ├── product/           # 商品业务逻辑
│   │   └── service.go     # 商品CRUD、状态管理、搜索等
│   ├── category/          # 分类业务逻辑（待实现）
│   ├── tag/               # 标签业务逻辑（待实现）
│   ├── recommend/         # 推荐系统业务逻辑（待实现）
│   └── admin/             # 后台管理业务逻辑（待实现）
│
├── repository/             # 仓储层（数据访问）
│   ├── user_repo.go       # 用户数据访问接口
│   ├── product_repo.go    # 商品数据访问接口
│   ├── category_repo.go   # 分类数据访问接口（待实现）
│   └── tag_repo.go        # 标签数据访问接口（待实现）
│
├── model/                  # 数据模型（实体类）
│   ├── user.go            # 用户模型（对应users表）
│   ├── product.go         # 商品模型（对应products表）
│   ├── category.go        # 分类模型（待实现）
│   ├── tag.go             # 标签模型（待实现）
│   └── ...                # 其他模型
│
├── common/                 # 公共工具
│   ├── resp/              # 统一响应格式
│   │   └── response.go    # Success/Error响应函数
│   ├── errors/            # 错误码定义
│   │   └── codes.go       # 业务错误码常量
│   ├── auth/              # 认证工具
│   │   ├── jwt.go         # JWT生成和解析
│   │   └── password.go    # 密码加密和验证（bcrypt）
│   └── util/              # 其他工具函数
│       ├── validator.go   # 参数验证
│       └── file.go        # 文件上传处理（待实现）
│
├── .env                    # 环境配置文件（包含敏感信息，不提交到Git）
├── go.mod                  # Go模块依赖
├── go.sum                  # 依赖校验和
└── README.md              # 项目说明文档

```

## 架构设计

### 分层架构

项目采用经典的三层架构模式：

```
HTTP请求
    ↓
Controller层（控制器）
    ├── 接收和验证HTTP请求参数
    ├── 调用Service层处理业务逻辑
    └── 构造统一格式的响应
    ↓
Service层（服务/业务逻辑）
    ├── 实现核心业务逻辑
    ├── 调用Repository层访问数据
    ├── 处理事务
    └── 返回业务结果
    ↓
Repository层（仓储/数据访问）
    ├── 封装数据库操作
    ├── 使用GORM执行SQL
    └── 返回数据模型
    ↓
Database（PostgreSQL）
```

### 职责划分

| 层级 | 职责 | 示例 |
|-----|------|------|
| **Controller** | HTTP请求处理、参数验证、响应构造 | 解析JSON请求体，调用service.Register() |
| **Service** | 业务逻辑、事务管理、权限校验 | 验证账号唯一性，哈希密码，生成JWT |
| **Repository** | 数据访问、SQL查询、ORM操作 | 插入用户记录，查询商品列表 |
| **Model** | 数据结构定义 | User、Product等结构体 |

## 核心功能模块

### 1. 用户模块
- 用户注册（账号密码 + 微信号）
- 用户登录（支持"记住我"）
- 个人信息查看和修改
- 头像上传
- 密码修改
- 昵称修改限制（30天一次）

### 2. 商品模块
- 商品发布（标题、描述、价格、图片、分类、标签、新旧程度）
- 商品搜索（关键词、分类、价格区间、标签）
- 商品详情查看
- 商品编辑
- 商品状态管理：
  - ForSale（在售）→ Delisted（下架）→ ForSale（重新上架）
  - ForSale（在售）→ Sold（已售，终态）
- 状态撤销（3秒窗口期，存储在Redis）
- 我的发布列表
- 联系卖家（获取微信号）
- 图片管理（上传、删除、设为主图）

### 3. 推荐系统
- 基于浏览历史的标签推荐
- 首页推荐商品（个性化 + 最新发布）
- 最近浏览记录（最多20条）

### 4. 分类和标签
- 分类管理（增删改查）
- 标签管理（增删改查）
- 删除前检查引用（防止误删）

### 5. 后台管理
- 用户管理（列表、搜索）
- 商品管理（列表、搜索、编辑）
- 分类/标签管理
- 统计数据

## API设计规范

### 统一响应格式

```json
{
  "code": 0,           // 0=成功，非0=错误码
  "message": "ok",     // 响应消息
  "data": {...}        // 响应数据（失败时为null）
}
```

### 错误码规范

- **1xxx**: 通用错误（参数、权限等）
  - 1001: 参数错误
  - 1002: 未认证
  - 1003: 无权限
- **2xxx**: 用户模块错误
- **3xxx**: 商品模块错误
- **4xxx**: 分类/标签模块错误

### 路由规范

- 基础路径：`/api/v1`
- RESTful风格：
  - `GET /products` - 列表
  - `GET /products/:id` - 详情
  - `POST /products` - 创建
  - `PUT /products/:id` - 更新
  - `DELETE /products/:id` - 删除

## 安全机制

### 1. 密码安全
- 使用bcrypt哈希算法（cost=10）
- 自动加盐，每次哈希结果不同
- 存储哈希值，不存储明文密码

### 2. JWT认证
- 无状态认证，不需要服务器存储session
- token包含用户ID和权限信息
- 支持短期token（1小时）和长期token（7天）
- 使用HMAC SHA256签名算法

### 3. 中间件保护
- AuthMiddleware: 验证登录状态
- AdminMiddleware: 验证管理员权限
- 防止未授权访问

### 4. 权限控制
- 用户只能编辑自己的商品
- 管理员可以编辑所有商品（状态不变时）
- 已售商品普通用户不可编辑

## 数据库设计

### 核心表

1. **users** - 用户表
   - 主键：id
   - 唯一索引：account
   - 外键关系：products.seller_id → users.id

2. **products** - 商品表
   - 主键：id
   - 索引：status, created_at, category_id
   - 状态枚举：ForSale, Delisted, Sold
   - 触发器：状态机验证

3. **product_images** - 商品图片表
   - 主键：id
   - 唯一索引：(product_id, is_primary)（确保只有一张主图）

4. **categories** - 分类表
5. **tags** - 标签表
6. **product_tags** - 商品标签关联表（多对多）
7. **user_recent_views** - 浏览记录表
   - 触发器：自动裁剪至最近20条

## 配置说明

### .env 文件示例

```env
APP_ENV=development
HTTP_PORT=8080
DB_DSN="host=localhost user=postgres password=123456 dbname=secondhand_dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"
REDIS_ADDR=localhost:6379
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_TTL=3600
JWT_REMEMBER_TTL=604800
FILE_STORAGE_DIR=./uploads
```

### 配置项说明

- **APP_ENV**: 运行环境（development/production）
- **HTTP_PORT**: HTTP服务器端口
- **DB_DSN**: PostgreSQL连接字符串
- **REDIS_ADDR**: Redis服务器地址
- **JWT_SECRET**: JWT签名密钥（生产环境必须修改）
- **JWT_ACCESS_TTL**: 普通登录token有效期（秒）
- **JWT_REMEMBER_TTL**: 记住我token有效期（秒）
- **FILE_STORAGE_DIR**: 文件上传存储目录

## 开发指南

### 启动服务

```bash
# 1. 进入backend目录
cd backend

# 2. 下载依赖
go mod tidy

# 3. 配置.env文件
# 编辑.env，填入数据库连接信息

# 4. 运行服务
go run ./cmd
```

### 添加新模块

1. 在 `model/` 下创建数据模型
2. 在 `repository/` 下定义数据访问接口
3. 在 `service/` 下实现业务逻辑
4. 在 `controller/` 下实现HTTP处理器
5. 在 `router/router.go` 中注册路由

### 代码规范

- 所有公开函数必须有文档注释
- 使用中文注释说明业务逻辑
- Controller层负责参数验证和响应构造
- Service层负责核心业务逻辑
- Repository层只负责数据访问
- 使用 `resp.Success` 和 `resp.Error` 返回统一响应

## 待完善功能

- [ ] JWT认证完整实现
- [ ] 管理员权限验证完整实现
- [ ] 用户注册/登录接口实现
- [ ] 商品CRUD接口实现
- [ ] 商品状态机和撤销功能
- [ ] 推荐系统实现
- [ ] 文件上传功能
- [ ] 分类和标签管理
- [ ] 后台管理接口
- [ ] 单元测试
- [ ] API文档（Swagger）

## 项目依赖

```
github.com/gin-gonic/gin          - Web框架
github.com/spf13/viper            - 配置管理
gorm.io/gorm                      - ORM框架
gorm.io/driver/postgres           - PostgreSQL驱动
github.com/redis/go-redis/v9      - Redis客户端
github.com/golang-jwt/jwt/v5      - JWT库
github.com/go-playground/validator/v10 - 参数验证
golang.org/x/crypto               - 密码加密
```

## 联系方式

如有问题，请参考 `docs/` 目录下的详细设计文档。
