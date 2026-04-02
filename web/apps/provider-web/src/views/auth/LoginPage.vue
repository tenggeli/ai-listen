<script setup lang="ts">
import { reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authService } from '../../application/auth'

const route = useRoute()
const router = useRouter()
const state = reactive({
  account: 'provider',
  password: 'provider123',
  loading: false,
  errorMessage: ''
})

async function submit(): Promise<void> {
  state.loading = true
  state.errorMessage = ''
  try {
    await authService.login(state.account, state.password)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/dashboard'
    await router.replace(redirect)
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '登录失败'
  } finally {
    state.loading = false
  }
}
</script>

<template>
  <main class="login-page">
    <section class="login-card">
      <h1>服务方后台</h1>
      <p class="sub-title">请输入服务方账号密码登录。</p>
      <form @submit.prevent="submit" class="login-form">
        <label>
          <span>账号</span>
          <input v-model="state.account" type="text" autocomplete="username" />
        </label>
        <label>
          <span>密码</span>
          <input v-model="state.password" type="password" autocomplete="current-password" />
        </label>
        <button type="submit" :disabled="state.loading">{{ state.loading ? '登录中...' : '登录' }}</button>
      </form>
      <p v-if="state.errorMessage" class="error">{{ state.errorMessage }}</p>
      <p class="hint">默认测试账号：provider / provider123</p>
    </section>
  </main>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  background: linear-gradient(160deg, #f8fafc 0%, #e2e8f0 100%);
  padding: 16px;
}
.login-card {
  width: min(420px, 100%);
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 10px 25px rgba(15, 23, 42, 0.08);
}
.sub-title {
  margin-top: 8px;
  color: #475569;
}
.login-form {
  margin-top: 16px;
  display: grid;
  gap: 12px;
}
label {
  display: grid;
  gap: 6px;
}
input {
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  padding: 10px 12px;
}
button {
  margin-top: 4px;
  border: 1px solid #0f172a;
  border-radius: 10px;
  padding: 10px 12px;
  background: #0f172a;
  color: #fff;
  font-weight: 600;
}
.error {
  margin-top: 12px;
  color: #b91c1c;
}
.hint {
  margin-top: 10px;
  color: #64748b;
  font-size: 13px;
}
</style>
