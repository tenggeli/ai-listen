import { SoundCategory } from './SoundCategory'
import { SoundTrack } from './SoundTrack'

export class SoundPageAggregate {
  constructor(
    public readonly title: string,
    public readonly subtitle: string,
    public readonly categories: SoundCategory[],
    public readonly recommendedTracks: SoundTrack[],
    public readonly currentTrackId: string,
    public readonly currentProgressText: string,
    public readonly totalDurationText: string,
    public readonly isPlaying: boolean
  ) {}

  hasTracks(): boolean {
    return this.recommendedTracks.length > 0
  }

  currentTrack(): SoundTrack | null {
    return this.recommendedTracks.find((track) => track.id === this.currentTrackId) ?? null
  }
}
