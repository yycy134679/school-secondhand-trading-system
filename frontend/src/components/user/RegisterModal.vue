<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Modal from '@/components/common/Modal.vue'

defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
  (e: 'switchToLogin'): void
}>()

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

const handleClose = () => {
  emit('update:visible', false)
  errorMessage.value = ''
  form.account = ''
  form.nickname = ''
  form.password = ''
  form.confirmPassword = ''
  form.wechatId = ''
}

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

    emit('update:visible', false)
    emit('success')

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
  <Modal
    :visible="visible"
    title="注册"
    width="400px"
    @update:visible="emit('update:visible', $event)"
    @close="handleClose"
  >
    <form @submit.prevent="handleSubmit" class="register-form">
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

      <div class="form-footer">
        已有账号? <a href="#" @click.prevent="emit('switchToLogin')">立即登录</a>
      </div>
    </form>
  </Modal>
</template>

<style scoped lang="scss">
.register-form {
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

    .required {
      color: #ff4d4f;
    }
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
