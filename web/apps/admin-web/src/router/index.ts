import { createRouter, createWebHistory } from 'vue-router'
import { authService } from '../application/auth'
import LoginPage from '../views/auth/LoginPage.vue'
import DashboardPage from '../views/dashboard/DashboardPage.vue'
import ProviderReviewPage from '../views/providers/ProviderReviewPage.vue'
import ServiceItemManagePage from '../views/services/ServiceItemManagePage.vue'
import SoundManagePage from '../views/sound/SoundManagePage.vue'
import OrderManagePage from '../views/orders/OrderManagePage.vue'
import ComplaintManagePage from '../views/complaints/ComplaintManagePage.vue'

export const router = createRouter({
  history: createWebHistory('/admin/'),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/login', component: LoginPage },
    { path: '/dashboard', component: DashboardPage, meta: { requiresAuth: true } },
    { path: '/providers/review', component: ProviderReviewPage, meta: { requiresAuth: true } },
    { path: '/services/manage', component: ServiceItemManagePage, meta: { requiresAuth: true } },
    { path: '/sounds/manage', component: SoundManagePage, meta: { requiresAuth: true } },
    { path: '/orders/manage', component: OrderManagePage, meta: { requiresAuth: true } },
    { path: '/complaints/manage', component: ComplaintManagePage, meta: { requiresAuth: true } }
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
