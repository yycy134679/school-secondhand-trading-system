<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
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
      <select
        v-model="form.categoryId"
        class="select"
        :class="{ error: errors.categoryId }"
      >
        <option :value="undefined" disabled>请选择分类</option>
        <option
          v-for="cat in appStore.categories"
          :key="cat.id"
          :value="cat.id"
        >
          {{ cat.name }}
        </option>
      </select>
      <span class="error-text" v-if="errors.categoryId">{{ errors.categoryId }}</span>
    </div>

    <div class="form-group">
      <label class="label">新旧程度</label>
      <div class="radio-group">
        <label
          v-for="cond in appStore.productConditions"
          :key="cond.id"
          class="radio-label"
        >
          <input
            type="radio"
            :value="cond.id"
            v-model="form.conditionId"
          />
          {{ cond.name }}
        </label>
      </div>
      <span class="error-text" v-if="errors.conditionId">{{ errors.conditionId }}</span>
    </div>

    <div class="form-group">
      <label class="label">标签</label>
      <div class="checkbox-group">
        <label v-for="tag in appStore.tags" :key="tag.id" class="checkbox-label">
          <input
            type="checkbox"
            :value="tag.id"
            v-model="form.tagIds"
          />
          {{ tag.name }}
        </label>
      </div>
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
        {{ loading ? '提交中...' : (mode === 'create' ? '发布' : '保存') }}
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

    .radio-group,
    .checkbox-group {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;

      .radio-label,
      .checkbox-label {
        display: flex;
        align-items: center;
        gap: 4px;
        cursor: pointer;
        font-size: 14px;
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
