import { UserIdentity } from '../../domain/identity/UserIdentity'
import { UserMe } from '../../domain/identity/UserMe'

const SESSION_KEY = 'listen_user_session'

export interface UserSession {
  userId: string
  accessToken: string
  refreshToken: string
  displayName: string
  profileCompleted: boolean
  personalityCompleted: boolean
}

export function loadSession(): UserSession | null {
  const raw = localStorage.getItem(SESSION_KEY)
  if (!raw) {
    return null
  }
  try {
    const parsed = JSON.parse(raw) as UserSession
    if (!parsed.userId || !parsed.accessToken) {
      return null
    }
    return parsed
  } catch {
    return null
  }
}

export function saveSession(identity: UserIdentity, me: UserMe): UserSession {
  const session: UserSession = {
    userId: identity.userId,
    accessToken: identity.accessToken,
    refreshToken: identity.refreshToken,
    displayName: identity.displayName,
    profileCompleted: me.profileCompleted,
    personalityCompleted: me.personalityCompleted
  }
  localStorage.setItem(SESSION_KEY, JSON.stringify(session))
  return session
}

export function updateSessionPatch(patch: Partial<UserSession>): UserSession | null {
  const current = loadSession()
  if (!current) {
    return null
  }
  const next = { ...current, ...patch }
  localStorage.setItem(SESSION_KEY, JSON.stringify(next))
  return next
}

export function clearSession(): void {
  localStorage.removeItem(SESSION_KEY)
}

export function currentUserIdOrDemo(): string {
  return loadSession()?.userId ?? 'demo-user-001'
}

export function nextOnboardingRoute(session: UserSession): string {
  if (!session.profileCompleted) {
    return '/profile/setup'
  }
  if (!session.personalityCompleted) {
    return '/personality/setup'
  }
  return '/home'
}

