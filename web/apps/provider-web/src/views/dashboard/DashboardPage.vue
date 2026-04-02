<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { HttpProviderOrderApi } from '../../api/ProviderOrderApi'
import { authService } from '../../application/auth'

const router = useRouter()
const orderApi = new HttpProviderOrderApi('/api/v1/provider')
const state = reactive({
  pendingCount: 0,
  loading: false,
  errorMessage: ''
})

onMounted(() => {
  void loadSummary()
})

async function loadSummary(): Promise<void> {
  state.loading = true
  state.errorMessage = ''
  try {
    const result = await orderApi.listOrders(authService.getAccessToken(), 1, 50)
    state.pendingCount = result.items.filter((item) => item.status === 'paid').length
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '加载失败'
  } finally {
    state.loading = false
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
        <strong>{{ state.loading ? '...' : `${state.pendingCount} 单` }}</strong>
      </article>
      <article class="card">
        <small>入口</small>
        <p><RouterLink to="/orders">进入订单管理</RouterLink></p>
        <p><RouterLink to="/profile">进入个人资料</RouterLink></p>
      </article>
    </section>
    <p v-if="state.errorMessage" class="error">{{ state.errorMessage }}</p>
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
@media (max-width: 960px) {
  .cards {
    grid-template-columns: 1fr;
  }
}
</style>
