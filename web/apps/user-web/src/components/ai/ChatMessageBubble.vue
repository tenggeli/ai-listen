<script setup lang="ts">
import type { AiMessage } from '../../domain/ai/AiMessage'

defineProps<{
  message: AiMessage
}>()

const logoSvg = `
  <path d="M37 37 C33 41 33 49 37 53" stroke="#7ec8dc" stroke-width="2.6" stroke-linecap="round" fill="none"/>
  <path d="M30 31 C24 37 24 53 30 59" stroke="#aadde9" stroke-width="2.2" stroke-linecap="round" fill="none"/>
  <path d="M64 62 C58 72 46 74 38 66 C28 56 28 40 36 31 C44 22 58 22 66 32 L64 36" stroke="#4aa8c4" stroke-width="3.2" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
  <path d="M64 36 C70 40 70 52 64 56" stroke="#4aa8c4" stroke-width="3" stroke-linecap="round" fill="none"/>
  <path d="M60 40 C62 50 56 60 50 62 C47 63 46 66 50 68 C53 70 57 68 57 64" stroke="#4aa8c4" stroke-width="2.6" stroke-linecap="round" fill="none"/>
`
</script>

<template>
  <div class="message-row" :class="message.senderType === 'user' ? 'is-user' : 'is-ai'">
    <div v-if="message.senderType !== 'user'" class="avatar" aria-hidden="true">
      <svg viewBox="0 0 90 90" fill="none" v-html="logoSvg" />
    </div>
    <div class="message-content">
      <div class="bubble" :class="message.senderType === 'user' ? 'user-bubble' : 'ai-bubble'">
        {{ message.content }}
      </div>
      <p class="time">{{ message.displayTime() }}</p>
    </div>
  </div>
</template>

<style scoped>
.message-row {
  display: flex;
  gap: 10px;
  align-items: flex-end;
}

.message-row.is-user {
  justify-content: flex-end;
}

.avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 12px;
  flex-shrink: 0;
  border: 1px solid rgba(74, 168, 196, 0.24);
  background: rgba(74, 168, 196, 0.14);
}

.avatar svg {
  width: 16px;
  height: 16px;
}

.message-content {
  max-width: min(78%, 540px);
}

.message-row.is-user .message-content {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.bubble {
  padding: 14px 16px;
  font-size: 14px;
  line-height: 1.9;
  white-space: pre-wrap;
  backdrop-filter: blur(12px);
}

.ai-bubble {
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px 20px 20px 20px;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(237, 247, 251, 0.82);
}

.user-bubble {
  border: 1px solid rgba(115, 213, 255, 0.28);
  border-radius: 20px 6px 20px 20px;
  background: rgba(115, 213, 255, 0.14);
  color: #f6fcff;
}

.time {
  margin: 4px 4px 0;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.18);
}
</style>
