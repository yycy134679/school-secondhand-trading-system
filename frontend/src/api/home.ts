import request from '@/utils/request'
import type { ApiResponse, PageResult } from '@common/types/api'
import type { Product } from '@common/types/product'

export interface HomeData {
  recommendations: Product[]
  latest: PageResult<Product>
}

export function getHomeData(params?: { page?: number; pageSize?: number }) {
  return request.get<ApiResponse<HomeData>>('/home', { params })
}
