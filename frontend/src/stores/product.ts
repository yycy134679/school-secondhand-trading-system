import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getHomeData } from '@/api/home'
import {
  searchProducts as searchProductsApi,
  getProductsByCategory as getProductsByCategoryApi,
  getMyProducts as getMyProductsApi,
  getProductDetail as getProductDetailApi,
  changeProductStatus as changeProductStatusApi,
} from '@/api/product'
import type { Product } from '@common/types/product'
import type { PageResult } from '@common/types/api'
import type { ProductDetail, ProductSearchParams } from '@/api/product'

export const useProductStore = defineStore('product', () => {
  // Home Data
  const homeRecommendations = ref<Product[]>([])
  const homeLatest = ref<PageResult<Product>>({ items: [], page: 1, pageSize: 20, total: 0 })

  // Search & List Cache
  const searchParams = ref<ProductSearchParams>({})
  const searchResults = ref<PageResult<Product>>({ items: [], page: 1, pageSize: 20, total: 0 })

  // Current Product Detail
  const currentProduct = ref<ProductDetail | null>(null)

  async function fetchHomeData(params?: { page?: number; pageSize?: number }) {
    try {
      const res = await getHomeData(params)
      homeRecommendations.value = res.data.data.recommendations
      homeLatest.value = res.data.data.latest
    } catch (error) {
      console.error('Failed to fetch home data:', error)
      throw error
    }
  }

  async function searchProducts(params: ProductSearchParams) {
    searchParams.value = params
    try {
      const res = await searchProductsApi(params)
      searchResults.value = res.data.data
      return res.data.data
    } catch (error) {
      console.error('Failed to search products:', error)
      throw error
    }
  }

  async function fetchCategoryProducts(categoryId: number, params: ProductSearchParams) {
    try {
      const res = await getProductsByCategoryApi(categoryId, params)
      return res.data.data
    } catch (error) {
      console.error('Failed to fetch category products:', error)
      throw error
    }
  }

  async function fetchMyProducts(params: { keyword?: string; page?: number; pageSize?: number }) {
    try {
      const res = await getMyProductsApi(params)
      return res.data.data
    } catch (error) {
      console.error('Failed to fetch my products:', error)
      throw error
    }
  }

  async function fetchProductDetail(id: number) {
    try {
      const res = await getProductDetailApi(id)
      currentProduct.value = res.data.data
      return res.data.data
    } catch (error) {
      console.error('Failed to fetch product detail:', error)
      throw error
    }
  }

  async function changeStatus(id: number, action: 'delist' | 'relist' | 'sold') {
    try {
      const res = await changeProductStatusApi(id, { action })
      // Update local state if current product matches
      if (currentProduct.value && currentProduct.value.id === id) {
        currentProduct.value.status = res.data.data.status
      }
      return res.data.data
    } catch (error) {
      console.error('Failed to change product status:', error)
      throw error
    }
  }

  return {
    homeRecommendations,
    homeLatest,
    searchParams,
    searchResults,
    currentProduct,
    fetchHomeData,
    searchProducts,
    fetchCategoryProducts,
    fetchMyProducts,
    fetchProductDetail,
    changeStatus,
  }
})
