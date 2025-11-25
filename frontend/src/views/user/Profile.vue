<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import AvatarUpload from '@/components/user/AvatarUpload.vue'
import ProductCard from '@/components/product/ProductCard.vue'
import Empty from '@/components/common/Empty.vue'
import Loading from '@/components/common/Loading.vue'
import { getRecentViews } from '@/api/user'
import type { Product } from '@common/types/product'

type TabType = 'profile' | 'password' | 'recent'

const userStore = useUserStore()
const activeTab = ref<TabType>('profile')

// 个人资料表单
const profileForm = ref({
  nickname: '',
  wechatId: '',
})
const profileSubmitting = ref(false)
const profileMessage = ref('')

// 修改密码表单
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})
const passwordSubmitting = ref(false)
const passwordMessage = ref('')

// 最近浏览
const recentProducts = ref<Product[]>([])
const recentLoading = ref(false)

const currentUser = computed(() => userStore.currentUser)

onMounted(() => {
  if (currentUser.value) {
    profileForm.value.nickname = currentUser.value.nickname
    profileForm.value.wechatId = currentUser.value.wechatId || ''
  }
})

// 处理头像上传成功
const handleAvatarSuccess = (url: string) => {
  console.log('头像上传成功:', url)
}

const handleAvatarError = (message: string) => {
  alert(message)
}

// 更新个人资料
const handleProfileSubmit = async () => {
  if (!profileForm.value.nickname.trim()) {
    profileMessage.value = '昵称不能为空'
    return
  }

  profileSubmitting.value = true
  profileMessage.value = ''
  try {
    await userStore.updateProfile({
      nickname: profileForm.value.nickname,
      wechatId: profileForm.value.wechatId || undefined,
    })
    profileMessage.value = '保存成功'
    setTimeout(() => {
      profileMessage.value = ''
    }, 3000)
  } catch (error: unknown) {
    const err = error as { response?: { data?: { code?: number; message?: string } } }
    const code = err?.response?.data?.code
    if (code === 1001) {
      profileMessage.value = '昵称修改过于频繁，请稍后再试'
    } else {
      profileMessage.value = err?.response?.data?.message || '保存失败'
    }
  } finally {
    profileSubmitting.value = false
  }
}

// 修改密码
const handlePasswordSubmit = async () => {
  passwordMessage.value = ''

  if (
    !passwordForm.value.oldPassword ||
    !passwordForm.value.newPassword ||
    !passwordForm.value.confirmPassword
  ) {
    passwordMessage.value = '请填写所有字段'
    return
  }

  if (passwordForm.value.newPassword.length < 8) {
    passwordMessage.value = '新密码至少8位'
    return
  }

  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    passwordMessage.value = '两次输入的密码不一致'
    return
  }

  passwordSubmitting.value = true
  try {
    await userStore.changePassword({
      oldPassword: passwordForm.value.oldPassword,
      newPassword: passwordForm.value.newPassword,
      confirmPassword: passwordForm.value.confirmPassword,
    })
    passwordMessage.value = '密码修改成功'
    // 清空表单
    passwordForm.value = {
      oldPassword: '',
      newPassword: '',
      confirmPassword: '',
    }
    setTimeout(() => {
      passwordMessage.value = ''
    }, 3000)
  } catch (error: unknown) {
    const err = error as { response?: { data?: { code?: number; message?: string } } }
    const code = err?.response?.data?.code
    if (code === 1001) {
      passwordMessage.value = '旧密码不正确或新密码不满足要求'
    } else {
      passwordMessage.value = err?.response?.data?.message || '密码修改失败'
    }
  } finally {
    passwordSubmitting.value = false
  }
}

// 切换 Tab 时加载最近浏览
const handleTabChange = async (tab: TabType) => {
  activeTab.value = tab
  if (tab === 'recent' && recentProducts.value.length === 0) {
    await loadRecentViews()
  }
}

// 加载最近浏览
const loadRecentViews = async () => {
  recentLoading.value = true
  try {
    const res = await getRecentViews()
    recentProducts.value = res.data.data
  } catch (error) {
    console.error('加载最近浏览失败:', error)
  } finally {
    recentLoading.value = false
  }
}
</script>

<template>
  <div class="profile-page">
    <div class="profile-container">
      <h1 class="page-title">个人中心</h1>

      <!-- Tab 导航 -->
      <div class="tabs">
        <button
          :class="['tab-item', { active: activeTab === 'profile' }]"
          @click="handleTabChange('profile')"
        >
          个人资料
        </button>
        <button
          :class="['tab-item', { active: activeTab === 'password' }]"
          @click="handleTabChange('password')"
        >
          修改密码
        </button>
        <button
          :class="['tab-item', { active: activeTab === 'recent' }]"
          @click="handleTabChange('recent')"
        >
          最近浏览
        </button>
      </div>

      <!-- Tab 内容 -->
      <div class="tab-content">
        <!-- 个人资料 Tab -->
        <div v-show="activeTab === 'profile'" class="profile-tab">
          <div class="form-section">
            <div class="form-item">
              <label class="form-label">头像</label>
              <div class="form-control">
                <AvatarUpload
                  :model-value="currentUser?.avatarUrl"
                  @success="handleAvatarSuccess"
                  @error="handleAvatarError"
                />
              </div>
            </div>

            <div class="form-item">
              <label class="form-label">账号</label>
              <div class="form-control">
                <input
                  type="text"
                  :value="currentUser?.account"
                  disabled
                  class="input-field disabled"
                />
                <span class="form-hint">账号不可修改</span>
              </div>
            </div>

            <div class="form-item">
              <label class="form-label">昵称</label>
              <div class="form-control">
                <input
                  v-model="profileForm.nickname"
                  type="text"
                  class="input-field"
                  placeholder="请输入昵称"
                />
                <span class="form-hint">昵称30天内只能修改一次</span>
              </div>
            </div>

            <div class="form-item">
              <label class="form-label">微信号</label>
              <div class="form-control">
                <input
                  v-model="profileForm.wechatId"
                  type="text"
                  class="input-field"
                  placeholder="请输入微信号"
                />
                <span class="form-hint">用于买家联系您</span>
              </div>
            </div>

            <div class="form-message" v-if="profileMessage">
              {{ profileMessage }}
            </div>

            <div class="form-actions">
              <button
                @click="handleProfileSubmit"
                :disabled="profileSubmitting"
                class="btn btn-primary"
              >
                {{ profileSubmitting ? '保存中...' : '保存修改' }}
              </button>
            </div>
          </div>
        </div>

        <!-- 修改密码 Tab -->
        <div v-show="activeTab === 'password'" class="password-tab">
          <div class="form-section">
            <div class="form-item">
              <label class="form-label">旧密码</label>
              <div class="form-control">
                <input
                  v-model="passwordForm.oldPassword"
                  type="password"
                  class="input-field"
                  placeholder="请输入旧密码"
                />
              </div>
            </div>

            <div class="form-item">
              <label class="form-label">新密码</label>
              <div class="form-control">
                <input
                  v-model="passwordForm.newPassword"
                  type="password"
                  class="input-field"
                  placeholder="请输入新密码（至少8位）"
                />
              </div>
            </div>

            <div class="form-item">
              <label class="form-label">确认新密码</label>
              <div class="form-control">
                <input
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  class="input-field"
                  placeholder="请再次输入新密码"
                />
              </div>
            </div>

            <div class="form-message" v-if="passwordMessage">
              {{ passwordMessage }}
            </div>

            <div class="form-actions">
              <button
                @click="handlePasswordSubmit"
                :disabled="passwordSubmitting"
                class="btn btn-primary"
              >
                {{ passwordSubmitting ? '修改中...' : '修改密码' }}
              </button>
            </div>
          </div>
        </div>

        <!-- 最近浏览 Tab -->
        <div v-show="activeTab === 'recent'" class="recent-tab">
          <Loading v-if="recentLoading" />
          <Empty
            v-else-if="recentProducts.length === 0"
            title="暂无浏览记录"
            description="去首页看看吧"
          />
          <div v-else class="products-grid">
            <ProductCard v-for="product in recentProducts" :key="product.id" :product="product" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.profile-page {
  min-height: calc(100vh - 64px - 80px);
  padding: 32px 0;
  background: #f7f8fa;
}

.profile-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

.page-title {
  font-size: 24px;
  font-weight: bold;
  color: #1a1a1a;
  margin-bottom: 24px;
}

.tabs {
  display: flex;
  gap: 8px;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 32px;
}

.tab-item {
  padding: 12px 24px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  color: #666666;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    color: #1a1a1a;
  }

  &.active {
    color: #0066ff;
    border-bottom-color: #0066ff;
  }
}

.tab-content {
  background: white;
  border-radius: 12px;
  padding: 32px;
  min-height: 400px;
}

.form-section {
  max-width: 600px;
}

.form-item {
  margin-bottom: 24px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.form-control {
  .input-field {
    width: 100%;
    height: 40px;
    padding: 0 16px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 14px;
    transition: all 0.2s;
    background: white;

    &:focus {
      outline: none;
      border-color: #0066ff;
      box-shadow: 0 0 0 3px rgba(0, 102, 255, 0.1);
    }

    &.disabled {
      background: #f7f8fa;
      color: #999999;
      cursor: not-allowed;
    }
  }

  .form-hint {
    display: block;
    font-size: 12px;
    color: #999999;
    margin-top: 4px;
  }
}

.form-message {
  padding: 12px 16px;
  background: #f0f9ff;
  border: 1px solid #0066ff;
  border-radius: 8px;
  color: #0066ff;
  font-size: 14px;
  margin-bottom: 16px;
}

.form-actions {
  margin-top: 32px;
}

.btn {
  height: 40px;
  padding: 0 32px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.btn-primary {
  background: #0066ff;
  color: white;

  &:hover:not(:disabled) {
    background: #0052cc;
  }
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 20px;
}
</style>
