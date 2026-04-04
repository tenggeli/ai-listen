<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { isUnauthorizedError } from '../../api/ApiError'
import { HttpProviderAuthApi } from '../../api/ProviderAuthApi'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'
import ProviderShell from '../../components/layout/ProviderShell.vue'

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
  <ProviderShell title="资料管理" subtitle="维护服务方资料与基础身份信息，保障展示信息一致性。" @logout="logout">
    <p v-if="state.pageState === PageLoadState.Idle" class="provider-sub">准备加载资料...</p>
    <p v-else-if="state.pageState === PageLoadState.Loading" class="provider-sub">加载中...</p>
    <section v-else-if="state.pageState === PageLoadState.Success" class="provider-card">
      <p><strong>Provider ID：</strong>{{ state.providerId }}</p>
      <p><strong>账号：</strong>{{ state.account }}</p>
      <p><strong>状态：</strong>{{ state.status }}</p>

      <div class="provider-form">
        <label class="provider-field">
          <span>昵称</span>
          <input v-model="state.displayName" type="text" placeholder="请输入昵称" maxlength="64" />
        </label>
        <label class="provider-field">
          <span>城市编码</span>
          <input v-model="state.cityCode" type="text" placeholder="如 310100" maxlength="16" />
        </label>
      </div>

      <button :disabled="state.saving" @click="saveProfile">
        {{ state.saving ? '保存中...' : '保存资料' }}
      </button>
    </section>
    <p v-if="state.successMessage" class="provider-success">{{ state.successMessage }}</p>
    <p v-if="state.pageState === PageLoadState.Error || state.errorMessage" class="provider-error">{{ state.errorMessage }}</p>
  </ProviderShell>
</template>
