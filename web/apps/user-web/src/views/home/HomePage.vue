<script setup lang="ts">
import { onMounted } from 'vue'
import BreathingCore from '../../components/ai/BreathingCore.vue'
import MatchCandidateCard from '../../components/ai/MatchCandidateCard.vue'
import { HomePageViewModel } from '../../application/ai/HomePageViewModel'
import { HttpAiApi, MockAiApi } from '../../api/AiApi'
import { PageLoadState } from '../../domain/ai/PageLoadState'

const useMock = import.meta.env.VITE_USE_MOCK !== 'false'
const api = useMock ? new MockAiApi() : new HttpAiApi('/api/v1')
const vm = new HomePageViewModel(api, 'demo-user-001')

onMounted(() => {
  void vm.initialize()
})

function onSubmit() {
  void vm.submitMatch()
}
</script>

<template>
  <main class="home">
    <section class="hero">
      <h1>Listen AI 主入口</h1>
      <p>告诉我你现在想被怎样陪伴，我会先帮你完成一轮智能匹配。</p>
      <BreathingCore />
    </section>

    <section class="entry-panel">
      <label for="query">文本输入</label>
      <textarea id="query" v-model="vm.state.query" rows="4" placeholder="例如：今晚有点焦虑，想找一个能耐心倾听的人聊一会儿" />
      <div class="actions">
        <button type="button" :disabled="vm.state.matchState === PageLoadState.Loading" @click="onSubmit">开始 AI 匹配</button>
        <button type="button" class="ghost">语音输入（占位）</button>
        <button type="button" class="ghost">高级筛选（占位）</button>
        <router-link class="ghost link-btn" to="/chat">进入 AI 对话页</router-link>
      </div>
      <p class="tips" v-if="vm.state.remainingState === PageLoadState.Success">今日剩余匹配次数：{{ vm.state.remaining }}</p>
      <p class="tips" v-else-if="vm.state.remainingState === PageLoadState.Loading">正在加载匹配次数...</p>
      <p class="tips error" v-else-if="vm.state.remainingState === PageLoadState.Error">{{ vm.state.errorMessage }}</p>
      <p class="tips" v-else>点击开始以初始化匹配状态</p>
    </section>

    <section class="result-panel">
      <h2>推荐承接区（固定 3 人）</h2>
      <p v-if="vm.state.matchState === PageLoadState.Idle">输入你的诉求后开始匹配。</p>
      <p v-else-if="vm.state.matchState === PageLoadState.Loading">AI 正在理解你的意图并生成推荐...</p>
      <p v-else-if="vm.state.matchState === PageLoadState.Empty">没有匹配到结果，建议尝试高级筛选。</p>
      <p v-else-if="vm.state.matchState === PageLoadState.Error" class="error">{{ vm.state.errorMessage }}</p>

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
  </main>
</template>

<style scoped>
.home {
  max-width: 860px;
  margin: 0 auto;
  padding: 24px 16px 40px;
}

.hero h1 {
  margin: 0;
  font-size: 28px;
}

.hero p {
  margin-top: 8px;
  color: #475569;
}

.entry-panel,
.result-panel {
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 16px;
  margin-top: 16px;
  background: #ffffff;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
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
  display: flex;
  gap: 8px;
  margin-top: 10px;
  flex-wrap: wrap;
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

.ghost {
  background: #eff6ff;
  color: #1e3a8a;
}

.link-btn {
  display: inline-block;
  text-decoration: none;
  border-radius: 8px;
  padding: 8px 14px;
}

.tips {
  margin: 10px 0 0;
  color: #1e3a8a;
}

.error {
  color: #b91c1c;
}

.candidate-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 10px;
}

@media (max-width: 640px) {
  .hero h1 {
    font-size: 24px;
  }
}
</style>
