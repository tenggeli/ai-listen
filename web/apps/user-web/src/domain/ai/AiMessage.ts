import { AiActionCard } from './AiActionCard'

export class AiMessage {
  constructor(
    public readonly senderType: string,
    public readonly content: string,
    public readonly createdAt: string,
    public readonly actionCard: AiActionCard | null = null,
    public readonly safetyLevel: string = 'normal'
  ) {}

  isUser(): boolean {
    return this.senderType === 'user'
  }

  displayTime(): string {
    const parsed = new Date(this.createdAt)
    if (Number.isNaN(parsed.getTime())) {
      return this.createdAt
    }

    return new Intl.DateTimeFormat('zh-CN', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    }).format(parsed)
  }
}
