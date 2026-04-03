<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { isUnauthorizedError } from '../../api/ApiError'
import { HttpProviderAuthApi } from '../../api/ProviderAuthApi'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'

const router = useRouter()
const authApi = new HttpProviderAuthApi('/api/v1/provider')
const state = reactive({
  pageState: PageLoadState.Idle,
  errorMessage: '',
  providerId: '',
  account: '',
  displayName: '',
  cityCode: '',
  status: ''
})

onMounted(() => {
  void loadProfile()
})

async function loadProfile(): Promise<void> {
  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  try {
    const me = await authApi.getMe(authService.getAccessToken())
    state.providerId = me.providerId
    state.account = me.account
    state.displayName = me.displayName
    state.cityCode = me.cityCode
    state.status = me.status
    state.pageState = PageLoadState.Success
  } catch (error) {
    if (isUnauthorizedError(error)) {
      authService.logout()
      await router.replace('/login')
      return
    }
    state.pageState = PageLoadState.Error
    state.errorMessage = error instanceof Error ? error.message : '加载失败'
  }
}

function logout(): void {
  authService.logout()
  void router.replace('/login')
}
</script>

<template>
  <main class="page">
    <nav class="top-nav">
      <RouterLink to="/dashboard">返回工作台</RouterLink>
      <button @click="logout">退出登录</button>
    </nav>
    <h1>个人资料</h1>
    <p v-if="state.pageState === PageLoadState.Idle">准备加载资料...</p>
    <p v-else-if="state.pageState === PageLoadState.Loading">加载中...</p>
    <section v-else-if="state.pageState === PageLoadState.Success" class="card">
      <p><strong>Provider ID：</strong>{{ state.providerId }}</p>
      <p><strong>账号：</strong>{{ state.account }}</p>
      <p><strong>昵称：</strong>{{ state.displayName }}</p>
      <p><strong>城市：</strong>{{ state.cityCode }}</p>
      <p><strong>状态：</strong>{{ state.status }}</p>
    </section>
    <p v-if="state.pageState === PageLoadState.Error" class="error">{{ state.errorMessage }}</p>
  </main>
</template>

<style scoped>
.page {
  padding: 20px;
}
.top-nav {
  margin-bottom: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  padding: 14px;
}
.error {
  color: #b91c1c;
}
</style>
