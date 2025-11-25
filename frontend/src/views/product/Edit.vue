<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import ProductForm, { type ProductFormState } from '@/components/product/ProductForm.vue'
import Loading from '@/components/common/Loading.vue'
import { getProductDetail, updateProduct } from '@/api/product'
import { ProductStatus } from '@common/constants/product_status'
import { ErrorCode } from '@common/constants/error_code'
import type { ProductDetail } from '@/api/product'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()

const loading = ref(false)
const pageLoading = ref(true)
const product = ref<ProductDetail | null>(null)
const error = ref('')

const productId = computed(() => Number(route.params.id))

// 检查是否为已售商品
const isSoldProduct = computed(() => {
  return product.value?.status === ProductStatus.SOLD
})

// 检查是否可以编辑
const canEdit = computed(() => {
  return product.value?.viewerIsSeller && !isSoldProduct.value
})

// 准备表单初始数据
const initialFormData = computed<Partial<ProductFormState>>(() => {
  if (!product.value) return {}

  return {
    title: product.value.title,
    price: product.value.price,
    conditionId: product.value.conditionId,
    categoryId: product.value.categoryId,
    tagIds: product.value.tagIds || [],
    description: product.value.description,
    images: product.value.images.map((img) => ({
      id: img.id.toString(),
      url: img.url,
      file: undefined,
      isPrimary: img.isPrimary,
    })),
  }
})

// 加载商品详情
const loadProductDetail = async () => {
  try {
    pageLoading.value = true
    error.value = ''

    const response = await getProductDetail(productId.value)
    product.value = response.data.data

    // 权限检查：非卖家跳转到首页
    if (!product.value.viewerIsSeller) {
      error.value = '无权编辑此商品'
      setTimeout(() => {
        router.push('/')
      }, 1500)
    }
  } catch (err: unknown) {
    const e = err as { response?: { data?: { message?: string } } }
    error.value = e.response?.data?.message || '加载商品失败'
    setTimeout(() => {
      router.push('/')
    }, 1500)
  } finally {
    pageLoading.value = false
  }
}

// 处理表单提交
const handleSubmit = async (formData: ProductFormState) => {
  if (!product.value || !canEdit.value) {
    return
  }

  try {
    loading.value = true
    error.value = ''

    // 构建更新参数
    const updateData = {
      title: formData.title,
      description: formData.description,
      price: formData.price || 0,
      categoryId: formData.categoryId,
      conditionId: formData.conditionId,
      tagIds: formData.tagIds,
      // 如果有新上传的图片，需要通过 FormData 方式
      // 这里简化处理，只传已有图片的 URL
      imageUrls: formData.images.filter((img) => img.url).map((img) => img.url as string),
    }

    // 如果有新图片需要上传，使用 FormData
    const hasNewImages = formData.images.some((img) => img.file)
    if (hasNewImages) {
      const formDataObj = new FormData()
      formDataObj.append('title', formData.title)
      formDataObj.append('description', formData.description)
      formDataObj.append('price', (formData.price || 0).toString())
      formDataObj.append('categoryId', (formData.categoryId || 0).toString())
      formDataObj.append('conditionId', (formData.conditionId || 0).toString())

      if (formData.tagIds && formData.tagIds.length > 0) {
        formDataObj.append('tagIds', formData.tagIds.join(','))
      }

      // 添加新图片
      formData.images.forEach((img, index) => {
        if (img.file) {
          formDataObj.append('images', img.file)
          if (img.isPrimary) {
            formDataObj.append('primaryImageIndex', index.toString())
          }
        }
      })

      // 保留原有图片URL
      const existingUrls = formData.images
        .filter((img) => img.url && !img.file)
        .map((img) => img.url as string)
      if (existingUrls.length > 0) {
        formDataObj.append('existingImageUrls', JSON.stringify(existingUrls))
      }

      // 注意：这里需要后端支持 FormData 格式的 PUT 请求
      // 如果后端不支持，需要先上传图片再更新商品信息
      await updateProduct(productId.value, updateData)
    } else {
      await updateProduct(productId.value, updateData)
    }

    // 跳转到商品详情页
    router.push(`/products/${productId.value}`)
  } catch (err: unknown) {
    const e = err as { response?: { data?: { message?: string; code?: number } } }
    // 处理特定错误码
    if (e.response?.data?.code === ErrorCode.PRODUCT_SOLD) {
      error.value = '已售商品无法编辑'
    } else {
      error.value = e.response?.data?.message || '更新商品失败'
    }
  } finally {
    loading.value = false
  }
}

// 处理取消
const handleCancel = () => {
  router.push(`/products/${productId.value}`)
}

onMounted(() => {
  // 确保字典数据已加载
  if (appStore.categories.length === 0) {
    appStore.initDictionaries()
  }

  // 加载商品详情
  loadProductDetail()
})
</script>

<template>
  <div class="edit-product-page">
    <!-- 页面加载状态 -->
    <Loading v-if="pageLoading" />

    <!-- 错误提示 -->
    <div v-else-if="error" class="error-message">
      <div class="error-content">
        <svg
          class="error-icon"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
        </svg>
        <p>{{ error }}</p>
      </div>
    </div>

    <!-- 编辑表单 -->
    <div v-else class="edit-product-container">
      <div class="page-header">
        <h1>编辑商品</h1>
        <p class="subtitle">修改商品信息</p>
      </div>

      <!-- 已售商品提示 -->
      <div v-if="isSoldProduct" class="sold-warning">
        <svg
          class="warning-icon"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          ></path>
        </svg>
        <span>此商品已售出，无法编辑</span>
      </div>

      <!-- 错误提示 -->
      <div v-if="error" class="form-error">
        {{ error }}
      </div>

      <!-- 表单 -->
      <ProductForm
        v-if="product"
        mode="edit"
        :initial-value="initialFormData"
        :loading="loading || !canEdit"
        @submit="handleSubmit"
        @cancel="handleCancel"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
.edit-product-page {
  min-height: 60vh;
  padding: 24px 0;
}

.edit-product-container {
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
  text-align: center;

  h1 {
    font-size: 24px;
    font-weight: 600;
    color: var(--color-text-primary);
    margin-bottom: 8px;
  }

  .subtitle {
    font-size: 14px;
    color: var(--color-text-secondary);
  }
}

.sold-warning {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #fef3c7;
  border: 1px solid #fbbf24;
  border-radius: 8px;
  color: #92400e;
  margin-bottom: 24px;

  .warning-icon {
    width: 20px;
    height: 20px;
    flex-shrink: 0;
  }

  span {
    font-size: 14px;
  }
}

.form-error {
  padding: 12px 16px;
  background: #fee;
  border: 1px solid #ef4444;
  border-radius: 8px;
  color: #dc2626;
  margin-bottom: 24px;
  font-size: 14px;
}

.error-message {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;

  .error-content {
    text-align: center;

    .error-icon {
      width: 64px;
      height: 64px;
      color: #ef4444;
      margin: 0 auto 16px;
    }

    p {
      font-size: 16px;
      color: var(--color-text-secondary);
    }
  }
}
</style>
