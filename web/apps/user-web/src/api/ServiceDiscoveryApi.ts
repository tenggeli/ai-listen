import { ProviderPublicProfile } from '../domain/service/ProviderPublicProfile'
import { ServiceCategory } from '../domain/service/ServiceCategory'
import { ServiceItem } from '../domain/service/ServiceItem'

export interface ListPublicProvidersInput {
  categoryId?: string
  keyword?: string
  page?: number
  pageSize?: number
}

export interface ListPublicProvidersOutput {
  items: ProviderPublicProfile[]
  total: number
}

export interface ServiceDiscoveryApi {
  getCategories(): Promise<ServiceCategory[]>
  listProviders(input: ListPublicProvidersInput): Promise<ListPublicProvidersOutput>
  getProvider(providerId: string): Promise<ProviderPublicProfile>
  getProviderServiceItems(providerId: string): Promise<ServiceItem[]>
}

export class HttpServiceDiscoveryApi implements ServiceDiscoveryApi {
  constructor(private readonly baseUrl = '/api/v1') {}

  async getCategories(): Promise<ServiceCategory[]> {
    const response = await fetch(`${this.baseUrl}/services/categories`)
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load categories failed')
    }
    return (payload.data.items as any[]).map((item) => new ServiceCategory(item.id, item.name, item.icon))
  }

  async listProviders(input: ListPublicProvidersInput): Promise<ListPublicProvidersOutput> {
    const query = new URLSearchParams()
    if (input.categoryId) {
      query.set('category_id', input.categoryId)
    }
    if (input.keyword) {
      query.set('keyword', input.keyword)
    }
    query.set('page', String(input.page ?? 1))
    query.set('page_size', String(input.pageSize ?? 10))

    const response = await fetch(`${this.baseUrl}/providers/public?${query.toString()}`)
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load providers failed')
    }
    return {
      items: (payload.data.items as any[]).map(buildProvider),
      total: payload.data.total
    }
  }

  async getProvider(providerId: string): Promise<ProviderPublicProfile> {
    const response = await fetch(`${this.baseUrl}/providers/public/${encodeURIComponent(providerId)}`)
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load provider detail failed')
    }
    return buildProvider(payload.data)
  }

  async getProviderServiceItems(providerId: string): Promise<ServiceItem[]> {
    const response = await fetch(`${this.baseUrl}/providers/public/${encodeURIComponent(providerId)}/service-items`)
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load service items failed')
    }
    return (payload.data.items as any[]).map(
      (item) =>
        new ServiceItem(
          item.id,
          item.provider_id,
          item.category_id,
          item.title,
          item.description,
          item.price_amount,
          item.price_unit,
          item.support_online
        )
    )
  }
}

export class MockServiceDiscoveryApi implements ServiceDiscoveryApi {
  private readonly categories = [
    new ServiceCategory('cat_all', '全部', 'sparkles'),
    new ServiceCategory('cat_food', '饭搭子', 'utensils'),
    new ServiceCategory('cat_movie', '电影搭子', 'film'),
    new ServiceCategory('cat_chat', '散步聊天', 'message-circle'),
    new ServiceCategory('cat_relax', '心理疏导', 'heart')
  ]

  private readonly providers = [
    new ProviderPublicProfile(
      'p_pub_001',
      '暖心倾听师 · 小林',
      'https://mock.listen.local/avatar/p_pub_001.png',
      '310100',
      '擅长情绪陪伴、轻聊天与晚间低压见面，适合刚下班、心里有点空的时候。',
      4.9,
      128,
      true,
      '实名认证',
      ['深夜聊天', '温柔倾听', '同城可约'],
      99,
      '30分钟'
    ),
    new ProviderPublicProfile(
      'p_pub_002',
      '电影散步搭子 · 念念',
      'https://mock.listen.local/avatar/p_pub_002.png',
      '440100',
      '适合今晚想出门透口气的人，擅长看电影、散步、轻社交。',
      4.8,
      76,
      true,
      '高评分',
      ['电影', '散步', 'ENFP'],
      158,
      '小时'
    ),
    new ProviderPublicProfile(
      'p_pub_003',
      '睡前放松向导 · 阿乔',
      'https://mock.listen.local/avatar/p_pub_003.png',
      '110100',
      '偏向声音陪伴、睡前放松和安抚式聊天，适合焦虑与失眠场景。',
      4.9,
      205,
      false,
      '夜间优先',
      ['助眠', '治愈声音', 'INFP'],
      99,
      '30分钟'
    )
  ]

  private readonly serviceItems = [
    new ServiceItem('si_001', 'p_pub_001', 'cat_chat', '线上倾听 30 分钟', '适合情绪安抚、想找人说话、睡前放松。', 99, '30分钟', true),
    new ServiceItem('si_002', 'p_pub_001', 'cat_chat', '同城散步聊天 1 小时', '适合下班后想透口气、想有人一起走一段路。', 158, '小时', false),
    new ServiceItem('si_003', 'p_pub_001', 'cat_movie', '电影陪伴 1 场', '偏轻社交陪伴型，不做高强度互动。', 188, '次', false),
    new ServiceItem('si_004', 'p_pub_002', 'cat_movie', '电影散步陪伴 1 场', '看电影后可散步聊天，低压轻社交。', 188, '次', false),
    new ServiceItem('si_005', 'p_pub_003', 'cat_relax', '睡前放松语音 30 分钟', '适合焦虑、失眠、需要被稳稳接住的时刻。', 99, '30分钟', true)
  ]

  async getCategories(): Promise<ServiceCategory[]> {
    await sleep(120)
    return [...this.categories]
  }

  async listProviders(input: ListPublicProvidersInput): Promise<ListPublicProvidersOutput> {
    await sleep(180)
    const keyword = (input.keyword ?? '').trim().toLowerCase()
    const categoryId = input.categoryId ?? 'cat_all'
    const filtered = this.providers.filter((provider) => {
      const matchedKeyword =
        keyword.length === 0 ||
        provider.displayName.toLowerCase().includes(keyword) ||
        provider.bio.toLowerCase().includes(keyword) ||
        provider.tags.some((tag) => tag.toLowerCase().includes(keyword))

      const matchedCategory =
        categoryId === 'cat_all' ||
        this.serviceItems.some((item) => item.providerId === provider.id && item.categoryId === categoryId)

      return matchedKeyword && matchedCategory
    })

    return { items: filtered, total: filtered.length }
  }

  async getProvider(providerId: string): Promise<ProviderPublicProfile> {
    await sleep(150)
    const provider = this.providers.find((item) => item.id === providerId)
    if (!provider) {
      throw new Error('服务方不存在')
    }
    return provider
  }

  async getProviderServiceItems(providerId: string): Promise<ServiceItem[]> {
    await sleep(180)
    return this.serviceItems.filter((item) => item.providerId === providerId)
  }
}

function buildProvider(item: any): ProviderPublicProfile {
  return new ProviderPublicProfile(
    item.id,
    item.display_name,
    item.avatar_url,
    item.city_code,
    item.bio,
    item.rating_avg,
    item.completed_orders,
    item.online,
    item.verification_label,
    item.tags ?? [],
    item.price_from,
    item.price_unit
  )
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}
