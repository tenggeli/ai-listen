import { createMemoryHistory, createRouter, createWebHistory, type RouterHistory } from 'vue-router'
import type { AuthService } from '../application/auth/AuthService'
import { authService } from '../application/auth'
import LoginPage from '../views/auth/LoginPage.vue'
import DashboardPage from '../views/dashboard/DashboardPage.vue'
import OrderListPage from '../views/orders/OrderListPage.vue'
import OrderDetailPage from '../views/orders/OrderDetailPage.vue'
import ProfilePage from '../views/profile/ProfilePage.vue'
import ServiceListPage from '../views/services/ServiceListPage.vue'

function createProviderRoutes() {
  return [
    { path: '/', redirect: '/dashboard' },
    { path: '/login', component: LoginPage },
    { path: '/dashboard', component: DashboardPage, meta: { requiresAuth: true } },
    { path: '/orders', component: OrderListPage, meta: { requiresAuth: true } },
    { path: '/orders/:id', component: OrderDetailPage, meta: { requiresAuth: true } },
    { path: '/profile', component: ProfilePage, meta: { requiresAuth: true } },
    { path: '/services', component: ServiceListPage, meta: { requiresAuth: true } }
  ]
}

export function createProviderRouter(
  service: Pick<AuthService, 'ensureSession'>,
  history: RouterHistory = createWebHistory('/provider/')
) {
  const router = createRouter({
    history,
    routes: createProviderRoutes()
  })

  router.beforeEach(async (to) => {
    const isAuthed = await service.ensureSession()
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

  return router
}

export function createProviderMemoryRouter(service: Pick<AuthService, 'ensureSession'>) {
  return createProviderRouter(service, createMemoryHistory('/provider/'))
}

export const router = createProviderRouter(authService)
