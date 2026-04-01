export interface AdminSession {
  accessToken: string
  adminId: string
  role: string
}

const SESSION_KEY = 'admin_web_auth_session'

export class AuthSessionStore {
  getSession(): AdminSession | null {
    const raw = localStorage.getItem(SESSION_KEY)
    if (!raw) {
      return null
    }

    try {
      const parsed = JSON.parse(raw) as Partial<AdminSession>
      if (!parsed.accessToken || !parsed.adminId || !parsed.role) {
        this.clear()
        return null
      }
      return {
        accessToken: parsed.accessToken,
        adminId: parsed.adminId,
        role: parsed.role
      }
    } catch {
      this.clear()
      return null
    }
  }

  saveSession(session: AdminSession): void {
    localStorage.setItem(SESSION_KEY, JSON.stringify(session))
  }

  clear(): void {
    localStorage.removeItem(SESSION_KEY)
  }

  getAccessToken(): string {
    const session = this.getSession()
    return session?.accessToken || ''
  }

  hasSession(): boolean {
    return this.getSession() !== null
  }
}
