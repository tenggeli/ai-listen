<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '../../application/auth'
import { HttpServiceItemAdminApi, MockServiceItemAdminApi } from '../../api/ServiceItemAdminApi'
import { ServiceItemManageViewModel } from '../../application/service_item/ServiceItemManageViewModel'
import { PageLoadState } from '../../domain/provider/PageLoadState'

const router = useRouter()
const useMock = import.meta.env.VITE_USE_MOCK === 'true'
const api = useMock
  ? new MockServiceItemAdminApi()
  : new HttpServiceItemAdminApi('/api/v1/admin', () => authService.getAccessToken())
const vm = new ServiceItemManageViewModel(api)

onMounted(() => {
  void vm.initialize()
})

function logout(): void {
  authService.logout()
  void router.replace('/login')
}
</script>

<template>
  <main class="page">
    <nav class="top-nav">
      <RouterLink to="/dashboard">返回仪表盘</RouterLink>
      <button @click="logout">退出登录</button>
    </nav>

    <header>
      <h1>服务项目管理</h1>
      <p>首版能力：列表筛选、详情查看、启用/停用。</p>
    </header>

    <section class="filters">
      <input v-model.trim="vm.state.providerId" placeholder="provider_id" />
      <input v-model.trim="vm.state.categoryId" placeholder="category_id" />
      <select v-model="vm.state.status">
        <option value="">全部状态</option>
        <option value="active">active</option>
        <option value="inactive">inactive</option>
      </select>
      <input v-model.trim="vm.state.keyword" placeholder="关键词" />
      <button @click="vm.applyFilter">查询</button>
    </section>

    <section v-if="vm.state.listState === PageLoadState.Loading">列表加载中...</section>
    <section v-else-if="vm.state.listState === PageLoadState.Empty">当前筛选条件下暂无服务项目。</section>
    <section v-else-if="vm.state.listState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</section>
    <section v-else-if="vm.state.listState === PageLoadState.Success" class="layout">
      <section class="list-panel">
        <button
          v-for="item in vm.state.items"
          :key="item.id"
          type="button"
          class="list-item"
          :class="{ active: item.id === vm.state.selectedServiceItemId }"
          @click="vm.selectServiceItem(item.id)"
        >
          <div class="title-row">
            <strong>{{ item.title }}</strong>
            <span class="status">{{ item.serviceStatus }}</span>
          </div>
          <p>{{ item.providerName }}</p>
          <p>{{ item.priceAmount }} / {{ item.priceUnit }}</p>
        </button>
      </section>

      <section class="detail-panel">
        <p v-if="vm.state.detailState === PageLoadState.Loading">详情加载中...</p>
        <p v-else-if="vm.state.detailState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</p>
        <template v-else-if="vm.state.detailState === PageLoadState.Success && vm.state.selectedServiceItemDetail">
          <h3>{{ vm.state.selectedServiceItemDetail.title }}</h3>
          <p>服务方：{{ vm.state.selectedServiceItemDetail.providerName }}</p>
          <p>分类：{{ vm.state.selectedServiceItemDetail.categoryId }}</p>
          <p>状态：{{ vm.state.selectedServiceItemDetail.serviceStatus }}</p>
          <p>价格：{{ vm.state.selectedServiceItemDetail.priceAmount }} / {{ vm.state.selectedServiceItemDetail.priceUnit }}</p>
          <p>线上支持：{{ vm.state.selectedServiceItemDetail.supportOnline ? '是' : '否' }}</p>
          <p>说明：{{ vm.state.selectedServiceItemDetail.description || '-' }}</p>
          <div class="actions">
            <button :disabled="vm.state.actionState === PageLoadState.Loading" @click="vm.updateStatus('activate')">启用</button>
            <button :disabled="vm.state.actionState === PageLoadState.Loading" @click="vm.updateStatus('deactivate')">停用</button>
          </div>
        </template>
      </section>
    </section>

    <p v-if="vm.state.actionState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</p>
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

h1 {
  margin: 0;
}

.filters {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 8px;
  margin: 12px 0;
}

.filters input,
.filters select,
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

.status {
  color: #334155;
  font-size: 12px;
}

.detail-panel {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 12px;
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
  .filters {
    grid-template-columns: 1fr 1fr;
  }

  .layout {
    grid-template-columns: 1fr;
  }
}
</style>

