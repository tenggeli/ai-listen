import { ProviderReviewStatus } from './ProviderReviewStatus'

export class ProviderSummary {
  constructor(
    public readonly id: string,
    public readonly displayName: string,
    public readonly cityCode: string,
    public readonly reviewStatus: ProviderReviewStatus
  ) {}
}
