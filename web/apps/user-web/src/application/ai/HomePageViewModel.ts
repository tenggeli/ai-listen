import { reactive } from 'vue'
import type { AiApi } from '../../api/AiApi'
import type { AiHomeDashboard } from '../../domain/ai/AiHomeDashboard'
import { PageLoadState } from '../../domain/ai/PageLoadState'
import type { MatchCandidate } from '../../domain/ai/MatchCandidate'

export interface HomePageState {
  homeState: PageLoadState
  query: string
  overview: AiHomeDashboard | null
  remaining: number
  remainingState: PageLoadState
  matchState: PageLoadState
  candidates: MatchCandidate[]
  errorMessage: string
}

export class HomePageViewModel {
  readonly state: HomePageState = reactive({
    homeState: PageLoadState.Idle,
    query: '',
    overview: null,
    remaining: 0,
    remainingState: PageLoadState.Idle,
    matchState: PageLoadState.Idle,
    candidates: [],
    errorMessage: ''
  })

  constructor(private readonly api: AiApi, private readonly userId: string) {}

  async initialize(): Promise<void> {
    this.state.homeState = PageLoadState.Loading
    this.state.remainingState = PageLoadState.Loading
    this.state.errorMessage = ''
    try {
      const [overview, remaining] = await Promise.all([
        this.api.getHomeDashboard(this.userId),
        this.api.getRemaining(this.userId)
      ])
      this.state.overview = overview.withRemainingCount(remaining)
      this.state.remaining = remaining
      this.state.homeState = overview.hasContent() ? PageLoadState.Success : PageLoadState.Empty
      this.state.remainingState = PageLoadState.Success
    } catch (error) {
      this.state.homeState = PageLoadState.Error
      this.state.remainingState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '加载剩余次数失败'
    }
  }

  applyQuickAction(prompt: string): void {
    this.state.query = prompt
  }

  async submitMatch(): Promise<void> {
    if (!this.state.query.trim()) {
      this.state.candidates = []
      this.state.matchState = PageLoadState.Empty
      this.state.errorMessage = ''
      return
    }

    this.state.matchState = PageLoadState.Loading
    this.state.errorMessage = ''
    try {
      const result = await this.api.match(this.userId, this.state.query)
      this.state.remaining = result.remainingCount
      if (this.state.overview) {
        this.state.overview = this.state.overview.withRemainingCount(result.remainingCount)
      }
      this.state.candidates = result.candidates
      this.state.matchState = result.hasCandidates() ? PageLoadState.Success : PageLoadState.Empty
    } catch (error) {
      this.state.matchState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '匹配失败，请稍后重试'
    }
  }
}
