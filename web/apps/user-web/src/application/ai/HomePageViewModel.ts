import { reactive } from 'vue'
import type { AiApi } from '../../api/AiApi'
import { PageLoadState } from '../../domain/ai/PageLoadState'
import type { MatchCandidate } from '../../domain/ai/MatchCandidate'

export interface HomePageState {
  query: string
  remaining: number
  remainingState: PageLoadState
  matchState: PageLoadState
  candidates: MatchCandidate[]
  errorMessage: string
}

export class HomePageViewModel {
  readonly state: HomePageState = reactive({
    query: '',
    remaining: 0,
    remainingState: PageLoadState.Idle,
    matchState: PageLoadState.Idle,
    candidates: [],
    errorMessage: ''
  })

  constructor(private readonly api: AiApi, private readonly userId: string) {}

  async initialize(): Promise<void> {
    this.state.remainingState = PageLoadState.Loading
    this.state.errorMessage = ''
    try {
      this.state.remaining = await this.api.getRemaining(this.userId)
      this.state.remainingState = PageLoadState.Success
    } catch (error) {
      this.state.remainingState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '加载剩余次数失败'
    }
  }

  async submitMatch(): Promise<void> {
    this.state.matchState = PageLoadState.Loading
    this.state.errorMessage = ''
    try {
      const result = await this.api.match(this.userId, this.state.query)
      this.state.remaining = result.remainingCount
      this.state.candidates = result.candidates
      this.state.matchState = result.hasCandidates() ? PageLoadState.Success : PageLoadState.Empty
    } catch (error) {
      this.state.matchState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '匹配失败，请稍后重试'
    }
  }
}
