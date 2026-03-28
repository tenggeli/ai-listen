<script setup lang="ts">
import type { SoundTrack } from '../../domain/ai/SoundTrack'

defineProps<{
  track: SoundTrack
  progressText: string
  totalDurationText: string
  isPlaying: boolean
}>()

defineEmits<{
  toggle: []
}>()

const waveformHeights = [8, 14, 20, 28, 32, 26, 18, 24, 30, 22, 16, 28, 32, 20, 12, 24, 30, 18, 26, 14, 20, 28, 16, 22, 30]
</script>

<template>
  <section class="now-playing">
    <p class="label"><span class="dot" />正在播放</p>
    <h2 class="title">{{ track.title }}</h2>
    <p class="author">{{ track.author }} · {{ track.category }}</p>

    <div class="waveform" aria-hidden="true">
      <span
        v-for="(height, index) in waveformHeights"
        :key="`${height}-${index}`"
        class="wave"
        :style="{ '--h': `${height}px`, '--delay': `${index * 0.05}s` }"
      />
    </div>

    <div class="progress-row">
      <span class="time-label">{{ progressText }}</span>
      <div class="progress-bar">
        <div class="progress-fill" />
        <div class="progress-dot" />
      </div>
      <span class="time-label">{{ totalDurationText }}</span>
    </div>

    <div class="controls">
      <button type="button" class="ctrl-btn" aria-label="上一首">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <polygon points="19 20 9 12 19 4 19 20" />
          <line x1="5" y1="19" x2="5" y2="5" />
        </svg>
      </button>
      <button type="button" class="play-btn" aria-label="播放控制" @click="$emit('toggle')">
        <svg v-if="isPlaying" width="22" height="22" viewBox="0 0 24 24" fill="none">
          <rect x="6" y="4" width="4" height="16" />
          <rect x="14" y="4" width="4" height="16" />
        </svg>
        <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none">
          <polygon points="5 3 19 12 5 21 5 3" />
        </svg>
      </button>
      <button type="button" class="ctrl-btn" aria-label="下一首">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <polygon points="5 4 15 12 5 20 5 4" />
          <line x1="19" y1="5" x2="19" y2="19" />
        </svg>
      </button>
    </div>
  </section>
</template>

<style scoped>
.now-playing {
  position: relative;
  overflow: hidden;
  padding: 20px;
  border: 1px solid rgba(74, 168, 196, 0.25);
  border-radius: 24px;
  background: linear-gradient(135deg, rgba(30, 55, 80, 0.8), rgba(15, 35, 55, 0.95));
}

.now-playing::before {
  content: '';
  position: absolute;
  top: -30px;
  right: -30px;
  width: 140px;
  height: 140px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(74, 168, 196, 0.1) 0%, transparent 70%);
  animation: glowPulse 4s ease-in-out infinite;
}

.label,
.author,
.time-label {
  position: relative;
  z-index: 1;
}

.label {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0 0 10px;
  font-size: 10px;
  letter-spacing: 0.2em;
  color: rgba(74, 168, 196, 0.64);
}

.dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #4aa8c4;
  animation: dotBlink 1.5s ease-in-out infinite;
}

.title {
  position: relative;
  z-index: 1;
  margin: 0 0 4px;
  font-size: 18px;
  font-weight: 400;
  color: rgba(255, 255, 255, 0.92);
}

.author {
  margin: 0 0 16px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.36);
}

.waveform {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: flex-end;
  gap: 3px;
  height: 32px;
  margin-bottom: 14px;
}

.wave {
  width: 3px;
  height: var(--h);
  border-radius: 2px;
  background: rgba(74, 168, 196, 0.5);
  animation: waveBounce 0.8s ease-in-out infinite alternate;
  animation-delay: var(--delay);
}

.progress-row {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
}

.progress-bar {
  position: relative;
  flex: 1;
  height: 3px;
  border-radius: 2px;
  background: rgba(255, 255, 255, 0.1);
}

.progress-fill {
  width: 35%;
  height: 100%;
  border-radius: 2px;
  background: linear-gradient(90deg, #4aa8c4, rgba(74, 168, 196, 0.5));
}

.progress-dot {
  position: absolute;
  top: -3.5px;
  left: 35%;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #4aa8c4;
  transform: translateX(-50%);
  box-shadow: 0 0 6px rgba(74, 168, 196, 0.6);
}

.time-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.25);
}

.controls {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
}

.ctrl-btn,
.play-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  cursor: pointer;
}

.ctrl-btn {
  background: transparent;
  color: rgba(255, 255, 255, 0.45);
}

.ctrl-btn svg,
.play-btn svg {
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.play-btn {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: rgba(74, 168, 196, 0.18);
  border: 1.5px solid rgba(74, 168, 196, 0.4);
  color: #4aa8c4;
  box-shadow: 0 0 20px rgba(74, 168, 196, 0.15);
}

@keyframes dotBlink {
  0%,
  100% {
    opacity: 0.3;
  }
  50% {
    opacity: 1;
  }
}

@keyframes glowPulse {
  0%,
  100% {
    transform: scale(0.9);
    opacity: 0.5;
  }
  50% {
    transform: scale(1.1);
    opacity: 1;
  }
}

@keyframes waveBounce {
  0% {
    transform: scaleY(0.28);
    opacity: 0.3;
  }
  100% {
    transform: scaleY(1);
    opacity: 0.84;
  }
}
</style>
