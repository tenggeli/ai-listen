import type { ProviderOrder } from '../domain/order/ProviderOrder'

export interface ProviderOrderApi {
  listOrders(accessToken: string, page: number, pageSize: number): Promise<{ items: ProviderOrder[]; total: number }>
  getOrder(accessToken: string, orderId: string): Promise<ProviderOrder>
  operate(accessToken: string, orderId: string, action: 'accept' | 'depart' | 'arrive' | 'start' | 'complete'): Promise<string>
}

export class HttpProviderOrderApi implements ProviderOrderApi {
  constructor(private readonly baseUrl = '/api/v1/provider') {}

  async listOrders(accessToken: string, page: number, pageSize: number): Promise<{ items: ProviderOrder[]; total: number }> {
    const response = await fetch(`${this.baseUrl}/orders?page=${page}&page_size=${pageSize}`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load orders failed')
    }
    const items = Array.isArray(payload.data?.items) ? payload.data.items.map(mapOrder) : []
    const total = Number(payload.data?.total ?? items.length)
    return { items, total }
  }

  async getOrder(accessToken: string, orderId: string): Promise<ProviderOrder> {
    const response = await fetch(`${this.baseUrl}/orders/${orderId}`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load order detail failed')
    }
    return mapOrder(payload.data)
  }

  async operate(accessToken: string, orderId: string, action: 'accept' | 'depart' | 'arrive' | 'start' | 'complete'): Promise<string> {
    const response = await fetch(`${this.baseUrl}/orders/${orderId}/${action}`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'order action failed')
    }
    return payload.data.status
  }
}

function mapOrder(data: any): ProviderOrder {
  const status = String(data.status ?? 'created')
  return {
    id: String(data.id ?? ''),
    userId: String(data.user_id ?? ''),
    providerId: String(data.provider_id ?? ''),
    providerName: String(data.provider_name ?? ''),
    serviceItemId: String(data.service_item_id ?? ''),
    serviceItemTitle: String(data.service_item_title ?? ''),
    amount: Number(data.amount ?? 0),
    currency: String(data.currency ?? 'CNY'),
    status:
      status === 'paid' ||
      status === 'accepted' ||
      status === 'on_the_way' ||
      status === 'arrived' ||
      status === 'in_service' ||
      status === 'completed'
        ? status
        : 'created',
    createdAt: String(data.created_at ?? ''),
    paidAt: data.paid_at ? String(data.paid_at) : null
  }
}
