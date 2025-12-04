<template>
  <div class="profile-page">
    <div class="profile-card">
      <header class="profile-header">
        <h1>个人信息</h1>
        <p>查看当前登录管理员的账户信息</p>
      </header>

      <section class="info-section">
        <div class="info-row">
          <span class="info-label">账号</span>
          <span class="info-value">{{ authStore.currentAccount?.account || '未知' }}</span>
        </div>
        <div class="info-row">
          <span class="info-label">昵称</span>
          <span class="info-value">{{ authStore.currentAccount?.nickname || '未设置' }}</span>
        </div>
        <div class="info-row">
          <span class="info-label">登录状态</span>
          <span class="info-value status" :class="{ online: authStore.isAuthenticated }">
            {{ authStore.isAuthenticated ? '已登录' : '未登录' }}
          </span>
        </div>
      </section>

      <section class="password-section">
        <h2>修改密码</h2>
        <p class="section-tip">请先输入当前密码，再设置新的登录密码。</p>
        <form class="password-form" @submit.prevent="submitPasswordChange">
          <label class="form-field">
            <span class="field-label">当前密码</span>
            <input
              v-model="passwordForm.oldPassword"
              type="password"
              class="field-input"
              placeholder="请输入当前密码"
              autocomplete="current-password"
            />
          </label>
          <label class="form-field">
            <span class="field-label">新密码</span>
            <input
              v-model="passwordForm.newPassword"
              type="password"
              class="field-input"
              placeholder="请输入至少 6 位新密码"
              autocomplete="new-password"
            />
          </label>
          <label class="form-field">
            <span class="field-label">确认新密码</span>
            <input
              v-model="passwordForm.confirmPassword"
              type="password"
              class="field-input"
              placeholder="请再次输入新密码"
              autocomplete="new-password"
            />
          </label>
          <p v-if="passwordForm.error" class="form-error">{{ passwordForm.error }}</p>
          <p v-if="passwordForm.success" class="form-success">{{ passwordForm.success }}</p>
          <button type="submit" class="btn submit" :disabled="passwordForm.submitting">
            {{ passwordForm.submitting ? '保存中…' : '保存新密码' }}
          </button>
        </form>
      </section>

      <section class="actions">
        <button type="button" class="btn primary" @click="backToDashboard">
          返回仪表盘
        </button>
        <button type="button" class="btn danger" @click="handleLogout">
          退出登录
        </button>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
  error: '',
  success: '',
  submitting: false,
})

const backToDashboard = () => {
  router.push('/admin')
}

const handleLogout = () => {
  authStore.logout()
  router.replace({
    path: '/login',
    query: {
      redirect: '/admin/profile',
    },
  })
}

const submitPasswordChange = () => {
  if (passwordForm.submitting) return
  passwordForm.error = ''
  passwordForm.success = ''

  if (!passwordForm.oldPassword.trim() || !passwordForm.newPassword.trim()) {
    passwordForm.error = '请输入当前密码和新的密码。'
    return
  }

  if (passwordForm.newPassword.trim().length < 6) {
    passwordForm.error = '新密码长度不能少于 6 位。'
    return
  }

  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    passwordForm.error = '两次输入的密码不一致。'
    return
  }

  passwordForm.submitting = true
  try {
    const result = authStore.changeOwnPassword(
      passwordForm.oldPassword,
      passwordForm.newPassword,
    )
    if (!result.success) {
      passwordForm.error = result.message
      return
    }

    passwordForm.success = result.message ?? '密码修改成功。'
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } finally {
    passwordForm.submitting = false
  }
}
</script>

<style scoped>
.profile-page {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 32px;
}

.profile-card {
  width: 100%;
  max-width: 560px;
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 12px 32px rgba(15, 23, 42, 0.12);
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.profile-header h1 {
  margin: 0 0 8px;
  font-size: 24px;
  color: #111827;
}

.profile-header p {
  margin: 0;
  color: #6b7280;
  font-size: 14px;
}

.info-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 12px 16px;
  border-radius: 12px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
}

.info-label {
  font-size: 14px;
  color: #6b7280;
}

.info-value {
  font-size: 16px;
  color: #1f2937;
}

.info-value.status {
  font-weight: 600;
}

.info-value.status.online {
  color: #059669;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.password-section {
  border-top: 1px solid #f3f4f6;
  padding-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.password-section h2 {
  margin: 0;
  font-size: 18px;
  color: #111827;
}

.section-tip {
  margin: 0;
  font-size: 13px;
  color: #6b7280;
}

.password-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
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
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
  outline: none;
}

.form-error {
  margin: -4px 0 0;
  color: #dc2626;
  font-size: 13px;
}

.form-success {
  margin: -4px 0 0;
  color: #059669;
  font-size: 13px;
}

.btn.submit {
  align-self: flex-start;
  border: none;
  background: #2563eb;
  color: #fff;
  padding: 10px 24px;
  border-radius: 999px;
}

.btn.submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn {
  border: 1px solid #d1d5db;
  background: #fff;
  color: #1f2937;
  padding: 10px 20px;
  border-radius: 999px;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 18px rgba(37, 99, 235, 0.15);
}

.btn.primary {
  border: none;
  background: #2563eb;
  color: #fff;
}

.btn.primary:hover {
  background: #1d4ed8;
}

.btn.danger {
  border: none;
  background: #dc2626;
  color: #fff;
}

.btn.danger:hover {
  background: #b91c1c;
}

@media (max-width: 640px) {
  .profile-page {
    padding: 16px;
  }

  .profile-card {
    padding: 24px;
  }

  .actions {
    flex-direction: column;
  }

  .btn {
    width: 100%;
  }
}
</style>
