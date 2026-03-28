<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { currentUserIdOrDemo } from '../../application/identity/AuthSession'
import HomeBottomNav, { type HomeNavItem } from '../../components/ai/HomeBottomNav.vue'
import SoundCategoryTabs from '../../components/ai/SoundCategoryTabs.vue'
import SoundNowPlayingCard from '../../components/ai/SoundNowPlayingCard.vue'
import SoundTrackList from '../../components/ai/SoundTrackList.vue'
import { SoundPageViewModel } from '../../application/ai/SoundPageViewModel'
import { HttpSoundApi, MockSoundApi } from '../../api/SoundApi'
import { PageLoadState } from '../../domain/ai/PageLoadState'

const useMock = import.meta.env.VITE_USE_MOCK !== 'false'
const api = useMock ? new MockSoundApi() : new HttpSoundApi('/api/v1')
const vm = new SoundPageViewModel(api, currentUserIdOrDemo())
const router = useRouter()
const headerTime = computed(() =>
  new Intl.DateTimeFormat('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  }).format(new Date())
)
const navItems: HomeNavItem[] = [
  { key: 'home', label: '首页', icon: 'home' },
  { key: 'square', label: '广场', icon: 'square' },
  { key: 'join', label: '加入', icon: 'join' },
  { key: 'voice', label: '声音', icon: 'voice', active: true },
  { key: 'profile', label: '我的', icon: 'profile' }
]

onMounted(() => {
  void vm.initialize()
})

function onSelectNav(key: string) {
  if (key === 'home') {
    void router.push('/home')
    return
  }
  if (key === 'join') {
    void router.push('/chat')
    return
  }
  if (key === 'voice') {
    void router.push('/sound')
  }
}
</script>

<template>
  <main class="sound-page">
    <div class="bg-glow" />

    <div class="status-bar">
      <span>{{ headerTime }}</span>
      <span>···</span>
    </div>

    <section class="screen-shell">
      <header class="page-header">
        <p class="hero-label">Night Companion</p>
        <h1 class="page-title">{{ vm.state.aggregate?.title ?? '声音' }}</h1>
        <p class="page-sub">{{ vm.state.aggregate?.subtitle ?? '用声音抚慰此刻的你' }}</p>
      </header>

      <SoundCategoryTabs
        v-if="vm.state.aggregate?.categories.length"
        class="tabs-section"
        :categories="vm.state.aggregate.categories"
        :selected-key="vm.state.selectedCategoryKey"
        @select="vm.selectCategory"
      />

      <section v-if="vm.state.pageState === PageLoadState.Loading" class="state-card">
        正在加载声音内容...
      </section>
      <section v-else-if="vm.state.pageState === PageLoadState.Error" class="state-card error">
        {{ vm.state.errorMessage }}
      </section>
      <section v-else-if="vm.state.pageState === PageLoadState.Empty" class="state-card">
        暂无可播放内容，稍后再来听听。
      </section>

      <SoundNowPlayingCard
        v-if="vm.state.pageState === PageLoadState.Success && vm.state.currentTrack"
        class="now-playing-section"
        :track="vm.state.currentTrack"
        :progress-text="vm.state.aggregate?.currentProgressText ?? '00:00'"
        :total-duration-text="vm.state.aggregate?.totalDurationText ?? vm.state.currentTrack.durationText"
        :is-playing="vm.state.isPlaying"
        @toggle="vm.togglePlayback"
      />

      <section v-if="vm.state.pageState === PageLoadState.Success" class="recommend-section">
        <p class="section-title">为你推荐</p>
        <SoundTrackList
          :tracks="vm.state.filteredTracks"
          :current-track-id="vm.state.currentTrack?.id ?? ''"
          :is-playing="vm.state.isPlaying"
          @select="vm.selectTrack"
        />
      </section>
    </section>

    <HomeBottomNav :items="navItems" @select="onSelectNav" />
  </main>
</template>

<style scoped>
.sound-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top, rgba(74, 168, 196, 0.08), transparent 32%),
    linear-gradient(180deg, #0b1926 0%, #0d1f2d 42%, #08131f 100%);
  color: rgba(255, 255, 255, 0.82);
}

.bg-glow {
  position: fixed;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse 50% 40% at 50% 25%, rgba(74, 168, 196, 0.08) 0%, transparent 70%),
    radial-gradient(ellipse 30% 50% at 15% 80%, rgba(74, 120, 196, 0.05) 0%, transparent 60%);
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
  padding: 0 24px;
  color: rgba(255, 255, 255, 0.3);
  transform: translateX(-50%);
  font-size: 12px;
}

.screen-shell {
  position: relative;
  z-index: 1;
  width: min(100%, 390px);
  min-height: 100vh;
  margin: 0 auto;
  padding: 52px 24px 96px;
}

.page-header {
  margin-top: 18px;
  padding: 24px;
  border-radius: 28px;
  border: 1px solid rgba(145, 220, 255, 0.14);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.07), rgba(255, 255, 255, 0.03));
  backdrop-filter: blur(18px);
}

.hero-label {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  padding: 7px 12px;
  border-radius: 999px;
  background: rgba(115, 213, 255, 0.08);
  color: #8fdfff;
  font-size: 12px;
}

.hero-label::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #8fdfff;
  box-shadow: 0 0 14px rgba(143, 223, 255, 0.8);
}

.page-title {
  margin: 14px 0 0;
  font-size: 28px;
  font-weight: 500;
  line-height: 1.22;
}

.page-sub {
  margin: 10px 0 0;
  font-size: 14px;
  color: rgba(237, 247, 251, 0.66);
  line-height: 1.8;
}

.tabs-section,
.state-card,
.now-playing-section,
.recommend-section {
  margin-top: 16px;
}

.state-card {
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.045);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 10px;
  font-size: 13px;
  letter-spacing: 0.16em;
  color: rgba(255, 255, 255, 0.35);
}

.section-title::after {
  content: '';
  flex: 1;
  height: 1px;
  background: rgba(255, 255, 255, 0.06);
}

.error {
  color: #fca5a5;
}

@media (max-width: 640px) {
  .screen-shell {
    padding-inline: 16px;
  }

  .status-bar {
    padding-inline: 20px;
  }
}
</style>
