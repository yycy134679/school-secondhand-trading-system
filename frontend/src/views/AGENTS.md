[根目录](../../CLAUDE.md) > **src/views**

# Views 模块（页面组件）

## 变更记录 (Changelog)

### 2025-11-25 00:02:47
- 初始化模块文档

---

## 模块职责

包含所有页面级 Vue 组件，对应路由表中的各个页面。每个页面负责页面布局、数据获取和用户交互。

## 入口与启动

通过 Vue Router 路由到各页面组件，在 `src/router/index.ts` 中配置。

### 路由懒加载

所有页面均使用动态导入实现懒加载：

```typescript
{
  path: '/products/:id',
  component: () => import('@/views/product/Detail.vue')
}
```

## 对外接口

### 页面结构

```
views/
├── home/
│   └── Index.vue          # 首页（推荐+最新商品）
├── product/
│   ├── Detail.vue         # 商品详情
│   ├── New.vue            # 发布商品（需登录）
│   ├── Edit.vue           # 编辑商品（需登录）
│   └── MyProducts.vue     # 我的商品列表（需登录）
├── user/
│   ├── Login.vue          # 登录
│   ├── Register.vue       # 注册
│   └── Profile.vue        # 个人中心（需登录）
├── category/
│   └── Index.vue          # 分类浏览
└── search/
    └── Index.vue          # 搜索结果
```

### 页面功能说明

| 页面 | 路径 | 需要登录 | 主要功能 |
|------|------|----------|----------|
| 首页 | `/` | 否 | 展示推荐商品和最新商品 |
| 商品详情 | `/products/:id` | 否 | 查看商品详情、联系卖家 |
| 发布商品 | `/products/new` | 是 | 发布新商品（多图上传、分类标签选择） |
| 编辑商品 | `/products/:id/edit` | 是 | 编辑商品信息 |
| 我的商品 | `/my/products` | 是 | 管理我的商品（上/下架、标记已售） |
| 登录 | `/login` | 否 | 用户登录（支持记住密码） |
| 注册 | `/register` | 否 | 用户注册 |
| 个人中心 | `/profile` | 是 | 查看和编辑个人资料、修改密码 |
| 分类浏览 | `/category/:id` | 否 | 按分类浏览商品列表 |
| 搜索结果 | `/search` | 否 | 搜索商品（支持多条件筛选） |

## 关键依赖与配置

### 页面组件通用依赖

- **Stores**: `useUserStore`, `useProductStore`, `useAppStore`
- **Router**: `useRouter`, `useRoute`
- **API**: 各业务 API 函数（如 `getProductDetail`, `searchProducts`）

### 路由元信息（meta）

在 `src/router/index.ts` 中配置：

```typescript
{
  path: '/products/new',
  meta: {
    requiresAuth: true,      // 需要登录
    title: '发布闲置',        // 页面标题
    breadcrumb: '发布闲置'    // 面包屑导航
  }
}
```

### 导航守卫

`router.beforeEach` 拦截未登录用户访问需认证页面：

```typescript
if (to.meta.requiresAuth && !userStore.isLoggedIn) {
  next({ path: '/login', query: { redirect: to.fullPath } })
}
```

## 数据模型

### 组件 Props

页面组件通常不接收 props（数据来自路由参数和 Pinia store）。

### 路由参数

- `products/:id` - 商品 ID
- `category/:id` - 分类 ID

### 页面 State

每个页面维护自己的响应式状态（如加载状态、表单数据、分页参数）：

```typescript
const loading = ref(false)
const product = ref<ProductDetail | null>(null)
```

## 测试与质量

### 测试建议

1. **组件渲染测试**: 使用 `@vue/test-utils` 挂载组件，验证 DOM 结构
2. **用户交互测试**: 模拟点击、表单提交等操作
3. **路由测试**: 验证路由跳转和参数传递
4. **状态集成测试**: 配合 Pinia store 测试数据流

示例：

```typescript
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import ProductDetail from './product/Detail.vue'

describe('ProductDetail', () => {
  it('should display product info', async () => {
    const wrapper = mount(ProductDetail, {
      global: {
        plugins: [createTestingPinia()],
        stubs: ['RouterLink']
      }
    })

    // 断言
    expect(wrapper.find('.product-title').text()).toBe('测试商品')
  })
})
```

## 常见问题 (FAQ)

### Q1: 如何在页面中获取路由参数？

使用 `useRoute()`:

```typescript
import { useRoute } from 'vue-router'

const route = useRoute()
const productId = Number(route.params.id)
```

### Q2: 如何设置页面标题？

在路由 meta 中配置 `title`，`router.beforeEach` 会自动设置：

```typescript
document.title = `${to.meta.title} - 校园二手交易平台`
```

### Q3: 如何处理页面加载状态？

使用本地 loading state：

```typescript
const loading = ref(false)

async function loadData() {
  loading.value = true
  try {
    await fetchProductDetail(productId)
  } finally {
    loading.value = false
  }
}
```

### Q4: 如何在页面间传递数据？

**推荐方式**:
1. 通过路由参数传递 ID，在目标页面重新获取数据
2. 使用 Pinia store 共享状态

**不推荐**: 通过 `query` 传递大量数据（URL 长度限制）

### Q5: 未登录访问需认证页面会发生什么？

路由守卫会拦截并重定向到登录页，同时记录原始目标地址：

```
/login?redirect=/products/new
```

登录成功后自动跳转回原地址。

## 相关文件清单

```
src/views/
├── home/
│   └── Index.vue
├── product/
│   ├── Detail.vue
│   ├── New.vue
│   ├── Edit.vue
│   └── MyProducts.vue
├── user/
│   ├── Login.vue
│   ├── Register.vue
│   └── Profile.vue
├── category/
│   └── Index.vue
└── search/
    └── Index.vue
```

**开发进度**: 所有页面骨架已创建，部分页面逻辑待实现（具体参考 `docs/前端开发任务.md`）。

---

**最后更新**: 2025-11-25 00:02:47
