import type { ProviderOrder } from '../domain/order/ProviderOrder'
import { ApiError } from './ApiError'

export interface ProviderOrderApi {
  listOrders(
    accessToken: string,
    page: number,
    pageSize: number,
    status?: ProviderOrder['status']
  ): Promise<{ items: ProviderOrder[]; total: number }>
  getOrder(accessToken: string, orderId: string): Promise<ProviderOrder>
  operate(accessToken: string, orderId: string, action: 'accept' | 'depart' | 'arrive' | 'start' | 'complete'): Promise<string>
}

export class HttpProviderOrderApi implements ProviderOrderApi {
  constructor(
    private readonly baseUrl = '/api/v1/provider',
    private readonly timeoutMs = 8000
  ) {}

  async listOrders(
    accessToken: string,
    page: number,
    pageSize: number,
    status?: ProviderOrder['status']
  ): Promise<{ items: ProviderOrder[]; total: number }> {
    const query = new URLSearchParams({
      page: String(page),
      page_size: String(pageSize)
    })
    if (status) {
      query.set('status', status)
    }
    const response = await this.fetchWithTimeout(`${this.baseUrl}/orders?${query.toString()}`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new ApiError(payload.message || 'load orders failed', response.status)
    }
    const items = Array.isArray(payload.data?.items) ? payload.data.items.map(mapOrder) : []
    const total = Number(payload.data?.total ?? items.length)
    return { items, total }
  }

  async getOrder(accessToken: string, orderId: string): Promise<ProviderOrder> {
    const response = await this.fetchWithTimeout(`${this.baseUrl}/orders/${orderId}`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new ApiError(payload.message || 'load order detail failed', response.status)
    }
    return mapOrder(payload.data)
  }

  async operate(accessToken: string, orderId: string, action: 'accept' | 'depart' | 'arrive' | 'start' | 'complete'): Promise<string> {
    const response = await this.fetchWithTimeout(`${this.baseUrl}/orders/${orderId}/${action}`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new ApiError(payload.message || 'order action failed', response.status)
    }
    return String(payload.data?.status ?? '')
  }

  private async fetchWithTimeout(input: string, init: RequestInit): Promise<Response> {
    const controller = new AbortController()
    const timer = setTimeout(() => controller.abort(), this.timeoutMs)
    try {
      return await fetch(input, {
        ...init,
        signal: controller.signal
      })
    } catch (error) {
      if (error instanceof DOMException && error.name === 'AbortError') {
        throw new ApiError('网络请求超时，请稍后重试', 408)
      }
      throw new ApiError('网络异常，请检查连接后重试', 0)
    } finally {
      clearTimeout(timer)
    }
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
