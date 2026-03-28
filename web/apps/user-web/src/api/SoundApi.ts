import { SoundCategory } from '../domain/ai/SoundCategory'
import { SoundPageAggregate } from '../domain/ai/SoundPageAggregate'
import { SoundTrack } from '../domain/ai/SoundTrack'

export interface SoundApi {
  getSoundPage(userId: string): Promise<SoundPageAggregate>
}

export class MockSoundApi implements SoundApi {
  async getSoundPage(_userId: string): Promise<SoundPageAggregate> {
    await sleep(220)

    const categories = [
      new SoundCategory('all', '全部'),
      new SoundCategory('nature', '自然白噪音'),
      new SoundCategory('sleep', '睡眠引导'),
      new SoundCategory('meditation', '正念冥想'),
      new SoundCategory('story', '治愈故事'),
      new SoundCategory('breath', '呼吸练习')
    ]

    const tracks = [
      new SoundTrack('track-rain', '深夜雨声 · 安眠版', '自然白噪音', '1,248 次播放', '35:00', '🌧', 'listen 治愈声音库'),
      new SoundTrack('track-wave', '海浪呼吸引导 · 4-7-8 法', '呼吸练习', '856 次播放', '12:00', '🌊', 'listen 治愈声音库'),
      new SoundTrack('track-forest', '森林清晨 · 入睡前冥想', '正念冥想', '2,034 次播放', '20:30', '🍃', 'listen 治愈声音库'),
      new SoundTrack('track-radio', '今晚陪你说话 · 深夜电台 Vol.12', '治愈故事', '3,671 次播放', '28:15', '🌙', 'listen 深夜电台'),
      new SoundTrack('track-fire', '壁炉暖意 · 冬夜陪伴', '自然白噪音', '987 次播放', '60:00', '🔥', 'listen 治愈声音库'),
      new SoundTrack('track-star', '星空冥想 · 放下今天的重量', '正念冥想', '4,122 次播放', '18:00', '⭐', 'listen 冥想室')
    ]

    return new SoundPageAggregate(
      '声音',
      '用声音抚慰此刻的你',
      categories,
      tracks,
      'track-rain',
      '12:34',
      '35:00',
      true
    )
  }
}

export class HttpSoundApi implements SoundApi {
  constructor(private readonly baseUrl = '/api/v1') {}

  async getSoundPage(userId: string): Promise<SoundPageAggregate> {
    const response = await fetch(`${this.baseUrl}/sounds?page=home&user_id=${encodeURIComponent(userId)}`)
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load sound page failed')
    }
    return buildSoundPageAggregate(payload.data)
  }
}

function buildSoundPageAggregate(payload: any): SoundPageAggregate {
  const categories = (payload.categories as any[]).map((item) => new SoundCategory(item.key, item.label))
  const tracks = (payload.recommended_tracks as any[]).map(
    (item) =>
      new SoundTrack(
        item.id,
        item.title,
        item.category,
        item.play_count_text,
        item.duration_text,
        item.emoji,
        item.author
      )
  )

  return new SoundPageAggregate(
    payload.title,
    payload.subtitle,
    categories,
    tracks,
    payload.current_track_id,
    payload.current_progress_text,
    payload.total_duration_text,
    payload.is_playing
  )
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}
