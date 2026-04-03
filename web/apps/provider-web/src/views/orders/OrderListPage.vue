<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { isUnauthorizedError } from '../../api/ApiError'
import { HttpProviderOrderApi } from '../../api/ProviderOrderApi'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'
import type { ProviderOrderStatus } from '../../domain/order/ProviderOrder'
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
  <main class="page">
    <nav class="top-nav">
      <RouterLink to="/dashboard">返回工作台</RouterLink>
      <button type="button" @click="logout">退出登录</button>
    </nav>

    <header class="header">
      <div>
        <h1>订单列表</h1>
        <p>支持按状态筛选、分页浏览，并进入订单详情页查看完整信息。</p>
      </div>
      <label class="filter">
        <span>状态筛选</span>
        <select :value="state.statusFilter" @change="updateFilter(($event.target as HTMLSelectElement).value)">
          <option value="all">全部状态</option>
          <option value="created">待支付</option>
          <option value="paid">待接单</option>
          <option value="accepted">已接单</option>
          <option value="on_the_way">出发中</option>
          <option value="arrived">已到达</option>
          <option value="in_service">服务中</option>
          <option value="completed">已完单</option>
        </select>
      </label>
    </header>

    <p v-if="state.pageState === PageLoadState.Idle">等待加载订单数据...</p>
    <p v-else-if="state.pageState === PageLoadState.Loading">订单加载中，请稍候...</p>
    <p v-else-if="state.pageState === PageLoadState.Error" class="error">{{ state.errorMessage }}</p>

    <section v-else>
      <div v-if="filteredItems.length === 0" class="empty">
        <p>当前筛选条件下暂无订单。</p>
        <p class="sub">可切换状态筛选或稍后重试。</p>
      </div>

      <div v-else class="list">
        <button v-for="item in filteredItems" :key="item.id" type="button" class="row" @click="openDetail(item.id)">
          <div class="row-head">
            <strong>#{{ item.id }}</strong>
            <span :class="['tag', getOrderStatusTagType(item.status)]">{{ getOrderStatusLabel(item.status) }}</span>
          </div>
          <p>服务项目：{{ item.serviceItemTitle }}</p>
          <p>用户 ID：{{ item.userId }}</p>
          <p>金额：¥{{ item.amount }}</p>
          <p>创建时间：{{ formatOrderTime(item.createdAt) }}</p>
          <p>支付时间：{{ formatOrderTime(item.paidAt) }}</p>
        </button>
      </div>

      <footer class="pager">
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
.header {
  margin-bottom: 12px;
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: flex-end;
}
.header h1 {
  margin: 0;
}
.header p {
  margin: 8px 0 0;
  color: #475569;
}
.filter {
  display: grid;
  gap: 6px;
}
.filter span {
  color: #475569;
  font-size: 13px;
}
.filter select {
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  padding: 8px 10px;
}
.list {
  display: grid;
  gap: 10px;
}
.row {
  text-align: left;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  padding: 12px;
}
.row-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}
.row p {
  margin: 8px 0 0;
  color: #334155;
}
.tag {
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
}
.tag.default {
  background: #e2e8f0;
  color: #1e293b;
}
.tag.warn {
  background: #fff7ed;
  color: #9a3412;
}
.tag.info {
  background: #e0f2fe;
  color: #075985;
}
.tag.success {
  background: #dcfce7;
  color: #166534;
}
.pager {
  margin-top: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}
.pager button {
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 8px 12px;
  background: #fff;
}
.empty {
  border: 1px dashed #cbd5e1;
  border-radius: 12px;
  padding: 20px;
  background: #fff;
}
.empty p {
  margin: 0;
}
.sub {
  margin-top: 8px !important;
  color: #64748b;
}
.error {
  color: #b91c1c;
}
@media (max-width: 900px) {
  .header {
    display: grid;
    align-items: initial;
  }
}
</style>
