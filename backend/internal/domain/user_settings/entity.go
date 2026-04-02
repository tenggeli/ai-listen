package user_settings

import (
	"errors"
	"strings"
)

var (
	ErrInvalidInput                   = errors.New("invalid input")
	ErrSettingsPersistenceUnavailable = errors.New("settings persistence only available in mysql mode")
)

type Preference struct {
	PreferSameCityProviders bool
	AutoPlaySoundPreview    bool
	HideOfflineProviders    bool
}

type Notification struct {
	OrderStatusUpdate     bool
	ComplaintResultNotice bool
	MarketingActivity     bool
}

type Privacy struct {
	ProfilePublicVisible       bool
	PersonalizedRecommendation bool
	RiskControlDataSharing     bool
}

type Settings struct {
	UserID       string
	Preference   Preference
	Notification Notification
	Privacy      Privacy
}

func NewDefault(userID string) (Settings, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return Settings{}, ErrInvalidInput
	}
	return Settings{
		UserID: userID,
		Preference: Preference{
			PreferSameCityProviders: true,
			AutoPlaySoundPreview:    true,
			HideOfflineProviders:    false,
		},
		Notification: Notification{
			OrderStatusUpdate:     true,
			ComplaintResultNotice: true,
			MarketingActivity:     false,
		},
		Privacy: Privacy{
			ProfilePublicVisible:       true,
			PersonalizedRecommendation: true,
			RiskControlDataSharing:     true,
		},
	}, nil
}
