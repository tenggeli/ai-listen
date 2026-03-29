<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <main class="settings-page">
    <section class="screen-shell">
      <header class="top-row">
        <button type="button" class="back" @click="router.push('/me')">返回我的</button>
        <span>设置入口</span>
      </header>

      <section class="card">
        <h1>设置（占位）</h1>
        <p>这里将承载账号、通知、隐私、偏好等设置项。</p>
        <button type="button" class="entry" @click="router.push('/profile/setup')">编辑基础资料</button>
        <button type="button" class="entry" @click="router.push('/personality/setup')">编辑性格标签</button>
      </section>
    </section>
  </main>
</template>

<style scoped>
.settings-page {
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

.entry {
  margin-top: 12px;
  width: 100%;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.86);
  padding: 13px 12px;
  text-align: left;
}
</style>
