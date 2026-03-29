export class ProviderPublicProfile {
  constructor(
    public readonly id: string,
    public readonly displayName: string,
    public readonly avatarUrl: string,
    public readonly cityCode: string,
    public readonly bio: string,
    public readonly ratingAvg: number,
    public readonly completedOrders: number,
    public readonly online: boolean,
    public readonly verificationLabel: string,
    public readonly tags: string[],
    public readonly priceFrom: number,
    public readonly priceUnit: string
  ) {}
}
