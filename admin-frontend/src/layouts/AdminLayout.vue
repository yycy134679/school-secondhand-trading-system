<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="brand">校园二手平台</div>
      <nav>
        <RouterLink
          v-for="link in links"
          :key="link.to"
          :to="link.to"
          class="nav-link"
          :class="{ active: route.path === link.to }"
        >
          {{ link.label }}
        </RouterLink>
      </nav>
    </aside>
    <div class="main">
      <header class="topbar">
        <h1>{{ currentTitle }}</h1>
        <div class="topbar-actions">
          <div class="user-pill" :class="{ viewer: isViewer }">
            <div class="user-avatar">{{ displayInitial }}</div>
            <div class="user-meta">
              <span class="user-name">{{ displayName }}</span>
              <span class="user-account">账号：{{ displayAccount }}</span>
              <span class="user-role">{{ roleLabel }}</span>
            </div>
          </div>
          <div v-if="isViewer" class="application-area">
            <button
              type="button"
              class="action-btn apply"
              :disabled="applyButtonDisabled"
              @click="openApplicationModal"
            >
              {{ applyButtonLabel }}
            </button>
            <span v-if="applicationStatusText" class="application-status">
              {{ applicationStatusText }}
            </span>
          </div>
          <button type="button" class="action-btn ghost" @click="returnToFrontend">
            返回首页
          </button>
          <button
            type="button"
            class="action-btn"
            :disabled="!canManage"
            @click="handleRestricted(openAnnouncementModal)"
          >
            发布公告
          </button>
          <button
            type="button"
            class="action-btn primary"
            :disabled="!canManage"
            @click="handleRestricted(openCreateAdminModal)"
          >
            新增管理员
          </button>
        </div>
      </header>
      <main class="content">
        <RouterView />
      </main>
    </div>
  </div>

  <BaseModal
    :open="showAnnouncementModal"
    title="发布公告"
    width="520px"
    @close="closeAnnouncementModal"
  >
    <form class="modal-form" @submit.prevent="submitAnnouncement">
      <label class="modal-label" for="announcement-input">公告内容</label>
      <textarea
        id="announcement-input"
        v-model="announcementForm.content"
        class="modal-textarea"
        rows="5"
        placeholder="请填写需要发布的公告内容"
      ></textarea>
      <p v-if="announcementForm.error" class="modal-error">{{ announcementForm.error }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeAnnouncementModal">取消</button>
        <button type="submit" class="btn primary">发布</button>
      </div>
    </form>
  </BaseModal>

  <BaseModal
    :open="showCreateAdminModal"
    title="新增管理员"
    width="520px"
    @close="closeCreateAdminModal"
  >
    <form class="modal-form" @submit.prevent="submitCreateAdmin">
      <label class="modal-label" for="admin-account-input">管理员账号</label>
      <input
        id="admin-account-input"
        v-model="createAdminForm.account"
        class="modal-input"
        placeholder="请输入管理员账号"
      />
      <label class="modal-label" for="admin-nickname-input">管理员昵称</label>
      <input
        id="admin-nickname-input"
        v-model="createAdminForm.nickname"
        class="modal-input"
        placeholder="请输入管理员昵称"
      />
      <label class="modal-label" for="admin-wechat-input">联系方式（可选）</label>
      <input
        id="admin-wechat-input"
        v-model="createAdminForm.wechat"
        class="modal-input"
        placeholder="请输入微信或手机号"
      />
      <p v-if="createAdminForm.error" class="modal-error">{{ createAdminForm.error }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeCreateAdminModal">取消</button>
        <button type="submit" class="btn primary">创建</button>
      </div>
    </form>
  </BaseModal>

  <BaseModal
    :open="applicationModal.open"
    title="申请成为管理员"
    width="520px"
    @close="closeApplicationModal"
  >
    <form class="modal-form" @submit.prevent="submitApplication">
      <p v-if="applicationStatusText" class="modal-tip">{{ applicationStatusText }}</p>
      <label class="modal-label" for="apply-reason-input">申请理由</label>
      <textarea
        id="apply-reason-input"
        v-model="applicationModal.reason"
        class="modal-textarea"
        rows="4"
        placeholder="请简单描述您希望成为管理员的原因或能够提供的帮助"
      ></textarea>
      <p v-if="applicationModal.error" class="modal-error">{{ applicationModal.error }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeApplicationModal">取消</button>
        <button type="submit" class="btn primary">提交申请</button>
      </div>
    </form>
  </BaseModal>

  <BaseModal
    :open="Boolean(feedbackModal)"
    title="操作提示"
    width="420px"
    @close="closeFeedback"
  >
    <p class="modal-message">{{ feedbackModal?.message }}</p>
    <div class="modal-actions">
      <button type="button" class="btn primary" @click="closeFeedback">知道了</button>
    </div>
  </BaseModal>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

import { useAdminDataStore } from '../stores/admin'
import BaseModal from '../components/BaseModal.vue'
import { useAuthStore } from '../stores/auth'

interface NavLink {
  to: string
  label: string
}

const route = useRoute()
const router = useRouter()
const adminStore = useAdminDataStore()
const authStore = useAuthStore()
const { users, adminApplications } = storeToRefs(adminStore)
const { canManage, isViewer } = storeToRefs(authStore)

const links: NavLink[] = [
  { to: '/admin', label: '仪表盘' },
  { to: '/admin/products', label: '商品管理' },
  { to: '/admin/users', label: '用户管理' },
  { to: '/admin/categories', label: '分类与标签' },
  { to: '/admin/profile', label: '个人信息' },
]

const currentTitle = computed(() => {
  const found = links.find((item) => route.path === item.to)
  return found?.label ?? '后台管理'
})

const showAnnouncementModal = ref(false)
const announcementForm = reactive({ content: '', error: '' })

const showCreateAdminModal = ref(false)
const createAdminForm = reactive({ account: '', nickname: '', wechat: '', error: '' })

const feedbackModal = ref<{ message: string } | null>(null)

const DEFAULT_FRONTEND_URL = 'http://localhost:5173'

const frontendBaseUrl = computed(() => {
  const envValue = import.meta.env.VITE_FRONTEND_URL as string | undefined
  const base = envValue?.trim() || DEFAULT_FRONTEND_URL
  const normalized = base.replace(/\/+$/, '')
  return normalized || DEFAULT_FRONTEND_URL
})

const displayName = computed(() => authStore.currentAccount?.nickname || authStore.currentAccount?.account || '管理员')
const displayAccount = computed(() => authStore.currentAccount?.account || '未登录')
const displayInitial = computed(() => {
  const source = authStore.currentAccount?.nickname || authStore.currentAccount?.account || '管'
  return source.charAt(0).toUpperCase()
})
const roleLabel = computed(() => (authStore.currentAccount?.role === 'admin' ? '角色：管理员' : '角色：普通成员（仅查看）'))

const viewerAccount = computed(() => authStore.currentAccount?.account ?? '')

const latestApplication = computed(() => {
  const account = viewerAccount.value
  if (!account) return null
  return adminApplications.value.find((item) => item.account === account) ?? null
})

const hasPendingApplication = computed(
  () => latestApplication.value?.status === 'pending',
)

const applicationStatusText = computed(() => {
  const application = latestApplication.value
  if (!application) return ''
  if (application.status === 'pending') {
    return `已于 ${application.createdAt} 提交申请，等待管理员审核。`
  }
  if (application.status === 'rejected') {
    return application.feedback
      ? `最近一次申请被驳回，原因：${application.feedback}`
      : '最近一次申请被驳回，请完善信息后再次尝试。'
  }
  return ''
})

const applyButtonDisabled = computed(() => {
  if (!isViewer.value) return true
  if (!authStore.isAuthenticated) return true
  return hasPendingApplication.value
})

const applyButtonLabel = computed(() =>
  hasPendingApplication.value ? '申请审核中' : '申请成为管理员',
)

const applicationModal = reactive({ open: false, reason: '', error: '' })

const returnToFrontend = () => {
  if (typeof window === 'undefined') return
  window.location.href = `${frontendBaseUrl.value}/`
}

const openApplicationModal = () => {
  if (applyButtonDisabled.value) {
    feedbackModal.value = {
      message: hasPendingApplication.value
        ? '您的申请正在审核中，请耐心等待管理员处理。'
        : '当前账号暂不支持申请管理员，请确认已登录。',
    }
    return
  }
  applicationModal.reason = ''
  applicationModal.error = ''
  applicationModal.open = true
}

const closeApplicationModal = () => {
  applicationModal.open = false
}

const submitApplication = () => {
  if (!authStore.currentAccount) {
    applicationModal.error = '未检测到登录信息，请重新登录后再试。'
    return
  }

  const reason = applicationModal.reason.trim()
  if (!reason) {
    applicationModal.error = '请填写申请理由。'
    return
  }

  const result = adminStore.submitAdminApplication({
    account: authStore.currentAccount.account,
    nickname: authStore.currentAccount.nickname || authStore.currentAccount.account,
    reason,
  })

  if (!result.success) {
    applicationModal.error = result.message
    return
  }

  applicationModal.open = false
  applicationModal.reason = ''
  feedbackModal.value = { message: result.message }
}

const handleRestricted = (fn: () => void) => {
  if (!canManage.value) {
    feedbackModal.value = { message: '当前账号仅支持查看数据，如需操作请联系管理员。' }
    return
  }
  fn()
}

const openAnnouncementModal = () => {
  announcementForm.content = ''
  announcementForm.error = ''
  showAnnouncementModal.value = true
}

const closeAnnouncementModal = () => {
  showAnnouncementModal.value = false
}

const submitAnnouncement = () => {
  const content = announcementForm.content.trim()
  if (!content) {
    announcementForm.error = '公告内容不能为空。'
    return
  }
  adminStore.addAnnouncement(content)
  showAnnouncementModal.value = false
  feedbackModal.value = { message: '公告已发布。' }
}

const openCreateAdminModal = () => {
  createAdminForm.account = ''
  createAdminForm.nickname = ''
  createAdminForm.wechat = ''
  createAdminForm.error = ''
  showCreateAdminModal.value = true
}

const closeCreateAdminModal = () => {
  showCreateAdminModal.value = false
}

const submitCreateAdmin = () => {
  const account = createAdminForm.account.trim()
  const nickname = (createAdminForm.nickname || account).trim()
  if (!account) {
    createAdminForm.error = '账号不能为空。'
    return
  }
  const exists = users.value.some((user) => user.account === account)
  if (exists) {
    createAdminForm.error = '账号已存在，请换一个账号。'
    return
  }

  adminStore.createAdminAccount({
    account,
    nickname,
    wechat: createAdminForm.wechat.trim() || undefined,
  })
  authStore.addAdminAccount(account, nickname, '123456')
  showCreateAdminModal.value = false
  feedbackModal.value = { message: '管理员账号已创建并设为启用状态。' }
  router.push('/admin/users')
}

const closeFeedback = () => {
  feedbackModal.value = null
}
</script>

<style scoped>
.layout {
  display: grid;
  grid-template-columns: 240px 1fr;
  min-height: 100vh;
  background: #f7f8fa;
  color: #1f2937;
}

.sidebar {
  background: #111827;
  color: #f9fafb;
  padding: 24px 16px;
  display: flex;
  flex-direction: column;
}

.brand {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 32px;
}

.nav-link {
  display: block;
  padding: 12px 14px;
  border-radius: 8px;
  color: #d1d5db;
  text-decoration: none;
  margin-bottom: 8px;
  transition: background 0.2s ease, color 0.2s ease;
}

.nav-link:hover {
  background: rgba(59, 130, 246, 0.12);
  color: #f9fafb;
}

.nav-link.active {
  background: #2563eb;
  color: #f9fafb;
}

.main {
  display: flex;
  flex-direction: column;
}

.topbar {
  background: #fff;
  padding: 20px 28px;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.topbar h1 {
  margin: 0;
  font-size: 20px;
}

.topbar-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.user-pill {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 14px;
  border-radius: 999px;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
}

.application-area {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  max-width: 220px;
}

.application-status {
  font-size: 12px;
  color: #6b7280;
  line-height: 1.3;
}

.user-pill.viewer {
  background: #fff7ed;
  border-color: #fdba74;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #2563eb, #60a5fa);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 16px;
}

.user-meta {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.user-name {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.user-account {
  font-size: 12px;
  color: #6b7280;
}

.user-role {
  font-size: 12px;
  color: #2563eb;
}

.action-btn {
  border: 1px solid #d1d5db;
  background: #fff;
  color: #1f2937;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.2s ease, color 0.2s ease;
}

.action-btn:hover {
  background: #f3f4f6;
}

.action-btn:disabled {
  cursor: not-allowed;
  opacity: 0.55;
  background: #f9fafb;
  color: #9ca3af;
}

.action-btn.apply {
  border-color: #2563eb;
  color: #2563eb;
}

.action-btn.apply:not(:disabled):hover {
  background: rgba(37, 99, 235, 0.12);
  color: #1d4ed8;
}

.action-btn.primary {
  border: none;
  background: #2563eb;
  color: #fff;
}

.action-btn.ghost {
  background: transparent;
}

.action-btn.primary:hover {
  background: #1d4ed8;
}

.content {
  flex: 1;
  padding: 28px;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.modal-label {
  font-size: 14px;
  color: #4b5563;
}

.modal-input,
.modal-textarea {
  width: 100%;
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px solid #d1d5db;
  font-size: 14px;
}

.modal-textarea {
  resize: vertical;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 8px;
}

.modal-error {
  color: #dc2626;
  font-size: 13px;
  margin: 4px 0 0;
}

.modal-message {
  margin: 0;
  font-size: 15px;
  color: #374151;
}

.modal-tip {
  margin: 0 0 8px;
  font-size: 13px;
  color: #6b7280;
}

.btn {
  border: 1px solid #d1d5db;
  background: #fff;
  color: #1f2937;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
}

.btn.primary {
  border: none;
  background: #2563eb;
  color: #fff;
}

.btn.ghost {
  background: transparent;
}

@media (max-width: 960px) {
  .layout {
    grid-template-columns: 72px 1fr;
  }

  .brand {
    display: none;
  }

  .nav-link {
    text-align: center;
    padding: 10px 8px;
  }
}
</style>
