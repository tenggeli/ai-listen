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
      <h1>给 Listen 一点性格线索</h1>
      <p class="desc">这会影响首页问候语、对话语气、服务推荐和声音内容偏好。</p>

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

        <h2 class="section-title">MBTI 倾向</h2>
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
    radial-gradient(circle at bottom right, rgba(21, 56, 93, 0.35), transparent 45%),
    linear-gradient(180deg, #08131c 0%, #08131f 100%);
}

.card {
  width: min(100%, 420px);
  margin-top: 18px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 28px;
  padding: 22px;
  background: rgba(255, 255, 255, 0.045);
  backdrop-filter: blur(16px);
  height: fit-content;
}

.step {
  margin: 0;
  color: #8fdfff;
  font-size: 12px;
  letter-spacing: 0.08em;
}

h1 {
  margin: 10px 0 8px;
  font-size: 30px;
  line-height: 1.2;
  font-weight: 500;
}

.desc {
  margin: 0;
  color: rgba(255, 255, 255, 0.66);
  line-height: 1.7;
  font-size: 14px;
}

h2 {
  margin: 14px 0 8px;
  font-size: 18px;
  font-weight: 500;
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
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.03);
  color: rgba(241, 248, 252, 0.62);
  padding: 12px 14px;
}

.tag.active,
.mbti.active {
  border-color: rgba(115, 213, 255, 0.4);
  background: rgba(115, 213, 255, 0.18);
  color: #9be4ff;
  box-shadow: 0 10px 24px rgba(89, 190, 231, 0.08);
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
  border-radius: 18px;
  padding: 16px 14px;
  font-weight: 600;
}

.secondary {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(241, 248, 252, 0.72);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.primary {
  background: linear-gradient(135deg, #9be4ff, #58bee8);
  color: #082133;
  box-shadow: 0 14px 30px rgba(88, 190, 232, 0.24);
}

button:disabled {
  opacity: 0.64;
}
</style>
