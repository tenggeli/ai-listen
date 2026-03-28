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
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 18px;
  padding: 14px;
  background: rgba(10, 35, 48, 0.74);
  backdrop-filter: blur(12px);
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
  color: rgba(237, 247, 251, 0.32);
}

h3 {
  margin: 0;
  font-size: 16px;
  color: rgba(255, 255, 255, 0.88);
}

span {
  color: #93e0ff;
  font-size: 12px;
}

p {
  margin: 8px 0 0;
  color: rgba(237, 247, 251, 0.6);
  font-size: 12px;
  line-height: 1.8;
}

footer {
  display: flex;
  gap: 10px;
  margin-top: 14px;
}

button {
  flex: 1;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 14px;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(237, 247, 251, 0.72);
  cursor: pointer;
}

.ghost {
  background: rgba(115, 213, 255, 0.12);
  border-color: rgba(115, 213, 255, 0.28);
  color: #93e0ff;
}
</style>
