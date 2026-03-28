<script setup lang="ts">
const logoSvg = `
  <path d="M37 37 C33 41 33 49 37 53" stroke="#7ec8dc" stroke-width="2.6" stroke-linecap="round" fill="none"/>
  <path d="M30 31 C24 37 24 53 30 59" stroke="#aadde9" stroke-width="2.2" stroke-linecap="round" fill="none"/>
  <path d="M64 62 C58 72 46 74 38 66 C28 56 28 40 36 31 C44 22 58 22 66 32 L64 36" stroke="#4aa8c4" stroke-width="3.2" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
  <path d="M64 36 C70 40 70 52 64 56" stroke="#4aa8c4" stroke-width="3" stroke-linecap="round" fill="none"/>
  <path d="M60 40 C62 50 56 60 50 62 C47 63 46 66 50 68 C53 70 57 68 57 64" stroke="#4aa8c4" stroke-width="2.6" stroke-linecap="round" fill="none"/>
`
</script>

<template>
  <div class="typing-row">
    <div class="avatar" aria-hidden="true">
      <svg viewBox="0 0 90 90" fill="none" v-html="logoSvg" />
    </div>
    <div class="dots">
      <span />
      <span />
      <span />
    </div>
  </div>
</template>

<style scoped>
.typing-row {
  display: flex;
  gap: 10px;
  align-items: flex-end;
}

.avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 12px;
  flex-shrink: 0;
  border: 1px solid rgba(74, 168, 196, 0.24);
  background: rgba(74, 168, 196, 0.14);
}

.avatar svg {
  width: 16px;
  height: 16px;
}

.dots {
  display: flex;
  gap: 4px;
  padding: 14px 16px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 6px 20px 20px 20px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(12px);
}

.dots span {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: rgba(74, 168, 196, 0.54);
  animation: typingBounce 1.2s ease-in-out infinite;
}

.dots span:nth-child(2) {
  animation-delay: 0.2s;
}

.dots span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typingBounce {
  0%,
  100% {
    transform: translateY(0);
    opacity: 0.4;
  }
  50% {
    transform: translateY(-4px);
    opacity: 1;
  }
}
</style>
