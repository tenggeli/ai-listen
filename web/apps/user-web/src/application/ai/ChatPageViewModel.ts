import { reactive } from 'vue'
import type { AiApi } from '../../api/AiApi'
import type { AiSession } from '../../domain/ai/AiSession'
import { PageLoadState } from '../../domain/ai/PageLoadState'

export interface ChatPageState {
  sessionId: string
  draft: string
  quickReplies: string[]
  sessionState: PageLoadState
  sendState: PageLoadState
  session: AiSession | null
  errorMessage: string
}

export class ChatPageViewModel {
  readonly state: ChatPageState = reactive({
    sessionId: '',
    draft: '',
    quickReplies: ['工作压力太大了', '想找人说说话', '最近睡不太好'],
    sessionState: PageLoadState.Idle,
    sendState: PageLoadState.Idle,
    session: null,
    errorMessage: ''
  })

  constructor(
    private readonly api: AiApi,
    private readonly userId: string,
    private readonly sceneType = 'companion'
  ) {}

  async initialize(): Promise<void> {
    this.state.sessionState = PageLoadState.Loading
    this.state.errorMessage = ''

    try {
      if (!this.state.sessionId) {
        this.state.sessionId = await this.api.createSession(this.userId, this.sceneType)
      }
      await this.reloadSession()
    } catch (error) {
      this.state.sessionState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '会话初始化失败'
    }
  }

  async submitMessage(): Promise<void> {
    const content = this.state.draft.trim()
    if (!content) {
      this.state.sendState = PageLoadState.Empty
      return
    }

    if (!this.state.sessionId) {
      this.state.sendState = PageLoadState.Error
      this.state.errorMessage = '会话未初始化'
      return
    }

    this.state.sendState = PageLoadState.Loading
    this.state.errorMessage = ''

    try {
      await this.api.appendMessage(this.state.sessionId, 'user', content)
      this.state.draft = ''
      await this.reloadSession()
      this.state.sendState = PageLoadState.Success
    } catch (error) {
      this.state.sendState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '消息发送失败'
    }
  }

  useQuickReply(reply: string): void {
    this.state.draft = reply
  }

  private async reloadSession(): Promise<void> {
    const session = await this.api.getSession(this.state.sessionId)
    this.state.session = session
    this.state.sessionState = session.hasMessages() ? PageLoadState.Success : PageLoadState.Empty
  }
}
