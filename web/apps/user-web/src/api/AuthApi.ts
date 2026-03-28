import { UserIdentity } from '../domain/identity/UserIdentity'
import { UserMe } from '../domain/identity/UserMe'

export interface AuthApi {
  loginBySms(phone: string, verifyCode: string, agreementAccepted: boolean): Promise<UserIdentity>
  loginByWechatMock(authCode: string, agreementAccepted: boolean): Promise<UserIdentity>
  getMe(accessToken: string): Promise<UserMe>
  saveProfile(
    accessToken: string,
    input: {
      nickname: string
      avatarUrl: string
      gender: string
      ageRange: string
      city: string
      bio: string
      genderChangeConfirmed: boolean
    }
  ): Promise<UserMe>
  savePersonality(
    accessToken: string,
    input: {
      mbti: string
      interestTags: string[]
    }
  ): Promise<UserMe>
  skipPersonality(accessToken: string): Promise<UserMe>
}

export class HttpAuthApi implements AuthApi {
  constructor(private readonly baseUrl = '/api/v1') {}

  async loginBySms(phone: string, verifyCode: string, agreementAccepted: boolean): Promise<UserIdentity> {
    const response = await fetch(`${this.baseUrl}/auth/login/sms`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        phone,
        verify_code: verifyCode,
        agreement_accepted: agreementAccepted
      })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'sms login failed')
    }
    return mapIdentity(payload.data)
  }

  async loginByWechatMock(authCode: string, agreementAccepted: boolean): Promise<UserIdentity> {
    const response = await fetch(`${this.baseUrl}/auth/login/wechat/mock`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        auth_code: authCode,
        agreement_accepted: agreementAccepted
      })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'wechat login failed')
    }
    return mapIdentity(payload.data)
  }

  async getMe(accessToken: string): Promise<UserMe> {
    const response = await fetch(`${this.baseUrl}/users/me`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'get me failed')
    }
    return mapUserMe(payload.data)
  }

  async saveProfile(
    accessToken: string,
    input: {
      nickname: string
      avatarUrl: string
      gender: string
      ageRange: string
      city: string
      bio: string
      genderChangeConfirmed: boolean
    }
  ): Promise<UserMe> {
    const response = await fetch(`${this.baseUrl}/users/me/profile`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${accessToken}`
      },
      body: JSON.stringify({
        nickname: input.nickname,
        avatar_url: input.avatarUrl,
        gender: input.gender,
        age_range: input.ageRange,
        city: input.city,
        bio: input.bio,
        gender_change_confirmed: input.genderChangeConfirmed
      })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'save profile failed')
    }
    return mapUserMe(payload.data)
  }

  async savePersonality(
    accessToken: string,
    input: {
      mbti: string
      interestTags: string[]
    }
  ): Promise<UserMe> {
    const response = await fetch(`${this.baseUrl}/users/me/personality`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${accessToken}`
      },
      body: JSON.stringify({
        mbti: input.mbti,
        interest_tags: input.interestTags
      })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'save personality failed')
    }
    return mapUserMe(payload.data)
  }

  async skipPersonality(accessToken: string): Promise<UserMe> {
    const response = await fetch(`${this.baseUrl}/users/me/personality/skip`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${accessToken}`
      }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'skip personality failed')
    }
    return mapUserMe(payload.data)
  }
}

function mapIdentity(payload: any): UserIdentity {
  return new UserIdentity(
    payload.user_id,
    payload.login_channel,
    payload.access_token,
    payload.refresh_token,
    payload.expires_in_seconds,
    payload.display_name,
    payload.avatar_url,
    payload.is_new_user,
    payload.profile_completed
  )
}

function mapUserMe(data: any): UserMe {
  return new UserMe(
    data.user_id,
    data.nickname ?? '',
    data.avatar_url ?? '',
    data.gender ?? '',
    data.age_range ?? '',
    data.city ?? '',
    data.bio ?? '',
    Array.isArray(data.interest_tags) ? data.interest_tags : [],
    data.mbti ?? '',
    Boolean(data.profile_completed),
    Boolean(data.personality_completed),
    Boolean(data.personality_skipped)
  )
}
