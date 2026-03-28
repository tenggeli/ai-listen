import { AiHomeQuickAction } from './AiHomeQuickAction'

export class AiHomeDashboard {
  constructor(
    public readonly greetingPeriod: string,
    public readonly greetingText: string,
    public readonly greetingSubText: string,
    public readonly moodEmoji: string,
    public readonly moodText: string,
    public readonly weatherText: string,
    public readonly companionDays: number,
    public readonly onlineCount: number,
    public readonly waitingCount: number,
    public readonly remainingCount: number,
    public readonly quickActions: AiHomeQuickAction[]
  ) {}

  hasContent(): boolean {
    return this.greetingText.length > 0 || this.quickActions.length > 0
  }

  withRemainingCount(remainingCount: number): AiHomeDashboard {
    return new AiHomeDashboard(
      this.greetingPeriod,
      this.greetingText,
      this.greetingSubText,
      this.moodEmoji,
      this.moodText,
      this.weatherText,
      this.companionDays,
      this.onlineCount,
      this.waitingCount,
      remainingCount,
      this.quickActions
    )
  }
}
