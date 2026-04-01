import type { UserSettings } from '../../domain/settings/UserSettings'
import { createDefaultUserSettings } from '../../domain/settings/UserSettings'

function storageKey(userId: string): string {
  return `listen_user_settings_${userId}`
}

export function loadUserSettings(userId: string): UserSettings {
  if (!userId) {
    return createDefaultUserSettings()
  }
  const raw = localStorage.getItem(storageKey(userId))
  if (!raw) {
    return createDefaultUserSettings()
  }
  try {
    const parsed = JSON.parse(raw) as UserSettings
    return {
      preference: {
        preferSameCityProviders: Boolean(parsed.preference?.preferSameCityProviders),
        autoPlaySoundPreview: Boolean(parsed.preference?.autoPlaySoundPreview),
        hideOfflineProviders: Boolean(parsed.preference?.hideOfflineProviders)
      },
      notification: {
        orderStatusUpdate: Boolean(parsed.notification?.orderStatusUpdate),
        complaintResultNotice: Boolean(parsed.notification?.complaintResultNotice),
        marketingActivity: Boolean(parsed.notification?.marketingActivity)
      },
      privacy: {
        profilePublicVisible: Boolean(parsed.privacy?.profilePublicVisible),
        personalizedRecommendation: Boolean(parsed.privacy?.personalizedRecommendation),
        riskControlDataSharing: Boolean(parsed.privacy?.riskControlDataSharing)
      }
    }
  } catch {
    return createDefaultUserSettings()
  }
}

export function saveUserSettings(userId: string, settings: UserSettings): void {
  if (!userId) {
    return
  }
  localStorage.setItem(storageKey(userId), JSON.stringify(settings))
}
