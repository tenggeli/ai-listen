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
          <option value="paid">待服务方接单</option>
          <option value="accepted">服务方已接单</option>
          <option value="on_the_way">服务方出发中</option>
          <option value="arrived">服务方已到达，待开始服务</option>
          <option value="in_service">服务进行中</option>
          <option value="completed">服务已完成</option>
          <option value="after_sale_processing">订单售后处理中</option>
          <option value="closed">订单已关闭</option>
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
            <span :class="['tag', getOrderStatusTagType(item.status)]">{{ getOrderStatusLabel(item) }}</span>
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
  color: rgba(234, 246, 251, 0.66);
}
.filter {
  display: grid;
  gap: 6px;
}
.filter span {
  color: rgba(234, 246, 251, 0.72);
  font-size: 13px;
}
.filter select {
  border: 1px solid rgba(148, 217, 255, 0.24);
  border-radius: 14px;
  padding: 8px 10px;
  background: rgba(255, 255, 255, 0.06);
  color: #eaf6fb;
}
.list {
  display: grid;
  gap: 10px;
}
.row {
  text-align: left;
  border: 1px solid rgba(148, 217, 255, 0.12);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.045);
  padding: 12px;
  color: #eaf6fb;
}
.row-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}
.row p {
  margin: 8px 0 0;
  color: rgba(234, 246, 251, 0.78);
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
.pager {
  margin-top: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}
.pager button {
  border: 1px solid rgba(148, 217, 255, 0.24);
  border-radius: 10px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.04);
  color: #eaf6fb;
}
.pager span {
  color: rgba(234, 246, 251, 0.78);
}
.empty {
  border: 1px dashed rgba(148, 217, 255, 0.28);
  border-radius: 16px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.04);
}
.empty p {
  margin: 0;
}
.sub {
  margin-top: 8px !important;
  color: rgba(234, 246, 251, 0.56);
}
.error {
  color: #ffd278;
}
@media (max-width: 900px) {
  .header {
    display: grid;
    align-items: initial;
  }
}
</style>
