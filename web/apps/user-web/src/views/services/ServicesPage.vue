<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ServicesPageViewModel } from '../../application/service/ServicesPageViewModel'
import HomeBottomNav, { type HomeNavItem } from '../../components/ai/HomeBottomNav.vue'
import { MockServiceDiscoveryApi, HttpServiceDiscoveryApi } from '../../api/ServiceDiscoveryApi'
import { PageLoadState } from '../../domain/ai/PageLoadState'

const useMock = import.meta.env.VITE_USE_MOCK === 'true'
const api = useMock ? new MockServiceDiscoveryApi() : new HttpServiceDiscoveryApi('/api/v1')
const vm = new ServicesPageViewModel(api)
const router = useRouter()
const navItems: HomeNavItem[] = [
  { key: 'home', label: '首页', icon: 'home' },
  { key: 'service', label: '服务', icon: 'square', active: true },
  { key: 'join', label: '加入', icon: 'join' },
  { key: 'voice', label: '声音', icon: 'voice' },
  { key: 'profile', label: '我的', icon: 'profile' }
]
const headerTime = computed(() =>
  new Intl.DateTimeFormat('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  }).format(new Date())
)

onMounted(() => {
  void vm.initialize()
})

function openDetail(providerId: string): void {
  void router.push(`/providers/${providerId}`)
}

function onSearch(): void {
  void vm.search()
}

function onSelectCategory(categoryId: string): void {
  void vm.selectCategory(categoryId)
}

function onSelectNav(key: string): void {
  if (key === 'home') {
    void router.push('/home')
    return
  }
  if (key === 'service') {
    void router.push('/services')
    return
  }
  if (key === 'join') {
    void router.push('/chat')
    return
  }
  if (key === 'voice') {
    void router.push('/sound')
    return
  }
  if (key === 'profile') {
    void router.push('/me')
  }
}
</script>

<template>
  <main class="services-page">
    <div class="status-bar">
      <span>{{ headerTime }}</span>
      <span>服务</span>
    </div>

    <section class="screen-shell">
      <header class="page-header">
        <h1>服务</h1>
        <p>查看 AI 推荐承接结果，也可以自己筛选今天更想要的陪伴方式。</p>
      </header>

      <div class="search-box">
        <input
          v-model="vm.state.searchKeyword"
          type="text"
          placeholder="搜索服务方、场景或你想要的陪伴方式"
          @keyup.enter="onSearch"
        />
        <button type="button" @click="onSearch">搜索</button>
      </div>

      <div class="category-row">
        <button
          v-for="category in vm.state.categories"
          :key="category.id"
          type="button"
          class="chip"
          :class="{ active: vm.state.selectedCategoryId === category.id }"
          @click="onSelectCategory(category.id)"
        >
          {{ category.name }}
        </button>
      </div>

      <section class="banner">
        <strong>今晚更适合“轻聊天 + 同城低压陪伴”</strong>
        <p>根据你的首页输入与性格偏好，优先为你展示响应快、标签温和、夜间在线的服务方。</p>
      </section>

      <section v-if="vm.state.pageState === PageLoadState.Loading" class="state-card">正在加载服务列表...</section>
      <section v-else-if="vm.state.pageState === PageLoadState.Error" class="state-card error">
        {{ vm.state.errorMessage }}
      </section>
      <section v-else-if="vm.state.pageState === PageLoadState.Empty" class="state-card">暂无服务方，换个筛选试试。</section>

      <section v-if="vm.state.pageState === PageLoadState.Success" class="provider-list">
        <article v-for="provider in vm.state.providers" :key="provider.id" class="provider-card">
          <div class="card-head">
            <div class="avatar">{{ provider.displayName.slice(0, 1) }}</div>
            <div class="info">
              <div class="name-row">
                <h3>{{ provider.displayName }}</h3>
                <span class="badge">{{ provider.verificationLabel }}</span>
              </div>
              <p class="meta">
                <span>{{ provider.ratingAvg.toFixed(1) }} 分</span>
                <span>{{ provider.completedOrders }} 次服务</span>
                <span>{{ provider.online ? '在线' : '离线' }}</span>
              </p>
              <p class="bio">{{ provider.bio }}</p>
              <div class="tags">
                <span v-for="tag in provider.tags" :key="tag" class="tag">{{ tag }}</span>
              </div>
            </div>
          </div>
          <div class="card-bottom">
            <p class="price">¥{{ provider.priceFrom }} <span>/ {{ provider.priceUnit }} 起</span></p>
            <div class="actions">
              <button type="button" class="ghost" @click="openDetail(provider.id)">查看详情</button>
              <button type="button" class="primary" @click="openDetail(provider.id)">发起订单</button>
            </div>
          </div>
        </article>
      </section>
    </section>

    <HomeBottomNav :items="navItems" @select="onSelectNav" />
  </main>
</template>

<style scoped>
.services-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top, rgba(74, 168, 196, 0.1), transparent 30%),
    linear-gradient(180deg, #0b1926 0%, #0d1f2d 42%, #08131f 100%);
  color: rgba(255, 255, 255, 0.88);
}

.status-bar {
  position: fixed;
  top: 0;
  left: 50%;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: min(100%, 390px);
  height: 44px;
  padding: 0 22px;
  color: rgba(255, 255, 255, 0.35);
  transform: translateX(-50%);
  font-size: 12px;
}

.screen-shell {
  width: min(100%, 390px);
  min-height: 100vh;
  margin: 0 auto;
  padding: 56px 16px 104px;
}

.page-header {
  padding: 18px;
  border-radius: 22px;
  border: 1px solid rgba(145, 220, 255, 0.14);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.07), rgba(255, 255, 255, 0.03));
}

.page-header h1 {
  margin: 0;
  font-size: 26px;
  font-weight: 500;
}

.page-header p {
  margin: 8px 0 0;
  color: rgba(239, 247, 251, 0.62);
  font-size: 13px;
  line-height: 1.7;
}

.search-box {
  margin-top: 14px;
  display: flex;
  gap: 8px;
}

.search-box input {
  flex: 1;
  border-radius: 14px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.88);
  padding: 12px 14px;
}

.search-box button {
  border: none;
  border-radius: 14px;
  padding: 0 14px;
  background: linear-gradient(135deg, #98e3ff, #59bee7);
  color: #082132;
  font-weight: 600;
}

.category-row {
  margin-top: 14px;
  display: flex;
  gap: 8px;
  overflow-x: auto;
  scrollbar-width: none;
}

.category-row::-webkit-scrollbar {
  display: none;
}

.chip {
  flex-shrink: 0;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.7);
  padding: 8px 13px;
  font-size: 12px;
}

.chip.active {
  background: rgba(115, 213, 255, 0.14);
  border-color: rgba(115, 213, 255, 0.3);
  color: #93e0ff;
}

.banner {
  margin-top: 14px;
  padding: 16px;
  border-radius: 20px;
  border: 1px solid rgba(115, 213, 255, 0.16);
  background: linear-gradient(135deg, rgba(20, 53, 74, 0.95), rgba(11, 28, 41, 0.96));
}

.banner strong {
  display: block;
  font-size: 15px;
}

.banner p {
  margin: 8px 0 0;
  color: rgba(239, 247, 251, 0.65);
  font-size: 13px;
  line-height: 1.7;
}

.state-card {
  margin-top: 14px;
  border-radius: 18px;
  padding: 14px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.04);
}

.provider-list {
  margin-top: 14px;
  display: grid;
  gap: 12px;
}

.provider-card {
  border-radius: 22px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.05);
  padding: 14px;
}

.card-head {
  display: flex;
  gap: 12px;
}

.avatar {
  width: 60px;
  height: 60px;
  border-radius: 18px;
  border: 1px solid rgba(115, 213, 255, 0.28);
  background: linear-gradient(135deg, rgba(115, 213, 255, 0.28), rgba(115, 213, 255, 0.08));
  display: grid;
  place-items: center;
  font-size: 22px;
  flex-shrink: 0;
}

.info {
  min-width: 0;
}

.name-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
}

.name-row h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

.badge {
  border-radius: 999px;
  padding: 5px 8px;
  font-size: 11px;
  color: #93e0ff;
  background: rgba(115, 213, 255, 0.12);
  white-space: nowrap;
}

.meta {
  margin: 6px 0 0;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  font-size: 12px;
  color: rgba(239, 247, 251, 0.55);
}

.bio {
  margin: 10px 0 0;
  font-size: 13px;
  line-height: 1.7;
  color: rgba(239, 247, 251, 0.7);
}

.tags {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  border-radius: 999px;
  padding: 5px 9px;
  font-size: 12px;
  color: rgba(239, 247, 251, 0.6);
  background: rgba(255, 255, 255, 0.06);
}

.card-bottom {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.price {
  margin: 0;
  font-size: 18px;
}

.price span {
  font-size: 12px;
  color: rgba(239, 247, 251, 0.45);
}

.actions {
  display: flex;
  gap: 8px;
}

.actions button {
  border: none;
  border-radius: 12px;
  padding: 10px 12px;
  font-size: 12px;
}

.ghost {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.75);
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.primary {
  background: linear-gradient(135deg, #98e3ff, #59bee7);
  color: #082132;
  font-weight: 600;
}

.error {
  color: #fca5a5;
}
</style>
