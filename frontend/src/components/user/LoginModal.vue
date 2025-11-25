<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useUserStore } from '@/stores/user'
import Modal from '@/components/common/Modal.vue'

defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
  (e: 'switchToRegister'): void
}>()

const userStore = useUserStore()

const form = reactive({
  account: '',
  password: '',
  rememberMe: false,
})

const loading = ref(false)
const errorMessage = ref('')

const handleClose = () => {
  emit('update:visible', false)
  errorMessage.value = ''
  form.account = ''
  form.password = ''
}

const handleSubmit = async () => {
  if (!form.account || !form.password) {
    errorMessage.value = '请输入账号和密码'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    await userStore.login({
      account: form.account,
      password: form.password,
      rememberMe: form.rememberMe,
    })

    emit('update:visible', false)
    emit('success')

    // Refresh page to ensure all states are updated
    window.location.reload()
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '登录失败'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Modal
    :visible="visible"
    title="登录"
    width="400px"
    @update:visible="emit('update:visible', $event)"
    @close="handleClose"
  >
    <form @submit.prevent="handleSubmit" class="login-form">
      <div class="form-item">
        <label>账号</label>
        <input
          v-model="form.account"
          type="text"
          placeholder="请输入账号"
          :disabled="loading"
          autofocus
        />
      </div>

      <div class="form-item">
        <label>密码</label>
        <input
          v-model="form.password"
          type="password"
          placeholder="请输入密码"
          :disabled="loading"
        />
      </div>

      <div class="form-actions">
        <label class="checkbox">
          <input type="checkbox" v-model="form.rememberMe" :disabled="loading" />
          <span>记住我</span>
        </label>
        <a href="#" class="forgot-password">忘记密码?</a>
      </div>

      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>

      <button type="submit" class="btn-submit" :disabled="loading">
        {{ loading ? '登录中...' : '登录' }}
      </button>

      <div class="form-footer">
        还没有账号? <a href="#" @click.prevent="emit('switchToRegister')">立即注册</a>
      </div>
    </form>
  </Modal>
</template>

<style scoped lang="scss">
.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 8px;

  label {
    font-size: 14px;
    color: #333;
    font-weight: 500;
  }

  input {
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
    transition: border-color 0.2s;

    &:focus {
      border-color: var(--color-primary, #0066ff);
      outline: none;
    }

    &:disabled {
      background-color: #f5f5f5;
      cursor: not-allowed;
    }
  }
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;

  .checkbox {
    display: flex;
    align-items: center;
    gap: 4px;
    cursor: pointer;
    user-select: none;
  }

  .forgot-password {
    color: #666;
    text-decoration: none;

    &:hover {
      color: var(--color-primary, #0066ff);
    }
  }
}

.error-message {
  color: #ff4d4f;
  font-size: 13px;
  background: #fff2f0;
  padding: 8px;
  border-radius: 4px;
  border: 1px solid #ffccc7;
}

.btn-submit {
  background-color: var(--color-primary, #0066ff);
  color: white;
  border: none;
  padding: 12px;
  border-radius: 4px;
  font-size: 16px;
  font-weight: 500;
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
  margin-top: 8px;

  a {
    color: var(--color-primary, #0066ff);
    text-decoration: none;
    font-weight: 500;

    &:hover {
      text-decoration: underline;
    }
  }
}
</style>
