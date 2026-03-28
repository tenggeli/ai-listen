<script setup lang="ts">
import { useRouter } from 'vue-router'
import { HttpAuthApi } from '../../api/AuthApi'
import { LoginPageViewModel } from '../../application/identity/LoginPageViewModel'

const router = useRouter()
const vm = new LoginPageViewModel(new HttpAuthApi(import.meta.env.VITE_AI_API_BASE_URL ?? '/api/v1'))

async function onSmsLogin() {
  try {
    const route = await vm.submitSmsLogin()
    await router.push(route)
  } catch (error) {
    vm.state.errorMessage = error instanceof Error ? error.message : '登录失败，请稍后重试'
  }
}

async function onWechatLogin() {
  try {
    const route = await vm.submitWechatLogin()
    await router.push(route)
  } catch (error) {
    vm.state.errorMessage = error instanceof Error ? error.message : '登录失败，请稍后重试'
  }
}
</script>

<template>
  <main class="auth-page">
    <section class="shell">
      <header class="hero">
        <p class="brand">LISTEN</p>
        <h1>今晚，不必一个人把情绪扛过去</h1>
        <p class="sub">用最低门槛进入，先让 AI 听见你，再逐步补齐资料与性格偏好。</p>
        <div class="feature-row">
          <span class="feature">微信一键进入</span>
          <span class="feature">手机号验证码</span>
          <span class="feature">轻资料补全</span>
        </div>
      </header>

      <section class="card">
        <div class="tab-row">
          <button type="button" class="tab" :class="{ active: vm.state.mode === 'wechat' }" @click="vm.setMode('wechat')">
            微信登录
          </button>
          <button type="button" class="tab" :class="{ active: vm.state.mode === 'sms' }" @click="vm.setMode('sms')">
            手机号登录
          </button>
        </div>

        <button type="button" class="wechat-btn" :disabled="vm.state.submitting" @click="onWechatLogin">
          微信授权登录
        </button>

        <div class="divider">或使用手机号登录</div>

        <label class="label" for="phone">手机号</label>
        <input id="phone" v-model="vm.state.phone" class="input" placeholder="13800113800" />

        <label class="label" for="code">验证码</label>
        <input id="code" v-model="vm.state.verifyCode" class="input" placeholder="输入 123" />

        <label class="agree" @click="vm.toggleAgreement">
          <span class="checkbox" :class="{ checked: vm.state.agreementAccepted }">{{ vm.state.agreementAccepted ? '✓' : '' }}</span>
          <span>我已阅读并同意《用户协议》《隐私政策》</span>
        </label>

        <button type="button" class="submit-btn" :disabled="vm.state.submitting || !vm.state.agreementAccepted" @click="onSmsLogin">
          登录后继续
        </button>

        <p v-if="vm.state.errorMessage" class="error">{{ vm.state.errorMessage }}</p>

        <div class="help">
          <h3>登录后引导</h3>
          <p>首次进入会先补全基础资料，再进入性格设置；两步都支持跳过，先体验后完善。</p>
        </div>
      </section>
    </section>
  </main>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(43, 123, 158, 0.35), transparent 40%),
    radial-gradient(circle at bottom right, rgba(22, 64, 102, 0.45), transparent 42%),
    linear-gradient(180deg, #07131d 0%, #0b1d2b 55%, #08131f 100%);
  color: #eaf6ff;
}

.shell {
  width: min(100%, 430px);
  margin: 0 auto;
  min-height: 100vh;
  padding: 54px 20px 30px;
}

.hero {
  position: relative;
  overflow: hidden;
  padding: 26px 22px 22px;
  border-radius: 28px;
  border: 1px solid rgba(145, 220, 255, 0.18);
  background: linear-gradient(160deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.03));
}

.hero::after {
  content: '';
  position: absolute;
  top: -40px;
  right: -40px;
  width: 180px;
  height: 180px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(115, 213, 255, 0.16), transparent 68%);
}

.brand {
  margin: 0;
  font-size: 30px;
  letter-spacing: 0.2em;
  font-weight: 300;
}

h1 {
  margin: 14px 0 0;
  font-size: 29px;
  line-height: 1.24;
  font-weight: 500;
}

.sub {
  margin: 10px 0 0;
  color: rgba(234, 246, 255, 0.75);
  line-height: 1.7;
  font-size: 14px;
}

.feature-row {
  display: flex;
  gap: 10px;
  margin-top: 18px;
  flex-wrap: wrap;
}

.feature {
  padding: 8px 12px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.05);
  font-size: 12px;
  color: rgba(239, 247, 251, 0.64);
}

.card {
  margin-top: 18px;
  padding: 22px;
  border-radius: 28px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(16px);
}

.tab-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.tab {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.04);
  color: rgba(234, 246, 255, 0.68);
  padding: 11px;
}

.tab.active {
  border-color: rgba(115, 213, 255, 0.35);
  background: rgba(115, 213, 255, 0.15);
  color: #9be7ff;
}

.wechat-btn,
.submit-btn {
  width: 100%;
  margin-top: 14px;
  border: none;
  border-radius: 16px;
  padding: 14px 16px;
  font-weight: 600;
}

.wechat-btn {
  background: linear-gradient(135deg, #80eba3, #43c773);
  color: #0a2716;
}

.submit-btn {
  background: linear-gradient(135deg, #9be4ff, #58bee8);
  color: #082133;
}

.submit-btn:disabled,
.wechat-btn:disabled {
  opacity: 0.65;
}

.divider {
  margin: 14px 0;
  text-align: center;
  color: rgba(234, 246, 255, 0.35);
  font-size: 12px;
}

.label {
  display: block;
  margin-top: 10px;
  margin-bottom: 6px;
  color: rgba(234, 246, 255, 0.62);
  font-size: 12px;
}

.input {
  width: 100%;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 14px;
  padding: 12px 14px;
  background: rgba(255, 255, 255, 0.05);
  color: #eaf6ff;
}

.agree {
  margin-top: 14px;
  display: flex;
  gap: 10px;
  align-items: center;
  color: rgba(234, 246, 255, 0.68);
  font-size: 12px;
  cursor: pointer;
}

.checkbox {
  width: 18px;
  height: 18px;
  border-radius: 6px;
  border: 1px solid rgba(149, 226, 255, 0.4);
  background: rgba(255, 255, 255, 0.06);
  display: grid;
  place-items: center;
  color: #96e6ff;
}

.checkbox.checked {
  background: rgba(115, 213, 255, 0.2);
}

.error {
  margin: 10px 0 0;
  color: #ffb6b6;
  font-size: 12px;
}

.help {
  margin-top: 16px;
  padding: 16px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(9, 35, 49, 0.7);
}

.help h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
}

.help p {
  margin: 8px 0 0;
  font-size: 13px;
  line-height: 1.7;
  color: rgba(239, 247, 251, 0.6);
}
</style>
