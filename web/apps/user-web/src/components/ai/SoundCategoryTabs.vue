<script setup lang="ts">
import type { SoundCategory } from '../../domain/ai/SoundCategory'

defineProps<{
  categories: SoundCategory[]
  selectedKey: string
}>()

defineEmits<{
  select: [key: string]
}>()
</script>

<template>
  <div class="tabs">
    <button
      v-for="category in categories"
      :key="category.key"
      type="button"
      class="tab"
      :class="{ active: category.key === selectedKey }"
      @click="$emit('select', category.key)"
    >
      {{ category.label }}
    </button>
  </div>
</template>

<style scoped>
.tabs {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  scrollbar-width: none;
}

.tabs::-webkit-scrollbar {
  display: none;
}

.tab {
  flex-shrink: 0;
  padding: 9px 14px;
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.04);
  color: rgba(237, 247, 251, 0.54);
  font-size: 12px;
  cursor: pointer;
}

.tab.active {
  border-color: rgba(115, 213, 255, 0.32);
  background: rgba(115, 213, 255, 0.12);
  color: #93e0ff;
}
</style>
