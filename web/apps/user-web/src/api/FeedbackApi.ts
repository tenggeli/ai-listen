import type { UserOrderFeedback } from '../domain/order/UserOrderFeedback'

export interface FeedbackApi {
  getOrderFeedback(accessToken: string, orderId: string): Promise<UserOrderFeedback>
  submitOrderFeedback(
    accessToken: string,
    orderId: string,
    input: {
      ratingScore: number
      reviewTags: string[]
      reviewContent: string
      complaintReason: string
      complaintContent: string
    }
  ): Promise<UserOrderFeedback>
}

export class HttpFeedbackApi implements FeedbackApi {
  constructor(private readonly baseUrl = '/api/v1') {}

  async getOrderFeedback(accessToken: string, orderId: string): Promise<UserOrderFeedback> {
    const response = await fetch(`${this.baseUrl}/orders/${orderId}/feedback`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'get feedback failed')
    }
    return mapFeedback(payload.data)
  }

  async submitOrderFeedback(
    accessToken: string,
    orderId: string,
    input: {
      ratingScore: number
      reviewTags: string[]
      reviewContent: string
      complaintReason: string
      complaintContent: string
    }
  ): Promise<UserOrderFeedback> {
    const response = await fetch(`${this.baseUrl}/orders/${orderId}/feedback`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${accessToken}`
      },
      body: JSON.stringify({
        rating_score: input.ratingScore,
        review_tags: input.reviewTags,
        review_content: input.reviewContent,
        complaint_reason: input.complaintReason,
        complaint_content: input.complaintContent
      })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'submit feedback failed')
    }
    return mapFeedback(payload.data)
  }
}

function mapFeedback(data: any): UserOrderFeedback {
  return {
    id: String(data.id ?? ''),
    orderId: String(data.order_id ?? ''),
    userId: String(data.user_id ?? ''),
    ratingScore: Number(data.rating_score ?? 0),
    reviewTags: Array.isArray(data.review_tags) ? data.review_tags : [],
    reviewContent: String(data.review_content ?? ''),
    hasComplaint: Boolean(data.has_complaint),
    complaintReason: String(data.complaint_reason ?? ''),
    complaintContent: String(data.complaint_content ?? ''),
    createdAt: String(data.created_at ?? '')
  }
}
