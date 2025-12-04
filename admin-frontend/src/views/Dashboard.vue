<template>
  <div class="dashboard">
    <section class="cards">
      <article class="stat-card" v-for="card in overviewCards" :key="card.title">
        <div class="stat-header">
          <span class="stat-title">{{ card.title }}</span>
          <span class="stat-trend" :class="card.trend >= 0 ? 'up' : 'down'">
            {{ card.trend >= 0 ? '+' : '' }}{{ card.trend }}%
          </span>
        </div>
        <p class="stat-value">{{ card.value }}</p>
        <p class="stat-subtitle">{{ card.subtitle }}</p>
      </article>
    </section>

    <section class="section sales-trend">
      <header class="section-head">
        <div>
          <h2>出售趋势</h2>
          <p class="section-subtitle">最近 7 天成交数量</p>
        </div>
        <div class="trend-meta">
          <span>累计：{{ saleTrendTotal }} 件</span>
          <span>日均：{{ saleTrendAverage }} 件</span>
        </div>
      </header>
      <div class="chart-wrapper">
        <svg
          class="chart"
          role="img"
          :viewBox="`0 0 ${chartDimensions.width} ${chartDimensions.height}`"
          aria-label="近 7 天出售趋势"
        >
          <defs>
            <linearGradient id="saleTrendGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#2563eb" stop-opacity="0.32" />
              <stop offset="100%" stop-color="#2563eb" stop-opacity="0" />
            </linearGradient>
          </defs>
          <path v-if="chartAreaPath" class="chart-area" :d="chartAreaPath" fill="url(#saleTrendGradient)" />
          <path v-if="chartLinePath" class="chart-line" :d="chartLinePath" />
          <g class="chart-points">
            <circle
              v-for="point in chartPoints"
              :key="point.label"
              class="chart-point"
              :cx="point.x"
              :cy="point.y"
              r="4"
            >
              <title>{{ point.label }}：{{ point.value }} 件</title>
            </circle>
          </g>
        </svg>
        <ul class="chart-axis">
          <li v-for="label in chartLabels" :key="label.label" :style="{ left: `${label.x}px` }">
            <span class="axis-dot"></span>
            <span class="axis-label">{{ label.label }}</span>
          </li>
        </ul>
      </div>
    </section>

    <section class="section">
      <header class="section-head">
        <h2>待处理审核</h2>
        <button type="button" class="link-btn" @click="handleViewAllPending">查看全部</button>
      </header>
      <ul class="review-list">
        <li v-for="item in pendingReviews" :key="item.id" class="review-item">
          <div>
            <p class="review-title">{{ item.title }}</p>
            <p class="review-meta">卖家：{{ item.seller }} · 提交时间：{{ item.submittedAt }}</p>
          </div>
          <div class="review-actions">
            <button
              type="button"
              class="action approve"
              :disabled="!canManage"
              @click="openApproveModal(item.id)"
            >
              通过
            </button>
            <button
              type="button"
              class="action reject"
              :disabled="!canManage"
              @click="openRejectModal(item.id)"
            >
              驳回
            </button>
          </div>
        </li>
      </ul>
    </section>

    <section class="section">
      <header class="section-head">
        <h2>活跃情况</h2>
        <span class="section-subtitle">最近 7 天</span>
      </header>
      <div class="activity-grid">
        <div class="activity-item" v-for="metric in activity" :key="metric.label">
          <span class="activity-label">{{ metric.label }}</span>
          <strong class="activity-value">{{ metric.value }}</strong>
          <span class="activity-note">{{ metric.note }}</span>
        </div>
      </div>
    </section>
  </div>

  <BaseModal
    :open="Boolean(reviewModal)"
    :title="reviewModal?.mode === 'approve' ? '通过审核' : '驳回审核'"
    width="480px"
    @close="closeReviewModal"
  >
    <template v-if="reviewModal">
      <p class="modal-message">
        {{ reviewModal.mode === 'approve' ? '确认通过该商品审核？' : '请填写驳回原因：' }}
      </p>
      <textarea
        v-if="reviewModal.mode === 'reject'"
        v-model="reviewModal.reason"
        class="modal-textarea"
        rows="4"
        placeholder="例如：商品图片模糊或描述不完整"
      ></textarea>
      <p v-if="reviewError" class="modal-error">{{ reviewError }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeReviewModal">取消</button>
        <button type="button" class="btn primary" @click="submitReviewModal">
          {{ reviewModal.mode === 'approve' ? '确认通过' : '确认驳回' }}
        </button>
      </div>
    </template>
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
import { computed, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'

import { useAdminDataStore } from '../stores/admin'
import { useAuthStore } from '../stores/auth'
import type { AdminProduct, AdminUser } from '../stores/admin'
import BaseModal from '../components/BaseModal.vue'

interface OverviewCard {
  title: string
  value: string
  subtitle: string
  trend: number
}

interface ActivityMetric {
  label: string
  value: string
  note: string
}

interface TrendPoint {
  label: string
  value: number
}

interface ChartPoint {
  x: number
  y: number
  label: string
  value: number
}

const chartDimensions = {
  width: 520,
  height: 220,
  paddingX: 32,
  paddingY: 24,
} as const

defineOptions({ name: 'AdminDashboardView' })

const router = useRouter()

const adminStore = useAdminDataStore()
const authStore = useAuthStore()
const { saleTrend, pendingReviews, products, users } = storeToRefs(adminStore)
const { canManage } = storeToRefs(authStore)

const saleTrendTotal = computed(() =>
  saleTrend.value.reduce((acc: number, item: TrendPoint) => acc + item.value, 0),
)

const saleTrendAverage = computed(() =>
  saleTrend.value.length === 0
    ? 0
    : Math.round((saleTrendTotal.value / saleTrend.value.length) * 10) / 10,
)

const chartPoints = computed<ChartPoint[]>(() => {
  if (saleTrend.value.length === 0) {
    return []
  }

  const usableWidth = chartDimensions.width - chartDimensions.paddingX * 2
  const usableHeight = chartDimensions.height - chartDimensions.paddingY * 2
  const stepX = saleTrend.value.length === 1 ? 0 : usableWidth / (saleTrend.value.length - 1)
  const maxValue = Math.max(...saleTrend.value.map((item: TrendPoint) => item.value)) || 1

  return saleTrend.value.map((item: TrendPoint, index: number) => {
    const x = chartDimensions.paddingX + index * stepX
    const normalized = item.value / maxValue
    const y =
      chartDimensions.height - chartDimensions.paddingY - normalized * usableHeight
    return { x, y, label: item.label, value: item.value }
  })
})

const chartLinePath = computed(() => {
  if (chartPoints.value.length === 0) {
    return ''
  }

  return chartPoints.value
    .map((point, index) => `${index === 0 ? 'M' : 'L'} ${point.x} ${point.y}`)
    .join(' ')
})

const chartAreaPath = computed(() => {
  const points = chartPoints.value
  if (points.length === 0) {
    return ''
  }

  const baseY = chartDimensions.height - chartDimensions.paddingY
  const startPoint = points[0]!
  const endPoint = points[points.length - 1]!

  const lines = points
    .map((point) => `L ${point.x} ${point.y}`)
    .join(' ')

  return `M ${startPoint.x} ${baseY} ${lines} L ${endPoint.x} ${baseY} Z`
})

const chartLabels = computed(() =>
  chartPoints.value.map((point) => ({ x: point.x, label: point.label, value: point.value })),
)

const recentWindow = 30 * 24 * 60 * 60 * 1000

const overviewCards = computed<OverviewCard[]>(() => {
  const now = Date.now()
  const pendingCount = pendingReviews.value.length
  const soldRevenue = products.value
    .filter((item: AdminProduct) => item.status === 'sold')
    .reduce((total: number, item: AdminProduct) => total + item.price, 0)
  const newUserCount = users.value.filter((user: AdminUser) => {
    const date = Date.parse(user.registeredAt)
    if (Number.isNaN(date)) return false
    return now - date <= recentWindow
  }).length
  const bannedCount = users.value.filter((user: AdminUser) => user.status === 'banned').length

  return [
    {
      title: '待审核商品',
      value: String(pendingCount),
      subtitle: '需要管理员跟进的商品',
      trend: pendingCount === 0 ? -100 : Math.min(99, pendingCount * 5),
    },
    {
      title: '近 30 天新用户',
      value: String(newUserCount),
      subtitle: '最近 30 天注册用户数',
      trend: newUserCount * 6,
    },
    {
      title: '成交总额',
      value: `￥${soldRevenue.toLocaleString('zh-CN')}`,
      subtitle: '已售商品成交额（示例数据）',
      trend: soldRevenue > 0 ? 12.6 : 0,
    },
    {
      title: '封禁账号',
      value: String(bannedCount),
      subtitle: '当前被封禁的用户数',
      trend: -bannedCount * 8,
    },
  ]
})

const handleViewAllPending = () => {
  router.push({ path: '/admin/products', query: { status: 'pending' } })
}

const reviewModal = ref<{ id: number; mode: 'approve' | 'reject'; reason: string } | null>(null)
const reviewError = ref('')
const feedbackModal = ref<{ message: string } | null>(null)

const showPermissionWarning = () => {
  feedbackModal.value = { message: '当前账号仅支持查看数据，如需操作请联系管理员。' }
}

const openApproveModal = (id: number) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  reviewError.value = ''
  reviewModal.value = { id, mode: 'approve', reason: '' }
}

const openRejectModal = (id: number) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  reviewError.value = ''
  reviewModal.value = { id, mode: 'reject', reason: '' }
}

const closeReviewModal = () => {
  reviewModal.value = null
}

const submitReviewModal = () => {
  const modal = reviewModal.value
  if (!modal) return

  if (!canManage.value) {
    showPermissionWarning()
    reviewModal.value = null
    return
  }

  if (modal.mode === 'reject') {
    const reason = modal.reason.trim()
    if (!reason) {
      reviewError.value = '请填写驳回原因。'
      return
    }
    const success = adminStore.rejectPendingReview(modal.id)
    reviewModal.value = null
    feedbackModal.value = {
      message: success ? `已驳回该商品：${reason}` : '未找到该审核项。',
    }
    return
  }

  const success = adminStore.approvePendingReview(modal.id)
  reviewModal.value = null
  feedbackModal.value = {
    message: success ? '已通过该商品审核，并已移至在售列表。' : '未找到该审核项。',
  }
}

const closeFeedback = () => {
  feedbackModal.value = null
}

const activity: ActivityMetric[] = [
  { label: '发布商品', value: '134', note: '较上一周期 +12%' },
  { label: '发起聊天', value: '482', note: '保持稳定' },
  { label: '平台成交', value: '68', note: '同比提升 8%' },
  { label: '主动下架', value: '27', note: '需关注原因' },
]
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.stat-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.08);
  border: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-title {
  font-size: 14px;
  color: #6b7280;
}

.stat-trend {
  font-size: 13px;
  padding: 4px 8px;
  border-radius: 9999px;
  background: rgba(16, 185, 129, 0.16);
  color: #047857;
}

.stat-trend.down {
  background: rgba(239, 68, 68, 0.16);
  color: #b91c1c;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  margin: 0;
}

.stat-subtitle {
  font-size: 13px;
  color: #6b7280;
  margin: 0;
}

.section {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #e5e7eb;
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
}

.section-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-head h2 {
  margin: 0;
  font-size: 18px;
}

.section-subtitle {
  font-size: 13px;
  color: #6b7280;
}

.sales-trend .section-head {
  align-items: flex-start;
}

.sales-trend .trend-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
  font-size: 13px;
  color: #6b7280;
}

@media (max-width: 640px) {
  .sales-trend .section-head {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .sales-trend .trend-meta {
    align-items: flex-start;
  }
}

.chart-wrapper {
  position: relative;
  max-width: 520px;
  width: 100%;
  padding-bottom: 36px;
}

.chart {
  width: 100%;
  height: auto;
}

.chart-area {
  fill: url(#saleTrendGradient);
}

.chart-line {
  fill: none;
  stroke: #2563eb;
  stroke-width: 3;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.chart-point {
  fill: #2563eb;
  stroke: #ffffff;
  stroke-width: 2;
}

.chart-axis {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  margin: 0;
  padding: 0;
  list-style: none;
  height: 36px;
}

.chart-axis li {
  position: absolute;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  transform: translateX(-50%);
}

.axis-dot {
  display: block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #93c5fd;
}

.axis-label {
  font-size: 12px;
  color: #6b7280;
}

.link-btn {
  border: none;
  background: transparent;
  color: #2563eb;
  cursor: pointer;
  font-size: 14px;
}

.review-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.review-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-radius: 10px;
  background: #f9fafb;
}

.review-title {
  margin: 0 0 4px;
  font-weight: 600;
}

.review-meta {
  margin: 0;
  font-size: 13px;
  color: #6b7280;
}

.review-actions {
  display: flex;
  gap: 8px;
}

.action {
  border: none;
  border-radius: 6px;
  padding: 8px 14px;
  cursor: pointer;
  font-size: 14px;
}

.action.approve {
  background: #22c55e;
  color: #fff;
}

.action.reject {
  background: #ef4444;
  color: #fff;
}

.action:disabled {
  cursor: not-allowed;
  opacity: 0.55;
  background: #e5e7eb;
  color: #9ca3af;
}

.activity-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
}

.activity-item {
  padding: 16px;
  border-radius: 10px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
}

.activity-label {
  display: block;
  font-size: 13px;
  color: #6b7280;
}

.activity-value {
  display: block;
  font-size: 22px;
  margin: 8px 0 4px;
}

.activity-note {
  font-size: 12px;
  color: #9ca3af;
}

.modal-message {
  margin: 0 0 12px;
  font-size: 15px;
  color: #374151;
}

.modal-textarea {
  width: 100%;
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px solid #d1d5db;
  font-size: 14px;
  margin-bottom: 8px;
  resize: vertical;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.modal-error {
  color: #dc2626;
  font-size: 13px;
  margin: -4px 0 8px;
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
</style>
