## 项目概述

校园二手交易推荐系统 - 学生端前台应用。基于 Vue 3 + TypeScript + Vite + Pinia 构建。

## 常用命令

```bash
# 开发服务器（默认端口 5173）
pnpm dev

# 类型检查 + 构建
pnpm build

# 仅类型检查
pnpm type-check

# ESLint 检查并修复
pnpm lint

# Prettier 格式化
pnpm format
```

## 技术栈

- **框架**: Vue 3.5 (Composition API + `<script setup>`)
- **状态管理**: Pinia 3.0
- **路由**: Vue Router 4.6
- **构建工具**: Vite 7.x
- **语言**: TypeScript 5.9

## 架构说明

### 路径别名

- `@/*` → `./src/*`
- `@common/*` → `../common/*`

### 共享模块 (common/)

位于 `../common/`，前后端共享的类型定义与常量：

- `types/`: API 响应、用户、商品、分类、标签等 TypeScript 接口
- `constants/`: 错误码、商品状态、新旧程度等枚举常量

### API 代理

开发环境下 `/api` 代理到 `http://localhost:8080`，生产环境通过 `VITE_API_BASE_URL` 配置。

## 核心业务规则

### 商品状态机

```
ForSale（在售）↔ Delisted（已下架）
ForSale → Sold（已售，终态不可逆）
```

- `Sold` 状态一旦设置，不可变更为其他状态
- 已售商品：普通卖家不可编辑，管理员可编辑非状态字段

### 一物一件

每条商品记录代表一件实物，无库存字段，线下交易无在线支付。

### 联系卖家

- 未登录：不返回微信号，提示登录
- 登录且非卖家本人：返回完整微信号
- 卖家本人：隐藏联系入口

### 推荐算法

基于用户最近 20 条浏览记录的标签高频匹配，推荐相似商品。

## API 规范

- **前缀**: `/api/v1`
- **认证**: JWT Bearer Token
- **响应结构**: `{ code: number, message: string, data: T }`
- **错误码**: 0 成功，1001 参数错误，1002 未授权，详见 `common/constants/error_code.ts`
- 文档参考：编写 api相关代码时必须阅读参考api.md(docs/api.md)

## UI 设计规范

- **主色**: `#0066FF`
- **配色原则**: 60-30-10（主体色-辅助色-强调色）
- **字体**: Inter, -apple-system, sans-serif
- **间距基准**: 4px 网格
- 文档参考：进行 UI 开发时必须根据前端设计说明.md(docs/前端设计说明.md)进行开发

## 数据库关键表

| 表名               | 说明                                  |
| ------------------ | ------------------------------------- |
| users              | 用户（学生/管理员）                   |
| products           | 商品（status: ForSale/Sold/Delisted） |
| product_images     | 商品图片（is_primary 标记主图）       |
| categories         | 商品分类                              |
| tags               | 标签                                  |
| product_tags       | 商品-标签多对多关联                   |
| user_recent_views  | 最近浏览（自动裁剪到 20 条）          |
| product_conditions | 新旧程度枚举                          |

---

- 修改文件后需要执行pnpm run lint进行代码格式检查和pnpm type-check进行类型检查，确保代码符合规范且无类型错误。
