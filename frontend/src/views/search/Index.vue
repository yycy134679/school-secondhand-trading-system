<template>
  <div class="search-page">
    <!-- 搜索头部 -->
    <div class="search-header">
      <h1 class="search-title">搜索结果</h1>
      <p v-if="keyword" class="search-keyword">关键词: "{{ keyword }}"</p>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar" :class="{ sticky: isFilterSticky }">
      <div class="filter-bar-content">
        <!-- 排序 -->
        <div class="filter-item">
          <label>排序</label>
          <select v-model="filters.sort" class="filter-select" @change="onFilterChange">
            <option value="latest">最新发布</option>
            <option value="priceAsc">价格从低到高</option>
            <option value="priceDesc">价格从高到低</option>
          </select>
        </div>

        <!-- 价格区间 -->
        <div class="filter-item">
          <label>价格区间</label>
          <div class="price-range">
            <input
              v-model.number="filters.minPrice"
              type="number"
              placeholder="最低价"
              class="filter-input"
              @change="onFilterChange"
            />
            <span class="price-separator">-</span>
            <input
              v-model.number="filters.maxPrice"
              type="number"
              placeholder="最高价"
              class="filter-input"
              @change="onFilterChange"
            />
          </div>
        </div>

        <!-- 新旧程度 -->
        <div class="filter-item">
          <label>新旧程度</label>
          <div class="filter-tags">
            <button
              v-for="condition in productConditions"
              :key="condition.id"
              :class="['filter-tag', { active: isConditionSelected(condition.id) }]"
              @click="toggleCondition(condition.id)"
            >
              {{ condition.name }}
            </button>
          </div>
        </div>

        <!-- 发布时间 -->
        <div class="filter-item">
          <label>发布时间</label>
          <select
            v-model="filters.publishedTimeRange"
            class="filter-select"
            @change="onFilterChange"
          >
            <option value="all">全部</option>
            <option value="last_7_days">近7天</option>
            <option value="last_30_days">近30天</option>
          </select>
        </div>

        <!-- 重置按钮 -->
        <button class="reset-btn" @click="resetFilters">重置</button>
      </div>
    </div>

    <!-- 商品列表 -->
    <div class="products-container">
      <!-- 加载状态 - 骨架屏 -->
      <div v-if="loading" class="products-grid">
        <ProductCardSkeleton v-for="i in 20" :key="i" />
      </div>

      <!-- 商品网格 -->
      <div v-else-if="products.length > 0" class="products-grid">
        <ProductCard v-for="product in products" :key="product.id" :product="product" />
      </div>

      <!-- 空状态 -->
      <Empty
        v-else
        title="没有找到相关商品"
        description="换个词试试？或调整筛选条件"
        action-label="返回首页"
        @action="goHome"
      />
    </div>

    <!-- 分页 -->
    <Pagination
      v-if="!loading && total > 0"
      v-model:page="currentPage"
      :page-size="pageSize"
      :total="total"
      @update:page="onPageChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProductStore } from '@/stores/product'
import { useAppStore } from '@/stores/app'
import ProductCard from '@/components/product/ProductCard.vue'
import ProductCardSkeleton from '@/components/product/ProductCardSkeleton.vue'
import Empty from '@/components/common/Empty.vue'
import Pagination from '@/components/common/Pagination.vue'
import type { Product } from '@common/types/product'

const route = useRoute()
const router = useRouter()
const productStore = useProductStore()
const appStore = useAppStore()

// 状态
const loading = ref(false)
const products = ref<Product[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const isFilterSticky = ref(false)

// 关键词
const keyword = computed(() => (route.query.q as string) || '')

// 筛选条件
const filters = ref({
  sort: (route.query.sort as string) || 'latest',
  minPrice: route.query.minPrice ? Number(route.query.minPrice) : undefined,
  maxPrice: route.query.maxPrice ? Number(route.query.maxPrice) : undefined,
  conditionIds: route.query.conditionIds ? String(route.query.conditionIds) : '',
  publishedTimeRange: (route.query.publishedTimeRange as string) || 'all',
})

// 新旧程度列表
const productConditions = computed(() => appStore.productConditions)

// 判断新旧程度是否选中
const isConditionSelected = (id: number) => {
  const ids = filters.value.conditionIds.split(',').filter(Boolean).map(Number)
  return ids.includes(id)
}

// 切换新旧程度
const toggleCondition = (id: number) => {
  let ids = filters.value.conditionIds.split(',').filter(Boolean).map(Number)

  if (ids.includes(id)) {
    ids = ids.filter((i) => i !== id)
  } else {
    ids.push(id)
  }

  filters.value.conditionIds = ids.join(',')
  onFilterChange()
}

// 筛选变更
const onFilterChange = () => {
  currentPage.value = 1
  updateRoute()
}

// 页码变更
const onPageChange = (page: number) => {
  currentPage.value = page
  updateRoute()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 更新路由
const updateRoute = () => {
  const query: Record<string, string> = {}

  if (keyword.value) {
    query.q = keyword.value
  }
  if (filters.value.sort && filters.value.sort !== 'latest') {
    query.sort = filters.value.sort
  }
  if (filters.value.minPrice) {
    query.minPrice = String(filters.value.minPrice)
  }
  if (filters.value.maxPrice) {
    query.maxPrice = String(filters.value.maxPrice)
  }
  if (filters.value.conditionIds) {
    query.conditionIds = filters.value.conditionIds
  }
  if (filters.value.publishedTimeRange && filters.value.publishedTimeRange !== 'all') {
    query.publishedTimeRange = filters.value.publishedTimeRange
  }
  if (currentPage.value > 1) {
    query.page = String(currentPage.value)
  }

  router.push({ query })
}

// 重置筛选
const resetFilters = () => {
  filters.value = {
    sort: 'latest',
    minPrice: undefined,
    maxPrice: undefined,
    conditionIds: '',
    publishedTimeRange: 'all',
  }
  currentPage.value = 1
  updateRoute()
}

// 返回首页
const goHome = () => {
  router.push('/')
}

// 加载商品
const loadProducts = async () => {
  loading.value = true
  try {
    const params = {
      q: keyword.value || undefined,
      sort: filters.value.sort,
      minPrice: filters.value.minPrice,
      maxPrice: filters.value.maxPrice,
      conditionIds: filters.value.conditionIds || undefined,
      publishedTimeRange:
        filters.value.publishedTimeRange === 'all'
          ? undefined
          : filters.value.publishedTimeRange,
      page: currentPage.value,
      pageSize: pageSize.value,
    }

    const result = await productStore.searchProducts(params)
    products.value = result.items
    total.value = result.total
    currentPage.value = result.page
  } catch (error) {
    console.error('搜索商品失败:', error)
    products.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 监听路由变化
watch(
  () => route.query,
  (newQuery) => {
    // 同步筛选条件
    filters.value = {
      sort: (newQuery.sort as string) || 'latest',
      minPrice: newQuery.minPrice ? Number(newQuery.minPrice) : undefined,
      maxPrice: newQuery.maxPrice ? Number(newQuery.maxPrice) : undefined,
      conditionIds: newQuery.conditionIds ? String(newQuery.conditionIds) : '',
      publishedTimeRange: (newQuery.publishedTimeRange as string) || 'all',
    }
    currentPage.value = newQuery.page ? Number(newQuery.page) : 1

    // 重新加载
    loadProducts()
  },
  { immediate: true }
)

// 监听滚动，实现筛选栏吸顶
onMounted(() => {
  window.addEventListener('scroll', () => {
    isFilterSticky.value = window.scrollY > 200
  })
})
</script>

<style scoped lang="scss">
.search-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--spacing-lg) var(--spacing-md);
}

.search-header {
  margin-bottom: var(--spacing-lg);

  .search-title {
    font-size: 24px;
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: var(--spacing-sm);
  }

  .search-keyword {
    font-size: 14px;
    color: var(--text-secondary);
  }
}

.filter-bar {
  background: var(--bg-white);
  border-radius: var(--radius-lg);
  padding: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease-in-out;

  &.sticky {
    position: sticky;
    top: 80px;
    z-index: 100;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  .filter-bar-content {
    display: flex;
    flex-wrap: wrap;
    gap: var(--spacing-md);
    align-items: flex-start;
  }

  .filter-item {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-xs);

    label {
      font-size: 12px;
      color: var(--text-secondary);
      font-weight: 500;
    }
  }

  .filter-select {
    padding: 8px 12px;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    background: var(--bg-surface);
    font-size: 14px;
    color: var(--text-primary);
    cursor: pointer;
    transition: all 0.2s ease-in-out;

    &:focus {
      outline: none;
      border-color: var(--color-primary);
      background: var(--bg-white);
    }
  }

  .price-range {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);

    .filter-input {
      width: 100px;
      padding: 8px 12px;
      border: 1px solid var(--border-color);
      border-radius: var(--radius-md);
      background: var(--bg-surface);
      font-size: 14px;
      color: var(--text-primary);
      transition: all 0.2s ease-in-out;

      &:focus {
        outline: none;
        border-color: var(--color-primary);
        background: var(--bg-white);
      }

      &::placeholder {
        color: var(--text-caption);
      }
    }

    .price-separator {
      color: var(--text-caption);
    }
  }

  .filter-tags {
    display: flex;
    flex-wrap: wrap;
    gap: var(--spacing-xs);

    .filter-tag {
      padding: 6px 12px;
      border: 1px solid var(--border-color);
      border-radius: var(--radius-md);
      background: var(--bg-surface);
      font-size: 12px;
      color: var(--text-secondary);
      cursor: pointer;
      transition: all 0.2s ease-in-out;

      &:hover {
        border-color: var(--color-primary);
        color: var(--color-primary);
      }

      &.active {
        background: var(--color-primary);
        border-color: var(--color-primary);
        color: var(--bg-white);
      }
    }
  }

  .reset-btn {
    padding: 8px 16px;
    background: var(--bg-surface);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    font-size: 14px;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.2s ease-in-out;
    align-self: flex-end;

    &:hover {
      background: var(--bg-white);
      border-color: var(--color-primary);
      color: var(--color-primary);
    }
  }
}

.products-container {
  min-height: 400px;
  margin-bottom: var(--spacing-xl);
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: var(--spacing-lg);
}

@media (max-width: 768px) {
  .search-page {
    padding: var(--spacing-md) var(--spacing-sm);
  }

  .filter-bar-content {
    flex-direction: column;
  }

  .filter-item {
    width: 100%;
  }

  .products-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: var(--spacing-md);
  }
}
</style>
