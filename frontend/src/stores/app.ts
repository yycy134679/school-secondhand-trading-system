import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCategories } from '@/api/category'
import { getTags } from '@/api/tag'
import { getProductConditions } from '@/api/product'
import type { Category } from '@common/types/category'
import type { Tag } from '@common/types/tag'
import type { ProductCondition } from '@common/types/product_condition'

export const useAppStore = defineStore('app', () => {
  const loading = ref(false)
  const error = ref<string | null>(null)
  const categories = ref<Category[]>([])
  const tags = ref<Tag[]>([])
  const productConditions = ref<ProductCondition[]>([])

  async function initDictionaries() {
    loading.value = true
    error.value = null
    try {
      const [categoriesRes, tagsRes, conditionsRes] = await Promise.allSettled([
        getCategories(),
        getTags(),
        getProductConditions(),
      ])
      const errors: string[] = []

      if (categoriesRes.status === 'fulfilled') {
        categories.value = categoriesRes.value.data.data
      } else {
        errors.push('分类加载失败')
        console.error('Failed to load categories:', categoriesRes.reason)
      }

      if (tagsRes.status === 'fulfilled') {
        tags.value = tagsRes.value.data.data
      } else {
        errors.push('标签加载失败')
        console.error('Failed to load tags:', tagsRes.reason)
      }

      if (conditionsRes.status === 'fulfilled') {
        productConditions.value = conditionsRes.value.data.data
      } else {
        errors.push('新旧程度加载失败')
        console.error('Failed to load product conditions:', conditionsRes.reason)
      }

      if (errors.length > 0) {
        error.value = errors.join('；')
      }
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to load dictionaries'
      error.value = message
      console.error('Failed to init dictionaries:', err)
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    categories,
    tags,
    productConditions,
    initDictionaries,
  }
})
