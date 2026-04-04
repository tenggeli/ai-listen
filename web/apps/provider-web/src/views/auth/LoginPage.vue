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
  <main class="provider-login-page">
    <section class="provider-login-card">
      <p class="provider-login-kicker">Listen Provider Console</p>
      <h1>服务方后台登录</h1>
      <p class="provider-page-subtitle">围绕接单履约、经营成长与收益管理的服务方控制台。</p>
      <form @submit.prevent="submit" class="provider-login-form">
        <label class="provider-field">
          <span>服务方账号</span>
          <input v-model="state.account" type="text" autocomplete="username" />
        </label>
        <label class="provider-field">
          <span>登录密码</span>
          <input v-model="state.password" type="password" autocomplete="current-password" />
        </label>
        <button type="submit" :disabled="state.pageState === PageLoadState.Loading">
          {{ state.pageState === PageLoadState.Loading ? '登录中...' : '进入工作台' }}
        </button>
      </form>
      <p v-if="state.pageState === PageLoadState.Error" class="provider-error">{{ state.errorMessage }}</p>
      <p v-else-if="state.pageState === PageLoadState.Success" class="provider-success">登录成功，正在进入工作台...</p>
      <p class="provider-sub">默认测试账号：provider / provider123</p>
    </section>
  </main>
</template>
