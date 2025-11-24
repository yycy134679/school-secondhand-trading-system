import axios, { type InternalAxiosRequestConfig, type AxiosResponse } from 'axios'
import { getToken, removeToken } from './auth'
import { ErrorCode } from '@common/constants/error_code'
import type { ApiResponse } from '@common/types/api'

// 创建 axios 实例
const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 在发送请求之前做些什么
    const token = getToken()
    if (token) {
      // 让每个请求携带 token
      // ['Authorization'] 是自定义头部 key
      // 请根据实际情况修改
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: unknown) => {
    // 处理请求错误
    console.error('Request Error:', error)
    return Promise.reject(error)
  },
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const res = response.data

    // 如果 code 不为 0，则判断为错误
    if (res.code !== ErrorCode.SUCCESS) {
      // 1002: 未登录或 token 无效
      if (res.code === ErrorCode.UNAUTHORIZED) {
        // 移除 token
        removeToken()
        // 触发登录弹窗 (这里可以通过事件总线或 store 触发，暂时先打印日志或重定向)
        // 由于 store 可能还未初始化完成，这里简单处理，实际项目中可能需要配合 Router 或 Store
        // 例如: window.location.href = '/login' (如果不是弹窗模式)
        // 或者 dispatch 一个全局事件
        window.dispatchEvent(new CustomEvent('auth:unauthorized'))
      } else if (res.code === ErrorCode.FORBIDDEN) {
        // 1003: 权限不足
        // 显示 Toast 提示 (这里假设有一个全局 Toast 工具，或者暂时用 console.warn)
        console.warn('权限不足:', res.message)
        // 也可以 dispatch 事件让 UI 层处理
        window.dispatchEvent(new CustomEvent('auth:forbidden', { detail: res.message }))
      } else {
        // 其他业务错误
        console.error('API Error:', res.message)
      }
      return Promise.reject(new Error(res.message || 'Error'))
    } else {
      return response
    }
  },
  (error: unknown) => {
    console.error('Response Error:', error)
    return Promise.reject(error)
  },
)

export default service
