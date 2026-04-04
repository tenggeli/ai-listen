<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { isUnauthorizedError } from '../../api/ApiError'
import { HttpProviderOrderApi } from '../../api/ProviderOrderApi'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'
import type { ProviderOrderStatus } from '../../domain/order/ProviderOrder'
import ProviderShell from '../../components/layout/ProviderShell.vue'
import {
  filterOrdersByStatus,
  formatOrderTime,
  getOrderStatusLabel,
  getOrderStatusTagType,
  getTotalPages,
  type OrderStatusFilter
} from './orderShared'

const router = useRouter()
const api = new HttpProviderOrderApi('/api/v1/provider')
const pageSize = 10
const state = reactive({
  pageState: PageLoadState.Idle,
  errorMessage: '',
  page: 1,
  total: 0,
  items: [] as Array<{
    id: string
    serviceItemTitle: string
    userId: string
    amount: number
    status: ProviderOrderStatus
    statusReason: string
    createdAt: string
    paidAt: string | null
  }>,
  statusFilter: 'all' as OrderStatusFilter
})

const filteredItems = computed(() => filterOrdersByStatus(state.items, state.statusFilter))
const totalPages = computed(() => getTotalPages(state.total, pageSize))

onMounted(() => {
  void loadOrders()
})

async function loadOrders(): Promise<void> {
  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  try {
    const status = state.statusFilter === 'all' ? undefined : state.statusFilter
    const result = await api.listOrders(authService.getAccessToken(), state.page, pageSize, status)
    state.items = result.items
    state.total = result.total
    state.pageState = PageLoadState.Success
  } catch (error) {
    if (isUnauthorizedError(error)) {
      authService.logout()
      await router.replace('/login')
      return
    }
    state.pageState = PageLoadState.Error
    state.errorMessage = error instanceof Error ? error.message : '加载订单失败'
  }
}

function updateFilter(raw: string): void {
  state.statusFilter = raw as OrderStatusFilter
  state.page = 1
  void loadOrders()
}

function goPrevPage(): void {
  if (state.page <= 1) {
    return
  }
  state.page -= 1
  void loadOrders()
}

function goNextPage(): void {
  if (state.page >= totalPages.value) {
    return
  }
  state.page += 1
  void loadOrders()
}

function openDetail(orderId: string): void {
  void router.push(`/orders/${orderId}`)
}

function logout(): void {
  authService.logout()
  void router.replace('/login')
}
</script>

<template>
  <ProviderShell title="订单列表" subtitle="支持按状态筛选、分页浏览，并进入订单详情页查看完整信息。" @logout="logout">
    <section class="provider-filter-row">
      <div class="provider-filter">
        <span>状态筛选</span>
        <select :value="state.statusFilter" @change="updateFilter(($event.target as HTMLSelectElement).value)">
          <option value="all">全部状态</option>
          <option value="created">待支付</option>
          <option value="paid">待服务方接单</option>
          <option value="accepted">服务方已接单</option>
          <option value="on_the_way">服务方出发中</option>
          <option value="arrived">服务方已到达，待开始服务</option>
          <option value="in_service">服务进行中</option>
          <option value="completed">服务已完成</option>
          <option value="after_sale_processing">订单售后处理中</option>
          <option value="closed">订单已关闭</option>
        </select>
      </div>
    </section>

    <p v-if="state.pageState === PageLoadState.Idle" class="provider-sub">等待加载订单数据...</p>
    <p v-else-if="state.pageState === PageLoadState.Loading" class="provider-sub">订单加载中，请稍候...</p>
    <p v-else-if="state.pageState === PageLoadState.Error" class="provider-error">{{ state.errorMessage }}</p>

    <section v-else>
      <div v-if="filteredItems.length === 0" class="provider-empty">
        <p>当前筛选条件下暂无订单。</p>
        <p class="provider-sub">可切换状态筛选或稍后重试。</p>
      </div>

      <div v-else class="provider-stack">
        <button v-for="item in filteredItems" :key="item.id" type="button" class="provider-row-button" @click="openDetail(item.id)">
          <div class="provider-row-headline">
            <strong>#{{ item.id }}</strong>
            <span :class="['provider-tag', getOrderStatusTagType(item.status)]">{{ getOrderStatusLabel(item) }}</span>
          </div>
          <p>服务项目：{{ item.serviceItemTitle }}</p>
          <p>用户 ID：{{ item.userId }}</p>
          <p>金额：¥{{ item.amount }}</p>
          <p>创建时间：{{ formatOrderTime(item.createdAt) }}</p>
          <p>支付时间：{{ formatOrderTime(item.paidAt) }}</p>
        </button>
      </div>

      <footer class="provider-pager">
        <button type="button" :disabled="state.page <= 1 || state.pageState === PageLoadState.Loading" @click="goPrevPage">
          上一页
        </button>
        <span>第 {{ state.page }} / {{ totalPages }} 页</span>
        <button
          type="button"
          :disabled="state.page >= totalPages || state.pageState === PageLoadState.Loading"
          @click="goNextPage"
        >
          下一页
        </button>
      </footer>
    </section>
  </ProviderShell>
</template>
