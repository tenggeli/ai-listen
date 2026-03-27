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
      <h3>{{ name }}</h3>
      <span>匹配分 {{ score.toFixed(2) }}</span>
    </header>
    <p>{{ reason }}</p>
  </article>
</template>

<style scoped>
.card {
  border: 1px solid #dbeafe;
  border-radius: 12px;
  padding: 12px;
  background: #ffffff;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 8px;
}

h3 {
  margin: 0;
  font-size: 15px;
}

span {
  color: #1d4ed8;
  font-size: 12px;
}

p {
  margin: 8px 0 0;
  color: #475569;
  font-size: 13px;
  line-height: 1.5;
}
</style>
