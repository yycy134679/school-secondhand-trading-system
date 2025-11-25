import request from '@/utils/request'
import type { ApiResponse } from '@common/types/api'

export interface UploadResponse {
  url: string
}

export function uploadFile(file: File) {
  const formData = new FormData()
  formData.append('file', file)

  return request.post<ApiResponse<UploadResponse>>('/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}
