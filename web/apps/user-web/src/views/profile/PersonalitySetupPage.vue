<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { HttpAuthApi } from '../../api/AuthApi'
import { PersonalitySetupViewModel } from '../../application/identity/PersonalitySetupViewModel'

const router = useRouter()
const vm = new PersonalitySetupViewModel(new HttpAuthApi(import.meta.env.VITE_AI_API_BASE_URL ?? '/api/v1'))
const interestOptions = [
  'deep_chat',
  'healing_audio',
  'movie',
  'city_walk',
  'coffee',
  'sleep_relax',
  'light_social',
  'emotion_support'
]
const mbtiOptions = ['INFJ', 'INFP', 'ENFP', 'ENFJ', 'INTJ', 'INTP', 'ISFJ', 'ESFP']

onMounted(() => {
  void vm.initialize().catch((error: unknown) => {
    vm.state.errorMessage = error instanceof Error ? error.message : '加载性格配置失败'
  })
})

async function onSubmit() {
  try {
    const route = await vm.submitSave()
    await router.push(route)
  } catch (error) {
    vm.state.errorMessage = error instanceof Error ? error.message : '保存失败，请稍后重试'
  }
}

async function onSkip() {
  try {
    const route = await vm.submitSkip()
    await router.push(route)
  } catch (error) {
    vm.state.errorMessage = error instanceof Error ? error.message : '跳过失败，请稍后重试'
  }
}
</script>

<template>
  <main class="personality-page">
    <section class="card">
      <p class="step">步骤 2 / 2</p>
      <h1>性格设置</h1>
      <p class="desc">告诉 Listen 你的偏好，后续推荐和对话会更贴近你。</p>

      <div v-if="vm.state.loading" class="status">正在加载你的偏好...</div>

      <template v-else>
        <h2>兴趣标签</h2>
        <div class="tag-grid">
          <button
            v-for="tag in interestOptions"
            :key="tag"
            type="button"
            class="tag"
            :class="{ active: vm.state.selectedTags.includes(tag) }"
            @click="vm.toggleTag(tag)"
          >
            {{ tag }}
          </button>
        </div>

        <h2 class="section-title">MBTI</h2>
        <div class="mbti-grid">
          <button
            v-for="item in mbtiOptions"
            :key="item"
            type="button"
            class="mbti"
            :class="{ active: vm.state.selectedMbti === item }"
            @click="vm.toggleMbti(item)"
          >
            {{ item }}
          </button>
        </div>

        <p v-if="vm.state.errorMessage" class="error">{{ vm.state.errorMessage }}</p>

        <div class="action-row">
          <button type="button" class="secondary" :disabled="vm.state.saving" @click="onSkip">以后再设置</button>
          <button type="button" class="primary" :disabled="vm.state.saving" @click="onSubmit">完成并进入首页</button>
        </div>
      </template>
    </section>
  </main>
</template>

<style scoped>
.personality-page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  padding: 20px;
  background:
    radial-gradient(circle at top left, rgba(34, 118, 150, 0.35), transparent 44%),
    linear-gradient(180deg, #08131c 0%, #08131f 100%);
}

.card {
  width: min(100%, 420px);
  margin-top: 18px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 26px;
  padding: 22px;
  background: rgba(255, 255, 255, 0.05);
  height: fit-content;
}

.step {
  margin: 0;
  color: #8fdfff;
  font-size: 12px;
}

h1 {
  margin: 8px 0;
}

.desc {
  margin: 0;
  color: rgba(255, 255, 255, 0.72);
  line-height: 1.7;
  font-size: 13px;
}

h2 {
  margin: 14px 0 8px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.9);
}

.section-title {
  margin-top: 16px;
}

.status {
  margin-top: 16px;
  color: rgba(255, 255, 255, 0.75);
}

.tag-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.tag,
.mbti {
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.8);
  padding: 10px 12px;
}

.tag.active,
.mbti.active {
  border-color: rgba(115, 213, 255, 0.4);
  background: rgba(115, 213, 255, 0.18);
  color: #9be4ff;
}

.mbti-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.error {
  margin-top: 10px;
  color: #ffb6b6;
  font-size: 12px;
}

.action-row {
  display: flex;
  gap: 10px;
  margin-top: 16px;
}

button {
  flex: 1;
  border: none;
  border-radius: 14px;
  padding: 12px 14px;
  font-weight: 600;
}

.secondary {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.85);
}

.primary {
  background: linear-gradient(135deg, #9be4ff, #58bee8);
  color: #082133;
}

button:disabled {
  opacity: 0.64;
}
</style>
