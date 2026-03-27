export class AiMessage {
  constructor(
    public readonly senderType: string,
    public readonly content: string,
    public readonly createdAt: string
  ) {}
}
