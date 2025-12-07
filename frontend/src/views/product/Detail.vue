<template>
  <div class="product-detail-page">
    <ProductDetailSkeleton v-if="loading" />
    <div v-else-if="error" class="error-state">{{ error }}</div>
    <div v-else-if="product" class="product-content">
      <!-- 状态横幅 -->
      <div v-if="product.status !== 'ForSale'" class="status-banner">
        该商品{{ product.status === 'Sold' ? '已售出' : '已下架' }}
      </div>

      <div class="detail-container">
        <!-- 左侧图片区域 -->
        <div class="image-section">
          <div class="main-image">
            <img
              :src="activeImageUrl || '/placeholder.png'"
              :alt="product.title"
              @click="openLightbox"
            />
          </div>
          <div class="thumbnail-list" v-if="product.images && product.images.length > 1">
            <div
              v-for="(img, index) in product.images"
              :key="img.id"
              class="thumbnail-item"
              :class="{ active: activeImageIndex === index }"
              @click="activeImageIndex = index"
            >
              <img :src="img.url" :alt="`缩略图 ${index + 1}`" />
            </div>
          </div>
        </div>

        <!-- 右侧信息区域 -->
        <div class="info-section">
          <h1 class="product-title">{{ product.title }}</h1>

          <div class="product-price">{{ formatPrice(product.price) }}</div>

          <div class="tags-row">
            <ProductStatus :status="product.status" />
            <span class="condition-tag">{{ product.conditionName }}</span>
            <!-- 这里假设 categoryId 对应的名称需要额外获取或者后端返回了 categoryName，
                 根据 api.md ProductDetail 只有 categoryId。
                 暂时只显示 conditionName 和 status。
                 如果有 tags，也显示 tags。
            -->
          </div>

          <div class="meta-info">
            <span class="publish-time">发布于 {{ formatRelativeTime(product.createdAt) }}</span>
            <!-- <span class="view-count">浏览 {{ product.viewCount || 0 }} 次</span> -->
          </div>

          <div class="seller-card">
            <img
              :src="product.seller.avatarUrl || '/default-avatar.png'"
              class="seller-avatar"
              alt="卖家头像"
            />
            <div class="seller-info">
              <div class="seller-name">{{ product.seller.nickname }}</div>
              <div class="seller-desc">信誉极好</div>
            </div>
          </div>

          <div class="product-description">
            <h3>商品描述</h3>
            <p>{{ product.description }}</p>
          </div>

          <div class="action-area">
            <template v-if="product.viewerIsSeller">
              <button class="btn btn-primary" @click="goToEdit">编辑商品</button>
            </template>
            <template v-else>
              <button
                class="btn btn-primary btn-large"
                :disabled="product.status !== 'ForSale'"
                @click="contactSeller"
              >
                {{ product.status === 'ForSale' ? '联系卖家' : '不可购买' }}
              </button>
            </template>
          </div>
        </div>
      </div>
    </div>

    <!-- 登录弹窗 -->
    <LoginModal v-model:visible="showLoginModal" @success="handleLoginSuccess" />

    <!-- 联系卖家弹窗 -->
    <Modal v-model:visible="showContactModal" title="联系卖家" width="400px">
      <div class="contact-content">
        <div v-if="contactInfo.wechat" class="wechat-info">
          <div class="wechat-label">卖家微信号</div>
          <div class="wechat-value">
            <span class="wechat-id">{{ contactInfo.wechat }}</span>
            <button class="copy-btn" @click="copyWechat">复制</button>
            <span v-if="copyStatus" class="copy-tip">{{ copyStatus }}</span>
          </div>
        </div>
        <div class="contact-tips">{{ contactInfo.tips }}</div>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  getProductDetail,
  getProductContact,
  recordProductView,
  type ProductDetail,
} from '@/api/product'
import { formatPrice, formatRelativeTime } from '@/utils/format'
import { ErrorCode } from '@common/constants/error_code'
import ProductStatus from '@/components/product/ProductStatus.vue'
import ProductDetailSkeleton from '@/components/product/ProductDetailSkeleton.vue'
import LoginModal from '@/components/user/LoginModal.vue'
import Modal from '@/components/common/Modal.vue'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const product = ref<ProductDetail | null>(null)
const loading = ref(true)
const error = ref('')
const activeImageIndex = ref(0)

// 登录和联系卖家相关状态
const showLoginModal = ref(false)
const showContactModal = ref(false)
const contactLoading = ref(false)
const contactInfo = ref<{ wechat: string; tips: string }>({ wechat: '', tips: '' })
const copyStatus = ref('')
let copyTimer: number | null = null

const activeImageUrl = computed(() => {
  if (!product.value?.images?.length) return product.value?.mainImageUrl
  const img = product.value.images[activeImageIndex.value]
  return img ? img.url : product.value.mainImageUrl
})

const loadData = async () => {
  const id = Number(route.params.id)
  if (isNaN(id)) {
    error.value = '无效的商品ID'
    loading.value = false
    return
  }

  try {
    loading.value = true
    const response = await getProductDetail(id)
    const res = response.data
    if (res.code === 0) {
      product.value = res.data
      // 默认选中主图，或者第一张图
      const primaryIndex = res.data.images.findIndex((img) => img.isPrimary)
      activeImageIndex.value = primaryIndex >= 0 ? primaryIndex : 0
      await recordView(res.data.id)
    } else {
      error.value = res.message || '加载失败'
    }
  } catch (err) {
    error.value = '网络错误，请稍后重试'
    console.error(err)
  } finally {
    loading.value = false
  }
}

const openLightbox = () => {
  // TODO: 实现 Lightbox
  console.log('Open lightbox')
}

const goToEdit = () => {
  if (product.value) {
    router.push(`/products/${product.value.id}/edit`)
  }
}

const contactSeller = async () => {
  // 5.4.4.1 未登录时弹出登录框
  if (!userStore.isLoggedIn) {
    showLoginModal.value = true
    return
  }

  // 5.4.4.2 已登录且非卖家：调用 contact 接口
  if (!product.value) return

  try {
    contactLoading.value = true
    const response = await getProductContact(product.value.id)
    const res = response.data

    if (res.code === 0) {
      const data = res.data
      if (data.canContact && data.sellerWechat) {
        // 5.4.4.2.1 展示微信号弹窗
        contactInfo.value = {
          wechat: data.sellerWechat,
          tips: '请线下交易，注意安全',
        }
        showContactModal.value = true
      } else {
        // 5.4.4.2.2 降级提示
        contactInfo.value = {
          wechat: '',
          tips: data.tips || '卖家联系方式暂不可用，请稍后再试',
        }
        showContactModal.value = true
      }
    } else {
      alert(res.message || '获取联系方式失败')
    }
  } catch (err) {
    const error = err as { response?: { data?: { code?: number; message?: string } } }
    const errorCode = error?.response?.data?.code
    const errorMsg = error?.response?.data?.message || '网络错误，请稍后重试'

    // 处理状态机相关错误
    if (errorCode === ErrorCode.INVALID_STATUS_TRANSITION) {
      alert('商品状态已改变，无法联系卖家')
    } else if (errorCode === ErrorCode.PRODUCT_SOLD) {
      alert('该商品已售出')
    } else {
      alert(errorMsg)
    }
    console.error(err)
  } finally {
    contactLoading.value = false
  }
}

const handleLoginSuccess = async () => {
  showLoginModal.value = false
  // 登录成功后自动尝试联系卖家
  await contactSeller()
}

const copyWechat = () => {
  if (contactInfo.value.wechat) {
    navigator.clipboard.writeText(contactInfo.value.wechat).then(() => {
      copyStatus.value = '已复制'
      if (copyTimer) {
        window.clearTimeout(copyTimer)
      }
      copyTimer = window.setTimeout(() => {
        copyStatus.value = ''
        copyTimer = null
      }, 2000)
    })
  }
}

const recordView = async (productId: number) => {
  if (!userStore.isLoggedIn) return
  try {
    await recordProductView(productId)
  } catch (err) {
    console.warn('记录浏览失败:', err)
  }
}

onMounted(() => {
  loadData()
})

onBeforeUnmount(() => {
  if (copyTimer) {
    window.clearTimeout(copyTimer)
    copyTimer = null
  }
})
</script>

<style scoped lang="scss">
.product-detail-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 50px;
  font-size: 18px;
  color: #666;
}

.error-state {
  color: #ff4d4f;
}

.status-banner {
  background-color: #f5f5f5;
  color: #999;
  padding: 10px 20px;
  text-align: center;
  margin-bottom: 20px;
  border-radius: 4px;
  font-weight: bold;
}

.detail-container {
  display: flex;
  gap: 40px;

  @media (max-width: 768px) {
    flex-direction: column;
  }
}

.image-section {
  flex: 0 0 60%;
  max-width: 60%;

  @media (max-width: 768px) {
    flex: 0 0 100%;
    max-width: 100%;
  }

  .main-image {
    width: 100%;
    aspect-ratio: 4/3; // 或者 1:1，根据设计
    background-color: #f0f0f0;
    border-radius: 8px;
    overflow: hidden;
    margin-bottom: 16px;
    cursor: zoom-in;

    img {
      width: 100%;
      height: 100%;
      object-fit: contain; // 保持图片比例
    }
  }

  .thumbnail-list {
    display: flex;
    gap: 10px;
    overflow-x: auto;
    padding-bottom: 5px;

    .thumbnail-item {
      width: 80px;
      height: 80px;
      border-radius: 4px;
      overflow: hidden;
      cursor: pointer;
      border: 2px solid transparent;
      flex-shrink: 0;

      &.active {
        border-color: var(--color-primary, #0066ff);
      }

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }
  }
}

.info-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;

  .product-title {
    font-size: 24px;
    font-weight: 600;
    color: #333;
    margin: 0;
    line-height: 1.4;
  }

  .product-price {
    font-size: 28px;
    font-weight: bold;
    color: #ff4d4f; // 价格颜色
  }

  .tags-row {
    display: flex;
    gap: 10px;
    align-items: center;
    flex-wrap: wrap;

    .condition-tag {
      background-color: #e6f7ff;
      color: #1890ff;
      padding: 2px 8px;
      border-radius: 4px;
      font-size: 12px;
    }
  }

  .meta-info {
    display: flex;
    gap: 20px;
    color: #999;
    font-size: 14px;
  }

  .seller-card {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background-color: #f9f9f9;
    border-radius: 8px;

    .seller-avatar {
      width: 48px;
      height: 48px;
      border-radius: 50%;
      object-fit: cover;
    }

    .seller-info {
      .seller-name {
        font-weight: 600;
        color: #333;
      }
      .seller-desc {
        font-size: 12px;
        color: #666;
        margin-top: 4px;
      }
    }
  }

  .product-description {
    h3 {
      font-size: 16px;
      font-weight: 600;
      margin-bottom: 8px;
    }
    p {
      font-size: 14px;
      line-height: 1.6;
      color: #666;
      white-space: pre-wrap; // 保留换行
    }
  }

  .action-area {
    margin-top: auto;
    padding-top: 20px;

    .btn {
      padding: 10px 20px;
      border-radius: 4px;
      border: none;
      cursor: pointer;
      font-weight: 500;
      transition: all 0.3s;

      &.btn-primary {
        background-color: var(--color-primary, #0066ff);
        color: white;

        &:hover {
          opacity: 0.9;
        }

        &:disabled {
          background-color: #ccc;
          cursor: not-allowed;
        }
      }

      &.btn-large {
        width: 100%;
        height: 48px;
        font-size: 16px;
      }
    }
  }
}

.contact-content {
  padding: 20px 0;

  .wechat-info {
    margin-bottom: 20px;

    .wechat-label {
      font-size: 14px;
      color: #666;
      margin-bottom: 8px;
    }

      .wechat-value {
        display: flex;
        align-items: center;
        gap: 12px;

      .wechat-id {
        font-size: 18px;
        font-weight: 600;
        color: #333;
        flex: 1;
      }

      .copy-btn {
        padding: 6px 16px;
        background-color: var(--color-primary, #0066ff);
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 14px;
        transition: all 0.3s;

        &:hover {
          opacity: 0.9;
        }
      }

      .copy-tip {
        font-size: 12px;
        color: #52c41a;
        white-space: nowrap;
      }
    }
  }

  .contact-tips {
    padding: 12px;
    background-color: #fff7e6;
    border-left: 3px solid #faad14;
    color: #666;
    font-size: 14px;
    line-height: 1.6;
  }
}
</style>
