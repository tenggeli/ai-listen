<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { isUnauthorizedError } from '../../api/ApiError'
import { HttpProviderOrderApi } from '../../api/ProviderOrderApi'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'
import type { ProviderOrder } from '../../domain/order/ProviderOrder'
import { toOrderDetailErrorMessage } from './orderError'
import { formatOrderTime, getOrderStatusLabel, getOrderStatusTagType } from './orderShared'

const route = useRoute()
const router = useRouter()
const api = new HttpProviderOrderApi('/api/v1/provider')
const state = reactive<{
  pageState: PageLoadState
  errorMessage: string
  order: ProviderOrder | null
}>({
  pageState: PageLoadState.Idle,
  errorMessage: '',
  order: null
})

onMounted(() => {
  void loadDetail()
})

async function loadDetail(): Promise<void> {
  const orderId = String(route.params.id || '').trim()
  if (!orderId) {
    state.pageState = PageLoadState.Error
    state.errorMessage = '订单 ID 无效'
    return
  }

  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  try {
    state.order = await api.getOrder(authService.getAccessToken(), orderId)
    state.pageState = PageLoadState.Success
  } catch (error) {
    if (isUnauthorizedError(error)) {
      authService.logout()
      await router.replace('/login')
      return
    }
    state.pageState = PageLoadState.Error
    state.errorMessage = toOrderDetailErrorMessage(error)
  }
}

function logout(): void {
  authService.logout()
  void router.replace('/login')
}
</script>

<template>
  <main class="page">
    <nav class="top-nav">
      <RouterLink to="/orders">返回订单列表</RouterLink>
      <button type="button" @click="logout">退出登录</button>
    </nav>

    <h1>订单详情</h1>

    <p v-if="state.pageState === PageLoadState.Idle">等待加载订单详情...</p>
    <p v-else-if="state.pageState === PageLoadState.Loading">订单详情加载中，请稍候...</p>
    <p v-else-if="state.pageState === PageLoadState.Error" class="error">{{ state.errorMessage }}</p>

    <section v-else-if="state.order" class="card">
      <header class="card-head">
        <strong>#{{ state.order.id }}</strong>
        <span :class="['tag', getOrderStatusTagType(state.order.status)]">{{ getOrderStatusLabel(state.order.status) }}</span>
      </header>
      <p><span>服务项目</span>{{ state.order.serviceItemTitle }}</p>
      <p><span>用户 ID</span>{{ state.order.userId }}</p>
      <p><span>服务方</span>{{ state.order.providerName }}</p>
      <p><span>订单金额</span>¥{{ state.order.amount }} {{ state.order.currency }}</p>
      <p><span>创建时间</span>{{ formatOrderTime(state.order.createdAt) }}</p>
      <p><span>支付时间</span>{{ formatOrderTime(state.order.paidAt) }}</p>
    </section>
  </main>
</template>

<style scoped>
.page {
  padding: 20px;
}
.top-nav {
  margin-bottom: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
h1 {
  margin: 0 0 12px;
}
.card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  padding: 14px;
}
.card-head {
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.card p {
  margin: 10px 0 0;
  display: flex;
  justify-content: space-between;
  gap: 16px;
}
.card p span {
  color: #64748b;
}
.tag {
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
}
.tag.default {
  background: #e2e8f0;
  color: #1e293b;
}
.tag.warn {
  background: #fff7ed;
  color: #9a3412;
}
.tag.info {
  background: #e0f2fe;
  color: #075985;
}
.tag.success {
  background: #dcfce7;
  color: #166534;
}
.error {
  color: #b91c1c;
}
</style>
