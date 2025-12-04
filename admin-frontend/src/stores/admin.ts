import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export type ProductStatus = 'pending' | 'for-sale' | 'sold' | 'delisted'

export type UserRole = 'admin' | 'user'
export type UserStatus = 'active' | 'banned' | 'pending'

export interface AdminProduct {
  id: number
  title: string
  categoryId: number
  price: number
  seller: string
  wechat: string
  status: ProductStatus
  publishedAt: string
}

export interface PendingReviewItem {
  id: number
  productId: number
  title: string
  seller: string
  submittedAt: string
}

export interface AdminUser {
  id: number
  account: string
  nickname: string
  role: UserRole
  status: UserStatus
  wechat?: string
  registeredAt: string
  lastLogin: string
}

export interface AdminCategory {
  id: number
  name: string
  description: string
}

export interface AdminTag {
  id: number
  name: string
  description: string
  usage: number
}

export interface Announcement {
  id: number
  content: string
  createdAt: string
}

export interface TrendPoint {
  label: string
  value: number
}

export type AdminApplicationStatus = 'pending' | 'approved' | 'rejected'

export interface AdminApplication {
  id: number
  account: string
  nickname: string
  reason: string
  status: AdminApplicationStatus
  createdAt: string
  processedAt?: string
  reviewer?: string
  feedback?: string
}

let nextProductId = 105
let nextReviewId = 4
let nextUserId = 5
let nextCategoryId = 4
let nextTagId = 105
let nextAnnouncementId = 1
let nextApplicationId = 1

export const useAdminDataStore = defineStore('adminData', () => {
  const products = ref<AdminProduct[]>([
    {
      id: 101,
      title: 'ThinkPad X1 Carbon 2023 顶配 32G/1T',
      categoryId: 2,
      price: 6299,
      seller: 'linxia',
      wechat: 'linxia233',
      status: 'pending',
      publishedAt: '2025-11-30 21:16',
    },
    {
      id: 102,
      title: '飞跃牌羽绒服 M 码 九成新',
      categoryId: 3,
      price: 198,
      seller: 'forest',
      wechat: 'forest_yy',
      status: 'for-sale',
      publishedAt: '2025-11-29 11:04',
    },
    {
      id: 103,
      title: '雅思写作高分范文 + 课程兑换卡',
      categoryId: 1,
      price: 129,
      seller: 'carol',
      wechat: 'carol_ielts',
      status: 'for-sale',
      publishedAt: '2025-11-28 08:52',
    },
    {
      id: 104,
      title: 'Apple Watch S9 GPS 版 蓝色表带',
      categoryId: 2,
      price: 2599,
      seller: 'robin',
      wechat: 'robin-watch',
      status: 'sold',
      publishedAt: '2025-11-25 15:21',
    },
  ])

  const pendingReviews = ref<PendingReviewItem[]>([
    {
      id: 1,
      productId: 101,
      title: 'Apple iPad Air 5 256G Wi-Fi 版 深空灰色',
      seller: '张三',
      submittedAt: '2025-11-30 22:18',
    },
    {
      id: 2,
      productId: 102,
      title: '戴森 V12 吸尘器 九成新 使用三个月',
      seller: '李四',
      submittedAt: '2025-11-30 21:02',
    },
    {
      id: 3,
      productId: 103,
      title: '考研全套资料 408 + 英语一 + 政治',
      seller: '王五',
      submittedAt: '2025-11-29 19:44',
    },
  ])

  const users = ref<AdminUser[]>([
    {
      id: 1,
      account: 'admin',
      nickname: '超级管理员',
      role: 'admin',
      status: 'active',
      wechat: 'wechat_admin',
      registeredAt: '2024-09-01',
      lastLogin: '2025-12-01 08:10',
    },
    {
      id: 2,
      account: 'alice',
      nickname: 'Alice',
      role: 'user',
      status: 'active',
      wechat: 'alice-contact',
      registeredAt: '2025-03-16',
      lastLogin: '2025-11-30 20:42',
    },
    {
      id: 3,
      account: 'bob',
      nickname: 'Bob',
      role: 'user',
      status: 'pending',
      registeredAt: '2025-11-12',
      lastLogin: '未登录',
    },
    {
      id: 4,
      account: 'eve',
      nickname: 'Eve',
      role: 'user',
      status: 'banned',
      wechat: 'eve-2025',
      registeredAt: '2025-05-02',
      lastLogin: '2025-11-29 09:34',
    },
  ])

  const categories = ref<AdminCategory[]>([
    { id: 1, name: '教材教辅', description: '考研·四六级·专业课等学习资料' },
    { id: 2, name: '数码产品', description: '电脑数码、摄影设备、智能穿戴' },
    { id: 3, name: '生活日用', description: '宿舍电器、收纳清洁、居家好物' },
  ])

  const tags = ref<AdminTag[]>([
    { id: 101, name: '几乎全新', description: '使用不超过 1 个月，成色完好', usage: 342 },
    { id: 102, name: '送货上门', description: '卖家可协商校园内送货', usage: 211 },
    { id: 103, name: '价格可议', description: '支持少量砍价或换购', usage: 189 },
    { id: 104, name: '附带保修', description: '仍在官方质保期内', usage: 97 },
  ])

  const saleTrend = ref<TrendPoint[]>([
    { label: '周一', value: 8 },
    { label: '周二', value: 11 },
    { label: '周三', value: 9 },
    { label: '周四', value: 13 },
    { label: '周五', value: 17 },
    { label: '周六', value: 15 },
    { label: '周日', value: 19 },
  ])

  const announcements = ref<Announcement[]>([])
  const adminApplications = ref<AdminApplication[]>([])

  const categoryMap = computed(() => {
    const map = new Map<number, AdminCategory>()
    categories.value.forEach((item) => map.set(item.id, item))
    return map
  })

  const productById = computed(() => {
    const map = new Map<number, AdminProduct>()
    products.value.forEach((item) => map.set(item.id, item))
    return map
  })

  const productCountsByCategory = computed(() => {
    const counts = new Map<number, number>()
    categories.value.forEach((category) => counts.set(category.id, 0))
    products.value.forEach((product) => {
      counts.set(product.categoryId, (counts.get(product.categoryId) ?? 0) + 1)
    })
    return counts
  })

  const approvePendingReview = (id: number) => {
    const index = pendingReviews.value.findIndex((item) => item.id === id)
    if (index === -1) return false

    const [review] = pendingReviews.value.splice(index, 1)
    if (!review) return false
    const product = productById.value.get(review.productId)
    if (product) {
      product.status = 'for-sale'
      product.publishedAt = new Date().toISOString().replace('T', ' ').slice(0, 16)
    }

    const lastPoint = saleTrend.value[saleTrend.value.length - 1]
    if (lastPoint) {
      lastPoint.value += 1
    }
    return true
  }

  const rejectPendingReview = (id: number) => {
    const index = pendingReviews.value.findIndex((item) => item.id === id)
    if (index === -1) return false

    const [review] = pendingReviews.value.splice(index, 1)
    if (!review) return false
    const product = productById.value.get(review.productId)
    if (product) {
      product.status = 'delisted'
    }
    return true
  }

  const exportProducts = (productIds: number[] | null = null) => {
    const dataset = productIds && productIds.length > 0
      ? products.value.filter((item) => productIds.includes(item.id))
      : products.value

    const rows = dataset.map((item) => {
      const categoryName = categoryMap.value.get(item.categoryId)?.name ?? '未知分类'
      return [
        item.id,
        item.title,
        categoryName,
        item.price,
        item.seller,
        item.wechat,
        item.status,
        item.publishedAt,
      ]
    })

    const header = ['ID', '标题', '分类', '价格', '卖家', '微信号', '状态', '发布时间']
    return [header, ...rows]
      .map((row) => row.map((cell) => `"${String(cell).replace(/"/g, '""')}"`).join(','))
      .join('\n')
  }

  const createProduct = (payload: {
    title: string
    categoryId: number
    price: number
    seller: string
    wechat: string
  }) => {
    const newProduct: AdminProduct = {
      id: nextProductId++,
      title: payload.title,
      categoryId: payload.categoryId,
      price: payload.price,
      seller: payload.seller,
      wechat: payload.wechat,
      status: 'pending',
      publishedAt: new Date().toISOString().replace('T', ' ').slice(0, 16),
    }

    products.value.unshift(newProduct)

    pendingReviews.value.unshift({
      id: nextReviewId++,
      productId: newProduct.id,
      title: newProduct.title,
      seller: newProduct.seller,
      submittedAt: newProduct.publishedAt,
    })

    return newProduct
  }

  const toggleProductStatus = (productId: number, status: ProductStatus) => {
    const product = productById.value.get(productId)
    if (!product) return false
    product.status = status
    return true
  }

  const addAnnouncement = (content: string) => {
    announcements.value.unshift({
      id: nextAnnouncementId++,
      content,
      createdAt: new Date().toISOString().replace('T', ' ').slice(0, 16),
    })
  }

  const promoteToAdmin = (userId: number, nickname?: string) => {
    const user = users.value.find((item) => item.id === userId)
    if (!user) return false
    user.role = 'admin'
    if (nickname) user.nickname = nickname
    user.status = 'active'
    return true
  }

  const updateUserNickname = (userId: number, nickname: string) => {
    const trimmed = nickname.trim()
    if (!trimmed) return false
    const user = users.value.find((item) => item.id === userId)
    if (!user) return false
    user.nickname = trimmed
    return true
  }

  const createAdminAccount = (payload: { account: string; nickname: string; wechat?: string }) => {
    users.value.unshift({
      id: nextUserId++,
      account: payload.account,
      nickname: payload.nickname,
      role: 'admin',
      status: 'active',
      wechat: payload.wechat,
      registeredAt: new Date().toISOString().slice(0, 10),
      lastLogin: '未登录',
    })
  }

  const addViewerAccount = (payload: { account: string; nickname: string }) => {
    const exists = users.value.some((item) => item.account === payload.account)
    if (exists) return
    users.value.unshift({
      id: nextUserId++,
      account: payload.account,
      nickname: payload.nickname,
      role: 'user',
      status: 'active',
      registeredAt: new Date().toISOString().slice(0, 10),
      lastLogin: '未登录',
    })
  }

  const submitAdminApplication = (payload: {
    account: string
    nickname: string
    reason: string
  }): { success: boolean; message: string } => {
    const account = payload.account.trim()
    const nickname = payload.nickname.trim() || account
    const reason = payload.reason.trim()

    if (!account) {
      return { success: false, message: '账号信息缺失，请重新登录后再提交。' }
    }

    const user = users.value.find((item) => item.account === account)
    if (user?.role === 'admin') {
      return { success: false, message: '您已经是管理员，无需重复申请。' }
    }

    const pending = adminApplications.value.find(
      (item) => item.account === account && item.status === 'pending',
    )
    if (pending) {
      return { success: false, message: '您已有待审核的申请，请耐心等待管理员处理。' }
    }

    if (!reason) {
      return { success: false, message: '请填写申请理由，帮助管理员了解您的诉求。' }
    }

    if (!user) {
      addViewerAccount({ account, nickname })
    }

    adminApplications.value.unshift({
      id: nextApplicationId++,
      account,
      nickname,
      reason,
      status: 'pending',
      createdAt: new Date().toISOString().replace('T', ' ').slice(0, 16),
    })

    return { success: true, message: '申请已提交，请等待管理员审核。' }
  }

  const approveAdminApplication = (
    id: number,
    reviewer: string,
  ): { success: boolean; message: string } => {
    const application = adminApplications.value.find((item) => item.id === id)
    if (!application) {
      return { success: false, message: '未找到该申请。' }
    }

    if (application.status !== 'pending') {
      return { success: false, message: '该申请已处理，请刷新后查看最新状态。' }
    }

    let target = users.value.find((item) => item.account === application.account)
    if (!target) {
      addViewerAccount({ account: application.account, nickname: application.nickname })
      target = users.value.find((item) => item.account === application.account)
    }

    if (!target) {
      return { success: false, message: '未找到申请人账号，请稍后重试。' }
    }

    const promoted = promoteToAdmin(target.id, application.nickname)
    if (!promoted) {
      return { success: false, message: '升级管理员失败，请稍后重试。' }
    }

    application.status = 'approved'
    application.processedAt = new Date().toISOString().replace('T', ' ').slice(0, 16)
    application.reviewer = reviewer
    application.feedback = '申请已通过'

    return { success: true, message: '已同意该申请，并设置其为管理员。' }
  }

  const rejectAdminApplication = (
    id: number,
    reviewer: string,
    feedback: string,
  ): { success: boolean; message: string } => {
    const application = adminApplications.value.find((item) => item.id === id)
    if (!application) {
      return { success: false, message: '未找到该申请。' }
    }

    if (application.status !== 'pending') {
      return { success: false, message: '该申请已处理，请刷新后查看最新状态。' }
    }

    const note = feedback.trim()
    if (!note) {
      return { success: false, message: '请填写驳回原因，方便申请人了解情况。' }
    }

    application.status = 'rejected'
    application.processedAt = new Date().toISOString().replace('T', ' ').slice(0, 16)
    application.reviewer = reviewer
    application.feedback = note

    return { success: true, message: '已驳回该申请。' }
  }

  const toggleUserBan = (userId: number) => {
    const user = users.value.find((item) => item.id === userId)
    if (!user) return false
    user.status = user.status === 'banned' ? 'active' : 'banned'
    return true
  }

  const addCategory = (payload: { name: string; description: string }) => {
    categories.value.push({ id: nextCategoryId++, ...payload })
  }

  const updateCategory = (id: number, payload: { name: string; description: string }) => {
    const category = categories.value.find((item) => item.id === id)
    if (!category) return false
    category.name = payload.name
    category.description = payload.description
    return true
  }

  const removeCategory = (id: number) => {
    const hasProduct = products.value.some((item) => item.categoryId === id)
    if (hasProduct) return false
    const index = categories.value.findIndex((item) => item.id === id)
    if (index === -1) return false
    categories.value.splice(index, 1)
    return true
  }

  const addTag = (payload: { name: string; description: string }) => {
    tags.value.push({ id: nextTagId++, usage: 0, ...payload })
  }

  const updateTag = (id: number, payload: { name: string; description: string }) => {
    const tag = tags.value.find((item) => item.id === id)
    if (!tag) return false
    tag.name = payload.name
    tag.description = payload.description
    return true
  }

  const removeTag = (id: number) => {
    const index = tags.value.findIndex((item) => item.id === id)
    if (index === -1) return false
    tags.value.splice(index, 1)
    return true
  }

  return {
    products,
    pendingReviews,
    users,
    categories,
    tags,
    saleTrend,
    announcements,
    adminApplications,
    categoryMap,
    productCountsByCategory,
    approvePendingReview,
    rejectPendingReview,
    exportProducts,
    createProduct,
    toggleProductStatus,
    addAnnouncement,
    promoteToAdmin,
    createAdminAccount,
    updateUserNickname,
    addViewerAccount,
    submitAdminApplication,
    approveAdminApplication,
    rejectAdminApplication,
    toggleUserBan,
    addCategory,
    updateCategory,
    removeCategory,
    addTag,
    updateTag,
    removeTag,
  }
})
