export class AiActionCard {
  constructor(
    public readonly action: string,
    public readonly title: string,
    public readonly description: string,
    public readonly route: string,
    public readonly buttonText: string
  ) {}
}
