export class SoundTrack {
  constructor(
    public readonly id: string,
    public readonly title: string,
    public readonly category: string,
    public readonly playCountText: string,
    public readonly durationText: string,
    public readonly emoji: string,
    public readonly author: string
  ) {}
}
