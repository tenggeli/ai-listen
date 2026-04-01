import { createRouter, createWebHistory } from 'vue-router'
import { loadSession, nextOnboardingRoute } from '../application/identity/AuthSession'
import AuthEntryPage from '../views/auth/AuthEntryPage.vue'
import ChatPage from '../views/chat/ChatPage.vue'
import HomePage from '../views/home/HomePage.vue'
import OrderDetailPage from '../views/order/OrderDetailPage.vue'
import OrderFeedbackPage from '../views/order/OrderFeedbackPage.vue'
import OrderListPage from '../views/order/OrderListPage.vue'
import PaymentConfirmPage from '../views/payment/PaymentConfirmPage.vue'
import MyPage from '../views/me/MyPage.vue'
import PersonalitySetupPage from '../views/profile/PersonalitySetupPage.vue'
import ProfileSetupPage from '../views/profile/ProfileSetupPage.vue'
import SettingsPage from '../views/settings/SettingsPage.vue'
import SoundPage from '../views/sound/SoundPage.vue'
import ProviderDetailPage from '../views/services/ProviderDetailPage.vue'
import ServicesPage from '../views/services/ServicesPage.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/auth' },
    { path: '/auth', component: AuthEntryPage, meta: { public: true } },
    { path: '/profile/setup', component: ProfileSetupPage },
    { path: '/personality/setup', component: PersonalitySetupPage },
    { path: '/home', component: HomePage },
    { path: '/me', component: MyPage },
    { path: '/services', component: ServicesPage },
    { path: '/providers/:id', component: ProviderDetailPage },
    { path: '/orders', component: OrderListPage },
    { path: '/payment/confirm', component: PaymentConfirmPage },
    { path: '/orders/:id', component: OrderDetailPage },
    { path: '/orders/:id/feedback', component: OrderFeedbackPage },
    { path: '/settings', component: SettingsPage },
    { path: '/chat', component: ChatPage },
    { path: '/sound', component: SoundPage }
  ]
})

router.beforeEach((to) => {
  const session = loadSession()
  const isPublic = Boolean(to.meta.public)

  if (!session && !isPublic) {
    return '/auth'
  }

  if (session && to.path === '/auth') {
    return nextOnboardingRoute(session)
  }

  if (session && !session.profileCompleted && to.path !== '/profile/setup') {
    return '/profile/setup'
  }

  if (session && session.profileCompleted && !session.personalityCompleted && to.path !== '/personality/setup') {
    return '/personality/setup'
  }

  return true
})
