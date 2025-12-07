<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import LoginModal from '@/components/user/LoginModal.vue'
import RegisterModal from '@/components/user/RegisterModal.vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const searchQuery = ref(route.query.q?.toString() || '')
const showLoginModal = ref(false)
const showRegisterModal = ref(false)
const normalizeAdminBaseUrl = (urlValue: string | undefined) => {
  const fallback = 'http://localhost:5174'
  if (!urlValue) return fallback

  const trimmed = urlValue.trim()
  if (!trimmed) return fallback

  if (/\/login\/?$/i.test(trimmed)) {
    return trimmed.replace(/\/login\/?$/i, '') || fallback
  }

  if (/\/admin\/?$/i.test(trimmed)) {
    return trimmed.replace(/\/admin\/?$/i, '') || fallback
  }

  return trimmed.replace(/\/+$/, '') || fallback
}

const adminBaseUrl = computed(() =>
  normalizeAdminBaseUrl(import.meta.env.VITE_ADMIN_URL as string | undefined),
)

const adminLoginUrl = computed(() => `${adminBaseUrl.value.replace(/\/+$/, '')}/login`)

const handleSearch = () => {
  if (!searchQuery.value.trim()) return
  router.push({
    path: '/search',
    query: { q: searchQuery.value },
  })
}

const handlePublish = () => {
  if (!userStore.isLoggedIn) {
    showLoginModal.value = true
  } else {
    router.push('/products/new')
  }
}

const handleLoginClick = () => {
  showLoginModal.value = true
}

const handleLogout = () => {
  userStore.logout()
  router.push('/')
}

const handleLoginSuccess = () => {
  // Login modal handles closing itself
}

const handleRegisterSuccess = () => {
  // Register modal handles closing itself and redirecting
}

const handleAdminEntry = () => {
  if (typeof window === 'undefined') return
  window.location.href = adminLoginUrl.value
}

const handleLogoClick = async () => {
  try {
    if (route.path === '/') {
      await router.replace('/')
    } else {
      await router.push('/')
    }
  } catch (err) {
    console.warn('返回首页导航失败，尝试刷新页面', err)
    if (typeof window !== 'undefined') {
      window.location.href = '/'
    }
  }
}

const switchToRegister = () => {
  showLoginModal.value = false
  showRegisterModal.value = true
}

const switchToLogin = () => {
  showRegisterModal.value = false
  showLoginModal.value = true
}
</script>

<template>
  <header class="app-header">
    <div class="container header-content">
      <div class="logo">
        <a href="/" @click.prevent="handleLogoClick">校园二手</a>
      </div>

      <div class="search-bar">
        <div class="search-wrapper">
          <input
            type="text"
            v-model="searchQuery"
            placeholder="搜索 教材 / 耳机 / 自行车..."
            @keyup.enter="handleSearch"
          />
          <button class="search-btn" @click="handleSearch">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="11" cy="11" r="8"></circle>
              <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
            </svg>
          </button>
        </div>
      </div>

      <div class="actions">
        <button class="btn btn-primary publish-btn" @click="handlePublish">发布闲置</button>

        <div v-if="userStore.isLoggedIn" class="user-menu">
          <div class="avatar">
            <img
              :src="
                userStore.currentUser?.avatarUrl ||
                'https://ui-avatars.com/api/?name=' + (userStore.currentUser?.nickname || 'User')
              "
              alt="Avatar"
            />
          </div>
          <div class="dropdown-menu">
            <div class="user-name">{{ userStore.currentUser?.nickname }}</div>
            <router-link to="/profile" class="menu-item">个人中心</router-link>
            <router-link to="/my/products" class="menu-item">我发布的</router-link>
            <div class="divider"></div>
            <button class="menu-item logout-btn" @click="handleLogout">退出登录</button>
          </div>
        </div>

        <div v-else class="auth-buttons">
          <button class="btn btn-text" @click="handleLoginClick">登录 / 注册</button>
        </div>
      </div>
    </div>

    <LoginModal
      v-model:visible="showLoginModal"
      @success="handleLoginSuccess"
      @switch-to-register="switchToRegister"
    />

    <RegisterModal
      v-model:visible="showRegisterModal"
      @success="handleRegisterSuccess"
      @switch-to-login="switchToLogin"
    />
  </header>
</template>

<style scoped lang="scss">
.app-header {
  height: 64px;
  background-color: var(--color-bg-white, #fff);
  border-bottom: 1px solid var(--color-border, #e5e7eb);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.logo a {
  font-size: 24px;
  font-weight: bold;
  color: var(--color-primary, #0066ff);
  text-decoration: none;
  letter-spacing: -0.5px;
}

.search-bar {
  flex: 1;
  max-width: 500px;
  margin: 0 40px;
}

.search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-wrapper input {
  width: 100%;
  padding: 10px 40px 10px 16px;
  border: none;
  border-radius: 20px;
  background-color: #f7f8fa;
  font-size: 14px;
  transition: all 0.2s;

  &:focus {
    background-color: #fff;
    box-shadow: 0 0 0 2px var(--color-primary, #0066ff);
    outline: none;
  }
}

.search-btn {
  position: absolute;
  right: 8px;
  background: none;
  border: none;
  cursor: pointer;
  color: #666;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;

  &:hover {
    color: var(--color-primary, #0066ff);
  }
}

.actions {
  display: flex;
  align-items: center;
  gap: 20px;
}

.admin-link {
  background: none;
  padding: 6px 16px;
  border-radius: 18px;
  border: 1px solid #d1d5db;
  color: #1f2937;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  transition:
    background-color 0.2s ease,
    color 0.2s ease;

  &:hover {
    background-color: #f3f4f6;
    color: var(--color-primary, #0066ff);
  }
}

.publish-btn {
  padding: 8px 20px;
  border-radius: 20px;
  font-weight: 500;
  background-color: var(--color-primary, #0066ff);
  color: white;
  border: none;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #0052cc;
  }
}

.btn-text {
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;

  &:hover {
    color: var(--color-primary, #0066ff);
  }
}

.user-menu {
  position: relative;
  cursor: pointer;

  &:hover .dropdown-menu {
    display: block;
  }
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  border: 1px solid #e5e7eb;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.dropdown-menu {
  display: none;
  position: absolute;
  top: 100%;
  right: 0;
  width: 160px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  padding: 8px 0;
  margin-top: 8px;
  border: 1px solid #e5e7eb;

  &::before {
    content: '';
    position: absolute;
    top: -8px;
    left: 0;
    right: 0;
    height: 8px;
  }
}

.user-name {
  padding: 8px 16px;
  font-weight: bold;
  color: #1a1a1a;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.menu-item {
  display: block;
  padding: 8px 16px;
  color: #666;
  text-decoration: none;
  font-size: 14px;
  text-align: left;
  width: 100%;
  background: none;
  border: none;
  cursor: pointer;

  &:hover {
    background-color: #f7f8fa;
    color: var(--color-primary, #0066ff);
  }
}

.divider {
  height: 1px;
  background-color: #f0f0f0;
  margin: 4px 0;
}

.logout-btn {
  color: #ef4444;

  &:hover {
    color: #dc2626;
    background-color: #fef2f2;
  }
}
</style>
