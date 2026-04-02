import type { ProviderAuthApi } from '../../api/ProviderAuthApi'
import type { ProviderSession } from './AuthSessionStore'
import { AuthSessionStore } from './AuthSessionStore'

export class AuthService {
  private verifiedToken = ''

  constructor(
    private readonly api: ProviderAuthApi,
    private readonly sessionStore: AuthSessionStore
  ) {}

  async login(account: string, password: string): Promise<ProviderSession> {
    const result = await this.api.loginMock(account, password)
    const session: ProviderSession = {
      accessToken: result.accessToken,
      providerId: result.providerId
    }
    this.sessionStore.saveSession(session)
    this.verifiedToken = result.accessToken
    return session
  }

  async ensureSession(): Promise<boolean> {
    const session = this.sessionStore.getSession()
    if (!session) {
      return false
    }
    if (session.accessToken === this.verifiedToken) {
      return true
    }
    try {
      const me = await this.api.getMe(session.accessToken)
      this.sessionStore.saveSession({ accessToken: session.accessToken, providerId: me.providerId })
      this.verifiedToken = session.accessToken
      return true
    } catch {
      this.logout()
      return false
    }
  }

  logout(): void {
    this.sessionStore.clear()
    this.verifiedToken = ''
  }

  getAccessToken(): string {
    return this.sessionStore.getAccessToken()
  }
}
