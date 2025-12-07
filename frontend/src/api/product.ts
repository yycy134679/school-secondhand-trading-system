import request from '@/utils/request'
import type { ApiResponse, PageResult } from '@common/types/api'
import type { Product } from '@common/types/product'
import type { ProductCondition } from '@common/types/product_condition'

export interface ProductImage {
  id: number
  url: string
  sortOrder: number
  isPrimary: boolean
}

export interface ProductDetail extends Product {
  conditionName: string
  images: ProductImage[]
  tagIds: number[]
  seller: {
    id: number
    nickname: string
    avatarUrl: string
  }
  viewerIsSeller: boolean
  sellerWechat: string | null
}

// 发布商品参数 (FormData)
// title, description, price, categoryId, tagIds, conditionId, images
// 这里只定义接口，实际调用时传 FormData

// 编辑商品参数
export interface UpdateProductParams {
  title?: string
  description?: string
  price?: number
  categoryId?: number
  tagIds?: number[] // 前端使用数组，在请求前转换为逗号分隔字符串
  conditionId?: number
  imageUrls?: string[]
}

// 搜索参数
export interface ProductSearchParams {
  q?: string
  categoryId?: number
  tagIds?: string
  conditionIds?: string
  minPrice?: number
  maxPrice?: number
  publishedTimeRange?: string
  sort?: string
  page?: number
  pageSize?: number
}

// 状态变更参数
export interface ProductStatusParams {
  action: 'delist' | 'relist' | 'sold'
}

// 联系卖家响应
export interface ContactSellerResponse {
  canContact: boolean
  sellerWechat?: string
  tips?: string
}

export function createProduct(data: FormData) {
  return request.post<ApiResponse<Product>>('/products', data, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

export function updateProduct(id: number, data: UpdateProductParams) {
  // 转换 tagIds 数组为逗号分隔字符串
  const requestData = {
    ...data,
    tagIds:
      data.tagIds && Array.isArray(data.tagIds) ? (data.tagIds.join(',') as any) : data.tagIds,
  }
  return request.put<ApiResponse<Product>>(`/products/${id}`, requestData)
}

export function changeProductStatus(id: number, data: ProductStatusParams) {
  return request.post<ApiResponse<Product>>(`/products/${id}/status`, data)
}

export function undoProductStatusChange(id: number) {
  return request.post<ApiResponse<Product>>(`/products/${id}/status/undo`)
}

export function getProductDetail(id: number) {
  return request.get<ApiResponse<ProductDetail>>(`/products/${id}`)
}

export function recordProductView(id: number) {
  return request.post<ApiResponse<{ recorded: boolean }>>(`/products/${id}/view`)
}

export function getMyProducts(params: { keyword?: string; page?: number; pageSize?: number }) {
  return request.get<ApiResponse<PageResult<Product>>>('/products/my', { params })
}

export function searchProducts(params: ProductSearchParams) {
  return request.get<ApiResponse<PageResult<Product>>>('/products/search', { params })
}

export function getProductsByCategory(categoryId: number, params: ProductSearchParams) {
  return request.get<ApiResponse<PageResult<Product>>>(`/products/category/${categoryId}`, {
    params,
  })
}

export function getProductContact(id: number) {
  return request.get<ApiResponse<ContactSellerResponse>>(`/products/${id}/contact`)
}

export function getProductConditions() {
  return request.get<ApiResponse<ProductCondition[]>>('/product-conditions')
}
