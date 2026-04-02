<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { HttpOrderApi } from '../../api/OrderApi'
import { loadSession } from '../../application/identity/AuthSession'
import { PageLoadState } from '../../domain/ai/PageLoadState'
import type { UserOrder } from '../../domain/order/UserOrder'

const router = useRouter()
const orderApi = new HttpOrderApi('/api/v1')
const state = reactive<{
  pageState: PageLoadState
  items: UserOrder[]
  errorMessage: string
}>({
  pageState: PageLoadState.Idle,
  items: [],
  errorMessage: ''
})

onMounted(() => {
  void loadOrders()
})

function statusText(status: UserOrder['status']): string {
  switch (status) {
    case 'paid':
      return '待服务方接单'
    case 'accepted':
      return '已接单'
    case 'on_the_way':
      return '服务方出发中'
    case 'arrived':
      return '服务方已到达'
    case 'in_service':
      return '服务中'
    case 'completed':
      return '已完成'
    case 'after_sale_processing':
      return '售后处理中'
    case 'closed':
      return '已关闭'
    default:
      return '待支付'
  }
}

async function loadOrders(): Promise<void> {
  const session = loadSession()
  if (!session) {
    await router.push('/auth')
    return
  }

  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  try {
    const result = await orderApi.listOrders(session.accessToken, 1, 20)
    state.items = result.items
    state.pageState = result.items.length === 0 ? PageLoadState.Empty : PageLoadState.Success
  } catch (error) {
    state.pageState = PageLoadState.Error
    state.errorMessage = error instanceof Error ? error.message : '订单加载失败'
  }
}
</script>

<template>
  <main class="order-list-page">
    <section class="screen-shell">
      <header class="top-row">
        <button type="button" class="back" @click="router.push('/me')">返回我的</button>
        <span>订单列表</span>
      </header>

      <section v-if="state.pageState === PageLoadState.Loading" class="card">
        <h1>订单加载中</h1>
        <p>正在同步你的订单记录。</p>
      </section>

      <section v-else-if="state.pageState === PageLoadState.Error" class="card">
        <h1>加载失败</h1>
        <p>{{ state.errorMessage || '请稍后重试。' }}</p>
        <button type="button" class="primary" @click="loadOrders">重新加载</button>
      </section>

      <section v-else-if="state.pageState === PageLoadState.Empty" class="card">
        <h1>还没有订单</h1>
        <p>去服务页下第一单吧。</p>
        <button type="button" class="primary" @click="router.push('/services')">去服务页下单</button>
      </section>

      <section v-else class="list">
        <article v-for="item in state.items" :key="item.id" class="order-card" @click="router.push(`/orders/${item.id}`)">
          <div class="row">
            <h2>{{ item.serviceItemTitle }}</h2>
            <strong class="price">¥{{ item.amount }}</strong>
          </div>
          <p class="provider">{{ item.providerName }}</p>
          <div class="row meta">
            <span>{{ statusText(item.status) }}</span>
            <span>{{ new Date(item.createdAt).toLocaleString('zh-CN') }}</span>
          </div>
        </article>
      </section>
    </section>
  </main>
</template>

<style scoped>
.order-list-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top, rgba(74, 168, 196, 0.1), transparent 30%),
    linear-gradient(180deg, #0b1926 0%, #0d1f2d 42%, #08131f 100%);
  color: rgba(255, 255, 255, 0.88);
}

.screen-shell {
  width: min(100%, 390px);
  min-height: 100vh;
  margin: 0 auto;
  padding: 52px 16px 24px;
}

.top-row {
  display: flex;
  align-items: center;
  gap: 10px;
  color: rgba(255, 255, 255, 0.45);
  font-size: 12px;
}

.back {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.8);
  padding: 6px 10px;
}

.card {
  margin-top: 14px;
  border-radius: 18px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.045);
  padding: 16px;
}

.card h1 {
  margin: 0;
  font-size: 24px;
}

.card p {
  margin: 10px 0 0;
  color: rgba(239, 247, 251, 0.62);
  font-size: 14px;
  line-height: 1.7;
}

.primary {
  margin-top: 16px;
  width: 100%;
  border: none;
  border-radius: 14px;
  padding: 14px 10px;
  font-size: 14px;
  font-weight: 600;
  background: linear-gradient(135deg, #97e3ff, #58bee8);
  color: #082132;
}

.list {
  margin-top: 14px;
  display: grid;
  gap: 10px;
}

.order-card {
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.045);
  padding: 14px;
  cursor: pointer;
}

.row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.order-card h2 {
  margin: 0;
  font-size: 16px;
}

.provider {
  margin: 8px 0 0;
  color: rgba(239, 247, 251, 0.62);
  font-size: 13px;
}

.meta {
  margin-top: 10px;
  color: rgba(239, 247, 251, 0.52);
  font-size: 12px;
}

.price {
  color: rgba(151, 227, 255, 0.95);
}
</style>
