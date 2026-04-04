<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '../../application/auth'
import { HttpServiceItemAdminApi, MockServiceItemAdminApi } from '../../api/ServiceItemAdminApi'
import { ServiceItemManageViewModel } from '../../application/service_item/ServiceItemManageViewModel'
import { PageLoadState } from '../../domain/provider/PageLoadState'
import AdminShell from '../../components/layout/AdminShell.vue'

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
  <AdminShell title="服务项目管理" subtitle="支持服务项目筛选、详情查看与启停管理。" @logout="logout">
    <section class="admin-filters">
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

    <section v-if="vm.state.listState === PageLoadState.Loading" class="admin-card admin-loading">列表加载中...</section>
    <section v-else-if="vm.state.listState === PageLoadState.Empty" class="admin-card admin-empty">当前筛选条件下暂无服务项目。</section>
    <section v-else-if="vm.state.listState === PageLoadState.Error" class="admin-card error">{{ vm.state.errorMessage }}</section>
    <section v-else-if="vm.state.listState === PageLoadState.Success" class="admin-panel-grid">
      <section class="admin-list-panel">
        <button
          v-for="item in vm.state.items"
          :key="item.id"
          type="button"
          class="admin-list-item"
          :class="{ active: item.id === vm.state.selectedServiceItemId }"
          @click="vm.selectServiceItem(item.id)"
        >
          <div class="admin-title-row">
            <strong>{{ item.title }}</strong>
            <span class="pill-soft">{{ item.serviceStatus }}</span>
          </div>
          <p>{{ item.providerName }}</p>
          <p>{{ item.priceAmount }} / {{ item.priceUnit }}</p>
        </button>
      </section>

      <section class="admin-card admin-detail-panel">
        <p v-if="vm.state.detailState === PageLoadState.Loading">详情加载中...</p>
        <p v-else-if="vm.state.detailState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</p>
        <template v-else-if="vm.state.detailState === PageLoadState.Success && vm.state.selectedServiceItemDetail">
          <div class="admin-section-head">
            <h3>{{ vm.state.selectedServiceItemDetail.title }}</h3>
            <span class="pill-soft">{{ vm.state.selectedServiceItemDetail.serviceStatus }}</span>
          </div>
          <p>服务方：{{ vm.state.selectedServiceItemDetail.providerName }}</p>
          <p>分类：{{ vm.state.selectedServiceItemDetail.categoryId }}</p>
          <p>价格：{{ vm.state.selectedServiceItemDetail.priceAmount }} / {{ vm.state.selectedServiceItemDetail.priceUnit }}</p>
          <p>线上支持：{{ vm.state.selectedServiceItemDetail.supportOnline ? '是' : '否' }}</p>
          <p>说明：{{ vm.state.selectedServiceItemDetail.description || '-' }}</p>
          <div class="admin-actions">
            <button :disabled="vm.state.actionState === PageLoadState.Loading" @click="vm.updateStatus('activate')">启用</button>
            <button :disabled="vm.state.actionState === PageLoadState.Loading" @click="vm.updateStatus('deactivate')">停用</button>
          </div>
        </template>
      </section>
    </section>

    <p v-if="vm.state.actionState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</p>
  </AdminShell>
</template>
