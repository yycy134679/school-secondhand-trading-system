[根目录](../../CLAUDE.md) > **src/api**

# API 模块

## 变更记录 (Changelog)

### 2025-11-25 00:02:47
- 初始化模块文档

---

## 模块职责

封装所有后端 API 调用接口，提供类型安全的 HTTP 请求函数。每个文件对应一个业务领域的接口集合。

## 入口与启动

本模块为纯函数库，无独立启动入口。所有函数通过 `@/utils/request.ts` 封装的 axios 实例发起请求。

### 依赖关系

```
api/*.ts
  ↓ 依赖
@/utils/request.ts (axios 实例)
  ↓ 依赖
@common/types/api.ts (ApiResponse, PageResult)
@common/types/user.ts, product.ts, category.ts, tag.ts 等
```

## 对外接口

### user.ts - 用户认证与管理

| 函数名 | 功能 | 请求方式 | 路径 |
|--------|------|----------|------|
| `register()` | 用户注册 | POST | `/users/register` |
| `login()` | 用户登录 | POST | `/users/login` |
| `getProfile()` | 获取当前用户信息 | GET | `/users/profile` |
| `updateProfile()` | 更新用户资料 | PUT | `/users/profile` |
| `changePassword()` | 修改密码 | PUT | `/users/password` |
| `getRecentViews()` | 获取最近浏览 | GET | `/users/recent-views` |

### product.ts - 商品管理

| 函数名 | 功能 | 请求方式 | 路径 |
|--------|------|----------|------|
| `createProduct()` | 发布商品 | POST | `/products` |
| `updateProduct()` | 更新商品 | PUT | `/products/:id` |
| `changeProductStatus()` | 变更商品状态 | POST | `/products/:id/status` |
| `undoProductStatusChange()` | 撤销状态变更 | POST | `/products/:id/status/undo` |
| `getProductDetail()` | 获取商品详情 | GET | `/products/:id` |
| `getMyProducts()` | 获取我的商品列表 | GET | `/products/my` |
| `searchProducts()` | 搜索商品 | GET | `/products/search` |
| `getProductsByCategory()` | 按分类获取商品 | GET | `/products/category/:id` |
| `getProductContact()` | 获取卖家联系方式 | GET | `/products/:id/contact` |
| `getProductConditions()` | 获取新旧程度枚举 | GET | `/product-conditions` |

### home.ts - 首页数据

| 函数名 | 功能 | 请求方式 | 路径 |
|--------|------|----------|------|
| `getHomeData()` | 获取首页推荐和最新商品 | GET | `/home` |

### category.ts - 商品分类

| 函数名 | 功能 | 请求方式 | 路径 |
|--------|------|----------|------|
| `getCategories()` | 获取分类列表 | GET | `/categories` |

### tag.ts - 商品标签

| 函数名 | 功能 | 请求方式 | 路径 |
|--------|------|----------|------|
| `getTags()` | 获取标签列表 | GET | `/tags` |

## 关键依赖与配置

### 请求拦截器（request.ts）

- **认证**: 自动在请求头添加 `Authorization: Bearer <token>`
- **超时**: 10 秒
- **Base URL**: 开发环境 `/api/v1`，生产环境通过 `VITE_API_BASE_URL` 配置

### 响应拦截器（request.ts）

- **成功**: code 为 0 时正常返回
- **未授权** (1002): 清除 token，触发 `auth:unauthorized` 事件
- **禁止访问** (1003): 触发 `auth:forbidden` 事件
- **其他错误**: 抛出异常供调用方处理

## 数据模型

### 类型定义

所有接口的请求参数和响应类型均在对应文件中定义：

- `RegisterParams`, `LoginParams`, `UpdateProfileParams`, `ChangePasswordParams` (user.ts)
- `UpdateProductParams`, `ProductSearchParams`, `ProductStatusParams`, `ContactSellerResponse` (product.ts)
- `HomeData` (home.ts)

### 共享类型

从 `@common/types/` 引入：
- `ApiResponse<T>`: 统一响应结构
- `PageResult<T>`: 分页数据结构
- `User`, `Product`, `Category`, `Tag`, `ProductCondition`: 实体类型

## 测试与质量

### 当前状态

- 无单元测试（建议使用 Vitest + MSW 模拟 HTTP 请求）
- 类型检查覆盖所有接口

### 测试建议

1. **单元测试**：验证参数序列化、响应解析
2. **集成测试**：使用 MSW 模拟后端响应
3. **类型测试**：确保响应类型与 `@common/types/` 一致

## 常见问题 (FAQ)

### Q1: 如何添加新的 API 接口？

1. 在对应文件（如 `user.ts`）中定义请求参数类型（如需要）
2. 使用 `request.get/post/put/delete<ApiResponse<T>>()` 发起请求
3. 确保响应类型与后端一致（参考 `docs/api.md`）
4. 在 Pinia store 或组件中调用

### Q2: 为什么要封装 API 调用而不是直接用 axios？

- **类型安全**: 所有接口有明确的 TS 类型定义
- **统一处理**: 认证、错误处理在拦截器中统一完成
- **易于维护**: API 变更只需修改对应文件
- **便于测试**: 可以轻松 mock 整个 API 模块

### Q3: 如何处理文件上传？

使用 `FormData` 并设置 `Content-Type: multipart/form-data`：

```typescript
const formData = new FormData()
formData.append('title', '商品标题')
formData.append('images', file)

await createProduct(formData)
```

## 相关文件清单

```
src/api/
├── user.ts           # 用户认证与管理
├── product.ts        # 商品管理
├── home.ts           # 首页数据
├── category.ts       # 商品分类
└── tag.ts            # 商品标签
```

---

**最后更新**: 2025-11-25 00:02:47
