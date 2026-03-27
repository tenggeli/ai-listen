import { createRouter, createWebHistory } from 'vue-router'
import DashboardPage from '../views/dashboard/DashboardPage.vue'
import ProviderReviewPage from '../views/providers/ProviderReviewPage.vue'

export const router = createRouter({
  history: createWebHistory('/admin/'),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/dashboard', component: DashboardPage },
    { path: '/providers/review', component: ProviderReviewPage }
  ]
})
