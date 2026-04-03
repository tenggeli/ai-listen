import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from './ApiError'
import { HttpProviderServiceApi } from './ProviderServiceApi'

describe('HttpProviderServiceApi', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('loads provider services list', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({
        code: 0,
        data: {
          items: [
            {
              item_id: 'si_001',
              provider_id: 'p_pub_001',
              category_id: 'cat_chat',
              title: '线上倾听 30 分钟',
              description: 'desc',
              price_amount: 99,
              price_unit: '30分钟',
              support_online: true,
              sort_order: 1
            }
          ]
        }
      })
    })
    vi.stubGlobal('fetch', fetchMock)
    const api = new HttpProviderServiceApi('/api/v1/provider')

    const items = await api.listServices('token-1')

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/provider/services', {
      headers: { Authorization: 'Bearer token-1' }
    })
    expect(items[0].itemId).toBe('si_001')
    expect(items[0].supportOnline).toBe(true)
  })

  it('returns 404 when provider services not found', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
        status: 404,
        json: async () => ({
          code: 404,
          message: 'provider not found',
          data: {}
        })
      })
    )
    const api = new HttpProviderServiceApi('/api/v1/provider')

    await expect(api.listServices('token-1')).rejects.toEqual(
      expect.objectContaining<ApiError>({
        statusCode: 404
      })
    )
  })
})
