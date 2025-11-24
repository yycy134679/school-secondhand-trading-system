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
    try {
      const [categoriesRes, tagsRes, conditionsRes] = await Promise.all([
        getCategories(),
        getTags(),
        getProductConditions(),
      ])
      categories.value = categoriesRes.data.data
      tags.value = tagsRes.data.data
      productConditions.value = conditionsRes.data.data
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
