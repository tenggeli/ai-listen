<script setup lang="ts">
import { onMounted } from 'vue'
import ProviderTable from '../../components/providers/ProviderTable.vue'
import { ProviderReviewViewModel } from '../../application/provider/ProviderReviewViewModel'
import { HttpProviderAdminApi, MockProviderAdminApi } from '../../api/ProviderAdminApi'
import { PageLoadState } from '../../domain/provider/PageLoadState'

const useMock = import.meta.env.VITE_USE_MOCK !== 'false'
const api = useMock ? new MockProviderAdminApi() : new HttpProviderAdminApi('/api/v1/admin')
const vm = new ProviderReviewViewModel(api)

onMounted(() => {
  void vm.initialize()
})
</script>

<template>
  <main class="page">
    <header>
      <h1>服务方审核</h1>
      <p>管理后台 P0 模块：列表、详情、审核动作最小闭环。</p>
    </header>

    <section class="filters">
      <button @click="vm.changeFilter('')">全部</button>
      <button @click="vm.changeFilter('submitted')">待提交审核</button>
      <button @click="vm.changeFilter('under_review')">审核中</button>
      <button @click="vm.changeFilter('supplement_required')">需补充</button>
    </section>

    <section v-if="vm.state.listState === PageLoadState.Loading">列表加载中...</section>
    <section v-else-if="vm.state.listState === PageLoadState.Empty">当前筛选条件下暂无服务方。</section>
    <section v-else-if="vm.state.listState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</section>
    <section v-else-if="vm.state.listState === PageLoadState.Success" class="layout">
      <ProviderTable :items="vm.state.providers" :active-id="vm.state.selectedProviderId" @select="vm.selectProvider" />

      <div class="detail-panel">
        <p v-if="vm.state.detailState === PageLoadState.Loading">详情加载中...</p>
        <p v-else-if="vm.state.detailState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</p>
        <template v-else-if="vm.state.detailState === PageLoadState.Success && vm.state.selectedProviderDetail">
          <h3>{{ vm.state.selectedProviderDetail.displayName }}</h3>
          <p>城市：{{ vm.state.selectedProviderDetail.cityCode }}</p>
          <p>状态：{{ vm.state.selectedProviderDetail.reviewStatus }}</p>
          <p>简介：{{ vm.state.selectedProviderDetail.bio }}</p>
          <div class="actions">
            <button :disabled="vm.state.actionState === PageLoadState.Loading" @click="vm.review('approve')">通过</button>
            <button :disabled="vm.state.actionState === PageLoadState.Loading" @click="vm.review('reject')">拒绝</button>
            <button :disabled="vm.state.actionState === PageLoadState.Loading" @click="vm.review('require-supplement')">补充资料</button>
          </div>
        </template>
      </div>
    </section>

    <p class="error" v-if="vm.state.actionState === PageLoadState.Error">{{ vm.state.errorMessage }}</p>
  </main>
</template>

<style scoped>
.page {
  padding: 20px;
}

h1 {
  margin: 0;
}

.filters {
  display: flex;
  gap: 8px;
  margin: 12px 0;
  flex-wrap: wrap;
}

.layout {
  display: grid;
  grid-template-columns: 1.3fr 1fr;
  gap: 16px;
}

.detail-panel {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 12px;
}

.actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

button {
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 8px 12px;
  background: #fff;
  cursor: pointer;
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
