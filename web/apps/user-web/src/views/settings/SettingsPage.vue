<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { HttpAuthApi } from '../../api/AuthApi'
import { clearSession, loadSession } from '../../application/identity/AuthSession'
import { loadUserSettings, saveUserSettings } from '../../application/settings/UserSettingsStore'
import { PageLoadState } from '../../domain/ai/PageLoadState'
import type { UserMe } from '../../domain/identity/UserMe'
import type { UserSettings } from '../../domain/settings/UserSettings'

const router = useRouter()
const authApi = new HttpAuthApi(import.meta.env.VITE_AI_API_BASE_URL ?? '/api/v1')

const state = reactive<{
  pageState: PageLoadState
  me: UserMe | null
  settings: UserSettings | null
  saveMessage: string
  errorMessage: string
}>({
  pageState: PageLoadState.Idle,
  me: null,
  settings: null,
  saveMessage: '',
  errorMessage: ''
})

onMounted(() => {
  void initialize()
})

const interestTagsText = computed(() => {
  if (!state.me) {
    return '-'
  }
  return state.me.interestTags.length > 0 ? state.me.interestTags.join('、') : '-'
})

async function initialize(): Promise<void> {
  const session = loadSession()
  if (!session) {
    await router.push('/auth')
    return
  }

  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  state.saveMessage = ''
  try {
    state.me = await authApi.getMe(session.accessToken)
    state.settings = loadUserSettings(session.userId)
    state.pageState = PageLoadState.Success
  } catch (error) {
    state.pageState = PageLoadState.Error
    state.errorMessage = error instanceof Error ? error.message : '加载设置页失败'
  }
}

function saveSettings(): void {
  const session = loadSession()
  if (!session || !state.settings) {
    return
  }
  saveUserSettings(session.userId, state.settings)
  state.saveMessage = '设置已保存'
  window.setTimeout(() => {
    state.saveMessage = ''
  }, 1600)
}

async function logout(): Promise<void> {
  clearSession()
  await router.push('/auth')
}
</script>

<template>
  <main class="settings-page">
    <section class="screen-shell">
      <header class="top-row">
        <button type="button" class="back" @click="router.push('/me')">返回我的</button>
        <span>设置</span>
      </header>

      <section v-if="state.pageState === PageLoadState.Loading" class="card">正在加载设置...</section>
      <section v-else-if="state.pageState === PageLoadState.Error" class="card error">{{ state.errorMessage }}</section>

      <template v-if="state.pageState === PageLoadState.Success && state.me && state.settings">
        <section class="card">
          <h1>账号资料</h1>
          <div class="line"><span>昵称</span><strong>{{ state.me.nickname || '-' }}</strong></div>
          <div class="line"><span>城市</span><strong>{{ state.me.city || '-' }}</strong></div>
          <div class="line"><span>MBTI</span><strong>{{ state.me.mbti || '-' }}</strong></div>
          <div class="line"><span>兴趣</span><strong>{{ interestTagsText }}</strong></div>
          <div class="actions">
            <button type="button" class="entry" @click="router.push('/profile/setup')">编辑基础资料</button>
            <button type="button" class="entry" @click="router.push('/personality/setup')">编辑性格偏好</button>
          </div>
        </section>

        <section class="card">
          <h2>偏好设置</h2>
          <label class="switch-row">
            <span>优先推荐同城服务方</span>
            <input v-model="state.settings.preference.preferSameCityProviders" type="checkbox" />
          </label>
          <label class="switch-row">
            <span>声音页自动播放预览</span>
            <input v-model="state.settings.preference.autoPlaySoundPreview" type="checkbox" />
          </label>
          <label class="switch-row">
            <span>隐藏当前离线服务方</span>
            <input v-model="state.settings.preference.hideOfflineProviders" type="checkbox" />
          </label>
        </section>

        <section class="card">
          <h2>通知设置</h2>
          <label class="switch-row">
            <span>订单状态更新提醒</span>
            <input v-model="state.settings.notification.orderStatusUpdate" type="checkbox" />
          </label>
          <label class="switch-row">
            <span>投诉处理结果通知</span>
            <input v-model="state.settings.notification.complaintResultNotice" type="checkbox" />
          </label>
          <label class="switch-row">
            <span>活动与运营消息</span>
            <input v-model="state.settings.notification.marketingActivity" type="checkbox" />
          </label>
        </section>

        <section class="card">
          <h2>隐私设置</h2>
          <label class="switch-row">
            <span>公开展示我的资料</span>
            <input v-model="state.settings.privacy.profilePublicVisible" type="checkbox" />
          </label>
          <label class="switch-row">
            <span>允许个性化推荐</span>
            <input v-model="state.settings.privacy.personalizedRecommendation" type="checkbox" />
          </label>
          <label class="switch-row">
            <span>允许风控安全数据联动</span>
            <input v-model="state.settings.privacy.riskControlDataSharing" type="checkbox" />
          </label>
          <p class="tip">说明：隐私设置为当前阶段最小可用结构，后续会补充更细粒度权限。</p>
        </section>

        <section class="footer-actions">
          <button type="button" class="save" @click="saveSettings">保存设置</button>
          <button type="button" class="logout" @click="logout">退出登录</button>
          <p v-if="state.saveMessage" class="save-msg">{{ state.saveMessage }}</p>
        </section>
      </template>
    </section>
  </main>
</template>

<style scoped>
.settings-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top, rgba(74, 168, 196, 0.1), transparent 30%),
    linear-gradient(180deg, #0b1926 0%, #0d1f2d 42%, #08131f 100%);
  color: rgba(255, 255, 255, 0.88);
}

.screen-shell {
  width: min(100%, 390px);
  min-height: 100vh;
  margin: 0 auto;
  padding: 52px 16px 24px;
}

.top-row {
  display: flex;
  align-items: center;
  gap: 10px;
  color: rgba(255, 255, 255, 0.45);
  font-size: 12px;
}

.back {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.8);
  padding: 6px 10px;
}

.card {
  margin-top: 14px;
  border-radius: 18px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.045);
  padding: 16px;
}

.card h1 {
  margin: 0;
  font-size: 22px;
}

.card h2 {
  margin: 0;
  font-size: 18px;
}

.line {
  margin-top: 12px;
  display: flex;
  justify-content: space-between;
  gap: 10px;
  font-size: 14px;
}

.line span {
  color: rgba(239, 247, 251, 0.56);
}

.line strong {
  text-align: right;
}

.actions {
  margin-top: 12px;
  display: grid;
  gap: 8px;
}

.entry {
  width: 100%;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.86);
  padding: 13px 12px;
  text-align: left;
}

.switch-row {
  margin-top: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  font-size: 14px;
  color: rgba(239, 247, 251, 0.86);
}

.switch-row input[type='checkbox'] {
  width: 18px;
  height: 18px;
}

.tip {
  margin-top: 12px;
  font-size: 12px;
  color: rgba(239, 247, 251, 0.54);
}

.footer-actions {
  margin-top: 14px;
  display: grid;
  gap: 8px;
  padding-bottom: 24px;
}

.save,
.logout {
  border: none;
  border-radius: 14px;
  padding: 13px 12px;
  font-size: 14px;
  font-weight: 600;
}

.save {
  background: linear-gradient(135deg, #97e3ff, #58bee8);
  color: #082132;
}

.logout {
  background: rgba(255, 255, 255, 0.07);
  color: rgba(239, 247, 251, 0.78);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.save-msg {
  margin: 4px 2px 0;
  font-size: 12px;
  color: #93f5c7;
}

.error {
  color: #fca5a5;
}
</style>
