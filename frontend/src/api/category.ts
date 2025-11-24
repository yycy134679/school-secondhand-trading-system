import request from '@/utils/request'
import type { ApiResponse } from '@common/types/api'
import type { Category } from '@common/types/category'

export function getCategories() {
  return request.get<ApiResponse<Category[]>>('/categories')
}
