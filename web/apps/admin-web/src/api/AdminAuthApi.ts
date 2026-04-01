export interface AdminLoginResult {
  accessToken: string
  adminId: string
  role: string
}

export interface AdminMeResult {
  adminId: string
  account: string
  role: string
  displayName: string
  status: string
}

export interface AdminAuthApi {
  loginMock(account: string, password: string): Promise<AdminLoginResult>
  getMe(accessToken: string): Promise<AdminMeResult>
}

export class HttpAdminAuthApi implements AdminAuthApi {
  constructor(private readonly baseUrl = '/api/v1/admin/auth') {}

  async loginMock(account: string, password: string): Promise<AdminLoginResult> {
    const response = await fetch(`${this.baseUrl}/login/mock`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ account, password })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'login failed')
    }
    return {
      accessToken: payload.data.access_token,
      adminId: payload.data.admin_id,
      role: payload.data.role
    }
  }

  async getMe(accessToken: string): Promise<AdminMeResult> {
    const response = await fetch(`${this.baseUrl}/me`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new Error(payload.message || 'load current admin failed')
    }
    return {
      adminId: payload.data.admin_id,
      account: payload.data.account,
      role: payload.data.role,
      displayName: payload.data.display_name,
      status: payload.data.status
    }
  }
}
