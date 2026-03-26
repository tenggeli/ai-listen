<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAdminAuthStore } from "./stores/auth";

const route = useRoute();
const router = useRouter();
const authStore = useAdminAuthStore();

const showLayout = computed(() => route.path !== "/login");

const navItems = [
  {
    key: "dashboard",
    label: "治理看板",
    path: "/dashboard"
  },
  {
    key: "content",
    label: "内容治理",
    path: "/governance/content"
  },
  {
    key: "complaints",
    label: "投诉处理",
    path: "/governance/complaints"
  }
];

function goPath(path: string) {
  if (route.path !== path) {
    void router.push(path);
  }
}

async function handleLogout() {
  await authStore.logout();
  void router.replace("/login");
}
</script>

<template>
  <div v-if="showLayout" class="admin-layout">
    <aside class="side-nav">
      <p class="brand-eyebrow">listen admin web</p>
      <h1>运营治理台</h1>
      <nav class="menu-list">
        <button
          v-for="item in navItems"
          :key="item.key"
          class="menu-item"
          :class="{ active: route.path === item.path }"
          @click="goPath(item.path)"
        >
          {{ item.label }}
        </button>
      </nav>
    </aside>
    <div class="main-shell">
      <header class="top-bar">
        <p class="operator">{{ authStore.adminUser?.nickname || authStore.adminUser?.username || "管理员" }}</p>
        <button class="btn ghost" @click="handleLogout">退出登录</button>
      </header>
      <router-view />
    </div>
  </div>
  <router-view v-else />
</template>
