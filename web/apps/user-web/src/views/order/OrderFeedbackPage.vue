<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { HttpFeedbackApi } from '../../api/FeedbackApi'
import { HttpOrderApi } from '../../api/OrderApi'
import { loadSession } from '../../application/identity/AuthSession'
import { PageLoadState } from '../../domain/ai/PageLoadState'
import type { UserOrder } from '../../domain/order/UserOrder'

const route = useRoute()
const router = useRouter()
const orderApi = new HttpOrderApi('/api/v1')
const feedbackApi = new HttpFeedbackApi('/api/v1')

const quickTags = ['很会接情绪', '不尴尬', '回复及时', '声音舒服', '节奏太快', '与描述不符']
const complaintOptions = ['服务未按约开始', '服务内容与描述不符', '需要平台介入']

const state = reactive<{
  pageState: PageLoadState
  order: UserOrder | null
  ratingScore: number
  selectedTags: string[]
  reviewContent: string
  complaintReason: string
  complaintContent: string
  submitting: boolean
  submitted: boolean
  errorMessage: string
}>({
  pageState: PageLoadState.Idle,
  order: null,
  ratingScore: 0,
  selectedTags: [],
  reviewContent: '',
  complaintReason: '',
  complaintContent: '',
  submitting: false,
  submitted: false,
  errorMessage: ''
})

const orderId = computed(() => String(route.params.id ?? ''))
const scoreCandidates = Array.from({ length: 10 }, (_, i) => i + 1)

onMounted(() => {
  void initialize()
})

async function initialize(): Promise<void> {
  const session = loadSession()
  if (!session) {
    await router.push('/auth')
    return
  }
  if (!orderId.value) {
    state.pageState = PageLoadState.Error
    state.errorMessage = '订单编号无效'
    return
  }

  state.pageState = PageLoadState.Loading
  state.errorMessage = ''
  try {
    const order = await orderApi.getOrder(session.accessToken, orderId.value)
    state.order = order
    try {
      const existing = await feedbackApi.getOrderFeedback(session.accessToken, orderId.value)
      state.ratingScore = existing.ratingScore
      state.selectedTags = [...existing.reviewTags]
      state.reviewContent = existing.reviewContent
      state.complaintReason = existing.complaintReason
      state.complaintContent = existing.complaintContent
      state.submitted = true
    } catch {
      state.submitted = false
    }
    state.pageState = PageLoadState.Success
  } catch (error) {
    state.pageState = PageLoadState.Error
    state.errorMessage = error instanceof Error ? error.message : '页面加载失败'
  }
}

function toggleTag(tag: string): void {
  if (state.submitted) {
    return
  }
  const idx = state.selectedTags.indexOf(tag)
  if (idx >= 0) {
    state.selectedTags.splice(idx, 1)
    return
  }
  state.selectedTags.push(tag)
}

function selectComplaint(reason: string): void {
  if (state.submitted) {
    return
  }
  state.complaintReason = reason
}

async function submit(): Promise<void> {
  if (state.submitted) {
    await router.push('/me')
    return
  }
  const session = loadSession()
  if (!session || !state.order) {
    return
  }

  const hasReview = state.ratingScore > 0 || state.selectedTags.length > 0 || state.reviewContent.trim() !== ''
  const hasComplaint = state.complaintReason.trim() !== '' || state.complaintContent.trim() !== ''
  if (!hasReview && !hasComplaint) {
    state.errorMessage = '请至少填写评价或投诉中的一项。'
    return
  }
  if (hasReview && state.ratingScore < 1) {
    state.errorMessage = '提交评价时请给出 1-10 分评分。'
    return
  }
  if (hasComplaint && state.complaintReason.trim() === '') {
    state.errorMessage = '提交投诉时请选择投诉原因。'
    return
  }

  state.errorMessage = ''
  state.submitting = true
  try {
    await feedbackApi.submitOrderFeedback(session.accessToken, state.order.id, {
      ratingScore: state.ratingScore,
      reviewTags: state.selectedTags,
      reviewContent: state.reviewContent.trim(),
      complaintReason: state.complaintReason.trim(),
      complaintContent: state.complaintContent.trim()
    })
    state.submitted = true
  } catch (error) {
    state.errorMessage = error instanceof Error ? error.message : '提交失败，请稍后再试。'
  } finally {
    state.submitting = false
  }
}
</script>

<template>
  <main class="feedback-page">
    <section class="screen-shell">
      <header class="top-row">
        <button type="button" class="back" @click="router.push(`/orders/${orderId}`)">返回订单</button>
        <span>评价投诉</span>
      </header>

      <section v-if="state.pageState === PageLoadState.Loading" class="card">
        <h1>加载中</h1>
        <p>正在同步订单信息...</p>
      </section>
      <section v-else-if="state.pageState === PageLoadState.Error" class="card">
        <h1>加载失败</h1>
        <p>{{ state.errorMessage || '请稍后重试。' }}</p>
      </section>

      <template v-if="state.pageState === PageLoadState.Success && state.order">
        <section class="card hero">
          <h1>把这次体验留下来，也给售后留出口</h1>
          <p>订单：{{ state.order.serviceItemTitle }} · {{ state.order.providerName }}</p>
          <p v-if="state.submitted" class="ok">本订单已提交评价/投诉，感谢反馈。</p>
        </section>

        <section class="card">
          <h2>服务评价（1-10 分）</h2>
          <div class="score-grid">
            <button
              v-for="score in scoreCandidates"
              :key="score"
              type="button"
              class="score-item"
              :class="{ active: state.ratingScore === score }"
              :disabled="state.submitted"
              @click="state.ratingScore = score"
            >
              {{ score }}
            </button>
          </div>
          <div class="tags">
            <button
              v-for="tag in quickTags"
              :key="tag"
              type="button"
              class="tag"
              :class="{ active: state.selectedTags.includes(tag) }"
              :disabled="state.submitted"
              @click="toggleTag(tag)"
            >
              {{ tag }}
            </button>
          </div>
          <textarea
            v-model="state.reviewContent"
            class="textarea"
            :disabled="state.submitted"
            placeholder="写一点你的真实感受，帮助下一位用户。"
          />
        </section>

        <section class="card">
          <h2>投诉 / 异常反馈</h2>
          <div class="issue-list">
            <button
              v-for="reason in complaintOptions"
              :key="reason"
              type="button"
              class="issue"
              :class="{ active: state.complaintReason === reason }"
              :disabled="state.submitted"
              @click="selectComplaint(reason)"
            >
              {{ reason }}
            </button>
          </div>
          <textarea
            v-model="state.complaintContent"
            class="textarea"
            :disabled="state.submitted"
            placeholder="可补充投诉细节，便于平台快速处理。"
          />
        </section>

        <p v-if="state.errorMessage" class="error">{{ state.errorMessage }}</p>

        <section class="actions">
          <button type="button" class="ghost" @click="router.push('/me')">回到我的</button>
          <button type="button" class="primary" :disabled="state.submitting" @click="submit">
            {{ state.submitting ? '提交中...' : state.submitted ? '已提交，去我的' : '提交评价/投诉' }}
          </button>
        </section>
      </template>
    </section>
  </main>
</template>

<style scoped>
.feedback-page {
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

.hero h1 {
  margin: 0;
  font-size: 22px;
}

.hero p {
  margin: 10px 0 0;
  color: rgba(239, 247, 251, 0.62);
  font-size: 13px;
}

.ok {
  color: #93f5c7;
}

.card h2 {
  margin: 0;
  font-size: 16px;
}

.score-grid {
  margin-top: 12px;
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 8px;
}

.score-item,
.tag,
.issue {
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.04);
  color: rgba(239, 247, 251, 0.82);
  font-size: 13px;
}

.score-item {
  padding: 10px 0;
}

.tags {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  padding: 8px 10px;
}

.issue-list {
  margin-top: 12px;
  display: grid;
  gap: 8px;
}

.issue {
  padding: 10px 12px;
  text-align: left;
}

.active {
  border-color: rgba(151, 227, 255, 0.4);
  background: rgba(115, 213, 255, 0.12);
}

.textarea {
  margin-top: 12px;
  width: 100%;
  min-height: 92px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.04);
  color: rgba(239, 247, 251, 0.9);
  padding: 10px 12px;
  resize: vertical;
}

.actions {
  margin-top: 14px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.ghost,
.primary {
  border: none;
  border-radius: 14px;
  padding: 14px 10px;
  font-size: 14px;
  font-weight: 600;
}

.ghost {
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(239, 247, 251, 0.74);
}

.primary {
  background: linear-gradient(135deg, #97e3ff, #58bee8);
  color: #082132;
}

.primary:disabled {
  opacity: 0.65;
}

.error {
  margin-top: 12px;
  color: #fca5a5;
  font-size: 13px;
}
</style>
