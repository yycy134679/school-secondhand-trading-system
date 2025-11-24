import request from '@/utils/request'
import type { ApiResponse } from '@common/types/api'
import type { User, LoginResponse } from '@common/types/user'

// 注册参数
export interface RegisterParams {
  account: string
  nickname: string
  password: string
  confirmPassword: string
  wechatId?: string
}

// 登录参数
export interface LoginParams {
  account: string
  password: string
  rememberMe?: boolean
}

// 更新资料参数
export interface UpdateProfileParams {
  nickname?: string
  avatarUrl?: string
  wechatId?: string
}

// 修改密码参数
export interface ChangePasswordParams {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

export function register(data: RegisterParams) {
  return request.post<ApiResponse<LoginResponse>>('/users/register', data)
}

export function login(data: LoginParams) {
  return request.post<ApiResponse<LoginResponse>>('/users/login', data)
}

export function getProfile() {
  return request.get<ApiResponse<User>>('/users/profile')
}

export function updateProfile(data: UpdateProfileParams) {
  return request.put<ApiResponse<void>>('/users/profile', data)
}

export function changePassword(data: ChangePasswordParams) {
  return request.put<ApiResponse<void>>('/users/password', data)
}

import type { Product } from '@common/types/product'

export function getRecentViews() {
  // 假设返回的是商品列表，具体类型需确认，暂时用 any 或 Product[]
  // 根据 api.md 3.4，最近浏览返回商品列表
  // 假设 Product 类型已定义
  return request.get<ApiResponse<Product[]>>('/users/recent-views')
}
