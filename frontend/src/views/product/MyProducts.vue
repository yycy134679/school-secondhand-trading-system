<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMyProducts, changeProductStatus, undoProductStatusChange } from '@/api/product'
import type { Product } from '@common/types/product'
import { ErrorCode } from '@common/constants/error_code'
import ProductStatus from '@/components/product/ProductStatus.vue'
import Pagination from '@/components/common/Pagination.vue'
import Loading from '@/components/common/Loading.vue'
import Empty from '@/components/common/Empty.vue'
import { formatPrice, formatRelativeTime } from '@/utils/format'

const router = useRouter()

// 状态管理
const loading = ref(false)
const products = ref<Product[]>([])
const keyword = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 撤销相关
const undoInfo = ref<{
  productId: number
  action: string
  timer: number | null
  countdown: number
}>({
  productId: 0,
  action: '',
  timer: null,
  countdown: 0,
})

// 加载商品列表
const loadProducts = async () => {
  loading.value = true
  try {
    const response = await getMyProducts({
      keyword: keyword.value || undefined,
      page: page.value,
      pageSize: pageSize.value,
    })

    if (response.data.code === 0 && response.data.data) {
      products.value = response.data.data.items
      total.value = response.data.data.total
    }
  } catch (error) {
    console.error('加载商品失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  page.value = 1
  loadProducts()
}

// 分页变化
const handlePageChange = (newPage: number) => {
  page.value = newPage
  loadProducts()
}

// 编辑商品
const handleEdit = (productId: number) => {
  router.push(`/products/${productId}/edit`)
}

// 状态变更
const handleStatusChange = async (productId: number, action: 'delist' | 'relist' | 'sold') => {
  try {
    const response = await changeProductStatus(productId, { action })

    if (response.data.code === 0) {
      // 更新本地列表
      await loadProducts()

      // 对于下架和重新上架操作，显示撤销提示
      if (action === 'delist' || action === 'relist') {
        showUndoNotification(productId, action)
      } else {
        // 标记已售成功提示
        alert('商品已标记为已售')
      }
    }
  } catch (error) {
    const errorMsg =
      (error as { response?: { data?: { message?: string } } })?.response?.data?.message ||
      '操作失败'
    alert(errorMsg)
  }
}

// 显示撤销通知
const showUndoNotification = (productId: number, action: string) => {
  // 清除之前的计时器
  if (undoInfo.value.timer) {
    clearInterval(undoInfo.value.timer)
  }

  undoInfo.value = {
    productId,
    action,
    countdown: 3,
    timer: null,
  }

  // 倒计时
  undoInfo.value.timer = window.setInterval(() => {
    if (undoInfo.value.countdown > 0) {
      undoInfo.value.countdown--
    } else {
      hideUndoNotification()
    }
  }, 1000)
}

// 隐藏撤销通知
const hideUndoNotification = () => {
  if (undoInfo.value.timer) {
    clearInterval(undoInfo.value.timer)
  }
  undoInfo.value = {
    productId: 0,
    action: '',
    timer: null,
    countdown: 0,
  }
}

// 执行撤销
const handleUndo = async () => {
  const productId = undoInfo.value.productId
  hideUndoNotification()

  try {
    const response = await undoProductStatusChange(productId)

    if (response.data.code === 0) {
      await loadProducts()
      alert('已撤销操作')
    }
  } catch (error) {
    const err = error as { response?: { data?: { code?: number; message?: string } } }
    const errorCode = err?.response?.data?.code
    const errorMsg = err?.response?.data?.message || '撤销失败'

    if (errorCode === ErrorCode.REVOKE_FAILED) {
      alert('撤销超时或状态已改变')
    } else {
      alert(errorMsg)
    }
  }
}

// 获取状态操作按钮
const getStatusActions = (status: Product['status']) => {
  switch (status) {
    case 'ForSale':
      return [
        { label: '下架', action: 'delist' as const, className: 'btn-warning' },
        { label: '标记已售', action: 'sold' as const, className: 'btn-primary' },
      ]
    case 'Delisted':
      return [{ label: '重新上架', action: 'relist' as const, className: 'btn-success' }]
    case 'Sold':
      return []
    default:
      return []
  }
}

// 初始化
onMounted(() => {
  loadProducts()
})
</script>

<template>
  <div class="my-products-page">
    <div class="page-header">
      <h1>我发布的商品</h1>
      <div class="search-bar">
        <input
          v-model="keyword"
          type="text"
          placeholder="搜索我的商品..."
          @keyup.enter="handleSearch"
        />
        <button @click="handleSearch" class="btn-search">搜索</button>
      </div>
    </div>

    <!-- 撤销通知 -->
    <div v-if="undoInfo.productId" class="undo-notification">
      <span>操作成功，{{ undoInfo.countdown }}秒内可撤销</span>
      <button @click="handleUndo" class="btn-undo">撤销</button>
      <button @click="hideUndoNotification" class="btn-close">×</button>
    </div>

    <!-- 加载状态 -->
    <Loading v-if="loading" />

    <!-- 空状态 -->
    <Empty
      v-else-if="!loading && products.length === 0"
      title="暂无商品"
      description="还没有发布过商品，去清空闲置吧！"
      actionLabel="去发布"
      @action="router.push('/products/new')"
    />

    <!-- 商品列表 -->
    <div v-else class="products-list">
      <div v-for="product in products" :key="product.id" class="product-item">
        <div class="product-info" @click="router.push(`/products/${product.id}`)">
          <img
            :src="product.mainImageUrl || '/placeholder.png'"
            :alt="product.title"
            class="product-thumbnail"
          />
          <div class="product-details">
            <h3 class="product-title">{{ product.title }}</h3>
            <div class="product-meta">
              <span class="product-price">{{ formatPrice(product.price) }}</span>
              <ProductStatus :status="product.status" />
              <span class="product-time">{{ formatRelativeTime(product.createdAt) }}</span>
            </div>
          </div>
        </div>
        <div class="product-actions">
          <button @click="handleEdit(product.id)" class="btn-action btn-edit">编辑</button>
          <button
            v-for="statusAction in getStatusActions(product.status)"
            :key="statusAction.action"
            @click="handleStatusChange(product.id, statusAction.action)"
            class="btn-action"
            :class="statusAction.className"
          >
            {{ statusAction.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="!loading && products.length > 0" class="pagination-wrapper">
      <Pagination
        :page="page"
        :pageSize="pageSize"
        :total="total"
        @update:page="handlePageChange"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
.my-products-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    h1 {
      font-size: 24px;
      font-weight: 600;
      color: var(--color-text-primary, #1a1a1a);
      margin: 0;
    }

    .search-bar {
      display: flex;
      gap: 8px;

      input {
        width: 300px;
        height: 40px;
        padding: 0 16px;
        border: 1px solid var(--color-border, #e5e7eb);
        border-radius: 8px;
        font-size: 14px;
        background-color: #f7f8fa;
        transition: all 0.2s;

        &:focus {
          outline: none;
          background-color: #fff;
          border-color: var(--color-primary, #0066ff);
          box-shadow: 0 0 0 3px rgba(0, 102, 255, 0.1);
        }
      }

      .btn-search {
        height: 40px;
        padding: 0 24px;
        background-color: var(--color-primary, #0066ff);
        color: white;
        border: none;
        border-radius: 8px;
        font-size: 14px;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s;

        &:hover {
          background-color: #0052cc;
        }
      }
    }
  }

  .undo-notification {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    margin-bottom: 16px;
    background-color: #fff3cd;
    border: 1px solid #ffc107;
    border-radius: 8px;
    color: #856404;
    font-size: 14px;

    .btn-undo {
      padding: 4px 16px;
      background-color: var(--color-primary, #0066ff);
      color: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      font-size: 14px;
      transition: all 0.2s;

      &:hover {
        background-color: #0052cc;
      }
    }

    .btn-close {
      padding: 0;
      width: 24px;
      height: 24px;
      background: none;
      border: none;
      font-size: 24px;
      line-height: 1;
      color: #856404;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        color: #000;
      }
    }
  }

  .products-list {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);

    .product-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16px;
      border-bottom: 1px solid var(--color-border, #e5e7eb);
      transition: background-color 0.2s;

      &:last-child {
        border-bottom: none;
      }

      &:hover {
        background-color: #f7f8fa;
      }

      .product-info {
        display: flex;
        gap: 16px;
        flex: 1;
        cursor: pointer;

        .product-thumbnail {
          width: 80px;
          height: 80px;
          object-fit: cover;
          border-radius: 8px;
          background-color: #f5f5f5;
        }

        .product-details {
          flex: 1;
          display: flex;
          flex-direction: column;
          justify-content: center;

          .product-title {
            margin: 0 0 8px;
            font-size: 16px;
            font-weight: 500;
            color: var(--color-text-primary, #1a1a1a);
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
          }

          .product-meta {
            display: flex;
            align-items: center;
            gap: 16px;

            .product-price {
              font-size: 18px;
              font-weight: 600;
              color: var(--color-primary, #0066ff);
            }

            .product-time {
              font-size: 12px;
              color: var(--color-text-secondary, #999);
            }
          }
        }
      }

      .product-actions {
        display: flex;
        gap: 8px;
        flex-shrink: 0;

        .btn-action {
          height: 32px;
          padding: 0 16px;
          border: none;
          border-radius: 6px;
          font-size: 14px;
          font-weight: 500;
          cursor: pointer;
          transition: all 0.2s;
          white-space: nowrap;

          &.btn-edit {
            background-color: #f7f8fa;
            color: var(--color-text-primary, #1a1a1a);

            &:hover {
              background-color: #e5e7eb;
            }
          }

          &.btn-primary {
            background-color: var(--color-primary, #0066ff);
            color: white;

            &:hover {
              background-color: #0052cc;
            }
          }

          &.btn-warning {
            background-color: #ef4444;
            color: white;

            &:hover {
              background-color: #dc2626;
            }
          }

          &.btn-success {
            background-color: #10b981;
            color: white;

            &:hover {
              background-color: #059669;
            }
          }
        }
      }
    }
  }

  .pagination-wrapper {
    margin-top: 24px;
    display: flex;
    justify-content: center;
  }
}
</style>
