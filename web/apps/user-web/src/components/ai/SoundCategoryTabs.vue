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
  gap: 6px;
  overflow-x: auto;
  scrollbar-width: none;
}

.tabs::-webkit-scrollbar {
  display: none;
}

.tab {
  flex-shrink: 0;
  padding: 7px 16px;
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.32);
  cursor: pointer;
}

.tab.active {
  border-color: rgba(74, 168, 196, 0.38);
  background: rgba(74, 168, 196, 0.14);
  color: #4aa8c4;
}
</style>
