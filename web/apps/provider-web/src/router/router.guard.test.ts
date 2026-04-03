import { describe, expect, it, vi } from 'vitest'
import { createProviderMemoryRouter } from './index'

describe('provider router guard', () => {
  it('redirects unauthenticated user to /login', async () => {
    const ensureSession = vi.fn().mockResolvedValue(false)
    const router = createProviderMemoryRouter({ ensureSession })

    await router.push('/dashboard')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/login')
    expect(router.currentRoute.value.query.redirect).toBe('/dashboard')
  })

  it('redirects authenticated user away from /login to /dashboard', async () => {
    const ensureSession = vi.fn().mockResolvedValue(true)
    const router = createProviderMemoryRouter({ ensureSession })

    await router.push('/login')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/dashboard')
  })
})
