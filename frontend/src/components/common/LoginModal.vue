<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'

defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}>()

const userStore = useUserStore()
const router = useRouter()

const form = reactive({
  account: '',
  password: '',
  rememberMe: false,
})

const loading = ref(false)
const errorMsg = ref('')

const handleClose = () => {
  emit('update:visible', false)
  // Reset form
  form.account = ''
  form.password = ''
  errorMsg.value = ''
}

const handleLogin = async () => {
  if (!form.account || !form.password) {
    errorMsg.value = '请输入账号和密码'
    return
  }

  loading.value = true
  errorMsg.value = ''

  try {
    await userStore.login({
      account: form.account,
      password: form.password,
      rememberMe: form.rememberMe,
    })
    emit('success')
    handleClose()
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '登录失败，请检查账号密码'
    errorMsg.value = message
  } finally {
    loading.value = false
  }
}

const goToRegister = () => {
  handleClose()
  router.push('/register')
}
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click.self="handleClose">
    <div class="modal-content">
      <button class="close-btn" @click="handleClose">&times;</button>

      <div class="modal-header">
        <h2>欢迎回来</h2>
        <p>登录校园二手交易平台</p>
      </div>

      <div class="modal-body">
        <div v-if="errorMsg" class="error-alert">{{ errorMsg }}</div>

        <div class="form-group">
          <label>账号</label>
          <input
            v-model="form.account"
            type="text"
            placeholder="请输入账号"
            @keyup.enter="handleLogin"
          />
        </div>

        <div class="form-group">
          <label>密码</label>
          <input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            @keyup.enter="handleLogin"
          />
        </div>

        <div class="form-actions">
          <label class="checkbox-label">
            <input v-model="form.rememberMe" type="checkbox" />
            <span>记住我</span>
          </label>
          <a href="#" class="forgot-pwd">忘记密码?</a>
        </div>

        <button class="btn btn-primary btn-block" :disabled="loading" @click="handleLogin">
          {{ loading ? '登录中...' : '登录' }}
        </button>

        <div class="register-link">
          还没有账号? <a href="#" @click.prevent="goToRegister">立即注册</a>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  background: white;
  border-radius: 16px;
  width: 400px;
  padding: 32px;
  position: relative;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;

  &:hover {
    color: #333;
  }
}

.modal-header {
  text-align: center;
  margin-bottom: 24px;

  h2 {
    font-size: 24px;
    font-weight: bold;
    color: #1a1a1a;
    margin-bottom: 8px;
  }

  p {
    color: #666;
    font-size: 14px;
  }
}

.form-group {
  margin-bottom: 16px;

  label {
    display: block;
    margin-bottom: 8px;
    font-size: 14px;
    color: #333;
  }

  input {
    width: 100%;
    height: 40px;
    padding: 0 12px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    background-color: #f7f8fa;
    transition: all 0.2s;

    &:focus {
      background-color: white;
      border-color: var(--color-primary, #0066ff);
      outline: none;
    }
  }
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  font-size: 14px;

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 4px;
    cursor: pointer;
    color: #666;
  }

  .forgot-pwd {
    color: var(--color-primary, #0066ff);
    text-decoration: none;
  }
}

.btn-block {
  width: 100%;
  height: 40px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  background-color: var(--color-primary, #0066ff);
  color: white;
  border: none;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #0052cc;
  }

  &:disabled {
    background-color: #ccc;
    cursor: not-allowed;
  }
}

.register-link {
  margin-top: 16px;
  text-align: center;
  font-size: 14px;
  color: #666;

  a {
    color: var(--color-primary, #0066ff);
    text-decoration: none;
    font-weight: 500;
  }
}

.error-alert {
  background-color: #fef2f2;
  color: #ef4444;
  padding: 8px 12px;
  border-radius: 6px;
  margin-bottom: 16px;
  font-size: 14px;
}
</style>
