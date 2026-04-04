<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const props = defineProps<{
  title: string
  subtitle: string
}>()

const emit = defineEmits<{
  logout: []
}>()

const route = useRoute()

const menuItems = [
  { label: '经营总览', path: '/dashboard' },
  { label: '服务方审核中心', path: '/providers/review' },
  { label: '服务项目管理', path: '/services/manage' },
  { label: '声音内容审核', path: '/sounds/manage' },
  { label: '订单监管', path: '/orders/manage' },
  { label: '投诉仲裁', path: '/complaints/manage' }
]

const activePath = computed(() => route.path)
</script>

<template>
  <div class="admin-layout">
    <aside class="admin-sidebar">
      <div class="admin-brand">
        <strong>Listen 平台管理后台</strong>
        <span>统一承接运营、审核、客服、财务与风控任务，聚焦风险预警与处理闭环。</span>
      </div>

      <nav class="admin-menu">
        <RouterLink
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="admin-menu-link"
          :class="{ active: activePath === item.path }"
        >
          {{ item.label }}
        </RouterLink>
      </nav>
    </aside>

    <main class="admin-content">
      <section class="admin-topbar">
        <div>
          <h1 class="admin-page-title">{{ props.title }}</h1>
          <p class="admin-page-subtitle">{{ props.subtitle }}</p>
        </div>
        <div class="admin-topbar-actions">
          <slot name="topbar-right" />
          <button type="button" @click="emit('logout')">退出登录</button>
        </div>
      </section>

      <slot />
    </main>
  </div>
</template>
