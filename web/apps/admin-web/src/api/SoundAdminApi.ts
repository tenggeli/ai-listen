export interface AdminSoundItem {
  id: string
  categoryKey: string
  title: string
  playCountText: string
  durationText: string
  emoji: string
  author: string
  sortOrder: number
  status: string
}

export interface AdminSoundListResult {
  items: AdminSoundItem[]
  total: number
}

export interface SoundAdminApi {
  list(input: { categoryKey: string; status: string; keyword: string; pageNo: number; pageSize: number }): Promise<AdminSoundListResult>
  create(input: Omit<AdminSoundItem, 'id'> & { id?: string }): Promise<AdminSoundItem>
  update(id: string, input: Omit<AdminSoundItem, 'id' | 'status'>): Promise<AdminSoundItem>
  changeStatus(id: string, action: 'activate' | 'deactivate'): Promise<string>
}

export class HttpSoundAdminApi implements SoundAdminApi {
  constructor(
    private readonly baseUrl = '/api/v1/admin',
    private readonly getAccessToken: () => string = () => ''
  ) {}

  async list(input: { categoryKey: string; status: string; keyword: string; pageNo: number; pageSize: number }): Promise<AdminSoundListResult> {
    const query = new URLSearchParams()
    if (input.categoryKey) query.set('category_key', input.categoryKey)
    if (input.status) query.set('status', input.status)
    if (input.keyword) query.set('keyword', input.keyword)
    if (input.pageNo > 0) query.set('page_no', String(input.pageNo))
    if (input.pageSize > 0) query.set('page_size', String(input.pageSize))
    const queryString = query.toString() ? `?${query.toString()}` : ''
    const response = await fetch(`${this.baseUrl}/sounds${queryString}`, { headers: this.buildAuthHeaders() })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load sounds failed')
    }
    return {
      total: payload.data.total,
      items: (payload.data.items as any[]).map(mapSoundItem)
    }
  }

  async create(input: Omit<AdminSoundItem, 'id'> & { id?: string }): Promise<AdminSoundItem> {
    const response = await fetch(`${this.baseUrl}/sounds`, {
      method: 'POST',
      headers: {
        ...this.buildAuthHeaders(),
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        track_id: input.id ?? '',
        category_key: input.categoryKey,
        title: input.title,
        play_count_text: input.playCountText,
        duration_text: input.durationText,
        emoji: input.emoji,
        author: input.author,
        sort_order: input.sortOrder,
        status: input.status
      })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'create sound failed')
    }
    return mapSoundItem(payload.data)
  }

  async update(id: string, input: Omit<AdminSoundItem, 'id' | 'status'>): Promise<AdminSoundItem> {
    const response = await fetch(`${this.baseUrl}/sounds/${encodeURIComponent(id)}`, {
      method: 'PUT',
      headers: {
        ...this.buildAuthHeaders(),
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        category_key: input.categoryKey,
        title: input.title,
        play_count_text: input.playCountText,
        duration_text: input.durationText,
        emoji: input.emoji,
        author: input.author,
        sort_order: input.sortOrder
      })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'update sound failed')
    }
    return mapSoundItem(payload.data)
  }

  async changeStatus(id: string, action: 'activate' | 'deactivate'): Promise<string> {
    const response = await fetch(`${this.baseUrl}/sounds/${encodeURIComponent(id)}/${action}`, {
      method: 'POST',
      headers: this.buildAuthHeaders()
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'update sound status failed')
    }
    return payload.data.status
  }

  private buildAuthHeaders(): Record<string, string> {
    const accessToken = this.getAccessToken()
    if (!accessToken) {
      return {}
    }
    return { Authorization: `Bearer ${accessToken}` }
  }
}

export class MockSoundAdminApi implements SoundAdminApi {
  private sounds: AdminSoundItem[] = [
    {
      id: 'track-rain',
      categoryKey: 'nature',
      title: '深夜雨声 · 安眠版',
      playCountText: '1,248 次播放',
      durationText: '35:00',
      emoji: '🌧',
      author: 'listen 治愈声音库',
      sortOrder: 1,
      status: 'active'
    }
  ]

  async list(input: { categoryKey: string; status: string; keyword: string; pageNo: number; pageSize: number }): Promise<AdminSoundListResult> {
    await sleep(160)
    const keyword = input.keyword.trim()
    let items = this.sounds.filter((item) => !input.categoryKey || item.categoryKey === input.categoryKey)
    items = items.filter((item) => !input.status || item.status === input.status)
    items = items.filter((item) => !keyword || item.id.includes(keyword) || item.title.includes(keyword) || item.author.includes(keyword))
    items.sort((a, b) => a.sortOrder - b.sortOrder || a.id.localeCompare(b.id))
    const pageNo = input.pageNo > 0 ? input.pageNo : 1
    const pageSize = input.pageSize > 0 ? input.pageSize : 20
    const start = (pageNo - 1) * pageSize
    return { total: items.length, items: items.slice(start, start + pageSize) }
  }

  async create(input: Omit<AdminSoundItem, 'id'> & { id?: string }): Promise<AdminSoundItem> {
    await sleep(160)
    const id = input.id?.trim() || `track_${Date.now()}`
    const item: AdminSoundItem = {
      id,
      categoryKey: input.categoryKey,
      title: input.title,
      playCountText: input.playCountText,
      durationText: input.durationText,
      emoji: input.emoji,
      author: input.author,
      sortOrder: input.sortOrder,
      status: input.status
    }
    this.sounds.unshift(item)
    return item
  }

  async update(id: string, input: Omit<AdminSoundItem, 'id' | 'status'>): Promise<AdminSoundItem> {
    await sleep(160)
    const index = this.sounds.findIndex((item) => item.id === id)
    if (index < 0) {
      throw new Error('sound not found')
    }
    this.sounds[index] = { ...this.sounds[index], ...input }
    return this.sounds[index]
  }

  async changeStatus(id: string, action: 'activate' | 'deactivate'): Promise<string> {
    await sleep(160)
    const index = this.sounds.findIndex((item) => item.id === id)
    if (index < 0) {
      throw new Error('sound not found')
    }
    const nextStatus = action === 'activate' ? 'active' : 'inactive'
    this.sounds[index] = { ...this.sounds[index], status: nextStatus }
    return nextStatus
  }
}

function mapSoundItem(item: any): AdminSoundItem {
  return {
    id: item.id,
    categoryKey: item.category_key,
    title: item.title,
    playCountText: item.play_count_text ?? '',
    durationText: item.duration_text ?? '',
    emoji: item.emoji ?? '',
    author: item.author ?? '',
    sortOrder: Number(item.sort_order ?? 0),
    status: item.status
  }
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}
