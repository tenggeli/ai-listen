import { ProviderDetail } from '../domain/provider/ProviderDetail'
import { ProviderReviewStatus } from '../domain/provider/ProviderReviewStatus'
import { ProviderSummary } from '../domain/provider/ProviderSummary'

export interface ListProvidersResult {
  items: ProviderSummary[]
  total: number
}

export interface ProviderAdminApi {
  listProviders(reviewStatus: string): Promise<ListProvidersResult>
  getProviderDetail(providerId: string): Promise<ProviderDetail>
  review(providerId: string, action: 'approve' | 'reject' | 'require-supplement', reason: string): Promise<ProviderReviewStatus>
}

export class HttpProviderAdminApi implements ProviderAdminApi {
  constructor(
    private readonly baseUrl = '/api/v1/admin',
    private readonly getAccessToken: () => string = () => ''
  ) {}

  async listProviders(reviewStatus: string): Promise<ListProvidersResult> {
    const query = reviewStatus ? `?review_status=${encodeURIComponent(reviewStatus)}` : ''
    const response = await fetch(`${this.baseUrl}/providers${query}`, {
      headers: this.buildAuthHeaders()
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load providers failed')
    }
    return {
      total: payload.data.total,
      items: (payload.data.items as any[]).map(
        (item) => new ProviderSummary(item.id, item.display_name, item.city_code, item.review_status)
      )
    }
  }

  async getProviderDetail(providerId: string): Promise<ProviderDetail> {
    const response = await fetch(`${this.baseUrl}/providers/${providerId}`, {
      headers: this.buildAuthHeaders()
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load provider detail failed')
    }
    const item = payload.data
    return new ProviderDetail(item.id, item.display_name, item.city_code, item.bio, item.review_status)
  }

  async review(providerId: string, action: 'approve' | 'reject' | 'require-supplement', reason: string): Promise<ProviderReviewStatus> {
    const response = await fetch(`${this.baseUrl}/providers/${providerId}/${action}`, {
      method: 'POST',
      headers: {
        ...this.buildAuthHeaders(),
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ reason })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'review failed')
    }
    return payload.data.review_status
  }

  private buildAuthHeaders(): Record<string, string> {
    const accessToken = this.getAccessToken()
    if (!accessToken) {
      return {}
    }
    return { Authorization: `Bearer ${accessToken}` }
  }
}

export class MockProviderAdminApi implements ProviderAdminApi {
  private providers = [
    new ProviderSummary('p_001', '暖心倾听师-小林', '310100', ProviderReviewStatus.Submitted),
    new ProviderSummary('p_002', '夜谈伙伴-阿泽', '110100', ProviderReviewStatus.UnderReview),
    new ProviderSummary('p_003', '电影散步搭子-念念', '440100', ProviderReviewStatus.SupplementRequired)
  ]

  async listProviders(reviewStatus: string): Promise<ListProvidersResult> {
    await sleep(220)
    const items = reviewStatus ? this.providers.filter((item) => item.reviewStatus === reviewStatus) : this.providers
    return { items, total: items.length }
  }

  async getProviderDetail(providerId: string): Promise<ProviderDetail> {
    await sleep(180)
    const target = this.providers.find((item) => item.id === providerId)
    if (!target) {
      throw new Error('provider not found')
    }
    return new ProviderDetail(target.id, target.displayName, target.cityCode, '这里是服务方资料简介（Mock）。', target.reviewStatus)
  }

  async review(providerId: string, action: 'approve' | 'reject' | 'require-supplement'): Promise<ProviderReviewStatus> {
    await sleep(220)
    const index = this.providers.findIndex((item) => item.id === providerId)
    if (index < 0) {
      throw new Error('provider not found')
    }

    const current = this.providers[index]
    const status =
      action === 'approve'
        ? ProviderReviewStatus.Approved
        : action === 'reject'
          ? ProviderReviewStatus.Rejected
          : ProviderReviewStatus.SupplementRequired
    this.providers[index] = new ProviderSummary(current.id, current.displayName, current.cityCode, status)
    return status
  }
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}
