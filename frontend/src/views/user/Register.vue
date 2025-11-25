<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const form = reactive({
  account: '',
  nickname: '',
  password: '',
  confirmPassword: '',
  wechatId: '',
})

const loading = ref(false)
const errorMessage = ref('')

const handleSubmit = async () => {
  if (!form.account || !form.nickname || !form.password || !form.confirmPassword) {
    errorMessage.value = '请填写所有必填项'
    return
  }

  if (form.password !== form.confirmPassword) {
    errorMessage.value = '两次输入的密码不一致'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    await userStore.register({
      account: form.account,
      nickname: form.nickname,
      password: form.password,
      confirmPassword: form.confirmPassword,
      wechatId: form.wechatId || undefined,
    })

    // Jump to home page
    router.push('/')
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '注册失败'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1 class="auth-title">注册</h1>

      <form @submit.prevent="handleSubmit" class="auth-form">
        <div class="form-item">
          <label>账号 <span class="required">*</span></label>
          <input v-model="form.account" type="text" placeholder="请输入账号" :disabled="loading" />
        </div>

        <div class="form-item">
          <label>昵称 <span class="required">*</span></label>
          <input v-model="form.nickname" type="text" placeholder="请输入昵称" :disabled="loading" />
        </div>

        <div class="form-item">
          <label>密码 <span class="required">*</span></label>
          <input
            v-model="form.password"
            type="password"
            placeholder="请输入密码 (至少8位)"
            :disabled="loading"
          />
        </div>

        <div class="form-item">
          <label>确认密码 <span class="required">*</span></label>
          <input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            :disabled="loading"
          />
        </div>

        <div class="form-item">
          <label>微信号 (可选)</label>
          <input
            v-model="form.wechatId"
            type="text"
            placeholder="方便买家联系您"
            :disabled="loading"
          />
        </div>

        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>

        <button type="submit" class="btn-submit" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>

        <div class="form-footer">已有账号? <router-link to="/login">立即登录</router-link></div>
      </form>
    </div>
  </div>
</template>

<style scoped lang="scss">
.auth-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 64px - 200px);
  padding: 40px 20px;
  background-color: #f7f8fa;
}

.auth-card {
  background: white;
  padding: 40px;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.05);
  width: 100%;
  max-width: 400px;
}

.auth-title {
  font-size: 24px;
  font-weight: bold;
  color: #1a1a1a;
  margin-bottom: 32px;
  text-align: center;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 8px;

  label {
    font-size: 14px;
    color: #333;
    font-weight: 500;

    .required {
      color: #ef4444;
    }
  }

  input {
    padding: 12px 16px;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 14px;
    transition: all 0.2s;
    background-color: #f9fafb;

    &:focus {
      border-color: var(--color-primary, #0066ff);
      background-color: white;
      outline: none;
      box-shadow: 0 0 0 3px rgba(0, 102, 255, 0.1);
    }

    &:disabled {
      background-color: #f5f5f5;
      cursor: not-allowed;
    }
  }
}

.error-message {
  color: #ef4444;
  font-size: 14px;
  background: #fef2f2;
  padding: 10px;
  border-radius: 8px;
  border: 1px solid #fee2e2;
  text-align: center;
}

.btn-submit {
  background-color: var(--color-primary, #0066ff);
  color: white;
  border: none;
  padding: 14px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-top: 8px;

  &:hover:not(:disabled) {
    background-color: #0052cc;
  }

  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
}

.form-footer {
  text-align: center;
  font-size: 14px;
  color: #666;
  margin-top: 16px;

  a {
    color: var(--color-primary, #0066ff);
    text-decoration: none;
    font-weight: 600;

    &:hover {
      text-decoration: underline;
    }
  }
}
</style>
