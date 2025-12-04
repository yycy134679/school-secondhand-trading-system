import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

type LoginResult = { success: true } | { success: false; message: string }
type RegisterResult = { success: true } | { success: false; message: string }
type PasswordResult = { success: true; message?: string } | { success: false; message: string }

type StoredState = {
  account: string | null
  nickname: string | null
  isAuthenticated: boolean
  role: AuthRole | null
  accounts?: DemoAccount[]
}

type DemoAccount = {
  account: string
  password: string
  nickname: string
  role: AuthRole
}

export type AuthRole = 'admin' | 'viewer'

const STORAGE_KEY = 'admin-auth-state'

const defaultAccounts: DemoAccount[] = [
  { account: 'admin', password: 'admin123', nickname: '超级管理员', role: 'admin' },
  { account: 'manager', password: 'manager123', nickname: '运营经理', role: 'admin' },
]

export const useAuthStore = defineStore('adminAuth', () => {
  const isAuthenticated = ref(false)
  const currentAccount = ref<{ account: string; nickname: string; role: AuthRole } | null>(null)
  const accounts = ref<DemoAccount[]>([...defaultAccounts])

  const persist = () => {
    if (typeof window === 'undefined') return
    const payload: StoredState = {
      account: currentAccount.value?.account ?? null,
      nickname: currentAccount.value?.nickname ?? null,
      isAuthenticated: isAuthenticated.value,
      role: currentAccount.value?.role ?? null,
      accounts: accounts.value,
    }
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(payload))
  }

  if (typeof window !== 'undefined') {
    try {
      const raw = window.localStorage.getItem(STORAGE_KEY)
      if (raw) {
        const parsed = JSON.parse(raw) as StoredState
        if (Array.isArray(parsed.accounts) && parsed.accounts.length) {
          accounts.value = parsed.accounts
        }
        if (parsed.isAuthenticated && parsed.account) {
          isAuthenticated.value = true
          currentAccount.value = {
            account: parsed.account,
            nickname: parsed.nickname ?? parsed.account,
            role: parsed.role ?? (
              defaultAccounts.find((item) => item.account === parsed.account)?.role || 'viewer'
            ),
          }
        }
      }
    } catch (error) {
      console.warn('Failed to restore auth state:', error)
      window.localStorage.removeItem(STORAGE_KEY)
    }
  }

  const ensureDefaultAccounts = () => {
    let mutated = false
    for (const item of defaultAccounts) {
      if (!accounts.value.some((existing) => existing.account === item.account)) {
        accounts.value.push(item)
        mutated = true
      }
    }
    if (mutated) {
      persist()
    }
  }

  const login = (accountInput: string, passwordInput: string): LoginResult => {
    const account = accountInput.trim()
    const password = passwordInput.trim()

    if (!account || !password) {
      return { success: false, message: '请输入账号和密码。' }
    }

    ensureDefaultAccounts()
    const matched = accounts.value.find(
      (item) => item.account === account && item.password === password,
    )

    if (!matched) {
      return { success: false, message: '账号或密码错误。' }
    }

    isAuthenticated.value = true
    currentAccount.value = {
      account: matched.account,
      nickname: matched.nickname,
      role: matched.role,
    }
    persist()
    return { success: true }
  }

  const logout = () => {
    isAuthenticated.value = false
    currentAccount.value = null
    persist()
  }

  const register = (accountInput: string, passwordInput: string, nicknameInput?: string): RegisterResult => {
    const account = accountInput.trim()
    const password = passwordInput.trim()
    const nickname = (nicknameInput || account).trim()

    if (!account || !password) {
      return { success: false, message: '账号和密码不能为空。' }
    }

    ensureDefaultAccounts()

    if (accounts.value.some((item) => item.account === account)) {
      return { success: false, message: '账号已存在，请直接登录。' }
    }

    accounts.value = [...accounts.value, { account, password, nickname, role: 'viewer' }]
    currentAccount.value = {
      account,
      nickname,
      role: 'viewer',
    }
    isAuthenticated.value = true
    persist()
    return { success: true }
  }

  const addAdminAccount = (accountInput: string, nicknameInput?: string, passwordInput = '123456') => {
    const account = accountInput.trim()
    const nickname = (nicknameInput || account).trim()
    const password = passwordInput.trim() || '123456'
    if (!account) return

    ensureDefaultAccounts()

    const index = accounts.value.findIndex((item) => item.account === account)
    const adminEntry: DemoAccount = {
      account,
      nickname,
      password,
      role: 'admin',
    }

    if (index >= 0) {
      accounts.value[index] = { ...accounts.value[index], ...adminEntry }
    } else {
      accounts.value = [...accounts.value, adminEntry]
    }

    if (currentAccount.value?.account === account) {
      currentAccount.value = {
        account,
        nickname,
        role: 'admin',
      }
      isAuthenticated.value = true
    }

    persist()
  }

  const updateAccountPassword = (
    account: string,
    newPasswordInput: string,
  ): { success: boolean; updatedNickname?: string } => {
    const newPassword = newPasswordInput.trim()
    if (!account || !newPassword) {
      return { success: false }
    }

    const index = accounts.value.findIndex((item) => item.account === account)
    if (index === -1) {
      return { success: false }
    }

    const entry = accounts.value[index]!
    const updated: DemoAccount = {
      account: entry.account,
      nickname: entry.nickname,
      role: entry.role,
      password: newPassword,
    }
    accounts.value.splice(index, 1, updated)

    if (currentAccount.value?.account === account) {
      currentAccount.value = {
        account: updated.account,
        nickname: updated.nickname,
        role: updated.role,
      }
      isAuthenticated.value = true
    }

    persist()
    return { success: true, updatedNickname: updated.nickname }
  }

  const changeOwnPassword = (oldPasswordInput: string, newPasswordInput: string): PasswordResult => {
    if (!currentAccount.value) {
      return { success: false, message: '请先登录后再修改密码。' }
    }

    const oldPassword = oldPasswordInput.trim()
    const newPassword = newPasswordInput.trim()

    if (!oldPassword || !newPassword) {
      return { success: false, message: '原密码和新密码均不能为空。' }
    }

    if (newPassword.length < 6) {
      return { success: false, message: '新密码长度不能少于 6 位。' }
    }

    ensureDefaultAccounts()

    const entry = accounts.value.find((item) => item.account === currentAccount.value?.account)
    if (!entry) {
      return { success: false, message: '未找到当前账号，请重新登录。' }
    }

    if (entry.password !== oldPassword) {
      return { success: false, message: '原密码不正确。' }
    }

    const result = updateAccountPassword(entry.account, newPassword)
    if (!result.success) {
      return { success: false, message: '修改密码失败，请稍后重试。' }
    }

    return { success: true, message: '密码修改成功。' }
  }

  const resetPassword = (accountInput: string, newPasswordInput: string): PasswordResult => {
    const account = accountInput.trim()
    const newPassword = newPasswordInput.trim()

    if (!account || !newPassword) {
      return { success: false, message: '账号与新密码均不能为空。' }
    }

    if (newPassword.length < 6) {
      return { success: false, message: '新密码长度不能少于 6 位。' }
    }

    ensureDefaultAccounts()

    const existed = accounts.value.some((item) => item.account === account)
    if (!existed) {
      return { success: false, message: '未找到该账号。' }
    }

    const result = updateAccountPassword(account, newPassword)
    if (!result.success) {
      return { success: false, message: '设置密码失败，请稍后重试。' }
    }

    return { success: true, message: `已为账号 ${account} 设置新密码。` }
  }

  const updateAccountNickname = (accountInput: string, nicknameInput: string): boolean => {
    const account = accountInput.trim()
    const nickname = nicknameInput.trim()
    if (!account || !nickname) {
      return false
    }

    ensureDefaultAccounts()

    const index = accounts.value.findIndex((item) => item.account === account)
    if (index === -1) {
      return false
    }

    const updated: DemoAccount = {
      ...accounts.value[index]!,
      nickname,
    }

    accounts.value.splice(index, 1, updated)

    if (currentAccount.value?.account === account) {
      currentAccount.value = {
        account: updated.account,
        nickname: updated.nickname,
        role: updated.role,
      }
      isAuthenticated.value = true
    }

    persist()
    return true
  }

  return {
    isAuthenticated,
    currentAccount,
    login,
    register,
    logout,
    addAdminAccount,
    updateAccountNickname,
    changeOwnPassword,
    resetPassword,
    canManage: computed(() => currentAccount.value?.role === 'admin'),
    isViewer: computed(() => currentAccount.value?.role !== 'admin'),
  }
})
