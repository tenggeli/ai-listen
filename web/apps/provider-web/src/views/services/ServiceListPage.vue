<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { isUnauthorizedError } from '../../api/ApiError'
import { HttpProviderServiceApi } from '../../api/ProviderServiceApi'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'
import type { ProviderService } from '../../domain/service/ProviderService'

const router = useRouter()
const api = new HttpProviderServiceApi('/api/v1/provider')
const state = reactive({
  pageState: PageLoadState.Idle,
  errorMessage: '',
  items: [] as ProviderService[]
})

onMounted(() => {
  void loadServices()
})

async function loadServices(): Promise<void> {
  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  try {
    state.items = await api.listServices(authService.getAccessToken())
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
      <div class="links">
        <RouterLink to="/profile">资料编辑</RouterLink>
        <button @click="logout">退出登录</button>
      </div>
    </nav>
    <h1>服务项目</h1>

    <p v-if="state.pageState === PageLoadState.Idle">准备加载服务项目...</p>
    <p v-else-if="state.pageState === PageLoadState.Loading">加载中...</p>

    <section v-else-if="state.pageState === PageLoadState.Success" class="card">
      <p v-if="state.items.length === 0">暂无服务项目</p>
      <table v-else class="table">
        <thead>
          <tr>
            <th>项目</th>
            <th>价格</th>
            <th>分类</th>
            <th>线上支持</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in state.items" :key="item.itemId">
            <td>
              <strong>{{ item.title }}</strong>
              <p>{{ item.description || '暂无描述' }}</p>
            </td>
            <td>¥{{ item.priceAmount }} / {{ item.priceUnit || '次' }}</td>
            <td>{{ item.categoryId }}</td>
            <td>{{ item.supportOnline ? '支持' : '仅线下' }}</td>
          </tr>
        </tbody>
      </table>
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
.links {
  display: flex;
  gap: 12px;
  align-items: center;
}
.card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  padding: 14px;
  overflow-x: auto;
}
.table {
  width: 100%;
  border-collapse: collapse;
}
.table th,
.table td {
  text-align: left;
  padding: 10px 8px;
  border-bottom: 1px solid #e2e8f0;
  vertical-align: top;
}
.table td p {
  margin: 6px 0 0;
  color: #64748b;
}
.error {
  color: #b91c1c;
}
</style>
