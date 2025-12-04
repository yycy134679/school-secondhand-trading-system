<template>
  <div class="page">
    <header class="page-head">
      <div>
        <h2>用户管理</h2>
        <p class="subtitle">监督平台用户行为、角色分配与账号状态。</p>
      </div>
      <div class="actions">
        <button type="button" class="btn" :disabled="!canManage" @click="openApplicationModal">
          管理员申请
        </button>
        <button type="button" class="btn" :disabled="!canManage" @click="openBatchBanModal">
          批量封禁
        </button>
        <button
          type="button"
          class="btn primary"
          :disabled="!canManage"
          @click="openCreateAdminModal"
        >
          创建管理员
        </button>
      </div>
    </header>

    <section class="filters">
      <input v-model="keyword" class="input" placeholder="搜索账号或昵称" />
      <select v-model="role" class="input">
        <option value="">全部角色</option>
        <option value="admin">管理员</option>
        <option value="user">普通用户</option>
      </select>
      <select v-model="status" class="input">
        <option value="">全部状态</option>
        <option value="active">正常</option>
        <option value="banned">封禁</option>
        <option value="pending">待验证</option>
      </select>
      <button type="button" class="btn" @click="applyFilters">筛选</button>
    </section>

    <table class="table">
      <thead>
        <tr>
          <th>用户</th>
          <th>角色</th>
          <th>状态</th>
          <th>微信号</th>
          <th>最近登录</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in filteredUsers" :key="user.id">
          <td>
            <p class="title">{{ user.account }}</p>
            <p class="meta">昵称：{{ user.nickname }} · 注册：{{ user.registeredAt }}</p>
          </td>
          <td>
            <span class="tag" :class="user.role">{{ roleMap[user.role] }}</span>
          </td>
          <td>
            <span class="status" :class="user.status">{{ statusMap[user.status] }}</span>
          </td>
          <td>{{ user.wechat ?? '—' }}</td>
          <td>{{ user.lastLogin }}</td>
          <td class="ops">
            <button type="button" class="link" @click="openDetail(user.id)">详情</button>
            <button
              type="button"
              class="link"
              :disabled="!canManage"
              @click="openRenameModal(user)"
            >
              修改昵称
            </button>
            <button
              type="button"
              class="link"
              :disabled="!canManage"
              @click="openBanModal(user)"
            >
              {{ user.status === 'banned' ? '解除封禁' : '封禁账号' }}
            </button>
            <button
              type="button"
              class="link"
              :disabled="!canManage"
              @click="openPasswordModal(user)"
            >
              修改密码
            </button>
            <button
              v-if="user.role !== 'admin'"
              type="button"
              class="link"
              :disabled="!canManage"
              @click="openPromoteModal(user)"
            >
              设为管理员
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <BaseModal
    :open="Boolean(detailModal)"
    title="用户详情"
    width="520px"
    @close="closeDetailModal"
  >
    <template v-if="detailModal">
      <dl class="detail-list">
        <div>
          <dt>账号</dt>
          <dd>{{ detailModal.account }}</dd>
        </div>
        <div>
          <dt>昵称</dt>
          <dd>{{ detailModal.nickname }}</dd>
        </div>
        <div>
          <dt>角色</dt>
          <dd>{{ roleMap[detailModal.role] }}</dd>
        </div>
        <div>
          <dt>状态</dt>
          <dd>{{ statusMap[detailModal.status] }}</dd>
        </div>
        <div>
          <dt>微信</dt>
          <dd>{{ detailModal.wechat ?? '未填写' }}</dd>
        </div>
        <div>
          <dt>注册时间</dt>
          <dd>{{ detailModal.registeredAt }}</dd>
        </div>
        <div>
          <dt>最近登录</dt>
          <dd>{{ detailModal.lastLogin }}</dd>
        </div>
      </dl>
      <div class="modal-actions">
        <button type="button" class="btn primary" @click="closeDetailModal">关闭</button>
      </div>
    </template>
  </BaseModal>

  <BaseModal
    :open="Boolean(banModal)"
    :title="banModal?.mode === 'ban' ? '封禁账号' : '解除封禁'"
    width="460px"
    @close="closeBanModal"
  >
    <template v-if="banModal">
      <p class="modal-message">
        {{
          banModal.mode === 'ban'
            ? `确定封禁用户 ${banModal.account}（${banModal.nickname}）吗？`
            : `确定解除封禁用户 ${banModal.account}（${banModal.nickname}）吗？`
        }}
      </p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeBanModal">取消</button>
        <button type="button" class="btn primary" :disabled="!canManage" @click="submitBanModal">
          {{ banModal.mode === 'ban' ? '确认封禁' : '解除封禁' }}
        </button>
      </div>
    </template>
  </BaseModal>

  <BaseModal
    :open="Boolean(promoteModal)"
    title="设为管理员"
    width="480px"
    @close="closePromoteModal"
  >
    <template v-if="promoteModal">
      <p class="modal-message">为用户 {{ promoteModal.account }} 设置管理员身份。</p>
      <div class="modal-field">
        <label class="modal-label" for="promote-nickname-input">管理员昵称</label>
        <input
          id="promote-nickname-input"
          v-model="promoteModal.nickname"
          class="modal-input"
          placeholder="请输入管理员昵称"
        />
      </div>
      <p v-if="promoteModal.error" class="modal-error">{{ promoteModal.error }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closePromoteModal">取消</button>
        <button type="button" class="btn primary" :disabled="!canManage" @click="submitPromoteModal">
          确认设置
        </button>
      </div>
    </template>
  </BaseModal>

  <BaseModal
    :open="passwordModal.open"
    title="设置新密码"
    width="480px"
    @close="closePasswordModal"
  >
    <form class="modal-form" @submit.prevent="submitPasswordModal">
      <p class="modal-message">
        为账号 {{ passwordModal.account }}（{{ passwordModal.nickname }}）设置新密码。
      </p>
      <div class="modal-field">
        <label class="modal-label" for="reset-password-input">新密码</label>
        <input
          id="reset-password-input"
          v-model="passwordModal.newPassword"
          class="modal-input"
          type="password"
          placeholder="请输入至少 6 位新密码"
          autocomplete="new-password"
        />
      </div>
      <div class="modal-field">
        <label class="modal-label" for="reset-password-confirm-input">确认新密码</label>
        <input
          id="reset-password-confirm-input"
          v-model="passwordModal.confirmPassword"
          class="modal-input"
          type="password"
          placeholder="请再次输入新密码"
          autocomplete="new-password"
        />
      </div>
      <p v-if="passwordModal.error" class="modal-error">{{ passwordModal.error }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closePasswordModal">取消</button>
        <button type="submit" class="btn primary" :disabled="!canManage">
          保存新密码
        </button>
      </div>
    </form>
  </BaseModal>

  <BaseModal
    :open="Boolean(renameModal)"
    title="修改昵称"
    width="460px"
    @close="closeRenameModal"
  >
    <template v-if="renameModal">
      <p class="modal-message">为账号 {{ renameModal.account }} 设置新的展示昵称。</p>
      <div class="modal-field">
        <label class="modal-label" for="rename-nickname-input">新昵称</label>
        <input
          id="rename-nickname-input"
          v-model="renameModal.nickname"
          class="modal-input"
          placeholder="请输入新的昵称"
        />
      </div>
      <p v-if="renameModal.error" class="modal-error">{{ renameModal.error }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeRenameModal">取消</button>
        <button
          type="button"
          class="btn primary"
          :disabled="!canManage"
          @click="submitRenameModal"
        >
          保存
        </button>
      </div>
    </template>
  </BaseModal>

  <BaseModal
    :open="Boolean(batchBanModal)"
    title="批量封禁"
    width="480px"
    @close="closeBatchBanModal"
  >
    <template v-if="batchBanModal">
      <p class="modal-message">确定封禁以下 {{ batchBanModal.count }} 位普通用户吗？</p>
      <ul class="modal-list">
        <li v-for="account in batchBanModal.accounts" :key="account">{{ account }}</li>
      </ul>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeBatchBanModal">取消</button>
        <button type="button" class="btn primary" :disabled="!canManage" @click="submitBatchBan">
          确认封禁
        </button>
      </div>
    </template>
  </BaseModal>

  <BaseModal
    :open="createAdminModal.open"
    title="创建管理员"
    width="520px"
    @close="closeCreateAdminModal"
  >
    <form class="modal-form" @submit.prevent="submitCreateAdmin">
      <div class="modal-field">
        <label class="modal-label" for="create-admin-account">管理员账号</label>
        <input
          id="create-admin-account"
          v-model="createAdminModal.account"
          class="modal-input"
          placeholder="请输入管理员账号"
        />
      </div>
      <div class="modal-field">
        <label class="modal-label" for="create-admin-nickname">管理员昵称</label>
        <input
          id="create-admin-nickname"
          v-model="createAdminModal.nickname"
          class="modal-input"
          placeholder="请输入管理员昵称"
        />
      </div>
      <div class="modal-field">
        <label class="modal-label" for="create-admin-wechat">管理员微信（可选）</label>
        <input
          id="create-admin-wechat"
          v-model="createAdminModal.wechat"
          class="modal-input"
          placeholder="请输入联系方式"
        />
      </div>
      <p v-if="createAdminModal.error" class="modal-error">{{ createAdminModal.error }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeCreateAdminModal">取消</button>
        <button type="submit" class="btn primary" :disabled="!canManage">创建</button>
      </div>
    </form>
  </BaseModal>

  <BaseModal
    :open="applicationModal.open"
    title="管理员申请审核"
    width="600px"
    @close="closeApplicationModal"
  >
    <section class="application-section">
      <header class="application-header">
        <h4>待审核申请</h4>
        <span>共 {{ pendingApplications.length }} 条</span>
      </header>
      <ul v-if="pendingApplications.length" class="application-list">
        <li v-for="item in pendingApplications" :key="item.id" class="application-item">
          <div class="application-info">
            <p class="application-title">{{ item.account }}（{{ item.nickname }}）</p>
            <p class="application-reason">{{ item.reason }}</p>
            <p class="application-time">提交时间：{{ item.createdAt }}</p>
          </div>
          <div class="application-actions">
            <button
              type="button"
              class="btn primary"
              :disabled="!canManage"
              @click="approveApplication(item.id)"
            >
              同意
            </button>
            <button
              type="button"
              class="btn danger"
              :disabled="!canManage"
              @click="openRejectApplicationModal(item)"
            >
              驳回
            </button>
          </div>
        </li>
      </ul>
      <p v-else class="application-empty">暂无待审核申请。</p>
    </section>

    <section v-if="processedApplications.length" class="application-section history">
      <header class="application-header">
        <h4>已处理记录</h4>
        <span>共 {{ processedApplications.length }} 条</span>
      </header>
      <ul class="history-list">
        <li v-for="item in processedApplications" :key="item.id" class="history-item">
          <div>
            <p class="history-title">{{ item.account }}（{{ item.nickname }}）</p>
            <p class="history-summary">
              {{ item.status === 'approved' ? '已通过申请' : '已驳回申请' }} · 处理人：{{ item.reviewer || '系统' }}
            </p>
            <p class="history-meta">
              处理时间：{{ item.processedAt || '待定' }}
              <span v-if="item.feedback"> · 备注：{{ item.feedback }}</span>
            </p>
          </div>
        </li>
      </ul>
    </section>
  </BaseModal>

  <BaseModal
    :open="rejectApplicationModal.open"
    title="驳回管理员申请"
    width="520px"
    @close="closeRejectApplicationModal"
  >
    <form class="modal-form" @submit.prevent="submitRejectApplication">
      <p class="modal-message">
        确认驳回 {{ rejectApplicationModal.account }}（{{ rejectApplicationModal.nickname }}）的管理员申请？
      </p>
      <label class="modal-label" for="reject-application-reason">驳回原因</label>
      <textarea
        id="reject-application-reason"
        v-model="rejectApplicationModal.reason"
        class="modal-textarea"
        rows="4"
        placeholder="请填写驳回理由，将反馈给申请人"
      ></textarea>
      <p v-if="rejectApplicationModal.error" class="modal-error">{{ rejectApplicationModal.error }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeRejectApplicationModal">取消</button>
        <button type="submit" class="btn primary" :disabled="!canManage">确认驳回</button>
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
import { storeToRefs } from 'pinia'

import BaseModal from '../components/BaseModal.vue'
import { useAdminDataStore } from '../stores/admin'
import { useAuthStore } from '../stores/auth'
import type { AdminUser, AdminApplication, UserRole, UserStatus } from '../stores/admin'

defineOptions({ name: 'AdminUserListView' })

const keyword = ref('')
const role = ref<UserRole | ''>('')
const status = ref<UserStatus | ''>('')

const adminStore = useAdminDataStore()
const authStore = useAuthStore()
const { users, adminApplications } = storeToRefs(adminStore)
const { canManage } = storeToRefs(authStore)

const roleMap: Record<UserRole, string> = {
  admin: '管理员',
  user: '普通成员',
}

const statusMap: Record<UserStatus, string> = {
  active: '正常',
  banned: '封禁',
  pending: '待验证',
}

const filteredUsers = computed(() => {
  const keywordValue = keyword.value.trim()
  const roleFilter = role.value
  const statusFilter = status.value

  return users.value.filter((user: AdminUser) => {
    const keywordMatch = keywordValue
      ? [user.account, user.nickname].some((text) => text.includes(keywordValue))
      : true
    const roleMatch = roleFilter ? user.role === roleFilter : true
    const statusMatch = statusFilter ? user.status === statusFilter : true
    return keywordMatch && roleMatch && statusMatch
  })
})

const pendingApplications = computed(() =>
  adminApplications.value.filter((item: AdminApplication) => item.status === 'pending'),
)

const processedApplications = computed(() =>
  adminApplications.value.filter((item: AdminApplication) => item.status !== 'pending'),
)

const feedbackModal = ref<{ message: string } | null>(null)

const showPermissionWarning = () => {
  feedbackModal.value = { message: '当前账号仅支持查看数据，如需操作请联系管理员。' }
}

const detailModal = ref<AdminUser | null>(null)

interface BanModalState {
  id: number
  mode: 'ban' | 'unban'
  account: string
  nickname: string
}

const banModal = ref<BanModalState | null>(null)

interface PromoteModalState {
  id: number
  account: string
  nickname: string
  error: string
}

const promoteModal = ref<PromoteModalState | null>(null)

interface RenameModalState {
  id: number
  account: string
  nickname: string
  error: string
}

const renameModal = ref<RenameModalState | null>(null)

interface BatchBanModalState {
  ids: number[]
  accounts: string[]
  count: number
}

const batchBanModal = ref<BatchBanModalState | null>(null)

const createAdminModal = reactive({
  open: false,
  account: '',
  nickname: '',
  wechat: '',
  error: '',
})

const passwordModal = reactive({
  open: false,
  account: '',
  nickname: '',
  newPassword: '',
  confirmPassword: '',
  error: '',
})

const applicationModal = reactive({
  open: false,
})

const rejectApplicationModal = reactive({
  open: false,
  id: 0,
  account: '',
  nickname: '',
  reason: '',
  error: '',
})

const applyFilters = () => {
  feedbackModal.value = { message: '筛选条件已应用。' }
}

const openApplicationModal = () => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  applicationModal.open = true
}

const closeApplicationModal = () => {
  applicationModal.open = false
}

const approveApplication = (id: number) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  const reviewer = authStore.currentAccount?.account || '管理员'
  const result = adminStore.approveAdminApplication(id, reviewer)
  feedbackModal.value = { message: result.message }
}

const openRejectApplicationModal = (application: AdminApplication) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  rejectApplicationModal.open = true
  rejectApplicationModal.id = application.id
  rejectApplicationModal.account = application.account
  rejectApplicationModal.nickname = application.nickname
  rejectApplicationModal.reason = ''
  rejectApplicationModal.error = ''
}

const closeRejectApplicationModal = () => {
  rejectApplicationModal.open = false
  rejectApplicationModal.reason = ''
  rejectApplicationModal.error = ''
}

const submitRejectApplication = () => {
  if (!canManage.value) {
    showPermissionWarning()
    rejectApplicationModal.open = false
    return
  }

  const note = rejectApplicationModal.reason.trim()
  if (!note) {
    rejectApplicationModal.error = '请填写驳回原因。'
    return
  }

  const reviewer = authStore.currentAccount?.account || '管理员'
  const result = adminStore.rejectAdminApplication(
    rejectApplicationModal.id,
    reviewer,
    note,
  )

  if (!result.success) {
    rejectApplicationModal.error = result.message
    return
  }

  feedbackModal.value = { message: result.message }
  closeRejectApplicationModal()
}

const openDetail = (id: number) => {
  const user = users.value.find((item: AdminUser) => item.id === id)
  if (!user) {
    feedbackModal.value = { message: '未找到该用户。' }
    return
  }
  detailModal.value = { ...user }
}

const closeDetailModal = () => {
  detailModal.value = null
}

const openBanModal = (user: AdminUser) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  banModal.value = {
    id: user.id,
    mode: user.status === 'banned' ? 'unban' : 'ban',
    account: user.account,
    nickname: user.nickname,
  }
}

const closeBanModal = () => {
  banModal.value = null
}

const submitBanModal = () => {
  const modal = banModal.value
  if (!modal) return
  if (!canManage.value) {
    showPermissionWarning()
    banModal.value = null
    return
  }
  const success = adminStore.toggleUserBan(modal.id)
  banModal.value = null
  feedbackModal.value = {
    message: success
      ? modal.mode === 'ban'
        ? `已封禁用户 ${modal.account}。`
        : `已解除用户 ${modal.account} 的封禁。`
      : '未找到该用户。',
  }
}

const openPromoteModal = (user: AdminUser) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  promoteModal.value = {
    id: user.id,
    account: user.account,
    nickname: user.nickname,
    error: '',
  }
}

const closePromoteModal = () => {
  promoteModal.value = null
}

const openRenameModal = (user: AdminUser) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  renameModal.value = {
    id: user.id,
    account: user.account,
    nickname: user.nickname,
    error: '',
  }
}

const closeRenameModal = () => {
  renameModal.value = null
}

const submitRenameModal = () => {
  const modal = renameModal.value
  if (!modal) return
  if (!canManage.value) {
    showPermissionWarning()
    renameModal.value = null
    return
  }

  const nextNickname = modal.nickname.trim()
  if (!nextNickname) {
    modal.error = '昵称不能为空。'
    return
  }

  const success = adminStore.updateUserNickname(modal.id, nextNickname)
  if (!success) {
    modal.error = '修改失败，请稍后重试。'
    return
  }

  authStore.updateAccountNickname(modal.account, nextNickname)
  renameModal.value = null
  feedbackModal.value = { message: `已更新用户 ${modal.account} 的昵称。` }
}

const submitPromoteModal = () => {
  const modal = promoteModal.value
  if (!modal) return
  if (!canManage.value) {
    showPermissionWarning()
    promoteModal.value = null
    return
  }
  const nickname = modal.nickname.trim()
  const success = adminStore.promoteToAdmin(modal.id, nickname || undefined)
  if (!success) {
    modal.error = '未找到该用户。'
    return
  }
  authStore.updateAccountNickname(modal.account, nickname || modal.account)
  promoteModal.value = null
  feedbackModal.value = { message: `用户 ${modal.account} 已升级为管理员。` }
}

const openPasswordModal = (user: AdminUser) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  passwordModal.open = true
  passwordModal.account = user.account
  passwordModal.nickname = user.nickname
  passwordModal.newPassword = ''
  passwordModal.confirmPassword = ''
  passwordModal.error = ''
}

const closePasswordModal = () => {
  passwordModal.open = false
}

const submitPasswordModal = () => {
  if (!canManage.value) {
    showPermissionWarning()
    passwordModal.open = false
    return
  }

  if (!passwordModal.newPassword.trim()) {
    passwordModal.error = '请输入新密码。'
    return
  }

  if (passwordModal.newPassword.trim().length < 6) {
    passwordModal.error = '新密码长度不能少于 6 位。'
    return
  }

  if (passwordModal.newPassword !== passwordModal.confirmPassword) {
    passwordModal.error = '两次输入的密码不一致。'
    return
  }

  const result = authStore.resetPassword(passwordModal.account, passwordModal.newPassword)
  if (!result.success) {
    passwordModal.error = result.message
    return
  }

  passwordModal.open = false
  feedbackModal.value = { message: result.message ?? `已更新账号 ${passwordModal.account} 的密码。` }
}

const openBatchBanModal = () => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  const targets = filteredUsers.value.filter(
    (user) => user.role !== 'admin' && user.status !== 'banned',
  )
  if (targets.length === 0) {
    feedbackModal.value = { message: '没有可封禁的普通用户。' }
    return
  }
  batchBanModal.value = {
    ids: targets.map((item) => item.id),
    accounts: targets.slice(0, 6).map((item) => item.account),
    count: targets.length,
  }
}

const closeBatchBanModal = () => {
  batchBanModal.value = null
}

const submitBatchBan = () => {
  const modal = batchBanModal.value
  if (!modal) return
  if (!canManage.value) {
    showPermissionWarning()
    batchBanModal.value = null
    return
  }
  modal.ids.forEach((id) => {
    const user = users.value.find((item) => item.id === id)
    if (user && user.status !== 'banned') {
      adminStore.toggleUserBan(id)
    }
  })
  batchBanModal.value = null
  feedbackModal.value = { message: '已封禁选定的普通用户。' }
}

const openCreateAdminModal = () => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  createAdminModal.open = true
  createAdminModal.account = ''
  createAdminModal.nickname = ''
  createAdminModal.wechat = ''
  createAdminModal.error = ''
}

const closeCreateAdminModal = () => {
  createAdminModal.open = false
}

const submitCreateAdmin = () => {
  if (!canManage.value) {
    showPermissionWarning()
    createAdminModal.open = false
    return
  }
  const account = createAdminModal.account.trim()
  if (!account) {
    createAdminModal.error = '账号不能为空。'
    return
  }

  const exists = users.value.some((user) => user.account === account)
  if (exists) {
    createAdminModal.error = '账号已存在，请换一个账号。'
    return
  }

  const nickname = createAdminModal.nickname.trim() || account
  const wechat = createAdminModal.wechat.trim()

  adminStore.createAdminAccount({ account, nickname, wechat: wechat || undefined })
  authStore.addAdminAccount(account, nickname, '123456')
  createAdminModal.open = false
  feedbackModal.value = { message: '管理员账号创建成功。' }
}

const closeFeedback = () => {
  feedbackModal.value = null
}
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.page-head h2 {
  margin: 0 0 4px;
}

.subtitle {
  margin: 0;
  color: #6b7280;
}

.actions {
  display: flex;
  gap: 12px;
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

.btn.danger {
  border: none;
  background: #ef4444;
  color: #fff;
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.55;
  border-color: #e5e7eb;
  color: #9ca3af;
  background: #f3f4f6;
}

.btn.primary:disabled {
  background: #93c5fd;
  color: #fff;
}

.btn.danger:disabled {
  background: #fca5a5;
  color: #fff;
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.55;
  border-color: #e5e7eb;
  color: #9ca3af;
  background: #f3f4f6;
}

.btn.primary:disabled {
  background: #93c5fd;
  color: #fff;
}

.filters {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
}

.input {
  width: 100%;
  padding: 10px 12px;
  border-radius: 6px;
  border: 1px solid #d1d5db;
}

.table {
  width: 100%;
  border-collapse: collapse;
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e5e7eb;
}

.table th,
.table td {
  padding: 14px 16px;
  border-bottom: 1px solid #e5e7eb;
  text-align: left;
  font-size: 14px;
}

.table thead th {
  background: #f9fafb;
  font-weight: 600;
}

.title {
  margin: 0 0 6px;
  font-weight: 600;
}

.meta {
  margin: 0;
  color: #6b7280;
  font-size: 12px;
}

.tag {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 9999px;
  background: rgba(59, 130, 246, 0.18);
  color: #1d4ed8;
}

.tag.admin {
  background: rgba(124, 58, 237, 0.18);
  color: #6d28d9;
}

.status {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 9999px;
  background: #e5e7eb;
  color: #374151;
}

.status.active {
  background: rgba(16, 185, 129, 0.18);
  color: #047857;
}

.status.pending {
  background: rgba(250, 204, 21, 0.18);
  color: #b45309;
}

.status.banned {
  background: rgba(239, 68, 68, 0.18);
  color: #b91c1c;
}

.ops {
  display: flex;
  gap: 12px;
}

.link {
  border: none;
  background: transparent;
  color: #2563eb;
  cursor: pointer;
  padding: 0;
}

.link:disabled {
  cursor: not-allowed;
  color: #9ca3af;
}

.application-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.application-section.history {
  border-top: 1px solid #e5e7eb;
  padding-top: 16px;
  margin-top: 8px;
}

.application-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  color: #6b7280;
  font-size: 13px;
}

.application-header h4 {
  margin: 0;
  font-size: 15px;
  color: #111827;
}

.application-list,
.history-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.application-item,
.history-item {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  border-radius: 10px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
}

.application-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.application-title,
.history-title {
  margin: 0;
  font-weight: 600;
  color: #1f2937;
}

.application-reason {
  margin: 0;
  font-size: 13px;
  color: #4b5563;
  line-height: 1.4;
}

.application-time,
.history-summary,
.history-meta {
  margin: 0;
  font-size: 12px;
  color: #6b7280;
}

.application-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.application-empty {
  text-align: center;
  color: #9ca3af;
  font-size: 13px;
  padding: 12px 0;
}

.btn.ghost {
  background: transparent;
  border-color: transparent;
  color: #4b5563;
}

.modal-form {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
}

.modal-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.modal-label {
  font-size: 14px;
  color: #374151;
  font-weight: 500;
}

.modal-input {
  width: 100%;
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px solid #d1d5db;
  font-size: 14px;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.modal-input:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
  outline: none;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  grid-column: 1 / -1;
  margin-top: 4px;
}

.modal-error {
  margin: -4px 0 0;
  color: #dc2626;
  font-size: 13px;
  grid-column: 1 / -1;
}

.modal-message {
  margin: 0 0 12px;
  font-size: 15px;
  color: #374151;
}

.modal-list {
  margin: 0 0 12px;
  padding-left: 18px;
  color: #4b5563;
  font-size: 14px;
  max-height: 160px;
  overflow-y: auto;
}

.detail-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px 18px;
  margin: 0 0 16px;
  padding: 0;
}

.detail-list div {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-list dt {
  font-size: 13px;
  color: #6b7280;
}

.detail-list dd {
  margin: 0;
  font-size: 15px;
  color: #111827;
  word-break: break-all;
}
</style>
