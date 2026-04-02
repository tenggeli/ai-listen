package user

import (
	"encoding/json"
	"errors"
	"net/http"

	settingsApp "listen/backend/internal/application/user_settings"
	identityDomain "listen/backend/internal/domain/identity"
	settingsDomain "listen/backend/internal/domain/user_settings"
)

type SettingsController struct {
	getUC  settingsApp.GetSettingsUseCase
	saveUC settingsApp.SaveSettingsUseCase
}

func NewSettingsController(getUC settingsApp.GetSettingsUseCase, saveUC settingsApp.SaveSettingsUseCase) SettingsController {
	return SettingsController{getUC: getUC, saveUC: saveUC}
}

func (c SettingsController) HandleGetSettings(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)
	if userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	output, err := c.getUC.Execute(r.Context(), settingsApp.GetSettingsInput{UserID: userID})
	if err != nil {
		writeSettingsError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildSettingsResponse(output.Settings))
}

func (c SettingsController) HandleSaveSettings(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)
	if userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body SaveSettingsRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.saveUC.Execute(r.Context(), settingsApp.SaveSettingsInput{
		UserID: userID,
		Preference: settingsDomain.Preference{
			PreferSameCityProviders: body.Preference.PreferSameCityProviders,
			AutoPlaySoundPreview:    body.Preference.AutoPlaySoundPreview,
			HideOfflineProviders:    body.Preference.HideOfflineProviders,
		},
		Notification: settingsDomain.Notification{
			OrderStatusUpdate:     body.Notification.OrderStatusUpdate,
			ComplaintResultNotice: body.Notification.ComplaintResultNotice,
			MarketingActivity:     body.Notification.MarketingActivity,
		},
		Privacy: settingsDomain.Privacy{
			ProfilePublicVisible:       body.Privacy.ProfilePublicVisible,
			PersonalizedRecommendation: body.Privacy.PersonalizedRecommendation,
			RiskControlDataSharing:     body.Privacy.RiskControlDataSharing,
		},
	})
	if err != nil {
		writeSettingsError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, buildSettingsResponse(output.Settings))
}

func buildSettingsResponse(settings settingsDomain.Settings) SettingsResponseDTO {
	return SettingsResponseDTO{
		Preference: SettingsPreferenceDTO{
			PreferSameCityProviders: settings.Preference.PreferSameCityProviders,
			AutoPlaySoundPreview:    settings.Preference.AutoPlaySoundPreview,
			HideOfflineProviders:    settings.Preference.HideOfflineProviders,
		},
		Notification: SettingsNotificationDTO{
			OrderStatusUpdate:     settings.Notification.OrderStatusUpdate,
			ComplaintResultNotice: settings.Notification.ComplaintResultNotice,
			MarketingActivity:     settings.Notification.MarketingActivity,
		},
		Privacy: SettingsPrivacyDTO{
			ProfilePublicVisible:       settings.Privacy.ProfilePublicVisible,
			PersonalizedRecommendation: settings.Privacy.PersonalizedRecommendation,
			RiskControlDataSharing:     settings.Privacy.RiskControlDataSharing,
		},
	}
}

func writeSettingsError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, settingsDomain.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, settingsDomain.ErrSettingsPersistenceUnavailable):
		writeJSONError(w, http.StatusNotImplemented, err.Error())
	case errors.Is(err, identityDomain.ErrUserNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
