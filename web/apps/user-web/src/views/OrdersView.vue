<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { clearAccessToken, getAccessToken } from "../api/client";
import { getMyOrders, type Order } from "../api/user";

const router = useRouter();
const loading = ref(false);
const errorText = ref("");
const orders = ref<Order[]>([]);

const orderStatusMap: Record<number, string> = {
  10: "待支付",
  20: "待接单",
  30: "待出发",
  40: "服务中",
  50: "待确认完成",
  60: "已完成",
  70: "已取消"
};

function statusText(status: number) {
  return orderStatusMap[status] ?? `状态${status}`;
}

function goHome() {
  void router.push("/");
}

function goLogin() {
  void router.push("/login");
}

function logout() {
  clearAccessToken();
  orders.value = [];
  errorText.value = "已退出登录";
}

async function loadOrders() {
  if (!getAccessToken()) {
    errorText.value = "请先登录后查看订单";
    return;
  }
  loading.value = true;
  errorText.value = "";
  try {
    const data = await getMyOrders();
    orders.value = data.list ?? [];
  } catch (error) {
    errorText.value = "订单加载失败，请检查登录状态";
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  void loadOrders();
});
</script>

<template>
  <main class="page-shell content-shell">
    <section class="hero-card">
      <p class="eyebrow">listen user web</p>
      <h1>我的订单</h1>
      <p class="desc">展示当前登录用户的订单列表（接口：`GET /api/v1/users/me/orders`）。</p>
      <div class="actions">
        <button class="btn" @click="goHome">返回首页</button>
        <button class="btn" @click="goLogin">去登录</button>
        <button class="btn" @click="logout">退出登录</button>
        <button class="btn primary" :disabled="loading" @click="loadOrders">
          {{ loading ? "加载中..." : "刷新订单" }}
        </button>
      </div>
      <p v-if="errorText" class="error-tip">{{ errorText }}</p>
      <ul v-if="orders.length > 0" class="order-list">
        <li v-for="order in orders" :key="order.id">
          <div class="order-head">
            <strong>{{ order.orderNo }}</strong>
            <span class="pill">{{ statusText(order.status) }}</span>
          </div>
          <p>场景：{{ order.sceneText || "未填写" }}</p>
          <p>地址：{{ order.addressText || "未填写" }}</p>
          <p>服务时长：{{ order.plannedDuration }} 分钟</p>
          <p>订单金额：¥{{ (order.payAmount / 100).toFixed(2) }}</p>
        </li>
      </ul>
      <p v-else-if="!loading" class="desc">暂无订单数据。</p>
    </section>
  </main>
</template>
