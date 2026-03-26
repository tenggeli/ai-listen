<template>
  <main class="page-shell content-shell">
    <section class="hero-card">
      <p class="eyebrow">listen user web</p>
      <h1>用户主页</h1>
      <p class="desc">已接入用户登录、我的订单、个人中心页面，支持基础接口联调。</p>
      <div class="status-line">
        <span>后端健康状态</span>
        <strong>{{ healthText }}</strong>
      </div>
      <div class="actions">
        <button class="btn primary" @click="checkHealth">刷新健康检查</button>
        <button class="btn" @click="goLogin">短信登录</button>
        <button class="btn" @click="goOrders">我的订单</button>
        <button class="btn" @click="goProfile">个人中心</button>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { getHealth } from "../api/user";

const router = useRouter();
const healthText = ref("未检查");

function goLogin() {
  void router.push("/login");
}

function goOrders() {
  void router.push("/orders");
}

function goProfile() {
  void router.push("/profile");
}

async function checkHealth() {
  try {
    const data = await getHealth();
    healthText.value = data.status;
  } catch (error) {
    healthText.value = "请求失败";
  }
}
</script>
