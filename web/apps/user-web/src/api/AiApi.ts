import { AiHomeAggregate } from '../domain/ai/AiHomeAggregate'
import { MatchCandidate } from '../domain/ai/MatchCandidate'

export interface AiApi {
  getRemaining(userId: string): Promise<number>
  match(userId: string, inputText: string): Promise<AiHomeAggregate>
}

export class HttpAiApi implements AiApi {
  constructor(private readonly baseUrl = '/api/v1') {}

  async getRemaining(userId: string): Promise<number> {
    const response = await fetch(`${this.baseUrl}/ai/match/remaining?user_id=${encodeURIComponent(userId)}`)
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load remaining failed')
    }
    return payload.data.remaining
  }

  async match(userId: string, inputText: string): Promise<AiHomeAggregate> {
    const response = await fetch(`${this.baseUrl}/ai/match`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: userId, input_text: inputText })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'match failed')
    }

    const candidates = (payload.data.candidates as any[]).map(
      (item) => new MatchCandidate(item.provider_id, item.display_name, item.reason_text, item.score)
    )

    return new AiHomeAggregate(payload.data.remaining, candidates)
  }
}

export class MockAiApi implements AiApi {
  async getRemaining(): Promise<number> {
    await sleep(240)
    return 5
  }

  async match(_userId: string, inputText: string): Promise<AiHomeAggregate> {
    await sleep(520)
    if (!inputText.trim()) {
      return new AiHomeAggregate(5, [])
    }
    return new AiHomeAggregate(4, [
      new MatchCandidate('p_001', '暖心倾听师-小林', '你提到情绪低落，优先匹配擅长情绪陪伴的服务方', 0.93),
      new MatchCandidate('p_002', '夜谈伙伴-阿泽', '你希望有人陪聊，匹配夜间响应较快的服务方', 0.89),
      new MatchCandidate('p_003', '电影散步搭子-念念', '你有轻度线下陪伴诉求，匹配同城轻社交场景', 0.84)
    ])
  }
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}
