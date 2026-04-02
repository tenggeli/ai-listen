<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '../../application/auth'
import { HttpComplaintAdminApi, type ComplaintSummary, type ComplaintDetail } from '../../api/ComplaintAdminApi'

const router = useRouter()
const api = new HttpComplaintAdminApi('/api/v1/admin', () => authService.getAccessToken())

const state = reactive<{
  loading: boolean
  actionLoading: boolean
  items: ComplaintSummary[]
  selectedOrderId: string
  detail: ComplaintDetail | null
  errorMessage: string
}>({
  loading: false,
  actionLoading: false,
  items: [],
  selectedOrderId: '',
  detail: null,
  errorMessage: ''
})

onMounted(() => {
  void loadComplaints()
})

async function loadComplaints(): Promise<void> {
  state.loading = true
  state.errorMessage = ''
  try {
    const result = await api.list()
    state.items = result.items
    if (state.items.length > 0) {
      const targetID = state.selectedOrderId || state.items[0].orderId
      await selectComplaint(targetID)
    } else {
      state.selectedOrderId = ''
      state.detail = null
    }
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '加载投诉失败'
  } finally {
    state.loading = false
  }
}

async function selectComplaint(orderId: string): Promise<void> {
  state.selectedOrderId = orderId
  try {
    state.detail = await api.detail(orderId)
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '加载投诉详情失败'
  }
}

async function action(actionName: 'intervene' | 'resolve'): Promise<void> {
  if (!state.selectedOrderId) {
    return
  }
  state.actionLoading = true
  state.errorMessage = ''
  try {
    await api.action(state.selectedOrderId, actionName, actionName === 'resolve' ? '投诉已处理' : '平台人工介入')
    await loadComplaints()
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

function complaintStatusText(status: string): string {
  switch (status) {
    case 'pending':
      return '待处理'
    case 'processing':
      return '处理中'
    case 'resolved':
      return '已处理'
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
      <h1>投诉管理</h1>
      <p>支持投诉列表、详情、介入与处理完成。</p>
    </header>

    <section v-if="state.loading">投诉加载中...</section>
    <section v-else class="layout">
      <section class="list-panel">
        <button
          v-for="item in state.items"
          :key="item.orderId"
          type="button"
          class="list-item"
          :class="{ active: state.selectedOrderId === item.orderId }"
          @click="selectComplaint(item.orderId)"
        >
          <div class="title-row">
            <strong>{{ item.complaintReason || '未填写原因' }}</strong>
            <span>{{ complaintStatusText(item.complaintStatus) }}</span>
          </div>
          <p>订单：{{ item.orderId }}</p>
          <p>用户：{{ item.userId }} · 服务方：{{ item.providerName }}</p>
        </button>
      </section>

      <section class="detail-panel" v-if="state.detail">
        <h3>投诉详情</h3>
        <p>投诉状态：{{ complaintStatusText(state.detail.complaintStatus) }}</p>
        <p>订单：{{ state.detail.order.id }}</p>
        <p>服务项目：{{ state.detail.order.serviceItemTitle }}</p>
        <p>投诉原因：{{ state.detail.feedback.complaintReason || '-' }}</p>
        <p>投诉描述：{{ state.detail.feedback.complaintContent || '-' }}</p>
        <p>评价分：{{ state.detail.feedback.ratingScore || '-' }}</p>
        <div class="actions">
          <button :disabled="state.actionLoading" @click="action('intervene')">介入处理</button>
          <button :disabled="state.actionLoading" @click="action('resolve')">处理完成</button>
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
  .layout {
    grid-template-columns: 1fr;
  }
}
</style>
