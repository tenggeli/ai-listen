import { AiMessage } from './AiMessage'

export class AiSession {
  constructor(
    public readonly id: string,
    public readonly userId: string,
    public readonly sceneType: string,
    public readonly status: string,
    public readonly lastMessageAt: string,
    public readonly messages: AiMessage[]
  ) {}

  hasMessages(): boolean {
    return this.messages.length > 0
  }
}
