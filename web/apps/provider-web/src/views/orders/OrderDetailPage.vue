<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { isUnauthorizedError } from '../../api/ApiError'
import { HttpProviderOrderApi } from '../../api/ProviderOrderApi'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'
import type { ProviderOrder } from '../../domain/order/ProviderOrder'
import ProviderShell from '../../components/layout/ProviderShell.vue'
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
  <ProviderShell title="订单详情" subtitle="展示订单关键信息，并支持按状态执行下一步履约动作。" @logout="logout">
    <p v-if="state.pageState === PageLoadState.Idle" class="provider-sub">等待加载订单详情...</p>
    <p v-else-if="state.pageState === PageLoadState.Loading" class="provider-sub">订单详情加载中，请稍候...</p>
    <p v-else-if="state.pageState === PageLoadState.Error" class="provider-error">{{ state.errorMessage }}</p>

    <section v-else-if="state.order" class="provider-card">
      <div class="provider-section-head">
        <h3>#{{ state.order.id }}</h3>
        <span :class="['provider-tag', getOrderStatusTagType(state.order.status)]">{{ getOrderStatusLabel(state.order) }}</span>
      </div>
      <p>服务项目：{{ state.order.serviceItemTitle }}</p>
      <p>用户 ID：{{ state.order.userId }}</p>
      <p>服务方：{{ state.order.providerName }}</p>
      <p>订单金额：¥{{ state.order.amount }} {{ state.order.currency }}</p>
      <p v-if="state.order.statusActionReason">状态原因：{{ state.order.statusActionReason }}</p>
      <p v-if="state.order.statusUpdatedAt">状态时间：{{ formatOrderTime(state.order.statusUpdatedAt) }}</p>
      <p>创建时间：{{ formatOrderTime(state.order.createdAt) }}</p>
      <p>支付时间：{{ formatOrderTime(state.order.paidAt) }}</p>
      <div class="provider-actions-grid">
        <button v-if="nextAction" type="button" :disabled="state.actionState === PageLoadState.Loading" @click="triggerNextAction">
          {{ state.actionState === PageLoadState.Loading ? '提交中...' : getOrderActionLabel(nextAction) }}
        </button>
        <p v-else class="provider-sub">当前状态无可执行履约动作。</p>
      </div>
      <p v-if="state.actionState === PageLoadState.Error" class="provider-error">{{ state.actionMessage }}</p>
      <p v-else-if="state.actionState === PageLoadState.Success" class="provider-success">{{ state.actionMessage }}</p>
    </section>
  </ProviderShell>
</template>
