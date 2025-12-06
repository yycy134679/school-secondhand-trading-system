import request from '@/utils/request'
import type { ApiResponse } from '@common/types/api'
import type { Product } from '@common/types/product'

// 对齐后端返回的首页商品卡片结构
export interface HomeProduct {
  id: number
  title: string
  price: number
  status: Product['status']
  mainImage?: string
  mainImageUrl?: string
  description?: string
  conditionId?: number
  sellerId?: number
  categoryId?: number
  createdAt?: string
  updatedAt?: string
}

export interface HomeData {
  recommendations: HomeProduct[]
  latest: HomeProduct[]
  totalCount: number
}

export function getHomeData(params?: { page?: number; pageSize?: number }) {
  return request.get<ApiResponse<HomeData>>('/home', { params })
}
