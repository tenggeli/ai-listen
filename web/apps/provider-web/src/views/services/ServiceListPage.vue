<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { isUnauthorizedError } from '../../api/ApiError'
import { HttpProviderServiceApi } from '../../api/ProviderServiceApi'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'
import type { ProviderService } from '../../domain/service/ProviderService'
import ProviderShell from '../../components/layout/ProviderShell.vue'

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
  <ProviderShell title="服务项目与价格" subtitle="展示服务项目、价格、分类与线上支持状态，支撑服务运营管理。" @logout="logout">
    <p v-if="state.pageState === PageLoadState.Idle" class="provider-sub">准备加载服务项目...</p>
    <p v-else-if="state.pageState === PageLoadState.Loading" class="provider-sub">加载中...</p>

    <section v-else-if="state.pageState === PageLoadState.Success" class="provider-card">
      <div v-if="state.items.length === 0" class="provider-empty">暂无服务项目</div>
      <div v-else class="provider-table">
        <div class="provider-row-head">
          <div>项目</div>
          <div>价格</div>
          <div>分类</div>
          <div>线上支持</div>
          <div>说明</div>
        </div>
        <div v-for="item in state.items" :key="item.itemId" class="provider-row">
          <div><strong>{{ item.title }}</strong></div>
          <div>¥{{ item.priceAmount }} / {{ item.priceUnit || '次' }}</div>
          <div>{{ item.categoryId }}</div>
          <div>{{ item.supportOnline ? '支持' : '仅线下' }}</div>
          <div>{{ item.description || '暂无描述' }}</div>
        </div>
      </div>
    </section>

    <p v-if="state.pageState === PageLoadState.Error" class="provider-error">{{ state.errorMessage }}</p>
  </ProviderShell>
</template>
