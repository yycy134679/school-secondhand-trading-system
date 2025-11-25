[根目录](../../CLAUDE.md) > **src/components**

# Components 模块（公共组件）

## 变更记录 (Changelog)

### 2025-11-25 00:02:47
- 初始化模块文档

---

## 模块职责

提供可复用的 UI 组件，包括布局组件、业务组件和基础组件。所有组件应保持单一职责、高复用性和良好的类型定义。

## 入口与启动

无独立入口，各组件按需在页面或其他组件中导入使用。

## 对外接口

### 当前组件

```
components/
└── common/
    ├── AppHeader.vue      # 全局顶部导航
    └── AppFooter.vue      # 全局底部栏
```

### AppHeader.vue - 顶部导航

**功能**:
- Logo 和站点名称
- 搜索框（待完善）
- 用户状态展示（登录/未登录）
- 发布按钮（跳转到发布页面）

**Props**: 无

**Emits**: 无（通过 `router-link` 实现导航）

**使用示例**:

```vue
<template>
  <AppHeader />
</template>

<script setup lang="ts">
import AppHeader from '@/components/common/AppHeader.vue'
</script>
```

**待完善功能**:
- 搜索框提交逻辑
- 用户头像和下拉菜单
- 未登录时点击"发布"跳转登录页

### AppFooter.vue - 底部栏

**功能**:
- 版权信息
- 友情链接
- 备案号等

**Props**: 无

**Emits**: 无

**使用示例**:

```vue
<template>
  <AppFooter />
</template>

<script setup lang="ts">
import AppFooter from '@/components/common/AppFooter.vue'
</script>
```

## 关键依赖与配置

### 样式规范

- 使用 `scoped` 样式隔离
- 引用全局 CSS 变量（见 `assets/styles/theme.scss`）
- 组件内部样式使用 SCSS

### 组件设计原则

1. **单一职责**: 一个组件只做一件事
2. **Props 驱动**: 通过 props 传递数据，避免直接访问 store
3. **类型安全**: 使用 TypeScript 定义 props 和 emits
4. **无副作用**: 组件不应直接修改外部状态（通过 emits 通知父组件）

## 数据模型

### Props 定义示例

```typescript
interface Props {
  title: string
  visible?: boolean
}

const props = defineProps<Props>()
```

### Emits 定义示例

```typescript
interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'submit', data: FormData): void
}

const emit = defineEmits<Emits>()
```

## 测试与质量

### 测试建议

1. **组件渲染**: 验证组件在不同 props 下的渲染结果
2. **事件触发**: 模拟用户操作，验证 emits 是否正确触发
3. **插槽测试**: 验证 slot 内容是否正确显示
4. **样式测试**: 使用快照测试验证样式一致性

示例（使用 Vitest + @vue/test-utils）:

```typescript
import { mount } from '@vue/test-utils'
import AppHeader from './AppHeader.vue'

describe('AppHeader', () => {
  it('should render logo', () => {
    const wrapper = mount(AppHeader)
    expect(wrapper.find('.logo').text()).toContain('校园二手')
  })

  it('should show search bar', () => {
    const wrapper = mount(AppHeader)
    expect(wrapper.find('.search-bar input').exists()).toBe(true)
  })
})
```

## 常见问题 (FAQ)

### Q1: 何时创建新组件？

当满足以下条件之一时：
- 同一段 UI 在多个页面中重复使用（如商品卡片）
- 单个页面组件代码超过 200 行，可拆分为多个子组件
- 独立的业务逻辑块（如评论区、购物车）

### Q2: 组件应该放在哪个目录？

- `common/`: 通用布局组件（Header, Footer, Sidebar）
- `product/`: 商品相关业务组件（ProductCard, ProductList）
- `user/`: 用户相关组件（UserAvatar, UserInfo）
- `base/`: 基础 UI 组件（Button, Input, Modal）

### Q3: 如何在组件中使用 Pinia store？

**不推荐**: 组件内部直接使用 store（降低复用性）

**推荐**: 通过 props 传递数据，由父组件从 store 获取：

```vue
<!-- 父组件 -->
<script setup lang="ts">
import { useProductStore } from '@/stores/product'

const productStore = useProductStore()
const products = productStore.homeRecommendations
</script>

<template>
  <ProductList :items="products" />
</template>

<!-- 子组件 -->
<script setup lang="ts">
interface Props {
  items: Product[]
}
const props = defineProps<Props>()
</script>
```

### Q4: 如何优化组件性能？

1. **使用 `v-once`**: 对不会变化的内容只渲染一次
2. **使用 `v-memo`**: 缓存列表项渲染结果
3. **异步组件**: 大型组件使用 `defineAsyncComponent` 懒加载
4. **避免深层响应式**: 对大型静态数据使用 `shallowRef`

## 相关文件清单

```
src/components/
└── common/
    ├── AppHeader.vue
    └── AppFooter.vue
```

**建议新增组件**（按开发需要）:
- `product/ProductCard.vue`: 商品卡片（用于列表展示）
- `product/ProductList.vue`: 商品列表容器
- `base/Button.vue`: 统一按钮样式
- `base/Modal.vue`: 弹窗组件
- `base/Pagination.vue`: 分页组件
- `user/UserAvatar.vue`: 用户头像组件

---

**最后更新**: 2025-11-25 00:02:47
