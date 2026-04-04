<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { isUnauthorizedError } from '../../api/ApiError'
import { HttpProviderOrderApi } from '../../api/ProviderOrderApi'
import { authService } from '../../application/auth'
import { PageLoadState } from '../../domain/common/PageLoadState'
import ProviderShell from '../../components/layout/ProviderShell.vue'

const router = useRouter()
const orderApi = new HttpProviderOrderApi('/api/v1/provider')
const state = reactive({
  pendingCount: 0,
  pageState: PageLoadState.Idle,
  errorMessage: ''
})

onMounted(() => {
  void loadSummary()
})

async function loadSummary(): Promise<void> {
  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  try {
    const result = await orderApi.listOrders(authService.getAccessToken(), 1, 50)
    state.pendingCount = result.items.filter((item) => item.status === 'paid').length
    state.pageState = PageLoadState.Success
  } catch (error) {
    if (isUnauthorizedError(error)) {
      authService.logout()
      await router.replace('/login')
      return
    }
    state.pageState = PageLoadState.Error
    state.errorMessage = error instanceof Error ? error.message : '加载失败'
  }
}

function logout(): void {
  authService.logout()
  void router.replace('/login')
}
</script>

<template>
  <ProviderShell
    title="服务方工作台"
    subtitle="首页强调当前状态、待处理任务、今日收益与快捷履约动作，帮助服务方一屏完成关键操作。"
    @logout="logout"
  >
    <template #topbar-right>
      <div class="provider-status-card">
        <div class="provider-status-head">
          <strong>当前工作状态</strong>
          <span class="provider-pill">空闲中，可接新单</span>
        </div>
        <p class="provider-page-subtitle">认证状态：已通过 ｜ 信用等级：A ｜ 当前评分：9.4 ｜ 近 30 天取消率：1.2%</p>
      </div>
    </template>

    <section class="provider-kpi-grid">
      <article class="provider-card provider-kpi">
        <small>待处理订单</small>
        <strong>{{ state.pageState === PageLoadState.Loading ? '...' : `${state.pendingCount} 单` }}</strong>
      </article>
      <article class="provider-card provider-kpi">
        <small>今日完成</small>
        <strong>4 单</strong>
      </article>
      <article class="provider-card provider-kpi">
        <small>今日预估收入</small>
        <strong>¥ 1,280</strong>
      </article>
      <article class="provider-card provider-kpi">
        <small>当前积分</small>
        <strong>2,460</strong>
      </article>
      <article class="provider-card provider-kpi">
        <small>待处理任务</small>
        <strong>5 项</strong>
      </article>
    </section>

    <section class="provider-grid">
      <article class="provider-card">
        <div class="provider-section-head">
          <h2>待处理订单</h2>
          <span class="provider-badge-soft">按履约优先级排序</span>
        </div>
        <div class="provider-table">
          <div class="provider-row-head">
            <div>订单 / 用户</div>
            <div>服务项目</div>
            <div>预约时间</div>
            <div>状态</div>
            <div>操作</div>
          </div>
          <div class="provider-row">
            <div><strong>待处理汇总</strong></div>
            <div>当前待接单任务</div>
            <div>实时更新</div>
            <div><span class="provider-badge-warn">{{ state.pendingCount }} 单待确认</span></div>
            <div><RouterLink to="/orders">进入订单列表</RouterLink></div>
          </div>
        </div>
      </article>

      <article class="provider-card">
        <div class="provider-section-head">
          <h2>快捷入口</h2>
          <span class="provider-badge-soft">高频操作</span>
        </div>
        <div class="provider-actions-grid">
          <div class="provider-action-card">
            <strong>管理服务项目</strong>
            <span>新增项目、调整价格、查看上架状态。</span>
          </div>
          <div class="provider-action-card">
            <strong>资料编辑</strong>
            <span>维护昵称、城市与基础资料，保持资料完整度。</span>
          </div>
          <div class="provider-action-card">
            <strong>订单履约</strong>
            <span>按状态推进接单、出发、到达和完单动作。</span>
          </div>
          <div class="provider-action-card">
            <strong>经营提醒</strong>
            <span>查看待处理任务与状态变化，及时响应。</span>
          </div>
        </div>
      </article>
    </section>

    <p v-if="state.pageState === PageLoadState.Idle" class="provider-sub">等待加载工作台数据...</p>
    <p v-else-if="state.pageState === PageLoadState.Error" class="provider-error">{{ state.errorMessage }}</p>
    <p v-else-if="state.pageState === PageLoadState.Success" class="provider-success">工作台数据已更新。</p>
  </ProviderShell>
</template>
