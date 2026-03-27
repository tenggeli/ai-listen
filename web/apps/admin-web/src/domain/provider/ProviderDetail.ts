import { ProviderReviewStatus } from './ProviderReviewStatus'

export class ProviderDetail {
  constructor(
    public readonly id: string,
    public readonly displayName: string,
    public readonly cityCode: string,
    public readonly bio: string,
    public readonly reviewStatus: ProviderReviewStatus
  ) {}
}
