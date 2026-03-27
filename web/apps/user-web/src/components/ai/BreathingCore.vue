<template>
  <div class="breath-wrapper" aria-hidden="true">
    <div class="core" />
  </div>
</template>

<style scoped>
.breath-wrapper {
  display: grid;
  place-items: center;
  width: 148px;
  height: 148px;
  margin: 16px auto;
  border-radius: 999px;
  background: radial-gradient(circle, rgba(43, 97, 255, 0.12), transparent 68%);
}

.core {
  width: 88px;
  height: 88px;
  border-radius: 999px;
  background: linear-gradient(160deg, #1d4ed8, #0ea5e9);
  animation: breathing 2.6s infinite ease-in-out;
  box-shadow: 0 12px 32px rgba(37, 99, 235, 0.35);
}

@keyframes breathing {
  0% { transform: scale(0.95); }
  50% { transform: scale(1.08); }
  100% { transform: scale(0.95); }
}
</style>
