import { ApiError } from './ApiError'

export interface ProviderLoginResult {
  accessToken: string
  providerId: string
}

export interface ProviderMeResult {
  providerId: string
  account: string
  displayName: string
  status: string
  cityCode: string
}

export interface SaveProviderProfileInput {
  displayName: string
  cityCode: string
}

export interface ProviderAuthApi {
  loginMock(account: string, password: string): Promise<ProviderLoginResult>
  getMe(accessToken: string): Promise<ProviderMeResult>
  saveProfile(accessToken: string, input: SaveProviderProfileInput): Promise<ProviderMeResult>
}

export class HttpProviderAuthApi implements ProviderAuthApi {
  constructor(private readonly baseUrl = '/api/v1/provider') {}

  async loginMock(account: string, password: string): Promise<ProviderLoginResult> {
    const response = await fetch(`${this.baseUrl}/auth/login/mock`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ account, password })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new ApiError(payload.message || 'login failed', response.status)
    }
    return {
      accessToken: payload.data.access_token,
      providerId: payload.data.provider_id
    }
  }

  async getMe(accessToken: string): Promise<ProviderMeResult> {
    const response = await fetch(`${this.baseUrl}/profile`, {
      headers: {
        Authorization: `Bearer ${accessToken}`
      }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new ApiError(payload.message || 'load profile failed', response.status)
    }
    return {
      providerId: payload.data.provider_id,
      account: payload.data.account,
      displayName: payload.data.display_name,
      status: payload.data.status,
      cityCode: payload.data.city_code
    }
  }

  async saveProfile(accessToken: string, input: SaveProviderProfileInput): Promise<ProviderMeResult> {
    const response = await fetch(`${this.baseUrl}/profile`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        display_name: input.displayName,
        city_code: input.cityCode
      })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new ApiError(payload.message || 'save profile failed', response.status)
    }
    return {
      providerId: payload.data.provider_id,
      account: payload.data.account,
      displayName: payload.data.display_name,
      status: payload.data.status,
      cityCode: payload.data.city_code
    }
  }
}
