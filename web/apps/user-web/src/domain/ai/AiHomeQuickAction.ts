export class AiHomeQuickAction {
  constructor(
    public readonly key: string,
    public readonly label: string,
    public readonly icon: string,
    public readonly route: string,
    public readonly prompt: string
  ) {}
}
