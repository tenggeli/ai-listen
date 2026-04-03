<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { isUnauthorizedError } from '../../api/ApiError'
import { HttpProviderOrderApi } from '../../api/ProviderOrderApi'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'
import type { ProviderOrder } from '../../domain/order/ProviderOrder'
import { toOrderDetailErrorMessage } from './orderError'
import { formatOrderTime, getNextOrderAction, getOrderActionLabel, getOrderStatusLabel, getOrderStatusTagType } from './orderShared'

const route = useRoute()
const router = useRouter()
const api = new HttpProviderOrderApi('/api/v1/provider')
const state = reactive<{
  pageState: PageLoadState
  actionState: PageLoadState
  errorMessage: string
  actionMessage: string
  order: ProviderOrder | null
}>({
  pageState: PageLoadState.Idle,
  actionState: PageLoadState.Idle,
  errorMessage: '',
  actionMessage: '',
  order: null
})
const nextAction = computed(() => (state.order ? getNextOrderAction(state.order.status) : null))

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
  state.actionMessage = ''
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

async function triggerNextAction(): Promise<void> {
  if (!state.order || !nextAction.value) {
    return
  }
  if (state.actionState === PageLoadState.Loading) {
    return
  }
  state.actionState = PageLoadState.Loading
  state.actionMessage = ''
  try {
    await api.operate(authService.getAccessToken(), state.order.id, nextAction.value)
    state.actionState = PageLoadState.Success
    state.actionMessage = '履约状态更新成功。'
    await loadDetail()
  } catch (error) {
    if (isUnauthorizedError(error)) {
      authService.logout()
      await router.replace('/login')
      return
    }
    state.actionState = PageLoadState.Error
    state.actionMessage = error instanceof Error ? error.message : '履约操作失败'
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
      <div class="action-block">
        <button
          v-if="nextAction"
          type="button"
          class="action-btn"
          :disabled="state.actionState === PageLoadState.Loading"
          @click="triggerNextAction"
        >
          {{ state.actionState === PageLoadState.Loading ? '提交中...' : getOrderActionLabel(nextAction) }}
        </button>
        <p v-else class="action-tip">当前状态无可执行履约动作。</p>
        <p v-if="state.actionState === PageLoadState.Error" class="error">{{ state.actionMessage }}</p>
        <p v-else-if="state.actionState === PageLoadState.Success" class="success">{{ state.actionMessage }}</p>
      </div>
    </section>
  </main>
</template>

<style scoped>
.page {
  min-height: 100vh;
  padding: 20px;
  color: #eaf6fb;
  background:
    radial-gradient(circle at top left, rgba(31, 109, 141, 0.2), transparent 30%),
    radial-gradient(circle at bottom right, rgba(21, 59, 93, 0.22), transparent 32%),
    #06111b;
}
.top-nav {
  margin-bottom: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.top-nav a {
  color: #8bd7ff;
}
.top-nav button {
  border: 1px solid rgba(148, 217, 255, 0.24);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.04);
  color: #eaf6fb;
  padding: 8px 12px;
}
h1 {
  margin: 0 0 12px;
}
.card {
  border: 1px solid rgba(148, 217, 255, 0.12);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.045);
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
  color: rgba(234, 246, 251, 0.62);
}
.tag {
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
}
.tag.default {
  background: rgba(255, 255, 255, 0.12);
  color: #dceaf2;
}
.tag.warn {
  background: rgba(255, 189, 89, 0.16);
  color: #ffd278;
}
.tag.info {
  background: rgba(115, 213, 255, 0.16);
  color: #8bd7ff;
}
.tag.success {
  background: rgba(91, 212, 154, 0.16);
  color: #7df0bc;
}
.error {
  color: #ffd278;
}
.success {
  color: #7df0bc;
}
.action-block {
  margin-top: 16px;
  border-top: 1px solid rgba(148, 217, 255, 0.14);
  padding-top: 12px;
}
.action-btn {
  border: 1px solid rgba(115, 213, 255, 0.4);
  border-radius: 14px;
  background: rgba(115, 213, 255, 0.16);
  color: #eaf6fb;
  padding: 10px 14px;
}
.action-tip {
  color: rgba(234, 246, 251, 0.62);
}
</style>
