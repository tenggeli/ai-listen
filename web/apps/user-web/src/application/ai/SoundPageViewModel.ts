import { reactive } from 'vue'
import type { SoundApi } from '../../api/SoundApi'
import type { SoundPageAggregate } from '../../domain/ai/SoundPageAggregate'
import { PageLoadState } from '../../domain/ai/PageLoadState'
import { SoundTrack } from '../../domain/ai/SoundTrack'

export interface SoundPageState {
  pageState: PageLoadState
  selectedCategoryKey: string
  aggregate: SoundPageAggregate | null
  currentTrack: SoundTrack | null
  filteredTracks: SoundTrack[]
  isPlaying: boolean
  errorMessage: string
}

export class SoundPageViewModel {
  readonly state: SoundPageState = reactive({
    pageState: PageLoadState.Idle,
    selectedCategoryKey: 'all',
    aggregate: null,
    currentTrack: null,
    filteredTracks: [],
    isPlaying: true,
    errorMessage: ''
  })

  constructor(
    private readonly api: SoundApi,
    private readonly userId: string
  ) {}

  async initialize(): Promise<void> {
    this.state.pageState = PageLoadState.Loading
    this.state.errorMessage = ''

    try {
      const aggregate = await this.api.getSoundPage(this.userId)
      this.state.aggregate = aggregate
      this.state.selectedCategoryKey = 'all'
      this.state.currentTrack = aggregate.currentTrack()
      this.state.filteredTracks = aggregate.recommendedTracks
      this.state.isPlaying = aggregate.isPlaying
      this.state.pageState = aggregate.hasTracks() ? PageLoadState.Success : PageLoadState.Empty
    } catch (error) {
      this.state.pageState = PageLoadState.Error
      this.state.errorMessage = error instanceof Error ? error.message : '加载声音页失败'
    }
  }

  selectCategory(categoryKey: string): void {
    this.state.selectedCategoryKey = categoryKey
    const aggregate = this.state.aggregate
    if (!aggregate) {
      this.state.filteredTracks = []
      return
    }

    if (categoryKey === 'all') {
      this.state.filteredTracks = aggregate.recommendedTracks
      return
    }

    this.state.filteredTracks = aggregate.recommendedTracks.filter((track) => mapTrackCategoryToKey(track.category) === categoryKey)
  }

  selectTrack(trackId: string): void {
    const aggregate = this.state.aggregate
    if (!aggregate) {
      return
    }

    const nextTrack = aggregate.recommendedTracks.find((track) => track.id === trackId) ?? null
    this.state.currentTrack = nextTrack
    this.state.isPlaying = true
  }

  togglePlayback(): void {
    this.state.isPlaying = !this.state.isPlaying
  }
}

function mapTrackCategoryToKey(category: string): string {
  if (category.includes('白噪音')) {
    return 'nature'
  }
  if (category.includes('睡眠')) {
    return 'sleep'
  }
  if (category.includes('冥想')) {
    return 'meditation'
  }
  if (category.includes('故事') || category.includes('电台')) {
    return 'story'
  }
  if (category.includes('呼吸')) {
    return 'breath'
  }
  return 'all'
}
