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
  successMessage: '',
  saving: false,
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
  state.successMessage = ''
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

async function saveProfile(): Promise<void> {
  if (state.saving) {
    return
  }
  state.errorMessage = ''
  state.successMessage = ''
  state.saving = true
  try {
    const saved = await authApi.saveProfile(authService.getAccessToken(), {
      displayName: state.displayName,
      cityCode: state.cityCode
    })
    state.displayName = saved.displayName
    state.cityCode = saved.cityCode
    state.successMessage = '资料已保存'
  } catch (error) {
    if (isUnauthorizedError(error)) {
      authService.logout()
      await router.replace('/login')
      return
    }
    state.errorMessage = error instanceof Error ? error.message : '保存失败'
  } finally {
    state.saving = false
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
      <div class="nav-links">
        <RouterLink to="/services">服务项目</RouterLink>
        <button @click="logout">退出登录</button>
      </div>
    </nav>
    <h1>资料编辑</h1>
    <p v-if="state.pageState === PageLoadState.Idle">准备加载资料...</p>
    <p v-else-if="state.pageState === PageLoadState.Loading">加载中...</p>
    <section v-else-if="state.pageState === PageLoadState.Success" class="card">
      <p><strong>Provider ID：</strong>{{ state.providerId }}</p>
      <p><strong>账号：</strong>{{ state.account }}</p>
      <p><strong>状态：</strong>{{ state.status }}</p>

      <label class="field">
        <span>昵称</span>
        <input v-model="state.displayName" type="text" placeholder="请输入昵称" maxlength="64" />
      </label>
      <label class="field">
        <span>城市编码</span>
        <input v-model="state.cityCode" type="text" placeholder="如 310100" maxlength="16" />
      </label>

      <button :disabled="state.saving" @click="saveProfile">
        {{ state.saving ? '保存中...' : '保存资料' }}
      </button>
    </section>
    <p v-if="state.successMessage" class="success">{{ state.successMessage }}</p>
    <p v-if="state.pageState === PageLoadState.Error || state.errorMessage" class="error">{{ state.errorMessage }}</p>
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
.nav-links {
  display: flex;
  gap: 12px;
  align-items: center;
}
.card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  padding: 14px;
}
.field {
  margin-top: 12px;
  display: grid;
  gap: 6px;
}
.field input {
  height: 40px;
  padding: 0 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
}
button {
  margin-top: 14px;
  height: 40px;
  padding: 0 16px;
  border: 1px solid #0f172a;
  background: #0f172a;
  color: #fff;
  border-radius: 8px;
  cursor: pointer;
}
button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.error {
  color: #b91c1c;
}
.success {
  color: #166534;
}
</style>
