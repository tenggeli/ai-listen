<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const currentDate = computed(() =>
  new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  }).format(new Date())
)
</script>

<template>
  <main class="order-list-page">
    <section class="screen-shell">
      <header class="top-row">
        <button type="button" class="back" @click="router.push('/me')">返回我的</button>
        <span>订单入口</span>
      </header>

      <section class="card">
        <h1>订单列表（占位）</h1>
        <p>当前日期：{{ currentDate }}</p>
        <p>后续将接入真实订单列表接口与筛选能力。</p>
        <button type="button" class="primary" @click="router.push('/services')">去服务页下单</button>
      </section>
    </section>
  </main>
</template>

<style scoped>
.order-list-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top, rgba(74, 168, 196, 0.1), transparent 30%),
    linear-gradient(180deg, #0b1926 0%, #0d1f2d 42%, #08131f 100%);
  color: rgba(255, 255, 255, 0.88);
}

.screen-shell {
  width: min(100%, 390px);
  min-height: 100vh;
  margin: 0 auto;
  padding: 52px 16px 24px;
}

.top-row {
  display: flex;
  align-items: center;
  gap: 10px;
  color: rgba(255, 255, 255, 0.45);
  font-size: 12px;
}

.back {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.8);
  padding: 6px 10px;
}

.card {
  margin-top: 14px;
  border-radius: 18px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.045);
  padding: 16px;
}

.card h1 {
  margin: 0;
  font-size: 24px;
}

.card p {
  margin: 10px 0 0;
  color: rgba(239, 247, 251, 0.62);
  font-size: 14px;
  line-height: 1.7;
}

.primary {
  margin-top: 16px;
  width: 100%;
  border: none;
  border-radius: 14px;
  padding: 14px 10px;
  font-size: 14px;
  font-weight: 600;
  background: linear-gradient(135deg, #97e3ff, #58bee8);
  color: #082132;
}
</style>
