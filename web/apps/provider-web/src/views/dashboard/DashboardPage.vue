<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { isUnauthorizedError } from '../../api/ApiError'
import { HttpProviderOrderApi } from '../../api/ProviderOrderApi'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'

const router = useRouter()
const orderApi = new HttpProviderOrderApi('/api/v1/provider')
const state = reactive({
  pendingCount: 0,
  pageState: PageLoadState.Idle,
  errorMessage: ''
})

onMounted(() => {
  void loadSummary()
})

async function loadSummary(): Promise<void> {
  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  try {
    const result = await orderApi.listOrders(authService.getAccessToken(), 1, 50)
    state.pendingCount = result.items.filter((item) => item.status === 'paid').length
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
      <span>服务方工作台</span>
      <button @click="logout">退出登录</button>
    </nav>
    <section class="cards">
      <article class="card">
        <small>待处理订单</small>
        <strong>{{ state.pageState === PageLoadState.Loading ? '...' : `${state.pendingCount} 单` }}</strong>
      </article>
      <article class="card">
        <small>入口</small>
        <p><RouterLink to="/orders">进入订单管理</RouterLink></p>
        <p><RouterLink to="/services">进入服务项目</RouterLink></p>
        <p><RouterLink to="/profile">进入个人资料</RouterLink></p>
      </article>
    </section>
    <p v-if="state.pageState === PageLoadState.Idle">等待加载工作台数据...</p>
    <p v-else-if="state.pageState === PageLoadState.Error" class="error">{{ state.errorMessage }}</p>
    <p v-else-if="state.pageState === PageLoadState.Success" class="success">工作台数据已更新。</p>
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
.cards {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}
.card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  padding: 14px;
}
.card small {
  color: #64748b;
}
.card strong {
  display: block;
  margin-top: 10px;
  font-size: 26px;
}
.error {
  color: #b91c1c;
}
.success {
  color: #166534;
}
@media (max-width: 960px) {
  .cards {
    grid-template-columns: 1fr;
  }
}
</style>
