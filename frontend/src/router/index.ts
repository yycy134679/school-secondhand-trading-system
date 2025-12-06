import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // 公开路由
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/home/Index.vue'),
      meta: { title: '首页' },
    },
    {
      path: '/search',
      name: 'search',
      component: () => import('@/views/search/Index.vue'),
      meta: { title: '搜索结果', breadcrumb: '搜索' },
    },
    {
      path: '/category/:id',
      name: 'category',
      component: () => import('@/views/category/Index.vue'),
      meta: { title: '分类浏览', breadcrumb: '分类' },
    },
    {
      path: '/products/:id',
      name: 'product-detail',
      component: () => import('@/views/product/Detail.vue'),
      meta: { title: '商品详情', breadcrumb: '商品详情' },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/user/Login.vue'),
      meta: { title: '登录' },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/user/Register.vue'),
      meta: { title: '注册' },
    },

    // 需登录路由
    {
      path: '/products/new',
      name: 'product-new',
      component: () => import('@/views/product/New.vue'),
      meta: { requiresAuth: true, title: '发布闲置', breadcrumb: '发布闲置' },
    },
    {
      path: '/products/:id/edit',
      name: 'product-edit',
      component: () => import('@/views/product/Edit.vue'),
      meta: { requiresAuth: true, title: '编辑商品', breadcrumb: '编辑商品' },
    },
    {
      path: '/my/products',
      name: 'my-products',
      component: () => import('@/views/product/MyProducts.vue'),
      meta: { requiresAuth: true, title: '我发布的', breadcrumb: '我发布的' },
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/user/Profile.vue'),
      meta: { requiresAuth: true, title: '个人中心', breadcrumb: '个人中心' },
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()

  if (userStore.isLoggedIn && !userStore.currentUser) {
    await userStore.fetchProfile()
  }

  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 校园二手交易平台`
  } else {
    document.title = '校园二手交易平台'
  }

  // 权限检查
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    // 如果需要登录但未登录，跳转到登录页，并记录重定向地址
    next({
      path: '/login',
      query: { redirect: to.fullPath },
    })
  } else {
    next()
  }
})

export default router
