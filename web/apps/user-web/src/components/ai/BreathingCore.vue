<script setup lang="ts">
import { computed } from 'vue'
import { PageLoadState } from '../../domain/ai/PageLoadState'

const props = withDefaults(
  defineProps<{
    state?: PageLoadState
    caption?: string
  }>(),
  {
    state: PageLoadState.Idle,
    caption: 'listen'
  }
)

defineEmits<{
  click: []
}>()

const stateText = computed(() => {
  switch (props.state) {
    case PageLoadState.Loading:
      return '正在倾听'
    case PageLoadState.Success:
      return '推荐完成'
    case PageLoadState.Empty:
      return '暂未匹配到合适对象'
    case PageLoadState.Error:
      return '网络有点拥挤'
    default:
      return '轻触开始表达'
  }
})
</script>

<template>
  <button class="breath-wrapper" type="button" @click="$emit('click')">
    <div class="glow-outer" />
    <div class="ripple ripple-1" />
    <div class="ripple ripple-2" />
    <div class="ripple ripple-3" />
    <div class="mid-ring" />
    <div class="inner-ring" />
    <div class="core">
      <svg class="logo-icon" viewBox="0 0 90 90" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
        <path d="M37 37 C33 41 33 49 37 53" stroke="#7ec8dc" stroke-width="2.6" stroke-linecap="round" />
        <path d="M30 31 C24 37 24 53 30 59" stroke="#aadde9" stroke-width="2.2" stroke-linecap="round" />
        <path
          d="M64 62 C58 72 46 74 38 66 C28 56 28 40 36 31 C44 22 58 22 66 32 L64 36"
          stroke="#4aa8c4"
          stroke-width="3.2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
        <path d="M64 36 C70 40 70 52 64 56" stroke="#4aa8c4" stroke-width="3" stroke-linecap="round" />
        <path d="M60 40 C62 50 56 60 50 62 C47 63 46 66 50 68 C53 70 57 68 57 64" stroke="#4aa8c4" stroke-width="2.6" stroke-linecap="round" />
      </svg>
      <span class="caption">{{ caption }}</span>
    </div>
    <span class="state-text">{{ stateText }}</span>
  </button>
</template>

<style scoped>
.breath-wrapper {
  position: relative;
  display: grid;
  place-items: center;
  width: 280px;
  height: 320px;
  margin: 0;
  padding: 0;
  border: none;
  background: transparent;
  color: inherit;
  cursor: pointer;
}

.glow-outer {
  position: absolute;
  width: 280px;
  height: 280px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(74, 168, 196, 0.08) 0%, transparent 70%);
  animation: glowPulse 6s ease-in-out infinite;
}

.ripple {
  position: absolute;
  border-radius: 50%;
  border: 1px solid rgba(74, 168, 196, 0.12);
  animation: rippleExpand 5s ease-out infinite;
}

.ripple-1 {
  animation-delay: 0s;
}

.ripple-2 {
  animation-delay: 1.25s;
}

.ripple-3 {
  animation-delay: 2.5s;
}

.mid-ring {
  position: absolute;
  width: 174px;
  height: 174px;
  border-radius: 50%;
  border: 1px solid rgba(74, 168, 196, 0.18);
  animation: glowPulse 6s ease-in-out infinite;
}

.inner-ring {
  position: absolute;
  width: 132px;
  height: 132px;
  border-radius: 50%;
  border: 1.5px solid rgba(74, 168, 196, 0.26);
  animation: innerRingPulse 6s ease-in-out infinite;
}

.core {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 104px;
  height: 104px;
  border-radius: 50%;
  background: radial-gradient(circle at 40% 40%, rgba(74, 168, 196, 0.22) 0%, rgba(20, 50, 70, 0.86) 60%, rgba(13, 31, 45, 0.96) 100%);
  border: 1.5px solid rgba(74, 168, 196, 0.35);
  animation: breathing 6s ease-in-out infinite;
  box-shadow: 0 0 65px rgba(74, 168, 196, 0.22), inset 0 0 35px rgba(74, 168, 196, 0.08);
}

.logo-icon {
  width: 44px;
  height: 44px;
}

.caption {
  margin-top: 6px;
  font-size: 9px;
  letter-spacing: 0.3em;
  text-transform: lowercase;
  color: rgba(74, 168, 196, 0.72);
}

.state-text {
  position: absolute;
  bottom: 8px;
  font-size: 12px;
  letter-spacing: 0.08em;
  color: rgba(255, 255, 255, 0.42);
}

@keyframes breathing {
  0%,
  100% {
    transform: scale(0.93);
    box-shadow: 0 0 30px rgba(74, 168, 196, 0.12), inset 0 0 20px rgba(74, 168, 196, 0.04);
  }
  50% {
    transform: scale(1.04);
    box-shadow: 0 0 65px rgba(74, 168, 196, 0.28), inset 0 0 35px rgba(74, 168, 196, 0.08);
  }
}

@keyframes glowPulse {
  0%,
  100% {
    transform: scale(0.92);
    opacity: 0.5;
  }
  50% {
    transform: scale(1.08);
    opacity: 1;
  }
}

@keyframes innerRingPulse {
  0%,
  100% {
    transform: scale(0.94);
    border-color: rgba(74, 168, 196, 0.15);
  }
  50% {
    transform: scale(1.05);
    border-color: rgba(74, 168, 196, 0.35);
  }
}

@keyframes rippleExpand {
  0% {
    width: 110px;
    height: 110px;
    opacity: 0.7;
    border-color: rgba(74, 168, 196, 0.26);
  }
  60% {
    opacity: 0.28;
  }
  100% {
    width: 270px;
    height: 270px;
    opacity: 0;
    border-color: rgba(74, 168, 196, 0.03);
  }
}
</style>
