<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { HttpOrderApi } from '../../api/OrderApi'
import { HttpServiceDiscoveryApi, MockServiceDiscoveryApi } from '../../api/ServiceDiscoveryApi'
import { loadSession } from '../../application/identity/AuthSession'
import { PageLoadState } from '../../domain/ai/PageLoadState'
import type { ProviderPublicProfile } from '../../domain/service/ProviderPublicProfile'
import type { ServiceItem } from '../../domain/service/ServiceItem'

interface PaymentConfirmState {
  pageState: PageLoadState
  provider: ProviderPublicProfile | null
  serviceItem: ServiceItem | null
  submitting: boolean
  errorMessage: string
}

const route = useRoute()
const router = useRouter()
const useMock = import.meta.env.VITE_USE_MOCK === 'true'
const api = useMock ? new MockServiceDiscoveryApi() : new HttpServiceDiscoveryApi('/api/v1')
const orderApi = new HttpOrderApi('/api/v1')
const state: PaymentConfirmState = reactive({
  pageState: PageLoadState.Idle,
  provider: null,
  serviceItem: null,
  submitting: false,
  errorMessage: ''
})

onMounted(() => {
  void initialize()
})

async function initialize(): Promise<void> {
  const providerId = String(route.query.provider_id ?? '')
  const serviceItemId = String(route.query.service_item_id ?? '')
  if (!providerId || !serviceItemId) {
    state.pageState = PageLoadState.Error
    state.errorMessage = '下单信息不完整，请返回服务详情重新选择。'
    return
  }

  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  try {
    const [provider, items] = await Promise.all([api.getProvider(providerId), api.getProviderServiceItems(providerId)])
    const selected = items.find((item) => item.id === serviceItemId)
    if (!selected) {
      state.pageState = PageLoadState.Error
      state.errorMessage = '未找到对应服务项目，请返回重试。'
      return
    }
    state.provider = provider
    state.serviceItem = selected
    state.pageState = PageLoadState.Success
  } catch (error) {
    state.pageState = PageLoadState.Error
    state.errorMessage = error instanceof Error ? error.message : '加载支付确认信息失败'
  }
}

async function confirmPay(): Promise<void> {
  if (!state.provider || !state.serviceItem) {
    return
  }
  const session = loadSession()
  if (!session) {
    await router.push('/auth')
    return
  }

  state.submitting = true
  state.errorMessage = ''
  try {
    const created = await orderApi.createOrder(session.accessToken, {
      providerId: state.provider.id,
      providerName: state.provider.displayName,
      serviceItemId: state.serviceItem.id,
      serviceItemTitle: state.serviceItem.title,
      amount: state.serviceItem.priceAmount,
      currency: 'CNY'
    })
    const paid = await orderApi.payOrderMockSuccess(session.accessToken, created.id)
    await router.push(`/orders/${paid.id}`)
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '创建订单失败，请稍后重试。'
  } finally {
    state.submitting = false
  }
}
</script>

<template>
  <main class="payment-page">
    <section class="screen-shell">
      <header class="top-row">
        <button type="button" class="back" @click="router.back()">返回</button>
        <span>支付确认</span>
      </header>

      <section v-if="state.pageState === PageLoadState.Loading" class="state-card">正在准备支付信息...</section>
      <section v-else-if="state.pageState === PageLoadState.Error" class="state-card error">
        {{ state.errorMessage }}
      </section>

      <template v-if="state.pageState === PageLoadState.Success && state.provider && state.serviceItem">
        <section class="panel">
          <h1>确认订单</h1>
          <p class="sub">当前为支付 mock 流程，后续可替换真实支付网关。</p>
          <div class="line"><span>服务方</span><strong>{{ state.provider.displayName }}</strong></div>
          <div class="line"><span>服务项目</span><strong>{{ state.serviceItem.title }}</strong></div>
          <div class="line"><span>服务说明</span><strong>{{ state.serviceItem.description }}</strong></div>
          <div class="line total"><span>应付金额</span><strong>¥{{ state.serviceItem.priceAmount }}</strong></div>
        </section>

        <section class="actions">
          <button type="button" class="ghost" @click="router.push(`/providers/${state.provider.id}`)">返回修改</button>
          <button type="button" class="primary" :disabled="state.submitting" @click="confirmPay">
            {{ state.submitting ? '正在处理...' : '确认支付（Mock Success）' }}
          </button>
        </section>
        <p v-if="state.errorMessage" class="submit-error">{{ state.errorMessage }}</p>
      </template>
    </section>
  </main>
</template>

<style scoped>
.payment-page {
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

.state-card,
.panel {
  margin-top: 14px;
  border-radius: 18px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.045);
  padding: 16px;
}

.panel h1 {
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
  margin-top: 18px;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  font-size: 18px;
}

.actions {
  margin-top: 16px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.actions button {
  border: none;
  border-radius: 14px;
  padding: 14px 10px;
  font-size: 14px;
  font-weight: 600;
}

.actions .ghost {
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.74);
}

.actions .primary {
  background: linear-gradient(135deg, #97e3ff, #58bee8);
  color: #082132;
}

.actions .primary:disabled {
  opacity: 0.65;
}

.error {
  color: #fca5a5;
}

.submit-error {
  margin-top: 12px;
  color: #fca5a5;
  font-size: 13px;
}
</style>
