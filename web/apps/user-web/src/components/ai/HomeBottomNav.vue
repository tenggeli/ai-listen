<script setup lang="ts">
export interface HomeNavItem {
  key: string
  label: string
  icon: string
  active?: boolean
}

defineProps<{
  items: HomeNavItem[]
}>()

defineEmits<{
  select: [key: string]
}>()

function iconPath(icon: string): string {
  switch (icon) {
    case 'home':
      return '<path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" /><polyline points="9 22 9 12 15 12 15 22" />'
    case 'square':
      return '<circle cx="12" cy="8" r="3" /><path d="M5 21v-2a4 4 0 0 1 4-4h6a4 4 0 0 1 4 4v2" />'
    case 'join':
      return '<path d="M12 5v14" /><path d="M5 12h14" />'
    case 'voice':
      return '<polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" /><path d="M15.54 8.46a5 5 0 0 1 0 7.07" />'
    default:
      return '<circle cx="12" cy="8" r="4" /><path d="M4 22c1.8-3.5 5-5 8-5s6.2 1.5 8 5" />'
  }
}
</script>

<template>
  <nav class="bottom-nav" aria-label="主导航">
    <button
      v-for="item in items"
      :key="item.key"
      type="button"
      class="nav-item"
      :class="{ active: item.active, center: item.key === 'join' }"
      @click="$emit('select', item.key)"
    >
      <svg class="nav-icon" width="22" height="22" viewBox="0 0 24 24" fill="none" v-html="iconPath(item.icon)" />
      <span class="nav-label">{{ item.label }}</span>
    </button>
  </nav>
</template>

<style scoped>
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 50%;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: space-around;
  width: min(100%, 430px);
  gap: 8px;
  padding: 8px 12px calc(16px + env(safe-area-inset-bottom));
  border-top: 1px solid rgba(74, 168, 196, 0.12);
  background: rgba(10, 22, 34, 0.96);
  backdrop-filter: blur(24px);
  transform: translateX(-50%);
}

.nav-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  min-width: 56px;
  padding: 8px 10px;
  border: none;
  border-radius: 14px;
  background: transparent;
  color: rgba(255, 255, 255, 0.24);
  cursor: pointer;
}

.nav-item.center {
  width: 56px;
  height: 56px;
  margin-top: -24px;
  justify-content: center;
  border-radius: 50%;
  background: linear-gradient(145deg, rgba(74, 168, 196, 0.22), rgba(30, 70, 100, 0.6));
  border: 1.5px solid rgba(74, 168, 196, 0.38);
  box-shadow: 0 0 24px rgba(74, 168, 196, 0.2);
  color: #4aa8c4;
}

.nav-item.active {
  color: #4aa8c4;
}

.nav-icon {
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.nav-label {
  font-size: 10px;
  letter-spacing: 0.05em;
}

.nav-item.center .nav-label {
  position: absolute;
  bottom: -18px;
  color: rgba(74, 168, 196, 0.74);
}
</style>
