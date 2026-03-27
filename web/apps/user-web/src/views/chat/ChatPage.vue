<script setup lang="ts">
import { onMounted } from 'vue'
import { ChatPageViewModel } from '../../application/ai/ChatPageViewModel'
import { HttpAiApi, MockAiApi } from '../../api/AiApi'
import { PageLoadState } from '../../domain/ai/PageLoadState'

const useMock = import.meta.env.VITE_USE_MOCK !== 'false'
const api = useMock ? new MockAiApi() : new HttpAiApi('/api/v1')
const vm = new ChatPageViewModel(api, 'demo-user-001')

onMounted(() => {
  void vm.initialize()
})

function onSubmit() {
  void vm.submitMessage()
}
</script>

<template>
  <main class="chat-page">
    <section class="chat-header">
      <h1>AI 对话页</h1>
      <p>会话 ID：{{ vm.state.sessionId || '初始化中...' }}</p>
    </section>

    <section class="chat-body">
      <p v-if="vm.state.sessionState === PageLoadState.Loading">正在初始化 AI 会话...</p>
      <p v-else-if="vm.state.sessionState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</p>
      <p v-else-if="vm.state.sessionState === PageLoadState.Empty">会话已建立，先发送第一条消息吧。</p>

      <ul v-if="vm.state.sessionState === PageLoadState.Success" class="message-list">
        <li v-for="(item, index) in vm.state.session?.messages" :key="`${item.createdAt}-${index}`" :class="item.senderType">
          <span class="sender">{{ item.senderType }}</span>
          <p>{{ item.content }}</p>
        </li>
      </ul>
    </section>

    <section class="chat-input">
      <textarea v-model="vm.state.draft" rows="3" placeholder="输入你此刻想聊的话题" />
      <div class="actions">
        <button type="button" :disabled="vm.state.sendState === PageLoadState.Loading" @click="onSubmit">发送消息</button>
      </div>
      <p v-if="vm.state.sendState === PageLoadState.Loading">发送中...</p>
      <p v-else-if="vm.state.sendState === PageLoadState.Empty" class="tips">请输入内容后再发送。</p>
      <p v-else-if="vm.state.sendState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</p>
      <p v-else-if="vm.state.sendState === PageLoadState.Success" class="tips">发送成功</p>
      <p v-else class="tips">可以开始新的对话了。</p>
    </section>
  </main>
</template>

<style scoped>
.chat-page {
  max-width: 860px;
  margin: 0 auto;
  padding: 24px 16px 40px;
}

.chat-header h1 {
  margin: 0;
  font-size: 28px;
}

.chat-header p {
  margin-top: 8px;
  color: #475569;
}

.chat-body,
.chat-input {
  margin-top: 16px;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  background: #ffffff;
  padding: 16px;
}

.message-list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: grid;
  gap: 10px;
}

.message-list li {
  border-radius: 10px;
  padding: 10px 12px;
  background: #eff6ff;
}

.message-list li.user {
  background: #dbeafe;
}

.sender {
  display: block;
  font-size: 12px;
  color: #1e3a8a;
}

.message-list p {
  margin: 6px 0 0;
}

textarea {
  width: 100%;
  resize: vertical;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  padding: 10px;
  font: inherit;
}

.actions {
  margin-top: 10px;
}

button {
  border: none;
  border-radius: 8px;
  padding: 8px 14px;
  font-weight: 600;
  background: #1d4ed8;
  color: #fff;
  cursor: pointer;
}

button:disabled {
  background: #93c5fd;
  cursor: wait;
}

.tips {
  margin: 10px 0 0;
  color: #1e3a8a;
}

.error {
  margin: 10px 0 0;
  color: #b91c1c;
}

@media (max-width: 640px) {
  .chat-header h1 {
    font-size: 24px;
  }
}
</style>
