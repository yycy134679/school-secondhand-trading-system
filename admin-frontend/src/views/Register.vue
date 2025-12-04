<template>
  <div class="register-page">
    <div class="register-card">
      <header class="register-header">
        <h1>管理员注册</h1>
        <p>创建新的后台管理员账号</p>
      </header>
      <form class="register-form" @submit.prevent="handleSubmit">
        <label class="form-field">
          <span class="field-label">账号</span>
          <input
            v-model="account"
            type="text"
            class="field-input"
            placeholder="请输入账号"
            autocomplete="username"
          />
        </label>
        <label class="form-field">
          <span class="field-label">昵称</span>
          <input
            v-model="nickname"
            type="text"
            class="field-input"
            placeholder="请输入昵称（可选）"
            autocomplete="nickname"
          />
        </label>
        <label class="form-field">
          <span class="field-label">密码</span>
          <input
            v-model="password"
            type="password"
            class="field-input"
            placeholder="请输入密码"
            autocomplete="new-password"
          />
        </label>
        <label class="form-field">
          <span class="field-label">确认密码</span>
          <input
            v-model="confirmPassword"
            type="password"
            class="field-input"
            placeholder="请再次输入密码"
            autocomplete="new-password"
          />
        </label>
        <p v-if="error" class="form-error">{{ error }}</p>
        <button type="submit" class="submit-btn" :disabled="submitting">
          {{ submitting ? '注册中…' : '注册并登录' }}
        </button>
      </form>
      <footer class="register-footer">
        <span>已有账号？</span>
        <button type="button" class="link-btn" @click="goToLogin">返回登录</button>
      </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useAuthStore } from '../stores/auth'
import { useAdminDataStore } from '../stores/admin'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const adminStore = useAdminDataStore()

const account = ref('')
const nickname = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const submitting = ref(false)

const handleSubmit = async () => {
  if (submitting.value) return
  error.value = ''

  if (password.value.trim().length < 6) {
    error.value = '密码至少需要 6 位字符。'
    return
  }

  if (password.value !== confirmPassword.value) {
    error.value = '两次输入的密码不一致。'
    return
  }

  submitting.value = true
  try {
    const registerResult = authStore.register(account.value, password.value, nickname.value)
    if (!registerResult.success) {
      error.value = registerResult.message
      return
    }

    if (!authStore.isAuthenticated) {
      const loginResult = authStore.login(account.value, password.value)
      if (!loginResult.success) {
        error.value = loginResult.message
        return
      }
    }

    const normalizedAccount = account.value.trim()
    const normalizedNickname = nickname.value.trim() || normalizedAccount
    adminStore.addViewerAccount({ account: normalizedAccount, nickname: normalizedNickname })

    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/admin'
    await router.replace(redirect)
  } finally {
    submitting.value = false
  }
}

const goToLogin = () => {
  router.push({
    path: '/login',
    query: route.query,
  })
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #7c3aed 0%, #60a5fa 100%);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 16px;
}

.register-card {
  width: 100%;
  max-width: 460px;
  background: #ffffff;
  border-radius: 20px;
  box-shadow: 0 24px 60px rgba(15, 23, 42, 0.18);
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.register-header h1 {
  margin: 0 0 8px;
  font-size: 26px;
  color: #111827;
}

.register-header p {
  margin: 0;
  color: #6b7280;
  font-size: 14px;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-size: 14px;
  color: #4b5563;
}

.field-input {
  border: 1px solid #d1d5db;
  border-radius: 10px;
  padding: 12px;
  font-size: 15px;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.field-input:focus {
  outline: none;
  border-color: #7c3aed;
  box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.12);
}

.form-error {
  margin: -6px 0 0;
  color: #dc2626;
  font-size: 13px;
}

.submit-btn {
  margin-top: 4px;
  background: #7c3aed;
  color: #fff;
  border: none;
  border-radius: 999px;
  padding: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s ease, transform 0.2s ease;
}

.submit-btn:hover:not(:disabled) {
  background: #6d28d9;
  transform: translateY(-1px);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.register-footer {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6b7280;
}

.link-btn {
  border: none;
  background: none;
  color: #7c3aed;
  cursor: pointer;
  font-weight: 600;
  padding: 0;
}

.link-btn:hover {
  text-decoration: underline;
}
</style>
