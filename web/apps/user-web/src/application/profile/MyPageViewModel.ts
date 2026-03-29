import { reactive } from 'vue'
import type { AuthApi } from '../../api/AuthApi'
import { loadSession } from '../identity/AuthSession'
import { PageLoadState } from '../../domain/ai/PageLoadState'
import type { UserMe } from '../../domain/identity/UserMe'

export interface MyPageState {
  pageState: PageLoadState
  me: UserMe | null
  displayName: string
  errorMessage: string
}

export class MyPageViewModel {
  readonly state: MyPageState = reactive({
    pageState: PageLoadState.Idle,
    me: null,
    displayName: '',
    errorMessage: ''
  })

  constructor(private readonly api: AuthApi) {}

  async initialize(): Promise<void> {
    const session = loadSession()
    if (!session) {
      this.state.pageState = PageLoadState.Error
      this.state.errorMessage = '请先登录'
      return
    }

    this.state.pageState = PageLoadState.Loading
    this.state.displayName = session.displayName
    this.state.errorMessage = ''

    try {
      const me = await this.api.getMe(session.accessToken)
      this.state.me = me
      this.state.pageState = PageLoadState.Success
    } catch (error) {
      this.state.pageState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '加载我的页失败'
    }
  }
}
