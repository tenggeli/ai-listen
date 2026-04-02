<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { HttpProviderOrderApi } from '../../api/ProviderOrderApi'
import { authService } from '../../application/auth'
import type { ProviderOrder } from '../../domain/order/ProviderOrder'

const router = useRouter()
const api = new HttpProviderOrderApi('/api/v1/provider')

const state = reactive<{
  loading: boolean
  actionLoading: boolean
  items: ProviderOrder[]
  selectedOrder: ProviderOrder | null
  errorMessage: string
}>({
  loading: false,
  actionLoading: false,
  items: [],
  selectedOrder: null,
  errorMessage: ''
})

onMounted(() => {
  void loadOrders()
})

function statusText(status: ProviderOrder['status']): string {
  switch (status) {
    case 'paid':
      return '待接单'
    case 'accepted':
      return '已接单'
    case 'on_the_way':
      return '出发中'
    case 'arrived':
      return '已到达'
    case 'in_service':
      return '服务中'
    case 'completed':
      return '已完单'
    default:
      return '待支付'
  }
}

async function loadOrders(): Promise<void> {
  state.loading = true
  state.errorMessage = ''
  try {
    const result = await api.listOrders(authService.getAccessToken(), 1, 50)
    state.items = result.items
    if (!state.selectedOrder && state.items.length > 0) {
      state.selectedOrder = state.items[0]
    }
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '加载失败'
  } finally {
    state.loading = false
  }
}

async function refreshDetail(orderId: string): Promise<void> {
  try {
    const latest = await api.getOrder(authService.getAccessToken(), orderId)
    const index = state.items.findIndex((item) => item.id === latest.id)
    if (index >= 0) {
      state.items[index] = latest
    }
    state.selectedOrder = latest
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '刷新详情失败'
  }
}

async function act(action: 'accept' | 'depart' | 'arrive' | 'start' | 'complete'): Promise<void> {
  if (!state.selectedOrder) {
    return
  }
  state.actionLoading = true
  state.errorMessage = ''
  try {
    await api.operate(authService.getAccessToken(), state.selectedOrder.id, action)
    await refreshDetail(state.selectedOrder.id)
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '操作失败'
  } finally {
    state.actionLoading = false
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
      <RouterLink to="/dashboard">返回工作台</RouterLink>
      <button @click="logout">退出登录</button>
    </nav>

    <header>
      <h1>订单管理</h1>
      <p>支持接单、出发、到达、开始服务、完单。</p>
    </header>

    <p v-if="state.loading">订单加载中...</p>
    <section v-else class="layout">
      <section class="list-panel">
        <button
          v-for="item in state.items"
          :key="item.id"
          type="button"
          class="list-item"
          :class="{ active: state.selectedOrder?.id === item.id }"
          @click="refreshDetail(item.id)"
        >
          <div class="title-row">
            <strong>{{ item.serviceItemTitle }}</strong>
            <span>{{ statusText(item.status) }}</span>
          </div>
          <p>用户：{{ item.userId }}</p>
          <p>金额：¥{{ item.amount }}</p>
        </button>
      </section>

      <section class="detail-panel">
        <template v-if="state.selectedOrder">
          <h3>{{ state.selectedOrder.serviceItemTitle }}</h3>
          <p>状态：{{ statusText(state.selectedOrder.status) }}</p>
          <p>用户：{{ state.selectedOrder.userId }}</p>
          <p>服务方：{{ state.selectedOrder.providerName }}</p>
          <p>创建时间：{{ new Date(state.selectedOrder.createdAt).toLocaleString('zh-CN') }}</p>
          <div class="actions">
            <button :disabled="state.actionLoading" @click="act('accept')">接单</button>
            <button :disabled="state.actionLoading" @click="act('depart')">出发</button>
            <button :disabled="state.actionLoading" @click="act('arrive')">到达</button>
            <button :disabled="state.actionLoading" @click="act('start')">开始服务</button>
            <button :disabled="state.actionLoading" @click="act('complete')">完单</button>
          </div>
        </template>
      </section>
    </section>

    <p v-if="state.errorMessage" class="error">{{ state.errorMessage }}</p>
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
.layout {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 16px;
}
.list-panel {
  display: grid;
  gap: 8px;
}
.list-item {
  text-align: left;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: #fff;
  padding: 10px;
}
.list-item.active {
  border-color: #0f172a;
}
.title-row {
  display: flex;
  justify-content: space-between;
  gap: 8px;
}
.detail-panel {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 12px;
  background: #fff;
}
.actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.actions button {
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 8px 12px;
  background: #fff;
}
.error {
  color: #b91c1c;
}
@media (max-width: 960px) {
  .layout {
    grid-template-columns: 1fr;
  }
}
</style>
