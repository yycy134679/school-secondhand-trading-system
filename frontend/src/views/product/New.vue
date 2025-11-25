<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import ProductForm, { type ProductFormState } from '@/components/product/ProductForm.vue'
import Modal from '@/components/common/Modal.vue'
import { createProduct } from '@/api/product'
import { updateProfile } from '@/api/user'

const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()

const loading = ref(false)
const showWechatModal = ref(false)
const wechatId = ref('')
const wechatError = ref('')
const pendingFormData = ref<ProductFormState | null>(null)

// 检查用户是否填写了微信号
const needsWechatId = computed(() => {
  return !userStore.currentUser?.wechatId
})

onMounted(() => {
  // 确保字典数据已加载
  if (appStore.categories.length === 0) {
    appStore.initDictionaries()
  }
})

// 处理表单提交
const handleSubmit = async (formData: ProductFormState) => {
  // 检查微信号
  if (needsWechatId.value) {
    pendingFormData.value = formData
    showWechatModal.value = true
    return
  }

  await submitProduct(formData)
}

// 提交微信号
const submitWechatId = async () => {
  wechatError.value = ''

  if (!wechatId.value.trim()) {
    wechatError.value = '请输入微信号'
    return
  }

  try {
    loading.value = true
    await updateProfile({ wechatId: wechatId.value.trim() })
    // 更新用户信息
    await userStore.fetchProfile()
    showWechatModal.value = false

    // 继续提交商品
    if (pendingFormData.value) {
      await submitProduct(pendingFormData.value)
      pendingFormData.value = null
    }
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    wechatError.value = err.response?.data?.message || '更新微信号失败'
  } finally {
    loading.value = false
  }
}

// 提交商品
const submitProduct = async (formData: ProductFormState) => {
  try {
    loading.value = true

    // 构建 FormData
    const data = new FormData()
    data.append('title', formData.title)
    data.append('description', formData.description)
    data.append('price', formData.price?.toString() || '0')
    data.append('categoryId', formData.categoryId?.toString() || '')
    data.append('conditionId', formData.conditionId?.toString() || '')

    // tagIds 转为逗号分隔的字符串
    if (formData.tagIds && formData.tagIds.length > 0) {
      data.append('tagIds', formData.tagIds.join(','))
    }

    // 添加图片文件
    formData.images.forEach((img, index) => {
      if (img.file) {
        data.append('images', img.file)
        // 标记主图
        if (img.isPrimary) {
          data.append('primaryImageIndex', index.toString())
        }
      }
    })

    const res = await createProduct(data)
    const product = res.data.data

    // 跳转到商品详情页
    router.push(`/products/${product.id}`)
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    const message = err.response?.data?.message || '发布失败，请重试'
    alert(message)
  } finally {
    loading.value = false
  }
}

// 取消发布
const handleCancel = () => {
  router.back()
}
</script>

<template>
  <div class="new-product-page">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">发布闲置</h1>
        <p class="page-subtitle">快速发布你的闲置物品，让更多人看到</p>
      </div>

      <div class="form-container">
        <ProductForm
          mode="create"
          :loading="loading"
          @submit="handleSubmit"
          @cancel="handleCancel"
        />
      </div>
    </div>

    <!-- 微信号补全模态框 -->
    <Modal v-model:visible="showWechatModal" title="完善联系方式">
      <div class="wechat-modal">
        <p class="modal-description">发布商品前需要填写微信号，以便买家联系你</p>

        <div class="form-group">
          <label class="label">微信号</label>
          <input
            v-model="wechatId"
            type="text"
            class="input"
            placeholder="请输入你的微信号"
            :class="{ error: wechatError }"
          />
          <span v-if="wechatError" class="error-text">{{ wechatError }}</span>
        </div>

        <div class="modal-actions">
          <button
            type="button"
            class="btn btn-secondary"
            @click="showWechatModal = false"
            :disabled="loading"
          >
            取消
          </button>
          <button type="button" class="btn btn-primary" @click="submitWechatId" :disabled="loading">
            {{ loading ? '保存中...' : '保存并继续' }}
          </button>
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped lang="scss">
.new-product-page {
  min-height: calc(100vh - 64px - 80px);
  padding: 40px 0;
  background-color: #f7f8fa;
}

.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 20px;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.page-subtitle {
  font-size: 14px;
  color: #666666;
}

.form-container {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.wechat-modal {
  padding: 8px;
}

.modal-description {
  font-size: 14px;
  color: #666666;
  margin-bottom: 24px;
  line-height: 1.6;
}

.form-group {
  margin-bottom: 24px;
}

.label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.input {
  width: 100%;
  height: 40px;
  padding: 0 12px;
  font-size: 14px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background-color: #f7f8fa;
  transition: all 0.2s ease-in-out;

  &:focus {
    outline: none;
    background-color: #ffffff;
    border-color: #0066ff;
    box-shadow: 0 0 0 3px rgba(0, 102, 255, 0.1);
  }

  &.error {
    border-color: #ef4444;
    background-color: #fef2f2;
  }
}

.error-text {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: #ef4444;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 24px;
}

.btn {
  height: 40px;
  padding: 0 24px;
  font-size: 14px;
  font-weight: 500;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease-in-out;

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

.btn-primary {
  background-color: #0066ff;
  color: #ffffff;

  &:hover:not(:disabled) {
    background-color: #0052cc;
  }
}

.btn-secondary {
  background-color: #f7f8fa;
  color: #1a1a1a;

  &:hover:not(:disabled) {
    background-color: #e5e7eb;
  }
}
</style>
