<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { HttpAuthApi } from '../../api/AuthApi'
import { ProfileSetupViewModel } from '../../application/identity/ProfileSetupViewModel'

const router = useRouter()
const vm = new ProfileSetupViewModel(new HttpAuthApi(import.meta.env.VITE_AI_API_BASE_URL ?? '/api/v1'))
const ageOptions = ['18-24', '25-29', '30-34', '35-39', '40+']
const genderOptions = [
  { label: '女', value: 'female' },
  { label: '男', value: 'male' },
  { label: '暂不透露', value: 'unknown' }
]

onMounted(() => {
  void vm.initialize().catch((error: unknown) => {
    vm.state.errorMessage = error instanceof Error ? error.message : '加载资料失败'
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
  <main class="profile-page">
    <section class="card">
      <p class="step">步骤 1 / 2</p>
      <h1>先认识你一下</h1>
      <p class="desc">为了让首页推荐和对话语气更贴近你，先完成基础资料，后续可在“我的”继续编辑。</p>

      <div v-if="vm.state.loading" class="status">正在加载你的资料...</div>

      <template v-else>
        <label class="label" for="nickname">昵称</label>
        <input id="nickname" v-model="vm.state.nickname" class="input" placeholder="输入昵称" />

        <label class="label" for="avatar">头像 URL（可选）</label>
        <input id="avatar" v-model="vm.state.avatarUrl" class="input" placeholder="https://..." />

        <p class="label">性别</p>
        <div class="gender-row">
          <button
            v-for="item in genderOptions"
            :key="item.value"
            type="button"
            class="chip"
            :class="{ active: vm.state.gender === item.value }"
            @click="vm.state.gender = item.value"
          >
            {{ item.label }}
          </button>
        </div>

        <label class="label" for="age">年龄段</label>
        <select id="age" v-model="vm.state.ageRange" class="input">
          <option value="">请选择年龄段</option>
          <option v-for="item in ageOptions" :key="item" :value="item">{{ item }}</option>
        </select>

        <label class="label" for="city">所在城市</label>
        <input id="city" v-model="vm.state.city" class="input" placeholder="例如：上海" />

        <label class="label" for="bio">个人简介（可选）</label>
        <textarea id="bio" v-model="vm.state.bio" class="input textarea" rows="3" placeholder="简单介绍一下你自己" />

        <p v-if="vm.state.errorMessage" class="error">{{ vm.state.errorMessage }}</p>

        <div class="action-row">
          <button type="button" class="ghost" :disabled="vm.state.saving" @click="onSkip">稍后再填</button>
          <button type="button" class="primary" :disabled="vm.state.saving" @click="onSubmit">继续设置性格</button>
        </div>
      </template>
    </section>
  </main>
</template>

<style scoped>
.profile-page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  padding: 20px;
  background:
    radial-gradient(circle at top left, rgba(31, 105, 136, 0.32), transparent 42%),
    radial-gradient(circle at bottom right, rgba(29, 120, 149, 0.3), transparent 44%),
    linear-gradient(180deg, #091621 0%, #08131f 100%);
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
  font-size: 28px;
  font-weight: 500;
  line-height: 1.2;
}

.desc {
  margin: 0;
  color: rgba(255, 255, 255, 0.66);
  line-height: 1.7;
  font-size: 14px;
}

.status {
  margin-top: 16px;
  color: rgba(255, 255, 255, 0.75);
}

.label {
  display: block;
  margin-top: 12px;
  margin-bottom: 6px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.66);
}

.input {
  width: 100%;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.04);
  color: #eff7fb;
}

.textarea {
  resize: vertical;
  min-height: 84px;
}

.gender-row {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.chip {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.03);
  color: rgba(255, 255, 255, 0.68);
  padding: 11px 14px;
}

.chip.active {
  border-color: rgba(115, 213, 255, 0.4);
  background: rgba(115, 213, 255, 0.18);
  color: #9be4ff;
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
  padding: 15px 14px;
  font-weight: 600;
}

.ghost {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.74);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.primary {
  background: linear-gradient(135deg, #9be4ff, #58bee8);
  color: #082133;
  box-shadow: 0 12px 28px rgba(86, 190, 232, 0.24);
}

button:disabled {
  opacity: 0.64;
}
</style>
