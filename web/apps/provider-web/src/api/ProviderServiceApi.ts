import type { ProviderService } from '../domain/service/ProviderService'
import { ApiError } from './ApiError'

export interface ProviderServiceApi {
  listServices(accessToken: string): Promise<ProviderService[]>
}

export class HttpProviderServiceApi implements ProviderServiceApi {
  constructor(private readonly baseUrl = '/api/v1/provider') {}

  async listServices(accessToken: string): Promise<ProviderService[]> {
    const response = await fetch(`${this.baseUrl}/services`, {
      headers: {
        Authorization: `Bearer ${accessToken}`
      }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new ApiError(payload.message || 'load services failed', response.status)
    }
    const items = Array.isArray(payload.data?.items) ? payload.data.items : []
    return items.map((item: any) => ({
      itemId: String(item.item_id ?? ''),
      providerId: String(item.provider_id ?? ''),
      categoryId: String(item.category_id ?? ''),
      title: String(item.title ?? ''),
      description: String(item.description ?? ''),
      priceAmount: Number(item.price_amount ?? 0),
      priceUnit: String(item.price_unit ?? ''),
      supportOnline: Boolean(item.support_online),
      sortOrder: Number(item.sort_order ?? 0)
    }))
  }
}
