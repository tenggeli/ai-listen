<script setup lang="ts">
import type { AiHomeQuickAction } from '../../domain/ai/AiHomeQuickAction'

defineProps<{
  actions: AiHomeQuickAction[]
}>()

defineEmits<{
  select: [action: AiHomeQuickAction]
}>()

function iconMarkup(icon: string): string {
  switch (icon) {
    case 'join':
      return '<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="1.8"/><line x1="19" y1="8" x2="19" y2="14" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><line x1="22" y1="11" x2="16" y2="11" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>'
    case 'square':
      return '<circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="1.8"/><path d="M12 2v3M12 19v3M2 12h3M19 12h3M4.93 4.93l2.12 2.12M16.95 16.95l2.12 2.12M4.93 19.07l2.12-2.12M16.95 7.05l2.12-2.12" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"/>'
    case 'voice':
      return '<polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><path d="M15.54 8.46a5 5 0 0 1 0 7.07" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>'
    default:
      return '<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>'
  }
}
</script>

<template>
  <div class="quick-actions">
    <button v-for="action in actions" :key="action.key" type="button" class="quick-btn" @click="$emit('select', action)">
      <span class="quick-icon" aria-hidden="true">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" v-html="iconMarkup(action.icon)" />
      </span>
      <span class="quick-label">{{ action.label }}</span>
    </button>
  </div>
</template>

<style scoped>
.quick-actions {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.quick-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 12px 8px;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
  color: rgba(126, 200, 220, 0.88);
  cursor: pointer;
  transition: transform 0.2s ease, border-color 0.2s ease, background 0.2s ease;
}

.quick-btn:hover {
  transform: translateY(-1px);
  background: rgba(74, 168, 196, 0.08);
  border-color: rgba(74, 168, 196, 0.22);
}

.quick-label {
  font-size: 11px;
  letter-spacing: 0.08em;
  color: rgba(255, 255, 255, 0.48);
}

@media (max-width: 480px) {
  .quick-actions {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
