<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import BreathingCore from '../../components/ai/BreathingCore.vue'
import HomeBottomNav, { type HomeNavItem } from '../../components/ai/HomeBottomNav.vue'
import HomeQuickActions from '../../components/ai/HomeQuickActions.vue'
import MatchCandidateCard from '../../components/ai/MatchCandidateCard.vue'
import { HomePageViewModel } from '../../application/ai/HomePageViewModel'
import { currentUserIdOrDemo } from '../../application/identity/AuthSession'
import { HttpAiApi } from '../../api/AiApi'
import type { AiHomeQuickAction } from '../../domain/ai/AiHomeQuickAction'
import { PageLoadState } from '../../domain/ai/PageLoadState'

const api = new HttpAiApi(import.meta.env.VITE_AI_API_BASE_URL ?? '/api/v1')
const vm = new HomePageViewModel(api, currentUserIdOrDemo())
const router = useRouter()
const navItems: HomeNavItem[] = [
  { key: 'home', label: '首页', icon: 'home', active: true },
  { key: 'service', label: '服务', icon: 'square' },
  { key: 'join', label: '加入', icon: 'join' },
  { key: 'voice', label: '声音', icon: 'voice' },
  { key: 'profile', label: '我的', icon: 'profile' }
]

onMounted(() => {
  void vm.initialize()
})

function onSubmit() {
  void vm.submitMatch()
}

function onSelectQuickAction(action: AiHomeQuickAction) {
  vm.applyQuickAction(action.prompt)
  if (action.route === '/chat' || action.route === '/sound') {
    void router.push(action.route)
  }
}

function onSelectNav(key: string) {
  if (key === 'join') {
    void router.push('/chat')
    return
  }
  if (key === 'home') {
    void router.push('/home')
    return
  }
  if (key === 'service') {
    void router.push('/services')
    return
  }
  if (key === 'voice') {
    void router.push('/sound')
  }
}
</script>

<template>
  <main class="home">
    <section class="screen-shell">
      <div class="bg-glow" />

      <header class="top-nav">
        <div class="brand-block">
          <p class="brand-name">listen</p>
          <p class="brand-sub">AI COMPANION ENTRY</p>
        </div>
        <div class="top-actions">
          <button type="button" class="icon-btn" aria-label="搜索">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <circle cx="11" cy="11" r="7" />
              <line x1="16.5" y1="16.5" x2="22" y2="22" />
            </svg>
          </button>
          <button type="button" class="icon-btn" aria-label="通知">
            <span class="badge" />
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
              <path d="M13.73 21a2 2 0 0 1-3.46 0" />
            </svg>
          </button>
        </div>
      </header>

      <section v-if="vm.state.homeState === PageLoadState.Success && vm.state.overview" class="greeting-section">
        <p class="greeting-time">{{ vm.state.overview.greetingPeriod }}</p>
        <h1 class="greeting-text">{{ vm.state.overview.greetingText }}</h1>
        <p class="greeting-sub">{{ vm.state.overview.greetingSubText }}</p>
      </section>
      <section v-else-if="vm.state.homeState === PageLoadState.Loading" class="status-message">
        正在加载首页内容...
      </section>
      <section v-else-if="vm.state.homeState === PageLoadState.Empty" class="status-message">
        暂无首页内容，请稍后重试。
      </section>
      <section v-else-if="vm.state.homeState === PageLoadState.Error" class="status-message error">
        {{ vm.state.errorMessage }}
      </section>

      <section v-if="vm.state.overview" class="status-card">
        <div class="status-left">
          <div class="mood-dot">{{ vm.state.overview.moodEmoji }}</div>
          <div>
            <p class="status-label">今日心情</p>
            <p class="status-value">{{ vm.state.overview.moodText }}</p>
          </div>
        </div>
        <div class="status-right">
          <p>{{ vm.state.overview.weatherText }}</p>
          <p>已陪伴 {{ vm.state.overview.companionDays }} 天</p>
        </div>
      </section>

      <HomeQuickActions
        v-if="vm.state.overview?.quickActions.length"
        class="quick-panel"
        :actions="vm.state.overview.quickActions"
        @select="onSelectQuickAction"
      />

      <div class="section-divider">今日陪伴</div>

      <section class="breath-section">
        <BreathingCore :state="vm.state.matchState" @click="onSubmit" />
        <p v-if="vm.state.overview" class="partner-hint">
          <span class="partner-dot" />
          <span>今晚在线陪伴者 <strong>{{ vm.state.overview.waitingCount }}</strong> 人</span>
        </p>
      </section>

      <section class="entry-panel">
        <label for="query">文本输入</label>
        <textarea
          id="query"
          v-model="vm.state.query"
          rows="3"
          placeholder="例如：今天加班到很晚，回家的路上突然很想找个人说说话。"
        />
        <div class="tips-row">
          <p v-if="vm.state.remainingState === PageLoadState.Success">今日剩余匹配次数：{{ vm.state.remaining }}</p>
          <p v-else-if="vm.state.remainingState === PageLoadState.Loading">正在同步匹配次数...</p>
          <p v-else-if="vm.state.remainingState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</p>
          <p v-else>点击呼吸球或主按钮开始。</p>
        </div>
      </section>

      <section class="result-panel">
        <div class="result-head">
          <div>
            <p class="result-eyebrow">推荐承接区</p>
            <h2>今晚可先从这几位开始</h2>
          </div>
          <router-link class="chat-link" to="/chat">进入 AI 对话</router-link>
        </div>

        <p v-if="vm.state.matchState === PageLoadState.Idle" class="result-message">输入你的诉求后开始匹配。</p>
        <p v-else-if="vm.state.matchState === PageLoadState.Loading" class="result-message">AI 正在理解你的意图并生成推荐...</p>
        <p v-else-if="vm.state.matchState === PageLoadState.Empty" class="result-message">没有匹配到结果，建议尝试高级筛选。</p>
        <p v-else-if="vm.state.matchState === PageLoadState.Error" class="result-message error">{{ vm.state.errorMessage }}</p>

        <div v-if="vm.state.matchState === PageLoadState.Success" class="candidate-list">
          <MatchCandidateCard
            v-for="item in vm.state.candidates"
            :key="item.providerId"
            :name="item.displayName"
            :reason="item.reason"
            :score="item.score"
          />
        </div>
      </section>

      <section class="action-bar">
        <button type="button" class="main-btn" :disabled="vm.state.matchState === PageLoadState.Loading" @click="onSubmit">
          开始 AI 匹配
        </button>
        <button type="button" class="filter-btn">高级筛选</button>
      </section>
    </section>

    <HomeBottomNav :items="navItems" @select="onSelectNav" />
  </main>
</template>

<style scoped>
.home {
  min-height: 100vh;
  background:
    radial-gradient(circle at top, rgba(74, 168, 196, 0.12), transparent 30%),
    linear-gradient(180deg, #0b1926 0%, #0d1f2d 42%, #08131f 100%);
  color: rgba(255, 255, 255, 0.88);
}

.screen-shell {
  position: relative;
  width: min(100%, 390px);
  min-height: 100vh;
  margin: 0 auto;
  padding: 22px 20px 110px;
  overflow: hidden;
}

.bg-glow {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse 60% 40% at 50% 18%, rgba(74, 168, 196, 0.1) 0%, transparent 70%),
    radial-gradient(ellipse 40% 60% at 20% 65%, rgba(40, 100, 140, 0.08) 0%, transparent 60%);
}

.top-nav,
.greeting-section,
.status-card,
.quick-panel,
.breath-section,
.entry-panel,
.result-panel,
.action-bar {
  position: relative;
  z-index: 1;
}

.top-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 18px;
}

.brand-name {
  margin: 0;
  font-size: 22px;
  font-weight: 300;
  letter-spacing: 0.24em;
  text-transform: lowercase;
}

.brand-sub {
  margin: 2px 0 0;
  font-size: 11px;
  letter-spacing: 0.1em;
  color: #8bdcff;
}

.top-actions {
  display: flex;
  gap: 10px;
}

.icon-btn {
  position: relative;
  display: grid;
  place-items: center;
  width: 36px;
  height: 36px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.56);
}

.icon-btn svg {
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.badge {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #8fdfff;
  box-shadow: 0 0 12px rgba(143, 223, 255, 0.8);
}

.greeting-section {
  margin-top: 18px;
  padding: 24px;
  border-radius: 28px;
  border: 1px solid rgba(145, 220, 255, 0.14);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.07), rgba(255, 255, 255, 0.03));
  backdrop-filter: blur(18px);
}

.greeting-time {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 4px;
  padding: 7px 12px;
  border-radius: 999px;
  background: rgba(115, 213, 255, 0.08);
  font-size: 12px;
  color: #8fdfff;
}

.greeting-time::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #8fdfff;
  box-shadow: 0 0 14px rgba(143, 223, 255, 0.8);
}

.greeting-text {
  margin: 14px 0 0;
  font-size: 30px;
  font-weight: 500;
  line-height: 1.22;
}

.greeting-sub {
  margin: 10px 0 0;
  font-size: 14px;
  color: rgba(237, 247, 251, 0.66);
  line-height: 1.8;
}

.status-message {
  margin-top: 18px;
  padding: 14px 16px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.04);
}

.status-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 18px;
  padding: 16px 18px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.045);
  backdrop-filter: blur(12px);
}

.status-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.mood-dot {
  display: grid;
  place-items: center;
  width: 48px;
  height: 48px;
  border-radius: 18px;
  background: linear-gradient(135deg, rgba(115, 213, 255, 0.26), rgba(115, 213, 255, 0.08));
  border: 1px solid rgba(115, 213, 255, 0.28);
}

.status-label,
.status-right p,
.result-eyebrow {
  margin: 0;
  font-size: 11px;
  letter-spacing: 0.06em;
  color: rgba(255, 255, 255, 0.3);
}

.status-value {
  margin: 3px 0 0;
  color: rgba(255, 255, 255, 0.78);
}

.status-right {
  text-align: right;
}

.status-right p + p {
  margin-top: 4px;
  color: rgba(237, 247, 251, 0.54);
}

.quick-panel {
  margin-top: 16px;
}

.error {
  color: #fca5a5;
}

.section-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 22px 0 10px;
  font-size: 10px;
  letter-spacing: 0.25em;
  color: rgba(255, 255, 255, 0.22);
}

.section-divider::before,
.section-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: rgba(255, 255, 255, 0.06);
}

.section-divider::before {
  margin-right: 14px;
}

.section-divider::after {
  margin-left: 14px;
}

.breath-section {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.partner-hint {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.42);
}

.partner-hint strong {
  color: #8fdfff;
}

.partner-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #8fdfff;
  box-shadow: 0 0 12px rgba(143, 223, 255, 0.7);
}

.entry-panel,
.result-panel {
  margin-top: 18px;
  padding: 18px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.045);
  backdrop-filter: blur(12px);
}

.entry-panel label {
  display: block;
  margin-bottom: 10px;
  font-size: 14px;
}

.entry-panel textarea {
  width: 100%;
  resize: vertical;
  min-height: 96px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 18px;
  padding: 14px 16px;
  background: rgba(10, 35, 49, 0.78);
  color: rgba(255, 255, 255, 0.88);
  font: inherit;
}

.entry-panel textarea::placeholder {
  color: rgba(255, 255, 255, 0.26);
}

.tips-row {
  margin-top: 10px;
  font-size: 12px;
  color: rgba(126, 200, 220, 0.76);
}

.tips-row p {
  margin: 0;
}

.result-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.result-head h2 {
  margin: 4px 0 0;
  font-size: 16px;
  font-weight: 500;
}

.chat-link {
  color: #93e0ff;
  text-decoration: none;
  font-size: 12px;
}

.result-message {
  margin: 14px 0 0;
  color: rgba(255, 255, 255, 0.6);
}

.candidate-list {
  display: grid;
  gap: 12px;
  margin-top: 14px;
}

.action-bar {
  position: sticky;
  bottom: 90px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 20px 0 18px;
  background: linear-gradient(180deg, rgba(13, 31, 45, 0), rgba(13, 31, 45, 0.96) 32%);
  z-index: 0;
  pointer-events: none;
}

.main-btn,
.filter-btn {
  border: none;
  cursor: pointer;
  font: inherit;
  pointer-events: auto;
}

.main-btn {
  width: 100%;
  height: 56px;
  border-radius: 18px;
  background: linear-gradient(135deg, #9be4ff, #5abfe8);
  color: #082031;
  font-weight: 700;
  letter-spacing: 0.08em;
  box-shadow: 0 12px 26px rgba(88, 190, 232, 0.24);
}

.main-btn:disabled {
  opacity: 0.6;
}

.filter-btn {
  padding: 10px 20px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: rgba(237, 247, 251, 0.74);
}

@media (max-width: 640px) {
  .screen-shell {
    padding-inline: 16px;
  }

  .greeting-text {
    font-size: 26px;
  }

  .status-card {
    align-items: flex-start;
    flex-direction: column;
  }

  .action-bar {
    bottom: 86px;
  }
}
</style>
