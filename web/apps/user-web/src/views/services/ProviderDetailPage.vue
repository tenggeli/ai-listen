<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ProviderDetailPageViewModel } from '../../application/service/ProviderDetailPageViewModel'
import { MockServiceDiscoveryApi, HttpServiceDiscoveryApi } from '../../api/ServiceDiscoveryApi'
import { PageLoadState } from '../../domain/ai/PageLoadState'

const route = useRoute()
const router = useRouter()
const useMock = import.meta.env.VITE_USE_MOCK === 'true'
const api = useMock ? new MockServiceDiscoveryApi() : new HttpServiceDiscoveryApi('/api/v1')
const vm = new ProviderDetailPageViewModel(api)
const providerId = computed(() => String(route.params.id ?? ''))

onMounted(() => {
  if (!providerId.value) {
    return
  }
  void vm.initialize(providerId.value)
})

function selectedServiceTitle(serviceItemId: string): string {
  return vm.state.serviceItems.find((item) => item.id === serviceItemId)?.title ?? '请选择服务项目'
}

function goPaymentConfirm(): void {
  if (!vm.state.provider || !vm.state.selectedServiceItemId) {
    return
  }
  void router.push({
    path: '/payment/confirm',
    query: {
      provider_id: vm.state.provider.id,
      service_item_id: vm.state.selectedServiceItemId
    }
  })
}
</script>

<template>
  <main class="provider-page">
    <section class="screen-shell">
      <header class="top-row">
        <button type="button" class="back" @click="router.push('/services')">返回</button>
        <span>服务方详情</span>
      </header>

      <section v-if="vm.state.pageState === PageLoadState.Loading" class="state-card">正在加载服务方详情...</section>
      <section v-else-if="vm.state.pageState === PageLoadState.Error" class="state-card error">
        {{ vm.state.errorMessage }}
      </section>
      <section v-else-if="vm.state.pageState === PageLoadState.Empty" class="state-card">暂无可展示的服务项。</section>

      <template v-if="vm.state.pageState === PageLoadState.Success && vm.state.provider">
        <section class="hero">
          <div class="head">
            <div class="avatar">{{ vm.state.provider.displayName.slice(0, 1) }}</div>
            <div class="info">
              <h1>{{ vm.state.provider.displayName }}</h1>
              <p class="meta">
                <span>{{ vm.state.provider.ratingAvg.toFixed(1) }} 分</span>
                <span>{{ vm.state.provider.completedOrders }} 次服务</span>
                <span>{{ vm.state.provider.online ? '在线' : '离线' }}</span>
              </p>
              <div class="badges">
                <span class="badge">{{ vm.state.provider.verificationLabel }}</span>
                <span v-for="tag in vm.state.provider.tags" :key="tag" class="badge ghost">{{ tag }}</span>
              </div>
            </div>
          </div>
          <p class="bio">{{ vm.state.provider.bio }}</p>
        </section>

        <section class="panel">
          <h2>服务项目</h2>
          <div class="service-list">
            <button
              v-for="item in vm.state.serviceItems"
              :key="item.id"
              type="button"
              class="service-item"
              :class="{ active: vm.state.selectedServiceItemId === item.id }"
              @click="vm.selectServiceItem(item.id)"
            >
              <div class="title-row">
                <span>{{ item.title }}</span>
                <strong>¥{{ item.priceAmount }}</strong>
              </div>
              <p>{{ item.description }}</p>
            </button>
          </div>
        </section>

        <section class="panel">
          <h2>历史评价</h2>
          <article class="review">
            <strong>“会认真听，也很会让人放松。”</strong>
            <p>和她散步的时候没有被硬聊，节奏很舒服。适合下班心很乱、又不想一个人待着的时候。</p>
          </article>
          <article class="review">
            <strong>“很稳，很会接情绪。”</strong>
            <p>线上聊了半小时，结束后整个人松下来很多。没有被说教，体验很舒服。</p>
          </article>
        </section>
      </template>
    </section>

    <footer class="sheet" v-if="vm.state.pageState === PageLoadState.Success && vm.state.provider">
      <div class="sheet-top">
        <strong>发起订单</strong>
        <span>下单前需先选择服务项目与时长</span>
      </div>
      <div class="selected">{{ selectedServiceTitle(vm.state.selectedServiceItemId) }}</div>
      <div class="actions">
        <button type="button" class="ghost">收藏</button>
        <button type="button" class="primary" @click="goPaymentConfirm">确认并去支付</button>
      </div>
    </footer>
  </main>
</template>

<style scoped>
.provider-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top, rgba(74, 168, 196, 0.1), transparent 30%),
    linear-gradient(180deg, #0b1926 0%, #0d1f2d 42%, #08131f 100%);
  color: rgba(255, 255, 255, 0.88);
}

.screen-shell {
  width: min(100%, 390px);
  min-height: 100vh;
  margin: 0 auto;
  padding: 52px 16px 130px;
}

.top-row {
  display: flex;
  align-items: center;
  gap: 10px;
  color: rgba(255, 255, 255, 0.45);
  font-size: 12px;
}

.back {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.8);
  padding: 6px 10px;
}

.state-card {
  margin-top: 14px;
  border-radius: 18px;
  padding: 14px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.04);
}

.hero,
.panel {
  margin-top: 14px;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.045);
  padding: 16px;
}

.head {
  display: flex;
  gap: 14px;
}

.avatar {
  width: 78px;
  height: 78px;
  border-radius: 22px;
  border: 1px solid rgba(115, 213, 255, 0.28);
  background: linear-gradient(135deg, rgba(115, 213, 255, 0.28), rgba(115, 213, 255, 0.08));
  display: grid;
  place-items: center;
  font-size: 30px;
  flex-shrink: 0;
}

.info {
  min-width: 0;
}

.info h1 {
  margin: 0;
  font-size: 22px;
}

.meta {
  margin: 8px 0 0;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  font-size: 12px;
  color: rgba(239, 247, 251, 0.56);
}

.badges {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.badge {
  border-radius: 999px;
  padding: 6px 10px;
  background: rgba(115, 213, 255, 0.12);
  color: #93e0ff;
  font-size: 11px;
}

.badge.ghost {
  background: rgba(255, 255, 255, 0.06);
  color: rgba(239, 247, 251, 0.72);
}

.bio {
  margin: 14px 0 0;
  font-size: 14px;
  line-height: 1.8;
  color: rgba(239, 247, 251, 0.72);
}

.panel h2 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

.service-list {
  margin-top: 12px;
  display: grid;
  gap: 10px;
}

.service-item {
  text-align: left;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(10, 34, 48, 0.7);
  color: rgba(239, 247, 251, 0.88);
  padding: 12px;
}

.service-item.active {
  border-color: rgba(115, 213, 255, 0.35);
  background: rgba(115, 213, 255, 0.12);
}

.title-row {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  font-size: 14px;
}

.service-item p {
  margin: 8px 0 0;
  font-size: 12px;
  line-height: 1.7;
  color: rgba(239, 247, 251, 0.65);
}

.review {
  margin-top: 12px;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(10, 34, 48, 0.72);
  padding: 12px;
}

.review strong {
  font-size: 13px;
}

.review p {
  margin: 8px 0 0;
  font-size: 13px;
  line-height: 1.7;
  color: rgba(239, 247, 251, 0.65);
}

.sheet {
  position: fixed;
  left: 50%;
  bottom: 0;
  transform: translateX(-50%);
  width: min(100%, 390px);
  padding: 14px 16px calc(20px + env(safe-area-inset-bottom));
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(7, 18, 28, 0.96);
  backdrop-filter: blur(24px);
}

.sheet-top {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.sheet-top strong {
  font-size: 16px;
}

.sheet-top span {
  font-size: 12px;
  color: rgba(239, 247, 251, 0.52);
}

.selected {
  margin-top: 10px;
  border-radius: 12px;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.06);
  font-size: 13px;
}

.actions {
  margin-top: 12px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.actions button {
  border: none;
  border-radius: 14px;
  padding: 13px 10px;
  font-size: 14px;
  font-weight: 600;
}

.actions .ghost {
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.74);
}

.actions .primary {
  background: linear-gradient(135deg, #97e3ff, #58bee8);
  color: #082132;
}

.error {
  color: #fca5a5;
}
</style>
