# AGENTS.md

## 项目总览
- 单仓三子项目：`frontend` 学生端、`admin-frontend` 管理端（皆 Vue3+TS+Vite+Pinia）、`backend` Go Gin 服务；共享类型与常量位于 `common/`，通过别名 `@common/*` 引入。
- API 统一前缀 `/api/v1`，响应结构 `{ code, message, data }`，核心错误码：`1002` 未登录、`1003` 权限不足、`3004` 已售终态不可变更、`3005` 撤销超时。
- 数据库 schema 在 `sql/schema.sql`（含商品状态触发器、最近浏览裁剪触发器）；新旧程度以 `product_conditions` 表为唯一事实来源。

## 关键业务规则
- 商品状态机：`ForSale` ↔ `Delisted`，`ForSale → Sold` 终态不可逆；上/下架成功缓存 3 秒撤销窗口（普通卖家不可撤销已售）。
- 一物一件：无库存字段；主图来自 `product_images.is_primary=true`，同步到 `products.main_image_url`。
- 联系卖家：`GET /products/:id/contact` 按登录态与是否本人决定是否返回微信号；详情接口不直接下发微信号。
- 推荐：基于 `user_recent_views` 最近 20 条的标签高频匹配，并与最新在售结果去重。

## 后端开发要点（`backend/`）
- 架构：Controller → Service → Repository，路由集中在 `router/router.go` 并按模块拆分；中间件 `auth` 注入 `userID/isAdmin`，`admin` 校验管理员权限，`cors/logger/recovery` 常驻。
- 配置：`.env` 读入 `DB_DSN/JWT_SECRET/JWT_ACCESS_TTL/JWT_REMEMBER_TTL/FILE_STORAGE_DIR` 等；`config/database.go` 初始化 GORM Postgres。默认端口 8080。
- 启动/构建：`cd backend; go mod tidy; go run ./cmd` 或 `go build -o backend.exe ./cmd`。健康检查 `GET /health`。
- 业务要点：商品服务负责状态流转与撤销缓存；管理员接口仅允许更正已售商品的非状态字段；图片/文件处理在 `common/util/file.go`。

## 学生端前端约定（`frontend/`）
- 技术栈：Vue3 `<script setup>` + TS + Pinia + Vue Router + Axios；路径别名 `@/*`、`@common/*`（见 `tsconfig.app.json`、`vite.config.ts`）。
- 请求：统一使用 `src/utils/request.ts` 封装（10s 超时，自动附带 Bearer Token，`code!=0` 统一处理 1002/1003）；API 函数按领域放在 `src/api/*.ts`，勿直接 new axios。
- 状态：Pinia stores in `src/stores/`（user/product/app）；`app.initDictionaries()` 启动时拉取分类/标签/新旧程度，不要硬编码 `product_conditions`。
- 路由：`src/router/index.ts` 使用 `meta.requiresAuth` 守卫；未登录访问受限页触发登录弹窗/跳转。联系卖家按钮调用独立 contact 接口。
- 常用命令：`pnpm install`；`pnpm dev`（5173）；`pnpm type-check`、`pnpm lint`、`pnpm build`。

## 联调与数据约束
- 接口契约以 `docs/api.md` 为准，尤其联系卖家、状态机、撤销窗口；发现设计冲突以该文档和 schema 触发器为权威。
- 需登录操作请确保前端已写 token（localStorage key `school_trading_token`），后台接口需管理员 token。
- 新增类型或错误码时同步维护 `common/types` 与 `common/constants`，避免前后端漂移。