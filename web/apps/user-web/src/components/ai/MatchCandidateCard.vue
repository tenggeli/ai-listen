<script setup lang="ts">
defineProps<{
  name: string
  reason: string
  score: number
}>()
</script>

<template>
  <article class="card">
    <header>
      <div>
        <p class="eyebrow">推荐对象</p>
        <h3>{{ name }}</h3>
      </div>
      <span>匹配度 {{ Math.round(score * 100) }}%</span>
    </header>
    <p>{{ reason }}</p>
    <footer>
      <button type="button">查看详情</button>
      <button type="button" class="ghost">发起邀约</button>
    </footer>
  </article>
</template>

<style scoped>
.card {
  border: 1px solid rgba(74, 168, 196, 0.14);
  border-radius: 18px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
}

header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 8px;
}

.eyebrow {
  margin: 0 0 6px;
  font-size: 11px;
  letter-spacing: 0.08em;
  color: rgba(255, 255, 255, 0.28);
}

h3 {
  margin: 0;
  font-size: 16px;
  color: rgba(255, 255, 255, 0.88);
}

span {
  color: #7ec8dc;
  font-size: 12px;
}

p {
  margin: 8px 0 0;
  color: rgba(255, 255, 255, 0.62);
  font-size: 13px;
  line-height: 1.5;
}

footer {
  display: flex;
  gap: 10px;
  margin-top: 14px;
}

button {
  flex: 1;
  border: 1px solid rgba(74, 168, 196, 0.2);
  border-radius: 999px;
  padding: 10px 12px;
  background: rgba(30, 65, 88, 0.92);
  color: rgba(255, 255, 255, 0.86);
  cursor: pointer;
}

.ghost {
  background: rgba(255, 255, 255, 0.03);
}
</style>
