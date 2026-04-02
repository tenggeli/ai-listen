import type { ServiceItemDetail } from '../domain/service_item/ServiceItemDetail'
import type { ServiceItemSummary } from '../domain/service_item/ServiceItemSummary'

export interface ListServiceItemsResult {
  items: ServiceItemSummary[]
  total: number
}

export interface ServiceItemAdminApi {
  listServiceItems(input: { providerId: string; categoryId: string; status: string; keyword: string }): Promise<ListServiceItemsResult>
  getServiceItemDetail(serviceItemId: string): Promise<ServiceItemDetail>
  updateStatus(serviceItemId: string, action: 'activate' | 'deactivate'): Promise<string>
}

export class HttpServiceItemAdminApi implements ServiceItemAdminApi {
  constructor(
    private readonly baseUrl = '/api/v1/admin',
    private readonly getAccessToken: () => string = () => ''
  ) {}

  async listServiceItems(input: {
    providerId: string
    categoryId: string
    status: string
    keyword: string
  }): Promise<ListServiceItemsResult> {
    const query = new URLSearchParams()
    if (input.providerId) {
      query.set('provider_id', input.providerId)
    }
    if (input.categoryId) {
      query.set('category_id', input.categoryId)
    }
    if (input.status) {
      query.set('status', input.status)
    }
    if (input.keyword) {
      query.set('keyword', input.keyword)
    }
    const queryString = query.toString() ? `?${query.toString()}` : ''
    const response = await fetch(`${this.baseUrl}/service-items${queryString}`, {
      headers: this.buildAuthHeaders()
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load service items failed')
    }
    return {
      total: payload.data.total,
      items: (payload.data.items as any[]).map((item) => ({
        id: item.id,
        providerId: item.provider_id,
        providerName: item.provider_name,
        categoryId: item.category_id,
        title: item.title,
        priceAmount: Number(item.price_amount ?? 0),
        priceUnit: item.price_unit ?? '',
        serviceStatus: item.service_status
      }))
    }
  }

  async getServiceItemDetail(serviceItemId: string): Promise<ServiceItemDetail> {
    const response = await fetch(`${this.baseUrl}/service-items/${serviceItemId}`, {
      headers: this.buildAuthHeaders()
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load service item detail failed')
    }
    const item = payload.data
    return {
      id: item.id,
      providerId: item.provider_id,
      providerName: item.provider_name,
      categoryId: item.category_id,
      title: item.title,
      description: item.description,
      priceAmount: Number(item.price_amount ?? 0),
      priceUnit: item.price_unit ?? '',
      supportOnline: Boolean(item.support_online),
      sortOrder: Number(item.sort_order ?? 0),
      serviceStatus: item.service_status
    }
  }

  async updateStatus(serviceItemId: string, action: 'activate' | 'deactivate'): Promise<string> {
    const response = await fetch(`${this.baseUrl}/service-items/${serviceItemId}/${action}`, {
      method: 'POST',
      headers: this.buildAuthHeaders()
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'update status failed')
    }
    return payload.data.service_status
  }

  private buildAuthHeaders(): Record<string, string> {
    const accessToken = this.getAccessToken()
    if (!accessToken) {
      return {}
    }
    return { Authorization: `Bearer ${accessToken}` }
  }
}

export class MockServiceItemAdminApi implements ServiceItemAdminApi {
  private items: ServiceItemDetail[] = [
    {
      id: 'si_001',
      providerId: 'p_pub_001',
      providerName: '暖心倾听师 · 小林',
      categoryId: 'cat_chat',
      title: '线上倾听 30 分钟',
      description: '适合情绪安抚与睡前放松。',
      priceAmount: 99,
      priceUnit: '30分钟',
      supportOnline: true,
      sortOrder: 1,
      serviceStatus: 'active'
    },
    {
      id: 'si_002',
      providerId: 'p_pub_001',
      providerName: '暖心倾听师 · 小林',
      categoryId: 'cat_chat',
      title: '同城散步聊天 1 小时',
      description: '轻社交陪伴，节奏舒适。',
      priceAmount: 158,
      priceUnit: '小时',
      supportOnline: false,
      sortOrder: 2,
      serviceStatus: 'inactive'
    }
  ]

  async listServiceItems(input: { providerId: string; categoryId: string; status: string; keyword: string }): Promise<ListServiceItemsResult> {
    await sleep(220)
    const keyword = input.keyword.trim().toLowerCase()
    const items = this.items
      .filter((item) => !input.providerId || item.providerId === input.providerId)
      .filter((item) => !input.categoryId || item.categoryId === input.categoryId)
      .filter((item) => !input.status || item.serviceStatus === input.status)
      .filter((item) => !keyword || item.title.toLowerCase().includes(keyword) || item.providerName.toLowerCase().includes(keyword))
      .map((item) => ({
        id: item.id,
        providerId: item.providerId,
        providerName: item.providerName,
        categoryId: item.categoryId,
        title: item.title,
        priceAmount: item.priceAmount,
        priceUnit: item.priceUnit,
        serviceStatus: item.serviceStatus
      }))
    return { items, total: items.length }
  }

  async getServiceItemDetail(serviceItemId: string): Promise<ServiceItemDetail> {
    await sleep(180)
    const item = this.items.find((entry) => entry.id === serviceItemId)
    if (!item) {
      throw new Error('service item not found')
    }
    return item
  }

  async updateStatus(serviceItemId: string, action: 'activate' | 'deactivate'): Promise<string> {
    await sleep(180)
    const index = this.items.findIndex((entry) => entry.id === serviceItemId)
    if (index < 0) {
      throw new Error('service item not found')
    }
    const nextStatus = action === 'activate' ? 'active' : 'inactive'
    this.items[index] = {
      ...this.items[index],
      serviceStatus: nextStatus
    }
    return nextStatus
  }
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

