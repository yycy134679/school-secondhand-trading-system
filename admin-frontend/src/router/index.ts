import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

import AdminLayout from '../layouts/AdminLayout.vue'
import Dashboard from '../views/Dashboard.vue'
import ProductList from '../views/ProductList.vue'
import UserList from '../views/UserList.vue'
import CategoryManage from '../views/CategoryManage.vue'
import Profile from '../views/Profile.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import NotFound from '../views/NotFound.vue'

import { useAuthStore } from '../stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/admin',
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: Dashboard,
      },
      {
        path: 'products',
        name: 'AdminProducts',
        component: ProductList,
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: UserList,
      },
      {
        path: 'categories',
        name: 'AdminCategories',
        component: CategoryManage,
      },
      {
        path: 'profile',
        name: 'AdminProfile',
        component: Profile,
      },
    ],
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return next({
      path: '/login',
      query: {
        redirect: to.fullPath,
      },
    })
  }

  if ((to.path === '/login' || to.path === '/register') && authStore.isAuthenticated) {
    const target = typeof to.query.redirect === 'string' ? to.query.redirect : '/admin'
    return next(target)
  }

  return next()
})

export default router
