import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../api/ApiError'
import type { ProviderAuthApi } from '../../api/ProviderAuthApi'
import { AuthService } from './AuthService'
import { AuthSessionStore } from './AuthSessionStore'

function createApiMock(): ProviderAuthApi {
  return {
    loginMock: vi.fn(),
    getMe: vi.fn(),
    saveProfile: vi.fn()
  }
}

describe('AuthService', () => {
  beforeEach(() => {
    const store = new Map<string, string>()
    Object.defineProperty(globalThis, 'localStorage', {
      configurable: true,
      value: {
        getItem: (key: string) => store.get(key) ?? null,
        setItem: (key: string, value: string) => {
          store.set(key, value)
        },
        removeItem: (key: string) => {
          store.delete(key)
        },
        clear: () => {
          store.clear()
        }
      }
    })
  })

  it('saves session when login succeeds', async () => {
    const api = createApiMock()
    vi.mocked(api.loginMock).mockResolvedValue({
      accessToken: 'token-1',
      providerId: 'provider-1'
    })
    const service = new AuthService(api, new AuthSessionStore())

    const session = await service.login('provider', 'provider123')

    expect(session).toEqual({ accessToken: 'token-1', providerId: 'provider-1' })
    expect(service.getAccessToken()).toBe('token-1')
  })

  it('keeps empty session when login fails', async () => {
    const api = createApiMock()
    vi.mocked(api.loginMock).mockRejectedValue(new ApiError('invalid credentials', 401))
    const service = new AuthService(api, new AuthSessionStore())

    await expect(service.login('provider', 'bad-password')).rejects.toThrow('invalid credentials')
    expect(service.getAccessToken()).toBe('')
  })

  it('restores session after refresh by validating token', async () => {
    const api = createApiMock()
    vi.mocked(api.getMe).mockResolvedValue({
      providerId: 'provider-2',
      account: 'provider',
      displayName: 'Provider',
      status: 'active',
      cityCode: '310000'
    })
    const store = new AuthSessionStore()
    store.saveSession({ accessToken: 'token-refresh', providerId: 'old-id' })
    const service = new AuthService(api, store)

    const isValid = await service.ensureSession()

    expect(isValid).toBe(true)
    expect(service.getAccessToken()).toBe('token-refresh')
    expect(store.getSession()).toEqual({
      accessToken: 'token-refresh',
      providerId: 'provider-2'
    })
  })

  it('clears session when token is unauthorized', async () => {
    const api = createApiMock()
    vi.mocked(api.getMe).mockRejectedValue(new ApiError('token expired', 401))
    const store = new AuthSessionStore()
    store.saveSession({ accessToken: 'expired-token', providerId: 'provider-1' })
    const service = new AuthService(api, store)

    const isValid = await service.ensureSession()

    expect(isValid).toBe(false)
    expect(store.getSession()).toBeNull()
  })
})
