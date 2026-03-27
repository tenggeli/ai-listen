import { MatchCandidate } from './MatchCandidate'

export class AiHomeAggregate {
  constructor(
    public readonly remainingCount: number,
    public readonly candidates: MatchCandidate[]
  ) {}

  hasCandidates(): boolean {
    return this.candidates.length > 0
  }
}
