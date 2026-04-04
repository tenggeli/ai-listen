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
  { label: '工作台总览', path: '/dashboard' },
  { label: '订单列表', path: '/orders' },
  { label: '服务项目与价格', path: '/services' },
  { label: '资料管理', path: '/profile' }
]

const activePath = computed(() => route.path)
</script>

<template>
  <div class="provider-layout">
    <aside class="provider-sidebar">
      <div class="provider-brand">
        <strong>Listen 服务方后台</strong>
        <span>围绕状态管理、接单履约、经营成长与收益结算设计的服务方控制台。</span>
      </div>

      <nav class="provider-menu">
        <RouterLink
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="provider-menu-link"
          :class="{ active: activePath === item.path || (item.path === '/orders' && activePath.startsWith('/orders/')) }"
        >
          {{ item.label }}
        </RouterLink>
      </nav>
    </aside>

    <main class="provider-content">
      <section class="provider-topbar">
        <div>
          <h1 class="provider-page-title">{{ props.title }}</h1>
          <p class="provider-page-subtitle">{{ props.subtitle }}</p>
        </div>
        <div class="provider-topbar-actions">
          <slot name="topbar-right" />
          <button type="button" @click="emit('logout')">退出登录</button>
        </div>
      </section>

      <slot />
    </main>
  </div>
</template>
