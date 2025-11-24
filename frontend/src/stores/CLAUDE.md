[根目录](../../CLAUDE.md) > **src/stores**

# Stores 模块（状态管理）

## 变更记录 (Changelog)

### 2025-11-25 00:02:47
- 初始化模块文档

---

## 模块职责

使用 Pinia 管理全局应用状态，包括用户认证、商品数据、全局字典等。所有状态变更逻辑集中在此模块。

## 入口与启动

在 `src/main.ts` 中初始化：

```typescript
import pinia from './stores'
app.use(pinia)
```

`stores/index.ts` 导出已配置的 Pinia 实例。

## 对外接口

### user.ts - 用户状态

**State:**
- `token`: JWT token (同步到 localStorage)
- `currentUser`: 当前登录用户信息
- `rememberMe`: 是否记住登录状态

**Getters:**
- `isLoggedIn`: 是否已登录（计算属性）
- `isAdmin`: 是否为管理员（计算属性）

**Actions:**
- `login(loginForm)`: 登录并保存 token
- `logout()`: 退出登录，清除状态
- `fetchProfile()`: 获取当前用户信息
- `updateProfile(data)`: 更新用户资料
- `changePassword(data)`: 修改密码

### product.ts - 商品状态

**State:**
- `homeRecommendations`: 首页推荐商品列表
- `homeLatest`: 首页最新商品分页数据
- `searchParams`: 搜索参数缓存
- `searchResults`: 搜索结果分页数据
- `currentProduct`: 当前查看的商品详情

**Actions:**
- `fetchHomeData(params)`: 获取首页数据
- `searchProducts(params)`: 搜索商品
- `fetchCategoryProducts(categoryId, params)`: 按分类获取商品
- `fetchMyProducts(params)`: 获取我的商品列表
- `fetchProductDetail(id)`: 获取商品详情
- `changeStatus(id, action)`: 变更商品状态

### app.ts - 全局状态

**State:**
- `loading`: 全局加载状态
- `error`: 错误信息
- `categories`: 商品分类列表（字典）
- `tags`: 商品标签列表（字典）
- `productConditions`: 新旧程度枚举（字典）

**Actions:**
- `initDictionaries()`: 初始化全局字典数据（建议在 App.vue mounted 时调用）

## 关键依赖与配置

### 依赖关系

```
stores/*.ts
  ↓ 调用
api/*.ts (API 接口)
  ↓ 依赖
@common/types/ (共享类型)
```

### Pinia 配置

默认配置，未启用插件（如持久化插件）。token 通过 `utils/auth.ts` 手动同步到 localStorage。

## 数据模型

所有 state 类型与 `@common/types/` 中的实体类型对应：
- `User` - 用户
- `Product` - 商品
- `Category` - 分类
- `Tag` - 标签
- `ProductCondition` - 新旧程度
- `PageResult<T>` - 分页数据

## 测试与质量

### 当前状态

- 无单元测试（建议使用 Vitest + pinia.testing）
- 类型检查覆盖所有 actions 和 getters

### 测试建议

1. **Actions 测试**: 使用 `createPinia()` 创建测试实例，验证状态变更逻辑
2. **API Mock**: 使用 Vitest 的 `vi.mock()` 模拟 API 调用
3. **Getters 测试**: 验证计算属性的正确性

示例：

```typescript
import { setActivePinia, createPinia } from 'pinia'
import { useUserStore } from './user'

describe('User Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should login successfully', async () => {
    const store = useUserStore()
    await store.login({ account: 'test', password: '123456' })
    expect(store.isLoggedIn).toBe(true)
  })
})
```

## 常见问题 (FAQ)

### Q1: 什么时候用 Store，什么时候用组件 state？

- **Store**: 需要跨组件共享的状态（如用户信息、商品列表）
- **组件 state**: 仅在单个组件内使用的状态（如表单输入、加载状态）

### Q2: 为什么 token 不在 Pinia 中持久化？

考虑到安全性和灵活性，token 通过 `utils/auth.ts` 手动管理，可以在未来轻松切换存储方式（如 Cookie、SessionStorage）。

### Q3: 如何在 store 之间共享数据？

直接引入其他 store：

```typescript
import { useUserStore } from './user'

export const useProductStore = defineStore('product', () => {
  const userStore = useUserStore()

  // 使用 userStore.currentUser
})
```

### Q4: 何时调用 `initDictionaries()`？

建议在 `App.vue` 的 `onMounted` 钩子中调用：

```typescript
import { useAppStore } from '@/stores/app'

onMounted(async () => {
  const appStore = useAppStore()
  await appStore.initDictionaries()
})
```

## 相关文件清单

```
src/stores/
├── index.ts       # Pinia 实例
├── user.ts        # 用户状态（认证、个人信息）
├── product.ts     # 商品状态（列表、详情）
├── app.ts         # 全局状态（字典、加载、错误）
└── counter.ts     # （模板示例，可删除）
```

---

**最后更新**: 2025-11-25 00:02:47
