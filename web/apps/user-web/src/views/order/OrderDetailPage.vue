<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getOrderById } from '../../application/order/MockOrderStore'

const route = useRoute()
const router = useRouter()

const orderId = computed(() => String(route.params.id ?? ''))
const order = computed(() => getOrderById(orderId.value))
</script>

<template>
  <main class="order-page">
    <section class="screen-shell">
      <header class="top-row">
        <button type="button" class="back" @click="router.push('/services')">返回服务</button>
        <span>订单详情</span>
      </header>

      <section v-if="!order" class="card">
        <h1>订单不存在</h1>
        <p>可能是页面刷新后本地 mock 记录丢失，请返回服务页重新下单。</p>
      </section>

      <section v-else class="card">
        <h1>支付成功</h1>
        <p class="sub">这是订单详情占位页，后续可接真实订单中心接口。</p>
        <div class="line"><span>订单号</span><strong>{{ order.id }}</strong></div>
        <div class="line"><span>状态</span><strong>已支付</strong></div>
        <div class="line"><span>服务方</span><strong>{{ order.providerName }}</strong></div>
        <div class="line"><span>服务项目</span><strong>{{ order.serviceItemTitle }}</strong></div>
        <div class="line total"><span>支付金额</span><strong>¥{{ order.amount }}</strong></div>
        <button type="button" class="primary" @click="router.push('/services')">继续浏览服务</button>
      </section>
    </section>
  </main>
</template>

<style scoped>
.order-page {
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

.sub {
  margin: 8px 0 0;
  color: rgba(239, 247, 251, 0.62);
  font-size: 13px;
}

.line {
  margin-top: 12px;
  display: flex;
  justify-content: space-between;
  gap: 10px;
  font-size: 14px;
}

.line span {
  color: rgba(239, 247, 251, 0.56);
}

.line strong {
  text-align: right;
}

.line.total {
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  font-size: 18px;
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
