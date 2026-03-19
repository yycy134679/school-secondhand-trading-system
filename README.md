# 校园二手交易平台

[![Vue 3](https://img.shields.io/badge/Vue-3.5-42b883?logo=vue.js&logoColor=white)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178c6?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.23-00add8?logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-1.11-008ecf)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org/)

面向校园场景的二手交易系统，采用前后端分离架构，当前聚焦学生端与 Go 后端服务。仓库同时维护共享类型、错误码和数据库脚本，适合作为课程项目、全栈实训或单仓协同开发示例。

> [!NOTE]
> 本 README 基于仓库代码、接口文档与 SQL 脚本整理生成，当前未在本机实际启动验证。若你的本地环境尚未准备好，可以先按本文了解结构与接入方式，再补齐依赖后联调。

## 项目结构

```text
.
├── frontend/          学生端，Vue 3 + TypeScript + Pinia + Vite
├── backend/           后端服务，Go + Gin + GORM + PostgreSQL
├── common/            前后端共享类型与常量
├── docs/              接口契约文档
└── sql/               数据库建表与初始化脚本
```

## 核心能力

- 学生端支持首页推荐、搜索筛选、发布商品、编辑商品、查看详情、联系卖家、个人中心和我的发布。
- 后端按 `Controller -> Service -> Repository` 分层，统一暴露 `/api/v1` REST 接口。
- `common/` 维护共享类型和常量，前端通过 `@common/*` 别名直接复用，减少前后端漂移。
- PostgreSQL 脚本内置状态机与最近浏览裁剪触发器，关键业务约束不只停留在应用层。

## 技术栈

| 模块 | 技术 |
| --- | --- |
| 学生端 | Vue 3、TypeScript、Vue Router、Pinia、Axios、Vite |
| 后端 | Go 1.23、Gin、GORM、Viper、JWT |
| 数据库 | PostgreSQL |
| 共享层 | `common/types/*`、`common/constants/*` |

## 关键业务规则

- 商品状态机为 `ForSale -> Delisted -> ForSale` 与 `ForSale -> Sold`，其中 `Sold` 是终态，禁止回退。
- 商品是一物一件，没有库存字段；主图由 `product_images.is_primary = true` 决定，并同步冗余到 `products.main_image_url`。
- 联系卖家采用独立接口 `GET /products/:id/contact`，商品详情接口不会直接下发微信号。
- 推荐基于用户最近浏览的标签偏好计算，并与最新在售商品结果去重。
- 新旧程度以 `product_conditions` 表为唯一事实来源，前端不应硬编码。
- 统一响应结构为 `{ code, message, data }`，共享错误码定义位于 `common/constants/error_code.ts`。

## 仓库导览

### `frontend/` 学生端

- 路由覆盖首页、搜索、分类页、商品详情、发布商品、编辑商品、我的发布、登录注册和个人中心。
- 统一请求封装位于 `src/utils/request.ts`，开发环境默认请求 `/api/v1`，通过 Vite 代理到 `http://localhost:8080`。
- `App.vue` 启动时会调用 `app.initDictionaries()` 拉取分类、标签和新旧程度字典。
- 联系卖家逻辑位于商品详情页，只有登录且非商品发布者时才会实际请求联系方式。

### `backend/` 后端

- 应用入口在 `backend/cmd/main.go`，默认端口 `8080`。
- 配置由 `backend/config/config.go` 读取，优先级为环境变量 > `.env` > 默认值。
- 路由集中在 `backend/router/`，包含用户、商品、上传、推荐、分类、标签和新旧程度接口。
- 上传文件静态托管在 `/uploads`，存储目录由 `FILE_STORAGE_DIR` 控制。

### `docs/` 与 `sql/`

- [docs/api.md](./docs/api.md) 是接口契约主文档，联调时应优先以它为准。
- [sql/school-secondhand-trading.sql](./sql/school-secondhand-trading.sql) 包含建表、枚举、触发器与部分初始化数据。

## 快速开始

### 1. 准备依赖

- Node.js：`^20.19.0 || >=22.12.0`
- pnpm：用于前端依赖管理
- Go：`1.23`
- PostgreSQL：建议先准备一个可连接的开发库

### 2. 初始化数据库

导入 SQL 脚本：

```bash
psql -U postgres -d school-secondhand-trading -f sql/school-secondhand-trading.sql
```

如果你使用的是其他数据库名或账号，请按本地环境调整命令。

### 3. 配置后端环境变量

`backend/.env` 中已包含一组示例配置，常用字段如下：

```dotenv
APP_ENV=development
HTTP_PORT=8080
DB_DSN="host=... user=... password=... dbname=... port=5432 sslmode=disable TimeZone=Asia/Shanghai"
JWT_SECRET=your-secret
JWT_ACCESS_TTL=3600
JWT_REMEMBER_TTL=604800
FILE_STORAGE_DIR=./uploads
```

### 4. 启动后端

```bash
cd backend
go mod tidy
go run ./cmd
```

健康检查：

```bash
curl http://localhost:8080/health
```

### 5. 启动学生端

```bash
cd frontend
pnpm install
pnpm dev
```

默认开发环境会读取：

```dotenv
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## 常用命令

### 学生端

```bash
cd frontend
pnpm dev
pnpm type-check
pnpm lint
pnpm build
```

### 后端

```bash
cd backend
go mod tidy
go run ./cmd
go build -o backend ./cmd
```

## 接口概览

统一前缀：`/api/v1`

| 模块 | 示例接口 |
| --- | --- |
| 用户认证 | `POST /users/register`、`POST /users/login`、`GET /users/profile` |
| 商品 | `POST /products`、`GET /products/:id`、`GET /products/search`、`POST /products/:id/status` |
| 联系卖家 | `GET /products/:id/contact` |
| 推荐与浏览 | `GET /home`、`GET /users/recent-views`、`POST /products/:id/view` |
| 字典数据 | `GET /categories`、`GET /tags`、`GET /product-conditions` |

> [!TIP]
> 若你要新增字段、错误码或状态枚举，请同步更新 `common/types` 与 `common/constants`，否则前端与后端很容易出现契约漂移。

## 目录重点说明

- `common/types/`：共享的 `api`、`product`、`user`、`category`、`tag`、`product_condition` 类型定义。
- `common/constants/`：共享错误码、商品状态、新旧程度常量。
- `backend/common/util/file.go`：图片保存与上传 URL 生成逻辑。
- `frontend/src/stores/app.ts`：全局字典初始化入口。
- `frontend/src/api/product.ts`：商品、联系卖家和新旧程度相关 API 封装。

## 开发约定

- 学生端通过别名 `@common/*` 访问共享代码，避免复制类型定义。
- 登录态 token 使用 `school_trading_token` 保存在浏览器本地存储。
- 课程项目下的接口设计以 [docs/api.md](./docs/api.md) 与 SQL 触发器约束为准，出现冲突时优先回到这两处核对。

## 已知前提

- 当前仓库根目录没有统一的 workspace 安装脚本，学生端需单独执行 `pnpm install`。
- README 已尽量以代码事实为准，但由于当前未实际启动，不包含运行截图和在线演示地址。
