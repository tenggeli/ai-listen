<script setup lang="ts">
import { useRouter } from 'vue-router'
import { authService } from '../../application/auth'
import AdminShell from '../../components/layout/AdminShell.vue'

const router = useRouter()

function logout(): void {
  authService.logout()
  void router.replace('/login')
}
</script>

<template>
  <AdminShell
    title="平台经营总览"
    subtitle="首页优先展示高风险与待处理事项，再给出平台经营规模，支撑运营、审核与风控的协同处理。"
    @logout="logout"
  >
    <template #topbar-right>
      <div class="admin-alert">
        <strong>今日待处理风险</strong>
        <p>服务方审核待处理、异常订单、投诉仲裁与高风险账号复审均有积压，请优先处理高优先级任务。</p>
      </div>
    </template>

    <section class="admin-metrics">
      <div class="metric">
        <small>服务方总数</small>
        <strong>1,286</strong>
      </div>
      <div class="metric">
        <small>在线可接单</small>
        <strong>368</strong>
      </div>
      <div class="metric">
        <small>今日订单量</small>
        <strong>542</strong>
      </div>
      <div class="metric">
        <small>待审核入驻</small>
        <strong>18</strong>
      </div>
      <div class="metric">
        <small>投诉处理中</small>
        <strong>7</strong>
      </div>
      <div class="metric">
        <small>待放款金额</small>
        <strong>¥ 86,420</strong>
      </div>
    </section>

    <section class="admin-section-grid">
      <article class="admin-card">
        <div class="admin-section-head">
          <h2>核心处理中心</h2>
          <span class="pill-soft">列表处理</span>
        </div>
        <p><RouterLink to="/providers/review">服务方审核中心</RouterLink>：快速处理准入审核、补件和复审。</p>
        <p><RouterLink to="/orders/manage">订单监管</RouterLink>：聚焦异常状态订单与人工介入。</p>
        <p><RouterLink to="/complaints/manage">投诉仲裁</RouterLink>：追踪投诉状态并执行处置动作。</p>
      </article>

      <article class="admin-card">
        <div class="admin-section-head">
          <h2>内容与配置</h2>
          <span class="pill">常态巡检</span>
        </div>
        <p><RouterLink to="/services/manage">服务项目管理</RouterLink>：维护服务项状态与基础信息。</p>
        <p><RouterLink to="/sounds/manage">声音内容管理</RouterLink>：管理声音素材的上/下线与编辑。</p>
      </article>
    </section>
  </AdminShell>
</template>
