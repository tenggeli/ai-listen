export class ServiceItem {
  constructor(
    public readonly id: string,
    public readonly providerId: string,
    public readonly categoryId: string,
    public readonly title: string,
    public readonly description: string,
    public readonly priceAmount: number,
    public readonly priceUnit: string,
    public readonly supportOnline: boolean
  ) {}
}
