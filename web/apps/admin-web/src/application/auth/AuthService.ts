import type { AdminAuthApi } from '../../api/AdminAuthApi'
import type { AdminSession } from './AuthSessionStore'
import { AuthSessionStore } from './AuthSessionStore'

export class AuthService {
  private verifiedToken = ''

  constructor(
    private readonly api: AdminAuthApi,
    private readonly sessionStore: AuthSessionStore
  ) {}

  async login(account: string, password: string): Promise<AdminSession> {
    const result = await this.api.loginMock(account, password)
    const session: AdminSession = {
      accessToken: result.accessToken,
      adminId: result.adminId,
      role: result.role
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
      const currentAdmin = await this.api.getMe(session.accessToken)
      this.sessionStore.saveSession({
        accessToken: session.accessToken,
        adminId: currentAdmin.adminId,
        role: currentAdmin.role
      })
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
