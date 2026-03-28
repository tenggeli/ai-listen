export class UserIdentity {
  constructor(
    public readonly userId: string,
    public readonly loginChannel: string,
    public readonly accessToken: string,
    public readonly refreshToken: string,
    public readonly expiresInSeconds: number,
    public readonly displayName: string,
    public readonly avatarUrl: string,
    public readonly isNewUser: boolean,
    public readonly profileCompleted: boolean
  ) {}
}

