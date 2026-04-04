<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '../../application/auth'
import { HttpComplaintAdminApi, type ComplaintSummary, type ComplaintDetail } from '../../api/ComplaintAdminApi'
import AdminShell from '../../components/layout/AdminShell.vue'

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
  <AdminShell title="投诉仲裁" subtitle="支持投诉列表、详情与介入处理，沉淀完整处理日志。" @logout="logout">
    <section v-if="state.loading" class="admin-card admin-loading">投诉加载中...</section>
    <section v-else class="admin-panel-grid">
      <section class="admin-list-panel">
        <button
          v-for="item in state.items"
          :key="item.orderId"
          type="button"
          class="admin-list-item"
          :class="{ active: state.selectedOrderId === item.orderId }"
          @click="selectComplaint(item.orderId)"
        >
          <div class="admin-title-row">
            <strong>{{ item.complaintReason || '未填写原因' }}</strong>
            <span class="pill-risk">{{ complaintStatusText(item.complaintStatus) }}</span>
          </div>
          <p>订单：{{ item.orderId }}</p>
          <p>用户：{{ item.userId }} · 服务方：{{ item.providerName }}</p>
        </button>
      </section>

      <section class="admin-card admin-detail-panel" v-if="state.detail">
        <h3>投诉详情</h3>
        <p>投诉状态：{{ complaintStatusText(state.detail.complaintStatus) }}</p>
        <p>订单：{{ state.detail.order.id }}</p>
        <p>服务项目：{{ state.detail.order.serviceItemTitle }}</p>
        <p>投诉原因：{{ state.detail.feedback.complaintReason || '-' }}</p>
        <p>投诉描述：{{ state.detail.feedback.complaintContent || '-' }}</p>
        <p>评价分：{{ state.detail.feedback.ratingScore || '-' }}</p>
        <div class="admin-actions">
          <button :disabled="state.actionLoading" @click="action('intervene')">介入处理</button>
          <button :disabled="state.actionLoading" @click="action('resolve')">处理完成</button>
        </div>
        <div class="admin-log-panel">
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
  </AdminShell>
</template>
