<script setup lang="ts">
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import { loginBySMS, sendSMS } from "../api/user";
import { setAccessToken } from "../api/client";

const router = useRouter();
const mobile = ref("13800138000");
const code = ref("123456");
const resultText = ref("未登录");
const debugCode = ref("");
const loading = ref(false);

const canSubmit = computed(() => mobile.value.trim() !== "" && code.value.trim() !== "");

function goHome() {
  void router.push("/");
}

async function handleSendSMS() {
  try {
    const data = await sendSMS(mobile.value.trim());
    debugCode.value = data.debugCode ?? "";
  } catch (error) {
    debugCode.value = "";
    resultText.value = "验证码发送失败";
  }
}

async function handleLogin() {
  if (!canSubmit.value || loading.value) {
    return;
  }
  loading.value = true;
  try {
    const data = await loginBySMS(mobile.value.trim(), code.value.trim());
    setAccessToken(data.accessToken);
    resultText.value = `登录成功：${data.user.nickname}`;
    void router.push("/profile");
  } catch (error) {
    resultText.value = "登录失败，请检查手机号与验证码";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <main class="page-shell content-shell">
    <section class="hero-card">
      <p class="eyebrow">listen user web</p>
      <h1>短信登录</h1>
      <p class="desc">开发环境可先点击发送验证码，再使用返回的 debugCode 登录。</p>
      <div class="form-grid">
        <label class="field">
          <span>手机号</span>
          <input v-model="mobile" placeholder="请输入手机号" />
        </label>
        <label class="field">
          <span>验证码</span>
          <input v-model="code" placeholder="请输入验证码" />
        </label>
      </div>
      <div class="actions">
        <button class="btn" @click="goHome">返回首页</button>
        <button class="btn" @click="handleSendSMS">发送验证码</button>
        <button class="btn primary" :disabled="!canSubmit || loading" @click="handleLogin">
          {{ loading ? "登录中..." : "登录" }}
        </button>
      </div>
      <p class="status-line"><span>登录状态</span><strong>{{ resultText }}</strong></p>
      <p v-if="debugCode" class="status-line"><span>debugCode</span><strong>{{ debugCode }}</strong></p>
    </section>
  </main>
</template>
