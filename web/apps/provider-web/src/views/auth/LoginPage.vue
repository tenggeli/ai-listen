<script setup lang="ts">
import { reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'

const route = useRoute()
const router = useRouter()
const state = reactive({
  account: 'provider',
  password: 'provider123',
  pageState: PageLoadState.Idle,
  errorMessage: ''
})

async function submit(): Promise<void> {
  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  try {
    await authService.login(state.account, state.password)
    state.pageState = PageLoadState.Success
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/dashboard'
    await router.replace(redirect)
  } catch (error) {
    state.pageState = PageLoadState.Error
    state.errorMessage = error instanceof Error ? error.message : '登录失败'
  }
}
</script>

<template>
  <main class="login-page">
    <section class="login-card">
      <p class="kicker">Listen Provider Console</p>
      <h1>服务方后台登录</h1>
      <p class="sub-title">围绕接单履约、经营成长与收益管理的服务方控制台。</p>
      <form @submit.prevent="submit" class="login-form">
        <label>
          <span>服务方账号</span>
          <input v-model="state.account" type="text" autocomplete="username" />
        </label>
        <label>
          <span>登录密码</span>
          <input v-model="state.password" type="password" autocomplete="current-password" />
        </label>
        <button type="submit" :disabled="state.pageState === PageLoadState.Loading">
          {{ state.pageState === PageLoadState.Loading ? '登录中...' : '进入工作台' }}
        </button>
      </form>
      <p v-if="state.pageState === PageLoadState.Error" class="error">{{ state.errorMessage }}</p>
      <p v-else-if="state.pageState === PageLoadState.Success" class="success">登录成功，正在进入工作台...</p>
      <p class="hint">默认测试账号：provider / provider123</p>
    </section>
  </main>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  background:
    radial-gradient(circle at top left, rgba(31, 109, 141, 0.2), transparent 30%),
    radial-gradient(circle at bottom right, rgba(21, 59, 93, 0.24), transparent 32%),
    #06111b;
  padding: 16px;
}
.login-card {
  width: min(420px, 100%);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(148, 217, 255, 0.14);
  border-radius: 22px;
  padding: 28px 24px;
  box-shadow: 0 18px 48px rgba(3, 8, 12, 0.38);
  color: #eaf6fb;
}
.kicker {
  margin: 0;
  color: rgba(234, 246, 251, 0.56);
  font-size: 12px;
  letter-spacing: 1.2px;
  text-transform: uppercase;
}
h1 {
  margin: 12px 0 0;
  font-size: 30px;
  font-weight: 500;
}
.sub-title {
  margin-top: 8px;
  color: rgba(234, 246, 251, 0.66);
  line-height: 1.7;
}
.login-form {
  margin-top: 18px;
  display: grid;
  gap: 14px;
}
label {
  display: grid;
  gap: 6px;
  color: rgba(234, 246, 251, 0.76);
  font-size: 14px;
}
input {
  border: 1px solid rgba(148, 217, 255, 0.22);
  border-radius: 14px;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.06);
  color: #eaf6fb;
}
input::placeholder {
  color: rgba(234, 246, 251, 0.48);
}
button {
  margin-top: 4px;
  border: 1px solid rgba(115, 213, 255, 0.48);
  border-radius: 14px;
  padding: 11px 12px;
  background: rgba(115, 213, 255, 0.18);
  color: #f3fbff;
  font-weight: 600;
  cursor: pointer;
}
button:disabled {
  opacity: 0.7;
  cursor: default;
}
.error {
  margin-top: 12px;
  color: #ffd278;
}
.success {
  margin-top: 12px;
  color: #7df0bc;
}
.hint {
  margin-top: 10px;
  color: rgba(234, 246, 251, 0.54);
  font-size: 13px;
}
</style>
