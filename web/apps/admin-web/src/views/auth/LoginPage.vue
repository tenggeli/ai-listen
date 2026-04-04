<script setup lang="ts">
import { reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authService } from '../../application/auth'

interface LoginState {
  account: string
  password: string
  loading: boolean
  errorMessage: string
}

const route = useRoute()
const router = useRouter()

const state = reactive<LoginState>({
  account: 'admin',
  password: 'admin123',
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
      <h1>平台管理后台</h1>
      <p class="sub-title">请输入管理员账号密码登录。</p>

      <form @submit.prevent="submit" class="login-form">
        <label>
          <span>账号</span>
          <input v-model="state.account" type="text" autocomplete="username" />
        </label>
        <label>
          <span>密码</span>
          <input v-model="state.password" type="password" autocomplete="current-password" />
        </label>
        <button :disabled="state.loading" type="submit">
          {{ state.loading ? '登录中...' : '登录' }}
        </button>
      </form>

      <p v-if="state.errorMessage" class="error">{{ state.errorMessage }}</p>
      <p class="hint">默认测试账号：admin / admin123</p>
    </section>
  </main>
</template>
