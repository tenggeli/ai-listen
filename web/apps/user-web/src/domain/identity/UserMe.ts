export class UserMe {
  constructor(
    public readonly userId: string,
    public readonly nickname: string,
    public readonly avatarUrl: string,
    public readonly gender: string,
    public readonly ageRange: string,
    public readonly city: string,
    public readonly bio: string,
    public readonly interestTags: string[],
    public readonly mbti: string,
    public readonly profileCompleted: boolean,
    public readonly personalityCompleted: boolean,
    public readonly personalitySkipped: boolean
  ) {}
}
