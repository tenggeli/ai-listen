export interface ProviderSession {
  accessToken: string
  providerId: string
}

const SESSION_KEY = 'provider_web_auth_session'

export class AuthSessionStore {
  getSession(): ProviderSession | null {
    const raw = localStorage.getItem(SESSION_KEY)
    if (!raw) {
      return null
    }

    try {
      const parsed = JSON.parse(raw) as Partial<ProviderSession>
      if (!parsed.accessToken || !parsed.providerId) {
        this.clear()
        return null
      }
      return {
        accessToken: parsed.accessToken,
        providerId: parsed.providerId
      }
    } catch {
      this.clear()
      return null
    }
  }

  saveSession(session: ProviderSession): void {
    localStorage.setItem(SESSION_KEY, JSON.stringify(session))
  }

  clear(): void {
    localStorage.removeItem(SESSION_KEY)
  }

  getAccessToken(): string {
    return this.getSession()?.accessToken ?? ''
  }
}
