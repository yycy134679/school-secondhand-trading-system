<script setup lang="ts">
import { computed, onMounted, reactive, watch } from 'vue'
import { useAppStore } from '@/stores/app'
import ProductImageUpload, { type UploadImage } from './ProductImageUpload.vue'

export interface ProductFormState {
  title: string
  price: number | undefined
  conditionId: number | undefined
  categoryId: number | undefined
  tagIds: number[]
  description: string
  images: UploadImage[]
}

const props = defineProps<{
  mode: 'create' | 'edit'
  initialValue?: Partial<ProductFormState>
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'submit', payload: ProductFormState): void
  (e: 'cancel'): void
}>()

const appStore = useAppStore()

const form = reactive<ProductFormState>({
  title: '',
  price: undefined,
  conditionId: undefined,
  categoryId: undefined,
  tagIds: [],
  description: '',
  images: [],
})

const errors = reactive<Partial<Record<keyof ProductFormState, string>>>({})

onMounted(() => {
  if (props.initialValue) {
    Object.assign(form, props.initialValue)
  }
  // Ensure dictionaries are loaded
  if (appStore.categories.length === 0) {
    appStore.initDictionaries()
  }
})

const filteredTags = computed(() => {
  if (!form.categoryId) {
    return []
  }

  return appStore.tags.filter((tag) => tag.categoryId === form.categoryId)
})

watch(
  filteredTags,
  (tags) => {
    if (!form.categoryId) {
      form.tagIds = []
      return
    }
    if (appStore.tags.length === 0) {
      return
    }
    const validIds = tags.map((tag) => tag.id)
    form.tagIds = form.tagIds.filter((id) => validIds.includes(id))
  },
  { immediate: true },
)

const validate = (): boolean => {
  let isValid = true
  errors.title = ''
  errors.price = ''
  errors.categoryId = ''
  errors.conditionId = ''
  errors.description = ''
  errors.images = ''

  if (!form.title.trim()) {
    errors.title = '请输入商品标题'
    isValid = false
  }
  if (!form.price || form.price <= 0) {
    errors.price = '请输入有效的价格'
    isValid = false
  }
  if (!form.categoryId) {
    errors.categoryId = '请选择分类'
    isValid = false
  }
  if (!form.conditionId) {
    errors.conditionId = '请选择新旧程度'
    isValid = false
  }
  if (!form.description.trim()) {
    errors.description = '请输入商品描述'
    isValid = false
  }
  if (form.images.length === 0) {
    errors.images = '请至少上传一张图片'
    isValid = false
  }

  return isValid
}

const toggleTag = (tagId: number) => {
  const index = form.tagIds.indexOf(tagId)
  if (index > -1) {
    form.tagIds.splice(index, 1)
  } else {
    form.tagIds.push(tagId)
  }
}

const handleSubmit = () => {
  if (!validate()) return
  emit('submit', { ...form })
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="product-form">
    <div class="form-group">
      <label class="label">标题</label>
      <input
        v-model="form.title"
        type="text"
        class="input"
        placeholder="请输入商品标题"
        :class="{ error: errors.title }"
      />
      <span class="error-text" v-if="errors.title">{{ errors.title }}</span>
    </div>

    <div class="form-group">
      <label class="label">价格</label>
      <input
        v-model.number="form.price"
        type="number"
        step="0.01"
        class="input"
        placeholder="0.00"
        :class="{ error: errors.price }"
      />
      <span class="error-text" v-if="errors.price">{{ errors.price }}</span>
    </div>

    <div class="form-group">
      <label class="label">分类</label>
      <select v-model="form.categoryId" class="select" :class="{ error: errors.categoryId }">
        <option :value="undefined" disabled>请选择分类</option>
        <option v-for="cat in appStore.categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </option>
      </select>
      <span class="error-text" v-if="errors.categoryId">{{ errors.categoryId }}</span>
    </div>

    <div class="form-group">
      <label class="label">标签</label>
      <div class="tags-container" :class="{ disabled: !form.categoryId }">
        <div v-if="!form.categoryId" class="empty-text">请先选择分类</div>
        <div v-else-if="filteredTags.length === 0" class="empty-text">该分类暂无可选标签</div>
        <div
          v-else
          v-for="tag in filteredTags"
          :key="tag.id"
          class="tag-item"
          :class="{ active: form.tagIds.includes(tag.id) }"
          @click="toggleTag(tag.id)"
        >
          {{ tag.name }}
        </div>
      </div>
    </div>

    <div class="form-group">
      <label class="label">新旧程度</label>
      <select v-model="form.conditionId" class="select" :class="{ error: errors.conditionId }">
        <option :value="undefined" disabled>请选择新旧程度</option>
        <option v-for="cond in appStore.productConditions" :key="cond.id" :value="cond.id">
          {{ cond.name }}
        </option>
      </select>
      <span class="error-text" v-if="errors.conditionId">{{ errors.conditionId }}</span>
    </div>

    <div class="form-group">
      <label class="label">描述</label>
      <textarea
        v-model="form.description"
        class="textarea"
        rows="5"
        placeholder="描述一下宝贝的细节..."
        :class="{ error: errors.description }"
      ></textarea>
      <span class="error-text" v-if="errors.description">{{ errors.description }}</span>
    </div>

    <div class="form-group">
      <label class="label">图片</label>
      <ProductImageUpload v-model="form.images" />
      <span class="error-text" v-if="errors.images">{{ errors.images }}</span>
    </div>

    <div class="form-actions">
      <button type="button" class="btn-secondary" @click="$emit('cancel')">取消</button>
      <button type="submit" class="btn-primary" :disabled="loading">
        {{ loading ? '提交中...' : mode === 'create' ? '发布' : '保存' }}
      </button>
    </div>
  </form>
</template>

<style scoped lang="scss">
.product-form {
  max-width: 800px;
  margin: 0 auto;
  padding: 24px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);

  .form-group {
    margin-bottom: 20px;

    .label {
      display: block;
      margin-bottom: 8px;
      font-weight: 500;
      color: var(--color-text-primary, #333);
    }

    .input,
    .select,
    .textarea {
      width: 100%;
      padding: 8px 12px;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 14px;
      transition: border-color 0.2s;

      &:focus {
        border-color: var(--color-primary, #0066ff);
        outline: none;
      }

      &.error {
        border-color: var(--color-error, #ff4d4f);
      }
    }

    .error-text {
      display: block;
      margin-top: 4px;
      font-size: 12px;
      color: var(--color-error, #ff4d4f);
    }

    .tags-container {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      padding: 4px 0;

      &.disabled {
        opacity: 0.6;
        pointer-events: none;
      }

      .empty-text {
        font-size: 14px;
        color: #999;
        padding: 8px 0;
      }

      .tag-item {
        padding: 6px 16px;
        background-color: #f5f5f5;
        border: 1px solid #e0e0e0;
        border-radius: 20px;
        font-size: 14px;
        color: #666;
        cursor: pointer;
        transition: all 0.2s;
        user-select: none;

        &:hover {
          background-color: #e6f7ff;
          border-color: #91d5ff;
          color: #1890ff;
        }

        &.active {
          background-color: #e6f7ff;
          border-color: #1890ff;
          color: #1890ff;
          font-weight: 500;
        }
      }
    }
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 32px;
    padding-top: 24px;
    border-top: 1px solid #eee;

    button {
      padding: 8px 24px;
      border-radius: 4px;
      font-size: 14px;
      cursor: pointer;
      border: none;
      transition: opacity 0.2s;

      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }

      &.btn-primary {
        background: var(--color-primary, #0066ff);
        color: #fff;

        &:hover:not(:disabled) {
          opacity: 0.9;
        }
      }

      &.btn-secondary {
        background: #f5f5f5;
        color: #666;

        &:hover:not(:disabled) {
          background: #e0e0e0;
        }
      }
    }
  }
}
</style>
