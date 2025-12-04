<template>
  <div class="login-page">
    <div class="login-card">
      <header class="login-header">
        <h1>后台登录</h1>
        <p>使用演示账号登录后台管理系统</p>
      </header>
      <form class="login-form" @submit.prevent="handleSubmit">
        <label class="form-field">
          <span class="field-label">账号</span>
          <input
            v-model="account"
            type="text"
            class="field-input"
            placeholder="请输入账号，如 admin"
            autocomplete="username"
          />
        </label>
        <label class="form-field">
          <span class="field-label">密码</span>
          <input
            v-model="password"
            type="password"
            class="field-input"
            placeholder="请输入密码，如 admin123"
            autocomplete="current-password"
          />
        </label>
        <p v-if="error" class="form-error">{{ error }}</p>
        <button type="submit" class="submit-btn" :disabled="submitting">
          {{ submitting ? '登录中…' : '登录' }}
        </button>
      </form>
      <section class="demo-info">
        <h2>演示账号</h2>
        <ul>
          <li>admin / admin123</li>
          <li>manager / manager123</li>
        </ul>
      </section>
      <footer class="login-footer">
        <span>还没有账号？</span>
        <button type="button" class="link-btn" @click="goToRegister">前往注册</button>
      </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const account = ref('')
const password = ref('')
const error = ref('')
const submitting = ref(false)

const handleSubmit = async () => {
  if (submitting.value) return
  error.value = ''
  submitting.value = true
  try {
    const result = authStore.login(account.value, password.value)
    if (!result.success) {
      error.value = result.message
      return
    }
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/admin'
    await router.replace(redirect)
  } finally {
    submitting.value = false
  }
}

const goToRegister = () => {
  router.push({
    path: '/register',
    query: route.query,
  })
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #2563eb 0%, #60a5fa 100%);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 16px;
}

.login-card {
  width: 100%;
  max-width: 420px;
  background: #ffffff;
  border-radius: 20px;
  box-shadow: 0 24px 60px rgba(15, 23, 42, 0.18);
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.login-header h1 {
  margin: 0 0 8px;
  font-size: 26px;
  color: #111827;
}

.login-header p {
  margin: 0;
  color: #6b7280;
  font-size: 14px;
}

.login-form {
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
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
}

.form-error {
  margin: -6px 0 0;
  color: #dc2626;
  font-size: 13px;
}

.submit-btn {
  margin-top: 4px;
  background: #2563eb;
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
  background: #1d4ed8;
  transform: translateY(-1px);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.demo-info {
  border-top: 1px solid #e5e7eb;
  padding-top: 16px;
}

.demo-info h2 {
  margin: 0 0 8px;
  font-size: 14px;
  color: #4b5563;
}

.demo-info ul {
  margin: 0;
  padding-left: 18px;
  color: #6b7280;
  font-size: 13px;
}

.login-footer {
  display: flex;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
  color: #6b7280;
}

.link-btn {
  border: none;
  background: none;
  color: #2563eb;
  cursor: pointer;
  font-weight: 600;
  padding: 0;
}

.link-btn:hover {
  text-decoration: underline;
}
</style>
