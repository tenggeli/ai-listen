<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { HttpOrderApi } from '../../api/OrderApi'
import { loadSession } from '../../application/identity/AuthSession'
import { PageLoadState } from '../../domain/ai/PageLoadState'
import type { UserOrder } from '../../domain/order/UserOrder'

const route = useRoute()
const router = useRouter()
const orderApi = new HttpOrderApi('/api/v1')

const orderId = computed(() => String(route.params.id ?? ''))

const state = reactive<{
  pageState: PageLoadState
  order: UserOrder | null
  errorMessage: string
}>({
  pageState: PageLoadState.Idle,
  order: null,
  errorMessage: ''
})

onMounted(() => {
  void loadOrder()
})

const statusText = computed(() => {
  if (!state.order) {
    return ''
  }
  switch (state.order.status) {
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
})

async function loadOrder(): Promise<void> {
  const session = loadSession()
  if (!session) {
    await router.push('/auth')
    return
  }
  if (!orderId.value) {
    state.pageState = PageLoadState.Error
    state.errorMessage = '订单号无效'
    return
  }

  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  try {
    state.order = await orderApi.getOrder(session.accessToken, orderId.value)
    state.pageState = PageLoadState.Success
  } catch (error) {
    state.pageState = PageLoadState.Error
    state.errorMessage = error instanceof Error ? error.message : '订单加载失败'
  }
}
</script>

<template>
  <main class="order-page">
    <section class="screen-shell">
      <header class="top-row">
        <button type="button" class="back" @click="router.push('/services')">返回服务</button>
        <span>订单详情</span>
      </header>

      <section v-if="state.pageState === PageLoadState.Loading" class="card">
        <h1>订单加载中</h1>
        <p>正在同步订单详情，请稍候。</p>
      </section>

      <section v-else-if="state.pageState === PageLoadState.Error" class="card">
        <h1>订单不存在</h1>
        <p>{{ state.errorMessage || '订单不存在或无查看权限。' }}</p>
      </section>

      <section v-else-if="state.order" class="card">
        <h1>{{ state.order.status === 'completed' ? '订单已完成' : state.order.status === 'paid' ? '支付成功' : '订单详情' }}</h1>
        <p class="sub">当前订单已接入真实订单接口。</p>
        <div class="line"><span>订单号</span><strong>{{ state.order.id }}</strong></div>
        <div class="line"><span>状态</span><strong>{{ statusText }}</strong></div>
        <div class="line"><span>服务方</span><strong>{{ state.order.providerName }}</strong></div>
        <div class="line"><span>服务项目</span><strong>{{ state.order.serviceItemTitle }}</strong></div>
        <div class="line total"><span>支付金额</span><strong>¥{{ state.order.amount }}</strong></div>
        <button type="button" class="ghost" @click="router.push(`/orders/${state.order.id}/feedback`)">
          去评价 / 投诉
        </button>
        <button type="button" class="primary" @click="router.push('/services')">继续浏览服务</button>
      </section>
    </section>
  </main>
</template>

<style scoped>
.order-page {
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

.sub {
  margin: 8px 0 0;
  color: rgba(239, 247, 251, 0.62);
  font-size: 13px;
}

.line {
  margin-top: 12px;
  display: flex;
  justify-content: space-between;
  gap: 10px;
  font-size: 14px;
}

.line span {
  color: rgba(239, 247, 251, 0.56);
}

.line strong {
  text-align: right;
}

.line.total {
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  font-size: 18px;
}

.primary {
  margin-top: 10px;
  width: 100%;
  border: none;
  border-radius: 14px;
  padding: 14px 10px;
  font-size: 14px;
  font-weight: 600;
  background: linear-gradient(135deg, #97e3ff, #58bee8);
  color: #082132;
}

.ghost {
  margin-top: 16px;
  width: 100%;
  border-radius: 14px;
  padding: 14px 10px;
  font-size: 14px;
  font-weight: 600;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.74);
}
</style>
