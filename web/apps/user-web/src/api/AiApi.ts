import { AiHomeDashboard } from '../domain/ai/AiHomeDashboard'
import { AiHomeAggregate } from '../domain/ai/AiHomeAggregate'
import { AiHomeQuickAction } from '../domain/ai/AiHomeQuickAction'
import { AiActionCard } from '../domain/ai/AiActionCard'
import { AiMessage } from '../domain/ai/AiMessage'
import { AiSession } from '../domain/ai/AiSession'
import { MatchCandidate } from '../domain/ai/MatchCandidate'

export interface AiApi {
  getHomeDashboard(userId: string): Promise<AiHomeDashboard>
  getRemaining(userId: string): Promise<number>
  match(userId: string, inputText: string): Promise<AiHomeAggregate>
  createSession(userId: string, sceneType: string): Promise<string>
  getSession(sessionId: string): Promise<AiSession>
  appendMessage(sessionId: string, senderType: string, content: string): Promise<void>
}

export class HttpAiApi implements AiApi {
  constructor(private readonly baseUrl = '/api/v1') {}

  async getHomeDashboard(userId: string): Promise<AiHomeDashboard> {
    const response = await fetch(`${this.baseUrl}/ai/home?user_id=${encodeURIComponent(userId)}`)
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load home failed')
    }
    return buildHomeDashboard(payload.data)
  }

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

  async createSession(userId: string, sceneType: string): Promise<string> {
    const response = await fetch(`${this.baseUrl}/ai/sessions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: userId, scene_type: sceneType })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'create session failed')
    }
    return payload.data.session_id
  }

  async getSession(sessionId: string): Promise<AiSession> {
    const response = await fetch(`${this.baseUrl}/ai/sessions/${encodeURIComponent(sessionId)}`)
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load session failed')
    }
    return buildSession(payload.data)
  }

  async appendMessage(sessionId: string, senderType: string, content: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/ai/sessions/${encodeURIComponent(sessionId)}/messages`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ sender_type: senderType, content })
    })

    let payload: any = null
    try {
      payload = await response.json()
    } catch {
      payload = null
    }

    if (!response.ok || payload?.code !== 0) {
      throw new Error(payload?.message || 'AI 回复暂时超时，请稍后重试')
    }
  }
}

export class MockAiApi implements AiApi {
  private readonly sessions = new Map<string, AiSession>()

  async getHomeDashboard(_userId: string): Promise<AiHomeDashboard> {
    await sleep(220)
    return new AiHomeDashboard(
      '晚上好 · 周六',
      '今晚，遇见懂你的人',
      '有 1,247 位搭子正在等待陪伴',
      '🌙',
      '平静 · 有点想聊聊',
      '上海 21°C 微风',
      28,
      1247,
      312,
      5,
      buildQuickActions()
    )
  }

  async getRemaining(): Promise<number> {
    await sleep(240)
    return 5
  }

  async match(_userId: string, inputText: string): Promise<AiHomeAggregate> {
    await sleep(520)
    if (!inputText.trim()) {
      return new AiHomeAggregate(5, [])
    }

    const pool = [
      new MatchCandidate('p_001', '暖心倾听师-小林', '你提到情绪低落，优先匹配擅长情绪陪伴的服务方', 0.93),
      new MatchCandidate('p_002', '夜谈伙伴-阿泽', '你希望有人陪聊，匹配夜间响应较快的服务方', 0.89),
      new MatchCandidate('p_003', '电影散步搭子-念念', '你有轻度线下陪伴诉求，匹配同城轻社交场景', 0.84),
      new MatchCandidate('p_004', '睡前放松向导-阿乔', '你提到睡眠困扰，优先匹配放松场景服务方', 0.86),
      new MatchCandidate('p_005', '关系梳理教练-木子', '你提到关系压力，优先匹配关系倾听方向', 0.9)
    ]

    const candidates = [...pool].sort(() => Math.random() - 0.5).slice(0, 3)
    return new AiHomeAggregate(4, candidates)
  }

  async createSession(userId: string, sceneType: string): Promise<string> {
    await sleep(180)
    const sessionId = `sess_${Date.now()}`
    const initialTime = new Date().toISOString()
    this.sessions.set(
      sessionId,
      new AiSession(sessionId, userId, sceneType, 'active', initialTime, [
        new AiMessage(
          'assistant',
          '一个人走在夜路上的时候，情绪会更容易被放大。\n\n如果你愿意，今晚的这段路我可以陪你一起走。',
          initialTime
        )
      ])
    )
    return sessionId
  }

  async getSession(sessionId: string): Promise<AiSession> {
    await sleep(180)
    const session = this.sessions.get(sessionId)
    if (!session) {
      throw new Error('session not found')
    }
    return session
  }

  async appendMessage(sessionId: string, senderType: string, content: string): Promise<void> {
    await sleep(180)
    const session = this.sessions.get(sessionId)
    if (!session) {
      throw new Error('session not found')
    }

    const now = new Date()
    const userTime = now.toISOString()
    const replyTime = new Date(now.getTime() + 1200).toISOString()
    const nextMessages = [
      ...session.messages,
      new AiMessage(senderType, content, userTime),
      new AiMessage(
        'assistant',
        buildAssistantReply(content),
        replyTime,
        buildActionCard(content),
        buildSafetyLevel(content)
      )
    ]
    this.sessions.set(
      sessionId,
      new AiSession(session.id, session.userId, session.sceneType, session.status, replyTime, nextMessages)
    )
  }
}

function buildSession(payload: any): AiSession {
  const messages = (payload.messages as any[]).map((item) => {
    const cardPayload = item.action_card
    const actionCard = cardPayload
      ? new AiActionCard(
          cardPayload.action ?? '',
          cardPayload.title ?? '',
          cardPayload.description ?? '',
          cardPayload.route ?? '',
          cardPayload.button_text ?? ''
        )
      : null
    return new AiMessage(item.sender_type, item.content, item.created_at, actionCard, item.safety_level ?? 'normal')
  })
  return new AiSession(payload.id, payload.user_id, payload.scene_type, payload.status, payload.last_message_at, messages)
}

function buildHomeDashboard(payload: any): AiHomeDashboard {
  const quickActions = (payload.quick_actions as any[]).map(
    (item) => new AiHomeQuickAction(item.key, item.label, item.icon, item.route, item.prompt)
  )

  return new AiHomeDashboard(
    payload.greeting_period,
    payload.greeting_text,
    payload.greeting_sub_text,
    payload.mood_emoji,
    payload.mood_text,
    payload.weather_text,
    payload.companion_days,
    payload.online_count,
    payload.waiting_count,
    payload.remaining,
    quickActions
  )
}

function buildQuickActions(): AiHomeQuickAction[] {
  return [
    new AiHomeQuickAction('quick-join', '快速加入', 'join', '/chat', '想快速加入一个轻松的聊天陪伴场景'),
    new AiHomeQuickAction('square', '热门广场', 'square', '/home', '想看看大家最近都在聊什么'),
    new AiHomeQuickAction('voice', '治愈声音', 'voice', '/sound', '我想听一点能让我放松下来的声音'),
    new AiHomeQuickAction('topic', '今日话题', 'topic', '/home', '给我一个今晚适合开启聊天的话题')
  ]
}

function buildAssistantReply(content: string): string {
  if (isSensitiveContent(content)) {
    return '你提到的内容可能和自我伤害风险有关，这很重要。请先保证你身边环境安全，并尽快联系身边可信任的人；如果风险正在发生，请立即联系当地紧急援助电话。'
  }
  if (content.includes('压力') || content.includes('累')) {
    return '听起来你已经撑了很久。\n\n现在不用急着把事情想清楚，我们先把那股压着你的劲慢慢放下来，好吗？'
  }
  if (content.includes('睡')) {
    return '睡不好的时候，白天那些没来得及安放的情绪，常常会在夜里一起涌上来。\n\n要不要先和我说说，最近最让你停不下来的那件事？'
  }
  if (content.includes('说说话') || content.includes('孤单')) {
    return '那种突然很想找个人说话的时刻，往往不是脆弱，而是你已经很努力了。\n\n我在，你可以慢慢说。'
  }
  return '我听见你了。\n\n不用急着组织语言，想到哪儿说到哪儿就好，我会一直接着你。'
}

function buildActionCard(content: string): AiActionCard | null {
  if (isSensitiveContent(content)) {
    return new AiActionCard('go_sound', '先做 3 分钟稳定情绪', '进入声音页，先让呼吸慢下来，再决定下一步。', '/sound', '去声音')
  }
  if (content.includes('睡不着') || content.includes('失眠') || content.includes('放松') || content.includes('声音')) {
    return new AiActionCard('go_sound', '试试声音放松', '白噪音和呼吸引导可以先帮你稳住状态。', '/sound', '去声音')
  }
  if (content.includes('下单') || content.includes('服务') || content.includes('咨询')) {
    return new AiActionCard('go_service', '看看可选服务', '可以先浏览服务方，再决定是否发起邀约。', '/services', '去服务')
  }
  if (content.includes('找人') || content.includes('匹配') || content.includes('推荐') || content.includes('搭子') || content.includes('陪聊') || content.includes('孤单')) {
    return new AiActionCard('go_match', '为你匹配陪伴对象', '回到首页快速匹配，先看最适合你的 3 位推荐。', '/home', '去匹配')
  }
  return null
}

function buildSafetyLevel(content: string): string {
  return isSensitiveContent(content) ? 'high' : 'normal'
}

function isSensitiveContent(content: string): boolean {
  return ['自杀', '轻生', '不想活', '活不下去', '结束生命', '伤害自己', '自残'].some((item) =>
    content.includes(item)
  )
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}
