package user

type SettingsResponseDTO struct {
	Preference   SettingsPreferenceDTO   `json:"preference"`
	Notification SettingsNotificationDTO `json:"notification"`
	Privacy      SettingsPrivacyDTO      `json:"privacy"`
}

type SettingsPreferenceDTO struct {
	PreferSameCityProviders bool `json:"prefer_same_city_providers"`
	AutoPlaySoundPreview    bool `json:"auto_play_sound_preview"`
	HideOfflineProviders    bool `json:"hide_offline_providers"`
}

type SettingsNotificationDTO struct {
	OrderStatusUpdate     bool `json:"order_status_update"`
	ComplaintResultNotice bool `json:"complaint_result_notice"`
	MarketingActivity     bool `json:"marketing_activity"`
}

type SettingsPrivacyDTO struct {
	ProfilePublicVisible       bool `json:"profile_public_visible"`
	PersonalizedRecommendation bool `json:"personalized_recommendation"`
	RiskControlDataSharing     bool `json:"risk_control_data_sharing"`
}

type SaveSettingsRequestDTO struct {
	Preference   SettingsPreferenceDTO   `json:"preference"`
	Notification SettingsNotificationDTO `json:"notification"`
	Privacy      SettingsPrivacyDTO      `json:"privacy"`
}
