<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '../../application/auth'
import ProviderTable from '../../components/providers/ProviderTable.vue'
import AdminShell from '../../components/layout/AdminShell.vue'
import { ProviderReviewViewModel } from '../../application/provider/ProviderReviewViewModel'
import { HttpProviderAdminApi, MockProviderAdminApi } from '../../api/ProviderAdminApi'
import { PageLoadState } from '../../domain/provider/PageLoadState'

const router = useRouter()
const useMock = import.meta.env.VITE_USE_MOCK === 'true'
const api = useMock
  ? new MockProviderAdminApi()
  : new HttpProviderAdminApi('/api/v1/admin', () => authService.getAccessToken())
const vm = new ProviderReviewViewModel(api)

onMounted(() => {
  void vm.initialize()
})

function logout(): void {
  authService.logout()
  void router.replace('/login')
}
</script>

<template>
  <AdminShell title="服务方审核中心" subtitle="处理服务方准入、补件与复审，确保服务供给质量与合规性。" @logout="logout">
    <section class="admin-filters">
      <button @click="vm.changeFilter('')">全部</button>
      <button @click="vm.changeFilter('submitted')">待提交审核</button>
      <button @click="vm.changeFilter('under_review')">审核中</button>
      <button @click="vm.changeFilter('supplement_required')">需补充</button>
    </section>

    <section v-if="vm.state.listState === PageLoadState.Loading" class="admin-card admin-loading">列表加载中...</section>
    <section v-else-if="vm.state.listState === PageLoadState.Empty" class="admin-card admin-empty">当前筛选条件下暂无服务方。</section>
    <section v-else-if="vm.state.listState === PageLoadState.Error" class="admin-card error">{{ vm.state.errorMessage }}</section>
    <section v-else-if="vm.state.listState === PageLoadState.Success" class="admin-panel-grid">
      <ProviderTable :items="vm.state.providers" :active-id="vm.state.selectedProviderId" @select="vm.selectProvider" />

      <div class="admin-card admin-detail-panel">
        <p v-if="vm.state.detailState === PageLoadState.Loading">详情加载中...</p>
        <p v-else-if="vm.state.detailState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</p>
        <template v-else-if="vm.state.detailState === PageLoadState.Success && vm.state.selectedProviderDetail">
          <div class="admin-section-head">
            <h3>{{ vm.state.selectedProviderDetail.displayName }}</h3>
            <span class="pill-soft">{{ vm.state.selectedProviderDetail.reviewStatus }}</span>
          </div>
          <p>城市：{{ vm.state.selectedProviderDetail.cityCode }}</p>
          <p>简介：{{ vm.state.selectedProviderDetail.bio }}</p>
          <div class="admin-actions">
            <button :disabled="vm.state.actionState === PageLoadState.Loading" @click="vm.review('approve')">通过</button>
            <button :disabled="vm.state.actionState === PageLoadState.Loading" @click="vm.review('reject')">拒绝</button>
            <button :disabled="vm.state.actionState === PageLoadState.Loading" @click="vm.review('require-supplement')">补充资料</button>
          </div>
        </template>
      </div>
    </section>

    <p class="error" v-if="vm.state.actionState === PageLoadState.Error">{{ vm.state.errorMessage }}</p>
  </AdminShell>
</template>
