export interface PreferenceSettings {
  preferSameCityProviders: boolean
  autoPlaySoundPreview: boolean
  hideOfflineProviders: boolean
}

export interface NotificationSettings {
  orderStatusUpdate: boolean
  complaintResultNotice: boolean
  marketingActivity: boolean
}

export interface PrivacySettings {
  profilePublicVisible: boolean
  personalizedRecommendation: boolean
  riskControlDataSharing: boolean
}

export interface UserSettings {
  preference: PreferenceSettings
  notification: NotificationSettings
  privacy: PrivacySettings
}

export function createDefaultUserSettings(): UserSettings {
  return {
    preference: {
      preferSameCityProviders: true,
      autoPlaySoundPreview: true,
      hideOfflineProviders: false
    },
    notification: {
      orderStatusUpdate: true,
      complaintResultNotice: true,
      marketingActivity: false
    },
    privacy: {
      profilePublicVisible: true,
      personalizedRecommendation: true,
      riskControlDataSharing: true
    }
  }
}
