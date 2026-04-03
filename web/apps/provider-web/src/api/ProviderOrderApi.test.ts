import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from './ApiError'
import { HttpProviderOrderApi } from './ProviderOrderApi'

describe('HttpProviderOrderApi', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('passes pagination and status query params', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({
        code: 0,
        data: { items: [], total: 0 }
      })
    })
    vi.stubGlobal('fetch', fetchMock)
    const api = new HttpProviderOrderApi('/api/v1/provider')

    await api.listOrders('token-1', 2, 10, 'paid')

    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(fetchMock.mock.calls[0][0]).toBe('/api/v1/provider/orders?page=2&page_size=10&status=paid')
  })

  it('maps abort timeout to 408 error', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockRejectedValue(new DOMException('request aborted', 'AbortError'))
    )
    const api = new HttpProviderOrderApi('/api/v1/provider', 5)

    await expect(api.listOrders('token-1', 1, 10)).rejects.toEqual(
      expect.objectContaining<ApiError>({
        statusCode: 408
      })
    )
  })

  it('returns 404 error for missing order detail', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
        status: 404,
        json: async () => ({
          code: 404,
          message: 'order not found',
          data: {}
        })
      })
    )
    const api = new HttpProviderOrderApi('/api/v1/provider')

    await expect(api.getOrder('token-1', 'missing-id')).rejects.toEqual(
      expect.objectContaining<ApiError>({
        statusCode: 404
      })
    )
  })

  it('calls provider order action endpoint', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({
        code: 0,
        data: { status: 'accepted' }
      })
    })
    vi.stubGlobal('fetch', fetchMock)
    const api = new HttpProviderOrderApi('/api/v1/provider')

    const status = await api.operate('token-1', 'ord-1', 'accept')

    expect(status).toBe('accepted')
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v1/provider/orders/ord-1/accept',
      expect.objectContaining({
        method: 'POST',
        headers: { Authorization: 'Bearer token-1' }
      })
    )
  })

  it('returns 409 when action transition is invalid', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
        status: 409,
        json: async () => ({
          code: 409,
          message: 'invalid order transition',
          data: {}
        })
      })
    )
    const api = new HttpProviderOrderApi('/api/v1/provider')

    await expect(api.operate('token-1', 'ord-1', 'accept')).rejects.toEqual(
      expect.objectContaining<ApiError>({
        statusCode: 409
      })
    )
  })
})
