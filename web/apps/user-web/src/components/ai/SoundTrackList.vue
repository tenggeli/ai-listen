<script setup lang="ts">
import type { SoundTrack } from '../../domain/ai/SoundTrack'

defineProps<{
  tracks: SoundTrack[]
  currentTrackId: string
  isPlaying: boolean
}>()

defineEmits<{
  select: [trackId: string]
}>()
</script>

<template>
  <div class="list">
    <button
      v-for="track in tracks"
      :key="track.id"
      type="button"
      class="card"
      :class="{ active: track.id === currentTrackId }"
      @click="$emit('select', track.id)"
    >
      <span class="thumb">{{ track.emoji }}</span>
      <span class="info">
        <span class="title">{{ track.title }}</span>
        <span class="meta">{{ track.category }} · {{ track.playCountText }}</span>
      </span>
      <span class="duration">{{ track.durationText }}</span>
      <span class="play-icon" aria-hidden="true">
        <svg v-if="track.id === currentTrackId && isPlaying" width="16" height="16" viewBox="0 0 24 24" fill="none">
          <rect x="6" y="4" width="4" height="16" />
          <rect x="14" y="4" width="4" height="16" />
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none">
          <polygon points="5 3 19 12 5 21 5 3" />
        </svg>
      </span>
    </button>
  </div>
</template>

<style scoped>
.list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.04);
  color: inherit;
  cursor: pointer;
  text-align: left;
}

.card.active {
  border-color: rgba(74, 168, 196, 0.35);
  background: rgba(74, 168, 196, 0.08);
}

.thumb {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 46px;
  height: 46px;
  border: 1px solid rgba(74, 168, 196, 0.2);
  border-radius: 12px;
  background: rgba(74, 168, 196, 0.1);
  font-size: 20px;
  flex-shrink: 0;
}

.info {
  display: flex;
  flex: 1;
  min-width: 0;
  flex-direction: column;
}

.title {
  overflow: hidden;
  margin-bottom: 3px;
  color: rgba(255, 255, 255, 0.82);
  font-size: 14px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.meta,
.duration {
  color: rgba(255, 255, 255, 0.26);
  font-size: 11px;
}

.duration,
.play-icon {
  flex-shrink: 0;
}

.play-icon {
  color: rgba(74, 168, 196, 0.54);
}

.play-icon svg {
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}
</style>
