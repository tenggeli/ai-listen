import { createRouter, createWebHistory } from 'vue-router'
import { authService } from '../application/auth'
import LoginPage from '../views/auth/LoginPage.vue'
import DashboardPage from '../views/dashboard/DashboardPage.vue'
import OrderManagePage from '../views/orders/OrderManagePage.vue'
import ProfilePage from '../views/profile/ProfilePage.vue'

export const router = createRouter({
  history: createWebHistory('/provider/'),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/login', component: LoginPage },
    { path: '/dashboard', component: DashboardPage, meta: { requiresAuth: true } },
    { path: '/orders', component: OrderManagePage, meta: { requiresAuth: true } },
    { path: '/profile', component: ProfilePage, meta: { requiresAuth: true } }
  ]
})

router.beforeEach(async (to) => {
  const isAuthed = await authService.ensureSession()
  if (to.meta.requiresAuth && !isAuthed) {
    return {
      path: '/login',
      query: { redirect: to.fullPath }
    }
  }
  if (to.path === '/login' && isAuthed) {
    return { path: '/dashboard' }
  }
  return true
})
