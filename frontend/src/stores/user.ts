import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  login as loginApi,
  register as registerApi,
  getProfile as getProfileApi,
  updateProfile as updateProfileApi,
  changePassword as changePasswordApi,
} from '@/api/user'
import { setToken, getToken, removeToken } from '@/utils/auth'
import type {
  LoginParams,
  RegisterParams,
  UpdateProfileParams,
  ChangePasswordParams,
} from '@/api/user'
import type { User } from '@common/types/user'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(getToken())
  const currentUser = ref<User | null>(null)
  const rememberMe = ref(false)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => currentUser.value?.isAdmin || false)

  async function login(loginForm: LoginParams) {
    try {
      const res = await loginApi(loginForm)
      const data = res.data.data
      token.value = data.token
      currentUser.value = data.user
      setToken(data.token)
      rememberMe.value = loginForm.rememberMe || false
      return data
    } catch (error) {
      throw error
    }
  }

  async function logout() {
    token.value = null
    currentUser.value = null
    removeToken()
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const res = await getProfileApi()
      currentUser.value = res.data.data
    } catch {
      // Token might be invalid
      logout()
    }
  }

  async function updateProfile(data: UpdateProfileParams) {
    try {
      await updateProfileApi(data)
      // Refresh profile after update
      await fetchProfile()
    } catch (error) {
      throw error
    }
  }

  async function changePassword(data: ChangePasswordParams) {
    try {
      await changePasswordApi(data)
    } catch (error) {
      throw error
    }
  }

  async function register(registerForm: RegisterParams) {
    try {
      const res = await registerApi(registerForm)
      const data = res.data.data
      token.value = data.token
      currentUser.value = data.user
      setToken(data.token)
      return data
    } catch (error) {
      throw error
    }
  }

  return {
    token,
    currentUser,
    rememberMe,
    isLoggedIn,
    isAdmin,
    login,
    register,
    logout,
    fetchProfile,
    updateProfile,
    changePassword,
  }
})
