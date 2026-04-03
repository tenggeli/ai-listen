<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { ChatPageViewModel } from '../../application/ai/ChatPageViewModel'
import { currentUserIdOrDemo } from '../../application/identity/AuthSession'
import { HttpAiApi, MockAiApi } from '../../api/AiApi'
import ChatMessageBubble from '../../components/ai/ChatMessageBubble.vue'
import ChatTypingIndicator from '../../components/ai/ChatTypingIndicator.vue'
import { PageLoadState } from '../../domain/ai/PageLoadState'

const useMock = import.meta.env.VITE_USE_MOCK !== 'false'
const api = useMock ? new MockAiApi() : new HttpAiApi('/api/v1')
const vm = new ChatPageViewModel(api, currentUserIdOrDemo())
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

function onSubmit() {
  void vm.submitMessage()
}

function onUseQuickReply(reply: string) {
  vm.useQuickReply(reply)
}
</script>

<template>
  <main class="chat-page">
    <div class="bg-glow" />

    <div class="status-bar">
      <span>{{ headerTime }}</span>
      <span>···</span>
    </div>

    <header class="companion-bar">
      <RouterLink class="back-btn" to="/home" aria-label="返回首页">
        <span>‹</span>
      </RouterLink>
      <div class="companion-avatar" aria-hidden="true">
        <svg viewBox="0 0 90 90" fill="none">
          <path d="M37 37 C33 41 33 49 37 53" stroke="#7ec8dc" stroke-width="2.6" stroke-linecap="round" fill="none" />
          <path d="M30 31 C24 37 24 53 30 59" stroke="#aadde9" stroke-width="2.2" stroke-linecap="round" fill="none" />
          <path d="M64 62 C58 72 46 74 38 66 C28 56 28 40 36 31 C44 22 58 22 66 32 L64 36" stroke="#4aa8c4" stroke-width="3.2" stroke-linecap="round" stroke-linejoin="round" fill="none" />
          <path d="M64 36 C70 40 70 52 64 56" stroke="#4aa8c4" stroke-width="3" stroke-linecap="round" fill="none" />
          <path d="M60 40 C62 50 56 60 50 62 C47 63 46 66 50 68 C53 70 57 68 57 64" stroke="#4aa8c4" stroke-width="2.6" stroke-linecap="round" fill="none" />
        </svg>
      </div>
      <div class="companion-info">
        <p class="companion-name">listen</p>
        <p class="companion-status"><span class="status-dot" />在线 · 随时倾听</p>
      </div>
      <button type="button" class="more-btn" aria-label="更多操作">⋯</button>
    </header>

    <section class="chat-area">
      <div class="date-divider">今晚 {{ headerTime }}</div>

      <div v-if="vm.state.sessionState === PageLoadState.Loading" class="state-card">
        正在初始化陪伴会话...
      </div>
      <div v-else-if="vm.state.sessionState === PageLoadState.Error" class="state-card error">
        {{ vm.state.errorMessage }}
      </div>
      <div v-else-if="vm.state.sessionState === PageLoadState.Empty" class="state-card">
        <p>会话已建立，先说说你现在最想被听见的那句话。</p>
        <div class="quick-replies">
          <button
            v-for="reply in vm.state.quickReplies"
            :key="reply"
            type="button"
            class="quick-btn"
            @click="onUseQuickReply(reply)"
          >
            {{ reply }}
          </button>
        </div>
      </div>
      <ChatTypingIndicator v-else-if="vm.state.sendState === PageLoadState.Loading" />

      <div v-if="vm.state.sessionState === PageLoadState.Success" class="message-list">
        <ChatMessageBubble
          v-for="(item, index) in vm.state.session?.messages"
          :key="`${item.createdAt}-${index}`"
          :message="item"
        />

        <div class="quick-replies" v-if="vm.state.quickReplies.length">
          <button
            v-for="reply in vm.state.quickReplies"
            :key="reply"
            type="button"
            class="quick-btn"
            @click="onUseQuickReply(reply)"
          >
            {{ reply }}
          </button>
        </div>

        <ChatTypingIndicator v-if="vm.state.sendState === PageLoadState.Loading" />
      </div>
    </section>

    <footer class="input-bar">
      <div class="input-inner">
        <button type="button" class="voice-btn" aria-label="语音输入">🎤</button>
        <textarea
          v-model="vm.state.draft"
          class="input-field"
          rows="1"
          placeholder="想说点什么……"
          @keydown.enter.exact.prevent="onSubmit"
        />
        <button
          type="button"
          class="send-btn"
          :disabled="vm.state.sendState === PageLoadState.Loading"
          @click="onSubmit"
        >
          ›
        </button>
      </div>
      <p v-if="vm.state.sendState === PageLoadState.Empty" class="footer-tips">请输入内容后再发送。</p>
      <p v-else-if="vm.state.sendState === PageLoadState.Error" class="footer-tips error">{{ vm.state.errorMessage }}</p>
      <p v-else-if="vm.state.sendState === PageLoadState.Success" class="footer-tips">listen 已回复。</p>
      <p v-else-if="vm.state.sendState === PageLoadState.Loading" class="footer-tips">listen 正在回复...</p>
      <p v-else class="footer-tips">按 Enter 可直接发送，Shift + Enter 换行。</p>
    </footer>
  </main>
</template>

<style scoped>
.chat-page {
  position: relative;
  min-height: 100vh;
  margin: 0 auto;
  background:
    radial-gradient(circle at top, rgba(74, 168, 196, 0.08), transparent 32%),
    linear-gradient(180deg, #0b1926 0%, #0d1f2d 42%, #08131f 100%);
  color: rgba(255, 255, 255, 0.88);
}

.bg-glow {
  position: fixed;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse 60% 40% at 50% 28%, rgba(74, 168, 196, 0.06) 0%, transparent 70%),
    radial-gradient(ellipse 30% 40% at 50% 84%, rgba(74, 120, 160, 0.04) 0%, transparent 60%);
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

.companion-bar {
  position: fixed;
  top: 44px;
  left: 50%;
  z-index: 9;
  display: flex;
  align-items: center;
  gap: 12px;
  width: min(100%, 390px);
  margin-top: 14px;
  padding: 16px 18px;
  border-radius: 24px;
  border: 1px solid rgba(145, 220, 255, 0.14);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.07), rgba(255, 255, 255, 0.03));
  backdrop-filter: blur(10px);
  transform: translateX(-50%);
}

.back-btn,
.more-btn {
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
}

.back-btn span {
  font-size: 22px;
  line-height: 1;
}

.companion-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 46px;
  height: 46px;
  border: 1px solid rgba(115, 213, 255, 0.3);
  border-radius: 18px;
  background: linear-gradient(135deg, rgba(115, 213, 255, 0.28), rgba(115, 213, 255, 0.08));
  flex-shrink: 0;
}

.companion-avatar svg {
  width: 22px;
  height: 22px;
}

.companion-info {
  flex: 1;
}

.companion-name {
  margin: 0;
  font-size: 14px;
  letter-spacing: 0.08em;
}

.companion-status {
  margin: 2px 0 0;
  font-size: 11px;
  color: rgba(74, 168, 196, 0.68);
  letter-spacing: 0.08em;
}

.status-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  margin-right: 4px;
  border-radius: 50%;
  background: #4aa8c4;
  animation: dotPulse 2s ease-in-out infinite;
}

.chat-area {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: min(100%, 390px);
  min-height: 100vh;
  margin: 0 auto;
  padding: 126px 20px 146px;
}

.date-divider {
  text-align: center;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.18);
  letter-spacing: 0.18em;
}

.message-list {
  display: grid;
  gap: 16px;
}

.state-card {
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.72);
  line-height: 1.7;
}

.quick-replies {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.quick-btn {
  padding: 7px 14px;
  border: 1px solid rgba(115, 213, 255, 0.22);
  border-radius: 999px;
  background: rgba(115, 213, 255, 0.1);
  color: #90dfff;
  cursor: pointer;
}

.input-bar {
  position: fixed;
  bottom: 0;
  left: 50%;
  z-index: 5;
  width: min(100%, 390px);
  padding: 14px 18px calc(24px + env(safe-area-inset-bottom));
  background: rgba(7, 18, 28, 0.96);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(24px);
  transform: translateX(-50%);
}

.input-inner {
  display: flex;
  gap: 10px;
  align-items: flex-end;
}

.input-field {
  width: 100%;
  min-height: 42px;
  max-height: 120px;
  padding: 10px 18px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 21px;
  background: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.86);
  font: inherit;
  line-height: 1.5;
  resize: vertical;
}

.input-field::placeholder {
  color: rgba(255, 255, 255, 0.22);
}

.voice-btn,
.send-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  border: none;
  border-radius: 16px;
  cursor: pointer;
  flex-shrink: 0;
}

.voice-btn {
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.5);
}

.send-btn {
  background: linear-gradient(135deg, #97e3ff, #58bee8);
  color: #082132;
  font-size: 28px;
  box-shadow: 0 12px 24px rgba(88, 190, 232, 0.24);
}

.send-btn:disabled {
  opacity: 0.4;
  cursor: default;
}

.footer-tips {
  margin: 10px 6px 0;
  font-size: 12px;
  color: rgba(126, 200, 220, 0.72);
}

.error {
  color: #fca5a5;
}

@keyframes dotPulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.3;
  }
}

@media (max-width: 640px) {
  .chat-area {
    padding-inline: 16px;
  }

  .status-bar {
    padding-inline: 20px;
  }
}
</style>
