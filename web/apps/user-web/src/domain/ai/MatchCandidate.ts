export class MatchCandidate {
  constructor(
    public readonly providerId: string,
    public readonly displayName: string,
    public readonly reason: string,
    public readonly score: number
  ) {}
}
