[根目录](../../CLAUDE.md) > **src/utils**

# Utils 模块（工具函数）

## 变更记录 (Changelog)

### 2025-11-25 00:02:47
- 初始化模块文档

---

## 模块职责

提供通用的工具函数和封装，包括 HTTP 请求、认证、格式化等。所有工具函数应保持纯函数特性（无副作用，除非明确标注）。

## 入口与启动

无独立入口，各文件按需导入使用。

## 对外接口

### request.ts - HTTP 请求封装

基于 Axios 封装的 HTTP 客户端，提供统一的请求/响应拦截。

**核心配置:**
- Base URL: `import.meta.env.VITE_API_BASE_URL || '/api/v1'`
- Timeout: 10 秒
- Content-Type: `application/json` (默认)

**请求拦截器:**
- 自动添加 `Authorization: Bearer <token>` 头（如果已登录）

**响应拦截器:**
- 统一处理业务错误码（0: 成功，1002: 未授权，1003: 禁止访问）
- 未授权时触发 `auth:unauthorized` 事件并清除 token
- 禁止访问时触发 `auth:forbidden` 事件

**导出:**
- `service`: 配置好的 axios 实例（默认导出）

**使用示例:**

```typescript
import request from '@/utils/request'

const response = await request.get<ApiResponse<User>>('/users/profile')
const user = response.data.data
```

### auth.ts - 认证工具

管理 JWT Token 的存储和读取。

**函数列表:**

| 函数名 | 功能 | 参数 | 返回值 |
|--------|------|------|--------|
| `getToken()` | 获取 token | - | `string \| null` |
| `setToken(token)` | 保存 token | `token: string` | `void` |
| `removeToken()` | 删除 token | - | `void` |

**存储位置:**
- localStorage, key: `school_trading_token`

**注意事项:**
- Token 不包含 `Bearer ` 前缀（由 request.ts 拦截器添加）
- 生产环境建议考虑安全性（如使用 HttpOnly Cookie 或 Secure Storage）

## 关键依赖与配置

### 依赖项

- `axios`: HTTP 客户端
- `@common/types/api`: ApiResponse 类型
- `@common/constants/error_code`: 错误码常量

### 环境变量

- `VITE_API_BASE_URL`: API 服务地址（生产环境）
  - 开发环境默认使用 `/api/v1`（通过 Vite 代理到 `http://localhost:8080`）

## 数据模型

### request.ts 类型定义

```typescript
import type { AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@common/types/api'

// 响应拦截器返回类型
type InterceptedResponse = AxiosResponse<ApiResponse>
```

## 测试与质量

### 测试建议

1. **request.ts 测试**:
   - 使用 `axios-mock-adapter` 或 MSW 模拟 HTTP 响应
   - 验证拦截器逻辑（认证头、错误处理）

2. **auth.ts 测试**:
   - 使用 `localStorage` mock 验证存取逻辑
   - 测试 token 过期场景

示例（使用 Vitest）:

```typescript
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { getToken, setToken, removeToken } from './auth'

describe('auth utils', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('should set and get token', () => {
    setToken('test-token')
    expect(getToken()).toBe('test-token')
  })

  it('should remove token', () => {
    setToken('test-token')
    removeToken()
    expect(getToken()).toBeNull()
  })
})
```

## 常见问题 (FAQ)

### Q1: 如何处理请求超时？

当前 timeout 设置为 10 秒，可在 `request.ts` 中修改：

```typescript
const service = axios.create({
  timeout: 15000, // 15 秒
})
```

### Q2: 如何添加全局错误提示？

在响应拦截器中集成 Toast 组件（如 Element Plus 的 Message）：

```typescript
import { ElMessage } from 'element-plus'

service.interceptors.response.use(
  (response) => {
    if (res.code !== ErrorCode.SUCCESS) {
      ElMessage.error(res.message)
      return Promise.reject(new Error(res.message))
    }
    return response
  }
)
```

### Q3: 为什么使用 localStorage 而不是 Cookie？

- **优点**: 前后端分离，无需服务端配置 Cookie
- **缺点**: 容易受 XSS 攻击，建议配合 CSP 策略

生产环境可考虑切换到 HttpOnly Cookie 或使用 Secure Storage 方案。

### Q4: 如何取消请求？

使用 Axios 的 CancelToken：

```typescript
import axios from 'axios'

const source = axios.CancelToken.source()

request.get('/products', {
  cancelToken: source.token
})

// 取消请求
source.cancel('Operation canceled by the user.')
```

## 相关文件清单

```
src/utils/
├── request.ts     # Axios 封装与拦截器
└── auth.ts        # Token 存储管理
```

**建议新增文件**（按需）:
- `format.ts`: 日期、金额、文本格式化
- `validate.ts`: 表单验证规则
- `storage.ts`: 统一的本地存储封装（localStorage/sessionStorage）

---

**最后更新**: 2025-11-25 00:02:47
