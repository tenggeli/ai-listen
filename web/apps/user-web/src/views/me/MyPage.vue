<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { HttpAuthApi } from '../../api/AuthApi'
import { MyPageViewModel } from '../../application/profile/MyPageViewModel'
import HomeBottomNav, { type HomeNavItem } from '../../components/ai/HomeBottomNav.vue'
import { PageLoadState } from '../../domain/ai/PageLoadState'

const router = useRouter()
const vm = new MyPageViewModel(new HttpAuthApi(import.meta.env.VITE_AI_API_BASE_URL ?? '/api/v1'))
const navItems: HomeNavItem[] = [
  { key: 'home', label: '首页', icon: 'home' },
  { key: 'service', label: '服务', icon: 'square' },
  { key: 'join', label: '加入', icon: 'join' },
  { key: 'voice', label: '声音', icon: 'voice' },
  { key: 'profile', label: '我的', icon: 'profile', active: true }
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
  <main class="my-page">
    <div class="status-bar">
      <span>{{ headerTime }}</span>
      <span>我的</span>
    </div>

    <section class="screen-shell">
      <section class="hero-card">
        <div class="avatar">{{ vm.state.me?.nickname?.slice(0, 1) ?? '我' }}</div>
        <div class="hero-info">
          <h1>{{ vm.state.me?.nickname || vm.state.displayName || 'listener' }}</h1>
          <p>ID: {{ vm.state.me?.userId || '-' }}</p>
        </div>
      </section>

      <section v-if="vm.state.pageState === PageLoadState.Loading" class="state-card">正在加载用户信息...</section>
      <section v-else-if="vm.state.pageState === PageLoadState.Error" class="state-card error">
        {{ vm.state.errorMessage }}
      </section>

      <template v-if="vm.state.pageState === PageLoadState.Success && vm.state.me">
        <section class="card">
          <h2>用户信息</h2>
          <div class="line"><span>昵称</span><strong>{{ vm.state.me.nickname || '-' }}</strong></div>
          <div class="line"><span>城市</span><strong>{{ vm.state.me.city || '-' }}</strong></div>
          <div class="line"><span>MBTI</span><strong>{{ vm.state.me.mbti || '-' }}</strong></div>
          <div class="line"><span>兴趣标签</span><strong>{{ vm.state.me.interestTags.join('、') || '-' }}</strong></div>
        </section>

        <section class="card">
          <h2>功能入口</h2>
          <button type="button" class="entry" @click="router.push('/orders')">
            <span>订单入口</span>
            <span>查看我的订单 ></span>
          </button>
          <button type="button" class="entry" @click="router.push('/settings')">
            <span>设置入口</span>
            <span>账号与偏好设置 ></span>
          </button>
        </section>
      </template>
    </section>

    <HomeBottomNav :items="navItems" @select="onSelectNav" />
  </main>
</template>

<style scoped>
.my-page {
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

.hero-card,
.card,
.state-card {
  margin-top: 14px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.045);
  padding: 16px;
}

.hero-card {
  display: flex;
  align-items: center;
  gap: 14px;
}

.avatar {
  width: 64px;
  height: 64px;
  border-radius: 18px;
  border: 1px solid rgba(115, 213, 255, 0.3);
  background: linear-gradient(135deg, rgba(115, 213, 255, 0.28), rgba(115, 213, 255, 0.08));
  display: grid;
  place-items: center;
  font-size: 24px;
  flex-shrink: 0;
}

.hero-info h1 {
  margin: 0;
  font-size: 22px;
}

.hero-info p {
  margin: 6px 0 0;
  font-size: 12px;
  color: rgba(239, 247, 251, 0.6);
}

.card h2 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

.line {
  margin-top: 12px;
  display: flex;
  justify-content: space-between;
  gap: 10px;
  font-size: 14px;
}

.line span {
  color: rgba(239, 247, 251, 0.56);
}

.line strong {
  text-align: right;
}

.entry {
  margin-top: 10px;
  width: 100%;
  border: none;
  border-radius: 14px;
  padding: 14px 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: rgba(239, 247, 251, 0.86);
  display: flex;
  justify-content: space-between;
  gap: 8px;
  font-size: 13px;
}

.error {
  color: #fca5a5;
}
</style>
