import type { UserSettings } from '../domain/settings/UserSettings'

export class SettingsApiError extends Error {
  constructor(message: string, readonly statusCode: number) {
    super(message)
  }
}

export interface SettingsApi {
  getMySettings(accessToken: string): Promise<UserSettings>
  saveMySettings(accessToken: string, settings: UserSettings): Promise<UserSettings>
}

export class HttpSettingsApi implements SettingsApi {
  constructor(private readonly baseUrl = '/api/v1') {}

  async getMySettings(accessToken: string): Promise<UserSettings> {
    const response = await fetch(`${this.baseUrl}/users/me/settings`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new SettingsApiError(payload.message || 'get settings failed', response.status)
    }
    return mapSettings(payload.data)
  }

  async saveMySettings(accessToken: string, settings: UserSettings): Promise<UserSettings> {
    const response = await fetch(`${this.baseUrl}/users/me/settings`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        preference: {
          prefer_same_city_providers: settings.preference.preferSameCityProviders,
          auto_play_sound_preview: settings.preference.autoPlaySoundPreview,
          hide_offline_providers: settings.preference.hideOfflineProviders
        },
        notification: {
          order_status_update: settings.notification.orderStatusUpdate,
          complaint_result_notice: settings.notification.complaintResultNotice,
          marketing_activity: settings.notification.marketingActivity
        },
        privacy: {
          profile_public_visible: settings.privacy.profilePublicVisible,
          personalized_recommendation: settings.privacy.personalizedRecommendation,
          risk_control_data_sharing: settings.privacy.riskControlDataSharing
        }
      })
    })
    const payload = await response.json()
    if (!response.ok || payload.code !== 0) {
      throw new SettingsApiError(payload.message || 'save settings failed', response.status)
    }
    return mapSettings(payload.data)
  }
}

function mapSettings(data: any): UserSettings {
  return {
    preference: {
      preferSameCityProviders: Boolean(data.preference?.prefer_same_city_providers),
      autoPlaySoundPreview: Boolean(data.preference?.auto_play_sound_preview),
      hideOfflineProviders: Boolean(data.preference?.hide_offline_providers)
    },
    notification: {
      orderStatusUpdate: Boolean(data.notification?.order_status_update),
      complaintResultNotice: Boolean(data.notification?.complaint_result_notice),
      marketingActivity: Boolean(data.notification?.marketing_activity)
    },
    privacy: {
      profilePublicVisible: Boolean(data.privacy?.profile_public_visible),
      personalizedRecommendation: Boolean(data.privacy?.personalized_recommendation),
      riskControlDataSharing: Boolean(data.privacy?.risk_control_data_sharing)
    }
  }
}

