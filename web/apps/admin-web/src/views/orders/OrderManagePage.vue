<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '../../application/auth'
import { HttpOrderAdminApi, type AdminOrderSummary, type AdminOrderDetail } from '../../api/OrderAdminApi'

const router = useRouter()
const api = new HttpOrderAdminApi('/api/v1/admin', () => authService.getAccessToken())

const state = reactive<{
  loading: boolean
  actionLoading: boolean
  status: string
  keyword: string
  items: AdminOrderSummary[]
  selectedOrderId: string
  detail: AdminOrderDetail | null
  errorMessage: string
}>({
  loading: false,
  actionLoading: false,
  status: '',
  keyword: '',
  items: [],
  selectedOrderId: '',
  detail: null,
  errorMessage: ''
})

onMounted(() => {
  void loadOrders()
})

async function loadOrders(): Promise<void> {
  state.loading = true
  state.errorMessage = ''
  try {
    const result = await api.list({ status: state.status, keyword: state.keyword })
    state.items = result.items
    if (state.items.length > 0) {
      const targetID = state.selectedOrderId || state.items[0].id
      await selectOrder(targetID)
    } else {
      state.detail = null
      state.selectedOrderId = ''
    }
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '加载订单失败'
  } finally {
    state.loading = false
  }
}

async function selectOrder(orderId: string): Promise<void> {
  state.selectedOrderId = orderId
  try {
    state.detail = await api.detail(orderId)
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '加载详情失败'
  }
}

async function action(actionName: 'intervene' | 'close'): Promise<void> {
  if (!state.selectedOrderId) {
    return
  }
  state.actionLoading = true
  state.errorMessage = ''
  try {
    await api.action(state.selectedOrderId, actionName, actionName === 'intervene' ? '人工介入' : '人工关闭')
    await loadOrders()
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

function statusText(status: string): string {
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
      return '已完成'
    case 'after_sale_processing':
      return '售后处理中'
    case 'closed':
      return '已关闭'
    default:
      return status
  }
}
</script>

<template>
  <main class="page">
    <nav class="top-nav">
      <RouterLink to="/dashboard">返回仪表盘</RouterLink>
      <button @click="logout">退出登录</button>
    </nav>

    <header>
      <h1>订单管理</h1>
      <p>支持订单列表、详情与人工介入/关闭动作。</p>
    </header>

    <section class="filters">
      <select v-model="state.status">
        <option value="">全部状态</option>
        <option value="paid">paid</option>
        <option value="accepted">accepted</option>
        <option value="on_the_way">on_the_way</option>
        <option value="arrived">arrived</option>
        <option value="in_service">in_service</option>
        <option value="completed">completed</option>
        <option value="after_sale_processing">after_sale_processing</option>
        <option value="closed">closed</option>
      </select>
      <input v-model.trim="state.keyword" placeholder="关键词: 订单号/用户/服务方/项目" />
      <button @click="loadOrders">查询</button>
    </section>

    <section v-if="state.loading">订单加载中...</section>
    <section v-else class="layout">
      <section class="list-panel">
        <button
          v-for="item in state.items"
          :key="item.id"
          type="button"
          class="list-item"
          :class="{ active: state.selectedOrderId === item.id }"
          @click="selectOrder(item.id)"
        >
          <div class="title-row">
            <strong>{{ item.serviceItemTitle }}</strong>
            <span>{{ statusText(item.status) }}</span>
          </div>
          <p>订单号：{{ item.id }}</p>
          <p>用户：{{ item.userId }} · 服务方：{{ item.providerName }}</p>
          <p>金额：¥{{ item.amount }}</p>
        </button>
      </section>

      <section class="detail-panel" v-if="state.detail">
        <h3>{{ state.detail.order.serviceItemTitle }}</h3>
        <p>状态：{{ statusText(state.detail.order.status) }}</p>
        <p>订单号：{{ state.detail.order.id }}</p>
        <p>用户：{{ state.detail.order.userId }}</p>
        <p>服务方：{{ state.detail.order.providerName }}</p>
        <p>金额：¥{{ state.detail.order.amount }}</p>
        <p v-if="state.detail.feedback">反馈：{{ state.detail.feedback.hasComplaint ? '有投诉' : '仅评价' }}</p>
        <div class="actions">
          <button :disabled="state.actionLoading" @click="action('intervene')">人工介入</button>
          <button :disabled="state.actionLoading" @click="action('close')">人工关闭</button>
        </div>
        <div class="log-panel">
          <h4>操作日志</h4>
          <p v-if="state.detail.actionLogs.length === 0">暂无操作记录</p>
          <ul v-else>
            <li v-for="log in state.detail.actionLogs" :key="log.actionId">
              <strong>{{ log.scope }}/{{ log.action }}</strong>
              <span> · {{ log.operator }} · {{ log.reason || '无原因' }} · {{ log.updatedAt }}</span>
            </li>
          </ul>
        </div>
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
.filters {
  display: grid;
  grid-template-columns: 1fr 2fr auto;
  gap: 8px;
  margin: 12px 0;
}
.filters select,
.filters input,
.filters button {
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 8px 10px;
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
}
.actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}
.actions button {
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 8px 12px;
  background: #fff;
}
.log-panel {
  margin-top: 12px;
  border-top: 1px solid #e2e8f0;
  padding-top: 10px;
}
.log-panel ul {
  margin: 8px 0 0;
  padding-left: 18px;
}
.log-panel li {
  margin: 6px 0;
  color: #334155;
}
.error {
  color: #b91c1c;
}
@media (max-width: 960px) {
  .filters {
    grid-template-columns: 1fr;
  }
  .layout {
    grid-template-columns: 1fr;
  }
}
</style>
