<template>
  <div class="page">
    <header class="page-head">
      <div>
        <h2>商品管理</h2>
        <p class="subtitle">查看在售商品与审核状态，支持快速筛选与批量操作。</p>
      </div>
      <div class="actions">
        <button type="button" class="btn secondary" @click="handleExport">导出报表</button>
        <button
          type="button"
          class="btn primary"
          @click="openCreateProductModal"
        >
          新增商品
        </button>
      </div>
    </header>

    <section class="filters">
      <input v-model="keyword" class="input" type="search" placeholder="搜索商品标题或卖家账号" />
      <select v-model="status" class="input">
        <option value="">全部状态</option>
        <option value="pending">待审核</option>
        <option value="for-sale">在售中</option>
        <option value="sold">已售出</option>
        <option value="delisted">已下架</option>
      </select>
      <select v-model="categoryId" class="input">
        <option value="">全部分类</option>
        <option v-for="item in categoryOptions" :key="item.id" :value="String(item.id)">
          {{ item.name }}
        </option>
      </select>
      <button type="button" class="btn" @click="applyFilters">筛选</button>
    </section>

    <table class="table">
      <thead>
        <tr>
          <th>
            <input type="checkbox" v-model="selectAll" @change="toggleSelectAll" />
          </th>
          <th>商品信息</th>
          <th>卖家</th>
          <th>状态</th>
          <th>发布时间</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in filteredProducts" :key="row.id">
          <td>
            <input type="checkbox" v-model="selected" :value="row.id" />
          </td>
          <td>
            <p class="title">{{ row.title }}</p>
            <p class="meta">分类：{{ row.categoryName }} · 价格：￥{{ row.price }}</p>
          </td>
          <td>
            <p class="title">{{ row.seller }}</p>
            <p class="meta">微信：{{ row.wechat }}</p>
          </td>
          <td>
            <span class="status" :class="row.status">{{ statusMap[row.status] }}</span>
          </td>
          <td>{{ row.publishedAt }}</td>
          <td class="ops">
            <button type="button" class="link" @click="openDetail(row.id)">详情</button>
            <button
              type="button"
              class="link"
              :disabled="!canManage"
              @click="openApproveModal(row.id)"
            >
              通过
            </button>
            <button
              type="button"
              class="link danger"
              :disabled="!canManage"
              @click="openRejectModal(row.id)"
            >
              驳回
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>

  <BaseModal
    :open="Boolean(detailModal)"
    title="商品详情"
    width="540px"
    @close="closeDetailModal"
  >
    <template v-if="detailModal">
      <dl class="detail-list">
        <div>
          <dt>标题</dt>
          <dd>{{ detailModal.title }}</dd>
        </div>
        <div>
          <dt>分类</dt>
          <dd>{{ detailModal.categoryName }}</dd>
        </div>
        <div>
          <dt>价格</dt>
          <dd>￥{{ detailModal.price }}</dd>
        </div>
        <div>
          <dt>卖家</dt>
          <dd>{{ detailModal.seller }}</dd>
        </div>
        <div>
          <dt>微信</dt>
          <dd>{{ detailModal.wechat }}</dd>
        </div>
        <div>
          <dt>状态</dt>
          <dd>{{ statusMap[detailModal.status] }}</dd>
        </div>
        <div>
          <dt>发布时间</dt>
          <dd>{{ detailModal.publishedAt }}</dd>
        </div>
      </dl>
      <div class="modal-actions">
        <button type="button" class="btn primary" @click="closeDetailModal">关闭</button>
      </div>
    </template>
  </BaseModal>

  <BaseModal
    :open="Boolean(actionModal)"
    :title="actionModal?.mode === 'approve' ? '通过商品审核' : '驳回商品审核'"
    width="500px"
    @close="closeActionModal"
  >
    <template v-if="actionModal">
      <p class="modal-message">
        {{ actionModal.mode === 'approve' ? '确认将该商品设为在售状态？' : '请填写驳回原因：' }}
      </p>
      <textarea
        v-if="actionModal.mode === 'reject'"
        v-model="actionModal.reason"
        class="modal-textarea"
        rows="4"
        placeholder="例如：商品描述不完整或图片不清晰"
      ></textarea>
      <p v-if="actionError" class="modal-error">{{ actionError }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeActionModal">取消</button>
        <button type="button" class="btn primary" :disabled="!canManage" @click="submitActionModal">
          {{ actionModal.mode === 'approve' ? '确认通过' : '确认驳回' }}
        </button>
      </div>
    </template>
  </BaseModal>

  <BaseModal
    :open="showCreateModal"
    title="新增商品"
    width="540px"
    @close="closeCreateProductModal"
  >
    <form class="modal-form" @submit.prevent="submitCreateProduct">
      <div class="modal-field">
        <label class="modal-label" for="create-title-input">商品标题</label>
        <input
          id="create-title-input"
          v-model="createProductForm.title"
          class="modal-input"
          placeholder="请输入商品标题"
        />
      </div>
      <div class="modal-field">
        <label class="modal-label" for="create-category-input">所属分类</label>
        <select
          id="create-category-input"
          v-model="createProductForm.categoryId"
          class="modal-input"
        >
          <option value="">请选择分类</option>
          <option v-for="item in categoryOptions" :key="item.id" :value="String(item.id)">
            {{ item.name }}
          </option>
        </select>
      </div>
      <div class="modal-field">
        <label class="modal-label" for="create-price-input">价格（元）</label>
        <input
          id="create-price-input"
          v-model="createProductForm.price"
          class="modal-input"
          type="number"
          min="0"
          step="0.01"
        />
      </div>
      <div class="modal-field">
        <label class="modal-label" for="create-seller-input">卖家账号</label>
        <input
          id="create-seller-input"
          v-model="createProductForm.seller"
          class="modal-input"
          :readonly="!canManage && Boolean(authStore.currentAccount?.account)"
          placeholder="请输入卖家账号"
        />
      </div>
      <div class="modal-field">
        <label class="modal-label" for="create-wechat-input">卖家微信（可选）</label>
        <input
          id="create-wechat-input"
          v-model="createProductForm.wechat"
          class="modal-input"
          placeholder="请输入微信号"
        />
      </div>
      <p v-if="createProductForm.error" class="modal-error">{{ createProductForm.error }}</p>
      <div class="modal-actions">
        <button type="button" class="btn ghost" @click="closeCreateProductModal">取消</button>
        <button type="submit" class="btn primary" :disabled="!canManage">创建</button>
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
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute } from 'vue-router'

import BaseModal from '../components/BaseModal.vue'
import { useAdminDataStore } from '../stores/admin'
import { useAuthStore } from '../stores/auth'
import type {
  AdminProduct,
  AdminCategory,
  PendingReviewItem,
  ProductStatus,
} from '../stores/admin'

defineOptions({ name: 'AdminProductListView' })

interface DisplayProductRow extends AdminProduct {
  categoryName: string
}

const keyword = ref('')
const status = ref<ProductStatus | ''>('')
const categoryId = ref('')
const selectAll = ref(false)
const selected = ref<number[]>([])

const adminStore = useAdminDataStore()
const authStore = useAuthStore()
const { products, categories, pendingReviews, categoryMap } = storeToRefs(adminStore)
const { canManage } = storeToRefs(authStore)

const allowedStatuses: ProductStatus[] = ['pending', 'for-sale', 'sold', 'delisted']
const isValidStatus = (value: string): value is ProductStatus =>
  allowedStatuses.includes(value as ProductStatus)

const categoryOptions = computed(() =>
  categories.value.map((item: AdminCategory) => ({ id: item.id, name: item.name })),
)

const statusMap: Record<ProductStatus, string> = {
  pending: '待审核',
  'for-sale': '在售中',
  sold: '已售出',
  delisted: '已下架',
}

const filteredProducts = computed<DisplayProductRow[]>(() => {
  const keywordValue = keyword.value.trim()
  const statusFilter = status.value
  const categoryFilter = categoryId.value

  return products.value
    .filter((product: AdminProduct) => {
      const matchKeyword = keywordValue
        ? [product.title, product.seller].some((text) => text.includes(keywordValue))
        : true
      const matchStatus = statusFilter ? product.status === statusFilter : true
      const matchCategory = categoryFilter ? String(product.categoryId) === categoryFilter : true
      return matchKeyword && matchStatus && matchCategory
    })
    .map((product: AdminProduct) => ({
      ...product,
      categoryName: categoryMap.value.get(product.categoryId)?.name ?? '未分类',
    }))
})

const route = useRoute()

onMounted(() => {
  if (typeof route.query.status === 'string' && isValidStatus(route.query.status)) {
    status.value = route.query.status
  }
})

watch(
  () => route.query.status,
  (value) => {
    if (typeof value === 'string' && isValidStatus(value)) {
      status.value = value
    } else if (value === undefined) {
      status.value = ''
    }
  },
)

watch([keyword, status, categoryId], () => {
  selectAll.value = false
  selected.value = []
})

const feedbackModal = ref<{ message: string } | null>(null)

const showPermissionWarning = () => {
  feedbackModal.value = { message: '当前账号仅支持查看数据，如需操作请联系管理员。' }
}

const closeFeedback = () => {
  feedbackModal.value = null
}

const detailModal = ref<DisplayProductRow | null>(null)
const actionModal = ref<{ id: number; mode: 'approve' | 'reject'; reason: string } | null>(null)
const actionError = ref('')
const showCreateModal = ref(false)
const createProductForm = reactive({
  title: '',
  categoryId: '',
  price: '',
  seller: '',
  wechat: '',
  error: '',
})

const applyFilters = () => {
  feedbackModal.value = { message: '筛选条件已应用，列表已更新。' }
}

const toggleSelectAll = () => {
  if (selectAll.value) {
    selected.value = filteredProducts.value.map((row) => row.id)
  } else {
    selected.value = []
  }
}

const openDetail = (id: number) => {
  const product = products.value.find((item: AdminProduct) => item.id === id)
  if (!product) {
    feedbackModal.value = { message: '未找到该商品。' }
    return
  }
  detailModal.value = {
    ...product,
    categoryName: categoryMap.value.get(product.categoryId)?.name ?? '未分类',
  }
}

const closeDetailModal = () => {
  detailModal.value = null
}

const openApproveModal = (id: number) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  const product = products.value.find((item: AdminProduct) => item.id === id)
  if (!product) {
    feedbackModal.value = { message: '未找到该商品。' }
    return
  }
  actionError.value = ''
  actionModal.value = { id: product.id, mode: 'approve', reason: '' }
}

const openRejectModal = (id: number) => {
  if (!canManage.value) {
    showPermissionWarning()
    return
  }
  const product = products.value.find((item: AdminProduct) => item.id === id)
  if (!product) {
    feedbackModal.value = { message: '未找到该商品。' }
    return
  }
  actionError.value = ''
  actionModal.value = { id: product.id, mode: 'reject', reason: '' }
}

const closeActionModal = () => {
  actionModal.value = null
  actionError.value = ''
}

const submitActionModal = () => {
  const modal = actionModal.value
  if (!modal) return

  if (!canManage.value) {
    showPermissionWarning()
    actionModal.value = null
    return
  }

  if (modal.mode === 'reject') {
    const reason = modal.reason.trim()
    if (!reason) {
      actionError.value = '请填写驳回原因。'
      return
    }
    const review = pendingReviews.value.find((item: PendingReviewItem) => item.productId === modal.id)
    const success = review
      ? adminStore.rejectPendingReview(review.id)
      : adminStore.toggleProductStatus(modal.id, 'delisted')
    actionModal.value = null
    feedbackModal.value = {
      message: success ? `已驳回该商品：${reason}` : '未找到可驳回的商品。',
    }
    return
  }

  const review = pendingReviews.value.find((item: PendingReviewItem) => item.productId === modal.id)
  const success = review
    ? adminStore.approvePendingReview(review.id)
    : adminStore.toggleProductStatus(modal.id, 'for-sale')
  actionModal.value = null
  feedbackModal.value = {
    message: success ? '商品已通过审核并上线。' : '未找到可审核的商品。',
  }
}

const handleExport = () => {
  const ids = selected.value.length > 0 ? selected.value : null
  const csv = adminStore.exportProducts(ids)
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  const timestamp = new Date().toISOString().replace(/[-:T]/g, '').slice(0, 12)
  link.download = `admin-products-${timestamp}.csv`
  link.click()
  URL.revokeObjectURL(url)
}

const openCreateProductModal = () => {
  if (categoryOptions.value.length === 0) {
    feedbackModal.value = { message: '暂无分类，请先创建分类。' }
    return
  }
  createProductForm.title = ''
  createProductForm.categoryId = String(categoryOptions.value[0]?.id ?? '')
  createProductForm.price = ''
  createProductForm.seller = canManage.value
    ? ''
    : authStore.currentAccount?.account ?? ''
  createProductForm.wechat = ''
  createProductForm.error = ''
  showCreateModal.value = true
}

const closeCreateProductModal = () => {
  showCreateModal.value = false
}

const submitCreateProduct = () => {
  const title = createProductForm.title.trim()
  if (!title) {
    createProductForm.error = '商品标题不能为空。'
    return
  }

  const categoryValue = Number(createProductForm.categoryId)
  if (!categoryOptions.value.some((item) => item.id === categoryValue)) {
    createProductForm.error = '请选择合法的分类。'
    return
  }

  const price = Number(createProductForm.price)
  if (!Number.isFinite(price) || price <= 0) {
    createProductForm.error = '请输入合法的价格。'
    return
  }

  const seller = createProductForm.seller.trim()
  if (!seller) {
    createProductForm.error = '卖家账号不能为空。'
    return
  }

  createProductForm.error = ''
  const wechat = createProductForm.wechat.trim()

  const product = adminStore.createProduct({
    title,
    categoryId: categoryValue,
    price,
    seller,
    wechat,
  })

  showCreateModal.value = false
  feedbackModal.value = { message: `已创建商品，编号：${product.id}，待管理员审核。` }
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
  align-items: flex-start;
}

.page-head h2 {
  margin: 0 0 4px;
}

.subtitle {
  margin: 0;
  color: #6b7280;
  font-size: 14px;
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

.btn.secondary {
  background: #f9fafb;
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
  align-items: center;
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

.status {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 9999px;
  font-size: 13px;
  background: #e5e7eb;
  color: #374151;
}

.status.pending {
  background: rgba(245, 158, 11, 0.18);
  color: #b45309;
}

.status.for-sale {
  background: rgba(59, 130, 246, 0.18);
  color: #1d4ed8;
}

.status.sold {
  background: rgba(16, 185, 129, 0.18);
  color: #047857;
}

.status.delisted {
  background: rgba(107, 114, 128, 0.18);
  color: #374151;
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

.link.danger {
  color: #dc2626;
}

.link:disabled {
  cursor: not-allowed;
  color: #9ca3af;
}

.link.danger:disabled {
  color: #fca5a5;
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
  display: block;
  font-size: 14px;
  color: #374151;
  font-weight: 500;
  margin-bottom: 4px;
}

.modal-input,
.modal-textarea,
.modal-select {
  width: 100%;
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px solid #d1d5db;
  background: #fff;
  font-size: 14px;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.modal-input:focus,
.modal-textarea:focus,
.modal-select:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
  outline: none;
}

.modal-textarea {
  resize: vertical;
  min-height: 100px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 4px;
  grid-column: 1 / -1;
}

.modal-error {
  margin: -6px 0 0;
  color: #dc2626;
  font-size: 13px;
  grid-column: 1 / -1;
}

.modal-message {
  margin: 0 0 12px;
  font-size: 15px;
  color: #374151;
}

.detail-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px 16px;
  margin: 0;
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
