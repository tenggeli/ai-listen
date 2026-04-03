import { getUserOrderStatusReason, isUserOrderStatus, type UserOrder } from '../domain/order/UserOrder'

export interface OrderApi {
  createOrder(
    accessToken: string,
    input: {
      providerId: string
      providerName: string
      serviceItemId: string
      serviceItemTitle: string
      amount: number
      currency: string
    }
  ): Promise<UserOrder>
  listOrders(accessToken: string, page: number, pageSize: number): Promise<{ items: UserOrder[]; total: number }>
  getOrder(accessToken: string, orderId: string): Promise<UserOrder>
  payOrderMockSuccess(accessToken: string, orderId: string): Promise<UserOrder>
}

export class HttpOrderApi implements OrderApi {
  constructor(private readonly baseUrl = '/api/v1') {}

  async createOrder(
    accessToken: string,
    input: {
      providerId: string
      providerName: string
      serviceItemId: string
      serviceItemTitle: string
      amount: number
      currency: string
    }
  ): Promise<UserOrder> {
    const response = await fetch(`${this.baseUrl}/orders`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${accessToken}`
      },
      body: JSON.stringify({
        provider_id: input.providerId,
        provider_name: input.providerName,
        service_item_id: input.serviceItemId,
        service_item_title: input.serviceItemTitle,
        amount: input.amount,
        currency: input.currency
      })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'create order failed')
    }
    return mapOrder(payload.data)
  }

  async listOrders(accessToken: string, page: number, pageSize: number): Promise<{ items: UserOrder[]; total: number }> {
    const response = await fetch(`${this.baseUrl}/orders?page=${page}&page_size=${pageSize}`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'list orders failed')
    }
    const items = Array.isArray(payload.data?.items) ? payload.data.items.map(mapOrder) : []
    const total = Number(payload.data?.total ?? items.length)
    return { items, total }
  }

  async getOrder(accessToken: string, orderId: string): Promise<UserOrder> {
    const response = await fetch(`${this.baseUrl}/orders/${orderId}`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'get order failed')
    }
    return mapOrder(payload.data)
  }

  async payOrderMockSuccess(accessToken: string, orderId: string): Promise<UserOrder> {
    const response = await fetch(`${this.baseUrl}/orders/${orderId}/pay/mock-success`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'pay order failed')
    }
    return mapOrder(payload.data)
  }
}

function mapOrder(data: any): UserOrder {
  const rawStatus = String(data.status ?? 'created')
  const status = isUserOrderStatus(rawStatus) ? rawStatus : 'created'
  const statusReasonRaw = String(data.status_reason ?? '').trim()
  return {
    id: String(data.id ?? ''),
    userId: String(data.user_id ?? ''),
    providerId: String(data.provider_id ?? ''),
    providerName: String(data.provider_name ?? ''),
    serviceItemId: String(data.service_item_id ?? ''),
    serviceItemTitle: String(data.service_item_title ?? ''),
    amount: Number(data.amount ?? 0),
    currency: String(data.currency ?? 'CNY'),
    status,
    statusReason: statusReasonRaw || getUserOrderStatusReason(status),
    statusActionReason: String(data.status_action_reason ?? ''),
    statusUpdatedAt: data.status_updated_at ? String(data.status_updated_at) : null,
    createdAt: String(data.created_at ?? ''),
    paidAt: data.paid_at ? String(data.paid_at) : null
  }
}
