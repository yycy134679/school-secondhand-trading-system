import request from '@/utils/request'
import type { ApiResponse } from '@common/types/api'
import type { Tag } from '@common/types/tag'

export function getTags() {
  return request.get<ApiResponse<Tag[]>>('/tags')
}
