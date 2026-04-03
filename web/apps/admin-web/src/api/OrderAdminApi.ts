import { getAdminOrderStatusReason, isAdminOrderStatus, type AdminOrderStatus } from '../domain/order/OrderStatus'

export interface AdminOrderSummary {
  id: string
  userId: string
  providerName: string
  serviceItemTitle: string
  amount: number
  status: AdminOrderStatus
  statusReason: string
  createdAt: string
}

export interface AdminOrderDetail {
  order: AdminOrderSummary & {
    providerId: string
    serviceItemId: string
    currency: string
    paidAt: string | null
  }
  feedback: {
    id: string
    ratingScore: number
    reviewTags: string[]
    reviewContent: string
    hasComplaint: boolean
    complaintReason: string
    complaintContent: string
  } | null
  actionLogs: Array<{
    actionId: string
    scope: string
    action: string
    operator: string
    reason: string
    statusBefore: string
    statusAfter: string
    statusAfterReason: string
    updatedAt: string
  }>
}

export interface AdminOrderApi {
  list(input: { status: string; keyword: string }): Promise<{ items: AdminOrderSummary[]; total: number }>
  detail(orderId: string): Promise<AdminOrderDetail>
  action(orderId: string, action: 'intervene' | 'close', reason: string): Promise<{ status: string }>
}

export class HttpOrderAdminApi implements AdminOrderApi {
  constructor(
    private readonly baseUrl = '/api/v1/admin',
    private readonly getAccessToken: () => string = () => ''
  ) {}

  async list(input: { status: string; keyword: string }): Promise<{ items: AdminOrderSummary[]; total: number }> {
    const query = new URLSearchParams()
    if (input.status) {
      query.set('status', input.status)
    }
    if (input.keyword) {
      query.set('keyword', input.keyword)
    }
    const suffix = query.toString() ? `?${query.toString()}` : ''
    const response = await fetch(`${this.baseUrl}/orders${suffix}`, {
      headers: this.buildAuthHeaders()
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load orders failed')
    }
    return {
      total: Number(payload.data.total ?? 0),
      items: (payload.data.items as any[]).map(mapOrderSummary)
    }
  }

  async detail(orderId: string): Promise<AdminOrderDetail> {
    const response = await fetch(`${this.baseUrl}/orders/${orderId}`, {
      headers: this.buildAuthHeaders()
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load order detail failed')
    }
    const order = mapOrderSummary(payload.data.order)
    return {
      order: {
        ...order,
        providerId: String(payload.data.order.provider_id ?? ''),
        serviceItemId: String(payload.data.order.service_item_id ?? ''),
        currency: String(payload.data.order.currency ?? 'CNY'),
        paidAt: payload.data.order.paid_at ? String(payload.data.order.paid_at) : null
      },
      feedback: payload.data.feedback
        ? {
            id: String(payload.data.feedback.id ?? ''),
            ratingScore: Number(payload.data.feedback.rating_score ?? 0),
            reviewTags: Array.isArray(payload.data.feedback.review_tags) ? payload.data.feedback.review_tags.map((x: any) => String(x)) : [],
            reviewContent: String(payload.data.feedback.review_content ?? ''),
            hasComplaint: Boolean(payload.data.feedback.has_complaint),
            complaintReason: String(payload.data.feedback.complaint_reason ?? ''),
            complaintContent: String(payload.data.feedback.complaint_content ?? '')
          }
        : null,
      actionLogs: Array.isArray(payload.data.action_logs)
        ? payload.data.action_logs.map((item: any) => ({
            actionId: String(item.action_id ?? ''),
            scope: String(item.scope ?? ''),
            action: String(item.action ?? ''),
            operator: String(item.operator ?? ''),
            reason: String(item.reason ?? ''),
            statusBefore: String(item.status_before ?? ''),
            statusAfter: String(item.status_after ?? ''),
            statusAfterReason: String(item.status_after_reason ?? ''),
            updatedAt: String(item.updated_at ?? '')
          }))
        : []
    }
  }

  async action(orderId: string, action: 'intervene' | 'close', reason: string): Promise<{ status: string }> {
    const response = await fetch(`${this.baseUrl}/orders/${orderId}/${action}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...this.buildAuthHeaders()
      },
      body: JSON.stringify({ reason })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'operate order failed')
    }
    return { status: String(payload.data?.order?.status ?? '') }
  }

  private buildAuthHeaders(): Record<string, string> {
    const accessToken = this.getAccessToken()
    if (!accessToken) {
      return {}
    }
    return { Authorization: `Bearer ${accessToken}` }
  }
}

function mapOrderSummary(item: any): AdminOrderSummary {
  const rawStatus = String(item.status ?? '')
  const status = isAdminOrderStatus(rawStatus) ? rawStatus : 'created'
  const statusReasonRaw = String(item.status_reason ?? '').trim()
  return {
    id: String(item.id ?? ''),
    userId: String(item.user_id ?? ''),
    providerName: String(item.provider_name ?? ''),
    serviceItemTitle: String(item.service_item_title ?? ''),
    amount: Number(item.amount ?? 0),
    status,
    statusReason: statusReasonRaw || getAdminOrderStatusReason(status),
    createdAt: String(item.created_at ?? '')
  }
}
