<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { clearAccessToken, getAccessToken } from "../api/client";
import { getMe, type User } from "../api/user";

const router = useRouter();
const loading = ref(false);
const errorText = ref("");
const profile = ref<User | null>(null);

function goHome() {
  void router.push("/");
}

function goLogin() {
  void router.push("/login");
}

function goOrders() {
  void router.push("/orders");
}

function logout() {
  clearAccessToken();
  profile.value = null;
  errorText.value = "已退出登录";
}

async function loadProfile() {
  if (!getAccessToken()) {
    errorText.value = "请先登录后查看个人信息";
    return;
  }
  loading.value = true;
  errorText.value = "";
  try {
    const data = await getMe();
    profile.value = data.user;
  } catch (error) {
    errorText.value = "个人信息加载失败，请重新登录";
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  void loadProfile();
});
</script>

<template>
  <main class="page-shell content-shell">
    <section class="hero-card">
      <p class="eyebrow">listen user web</p>
      <h1>个人中心</h1>
      <p class="desc">展示用户基础信息（接口：`GET /api/v1/users/me`）。</p>
      <div class="actions">
        <button class="btn" @click="goHome">返回首页</button>
        <button class="btn" @click="goLogin">去登录</button>
        <button class="btn" @click="goOrders">我的订单</button>
        <button class="btn" @click="logout">退出登录</button>
        <button class="btn primary" :disabled="loading" @click="loadProfile">
          {{ loading ? "加载中..." : "刷新资料" }}
        </button>
      </div>
      <p v-if="errorText" class="error-tip">{{ errorText }}</p>
      <div v-if="profile" class="profile-grid">
        <p><span>ID</span><strong>{{ profile.id }}</strong></p>
        <p><span>昵称</span><strong>{{ profile.nickname || "未设置" }}</strong></p>
        <p><span>手机号</span><strong>{{ profile.mobile }}</strong></p>
        <p><span>城市编码</span><strong>{{ profile.cityCode || "未设置" }}</strong></p>
        <p><span>生日</span><strong>{{ profile.birthday || "未设置" }}</strong></p>
      </div>
      <p v-else-if="!loading" class="desc">暂无用户资料。</p>
    </section>
  </main>
</template>
