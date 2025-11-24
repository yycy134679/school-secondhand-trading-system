export interface User {
  id: number
  account: string
  nickname: string
  avatarUrl?: string
  wechatId?: string
  isAdmin: boolean
  lastNicknameChangedAt?: string
  createdAt: string
  updatedAt: string
}

export interface LoginResponse {
  token: string
  user: User
}
