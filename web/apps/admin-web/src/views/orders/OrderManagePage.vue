<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '../../application/auth'
import { HttpOrderAdminApi, type AdminOrderSummary, type AdminOrderDetail } from '../../api/OrderAdminApi'
import AdminShell from '../../components/layout/AdminShell.vue'

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
</script>

<template>
  <AdminShell title="订单监管" subtitle="支持订单筛选、详情排查与人工介入处置。" @logout="logout">
    <section class="admin-filters admin-filters-compact">
      <select v-model="state.status">
        <option value="">全部状态</option>
        <option value="created">待支付（created）</option>
        <option value="paid">待服务方接单（paid）</option>
        <option value="accepted">服务方已接单（accepted）</option>
        <option value="on_the_way">服务方出发中（on_the_way）</option>
        <option value="arrived">服务方已到达，待开始服务（arrived）</option>
        <option value="in_service">服务进行中（in_service）</option>
        <option value="completed">服务已完成（completed）</option>
        <option value="after_sale_processing">订单售后处理中（after_sale_processing）</option>
        <option value="closed">订单已关闭（closed）</option>
      </select>
      <input v-model.trim="state.keyword" placeholder="关键词: 订单号/用户/服务方/项目" />
      <button @click="loadOrders">查询</button>
    </section>

    <section v-if="state.loading" class="admin-card admin-loading">订单加载中...</section>
    <section v-else class="admin-panel-grid">
      <section class="admin-list-panel">
        <button
          v-for="item in state.items"
          :key="item.id"
          type="button"
          class="admin-list-item"
          :class="{ active: state.selectedOrderId === item.id }"
          @click="selectOrder(item.id)"
        >
          <div class="admin-title-row">
            <strong>{{ item.serviceItemTitle }}</strong>
            <span class="pill-soft">{{ item.statusReason }}</span>
          </div>
          <p>订单号：{{ item.id }}</p>
          <p>用户：{{ item.userId }} · 服务方：{{ item.providerName }}</p>
          <p>金额：¥{{ item.amount }}</p>
        </button>
      </section>

      <section class="admin-card admin-detail-panel" v-if="state.detail">
        <h3>{{ state.detail.order.serviceItemTitle }}</h3>
        <p>状态：{{ state.detail.order.statusReason }}</p>
        <p>订单号：{{ state.detail.order.id }}</p>
        <p>用户：{{ state.detail.order.userId }}</p>
        <p>服务方：{{ state.detail.order.providerName }}</p>
        <p>金额：¥{{ state.detail.order.amount }}</p>
        <p v-if="state.detail.feedback">反馈：{{ state.detail.feedback.hasComplaint ? '有投诉' : '仅评价' }}</p>
        <div class="admin-actions">
          <button :disabled="state.actionLoading" @click="action('intervene')">人工介入</button>
          <button :disabled="state.actionLoading" @click="action('close')">人工关闭</button>
        </div>
        <div class="admin-log-panel">
          <h4>操作日志</h4>
          <p v-if="state.detail.actionLogs.length === 0">暂无操作记录</p>
          <ul v-else>
            <li v-for="log in state.detail.actionLogs" :key="log.actionId">
              <strong>{{ log.scope }}/{{ log.action }}</strong>
              <span> · {{ log.operator }} · {{ log.reason || '无原因' }} · {{ log.updatedAt }}</span>
              <span v-if="log.statusAfterReason"> · 状态：{{ log.statusBefore || 'unknown' }} → {{ log.statusAfter }}（{{ log.statusAfterReason }}）</span>
            </li>
          </ul>
        </div>
      </section>
    </section>

    <p v-if="state.errorMessage" class="error">{{ state.errorMessage }}</p>
  </AdminShell>
</template>
