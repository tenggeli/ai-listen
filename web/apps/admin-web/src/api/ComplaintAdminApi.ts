export interface ComplaintSummary {
  orderId: string
  complaintStatus: string
  providerName: string
  userId: string
  complaintReason: string
  createdAt: string
}

export interface ComplaintDetail {
  complaintStatus: string
  order: {
    id: string
    userId: string
    providerName: string
    serviceItemTitle: string
    status: string
  }
  feedback: {
    complaintReason: string
    complaintContent: string
    reviewContent: string
    reviewTags: string[]
    ratingScore: number
  }
  actionLogs: Array<{
    actionId: string
    scope: string
    action: string
    operator: string
    reason: string
    updatedAt: string
  }>
}

export interface ComplaintAdminApi {
  list(): Promise<{ items: ComplaintSummary[]; total: number }>
  detail(orderId: string): Promise<ComplaintDetail>
  action(orderId: string, action: 'intervene' | 'resolve', reason: string): Promise<{ complaintStatus: string }>
}

export class HttpComplaintAdminApi implements ComplaintAdminApi {
  constructor(
    private readonly baseUrl = '/api/v1/admin',
    private readonly getAccessToken: () => string = () => ''
  ) {}

  async list(): Promise<{ items: ComplaintSummary[]; total: number }> {
    const response = await fetch(`${this.baseUrl}/complaints`, {
      headers: this.buildAuthHeaders()
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load complaints failed')
    }
    return {
      total: Number(payload.data.total ?? 0),
      items: (payload.data.items as any[]).map((item) => ({
        orderId: String(item.order?.id ?? ''),
        complaintStatus: String(item.complaint_status ?? ''),
        providerName: String(item.order?.provider_name ?? ''),
        userId: String(item.order?.user_id ?? ''),
        complaintReason: String(item.feedback?.complaint_reason ?? ''),
        createdAt: String(item.feedback?.created_at ?? '')
      }))
    }
  }

  async detail(orderId: string): Promise<ComplaintDetail> {
    const response = await fetch(`${this.baseUrl}/complaints/${orderId}`, {
      headers: this.buildAuthHeaders()
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load complaint detail failed')
    }
    const item = payload.data
    return {
      complaintStatus: String(item.complaint_status ?? ''),
      order: {
        id: String(item.order?.id ?? ''),
        userId: String(item.order?.user_id ?? ''),
        providerName: String(item.order?.provider_name ?? ''),
        serviceItemTitle: String(item.order?.service_item_title ?? ''),
        status: String(item.order?.status ?? '')
      },
      feedback: {
        complaintReason: String(item.feedback?.complaint_reason ?? ''),
        complaintContent: String(item.feedback?.complaint_content ?? ''),
        reviewContent: String(item.feedback?.review_content ?? ''),
        reviewTags: Array.isArray(item.feedback?.review_tags) ? item.feedback.review_tags.map((x: any) => String(x)) : [],
        ratingScore: Number(item.feedback?.rating_score ?? 0)
      },
      actionLogs: Array.isArray(item.action_logs)
        ? item.action_logs.map((log: any) => ({
            actionId: String(log.action_id ?? ''),
            scope: String(log.scope ?? ''),
            action: String(log.action ?? ''),
            operator: String(log.operator ?? ''),
            reason: String(log.reason ?? ''),
            updatedAt: String(log.updated_at ?? '')
          }))
        : []
    }
  }

  async action(orderId: string, action: 'intervene' | 'resolve', reason: string): Promise<{ complaintStatus: string }> {
    const response = await fetch(`${this.baseUrl}/complaints/${orderId}/${action}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...this.buildAuthHeaders()
      },
      body: JSON.stringify({ reason })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'operate complaint failed')
    }
    return { complaintStatus: String(payload.data?.complaint?.complaint_status ?? '') }
  }

  private buildAuthHeaders(): Record<string, string> {
    const accessToken = this.getAccessToken()
    if (!accessToken) {
      return {}
    }
    return { Authorization: `Bearer ${accessToken}` }
  }
}
