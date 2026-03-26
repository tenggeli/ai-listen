<script setup lang="ts">
import { computed, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAdminAuthStore } from "../stores/auth";

const route = useRoute();
const router = useRouter();
const authStore = useAdminAuthStore();

const username = ref("admin");
const password = ref("admin123456");
const statusText = ref("请输入管理员账号登录");

const canSubmit = computed(() => username.value.trim() !== "" && password.value.trim() !== "");

async function handleLogin() {
  if (!canSubmit.value || authStore.loading) {
    return;
  }

  try {
    const data = await authStore.login(username.value.trim(), password.value.trim());
    statusText.value = `登录成功：${data.adminUser.nickname || data.adminUser.username}`;

    const redirectPath = typeof route.query.redirect === "string" ? route.query.redirect : "/dashboard";
    void router.replace(redirectPath);
  } catch (error) {
    statusText.value = "登录失败，请检查账号密码或权限";
  }
}
</script>

<template>
  <main class="login-shell">
    <section class="login-card">
      <p class="brand-eyebrow">listen admin web</p>
      <h1>后台治理登录</h1>
      <p class="desc">已接入 `/api/v1/admin/auth/login` 与登录态持久化，支持受保护路由自动拦截。</p>
      <label class="field">
        <span>账号</span>
        <input v-model="username" autocomplete="username" placeholder="请输入管理员账号" />
      </label>
      <label class="field">
        <span>密码</span>
        <input v-model="password" type="password" autocomplete="current-password" placeholder="请输入密码" />
      </label>
      <button class="btn primary block" :disabled="!canSubmit || authStore.loading" @click="handleLogin">
        {{ authStore.loading ? "登录中..." : "登录" }}
      </button>
      <p class="status-line">
        <span>状态</span>
        <strong>{{ statusText }}</strong>
      </p>
    </section>
  </main>
</template>
