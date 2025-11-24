export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PageResult<T> {
  items: T[]
  page: number
  pageSize: number
  total: number
}
