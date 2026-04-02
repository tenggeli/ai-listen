package user_settings

import (
	"context"
	"strings"

	identityDomain "listen/backend/internal/domain/identity"
	domain "listen/backend/internal/domain/user_settings"
)

type GetSettingsInput struct {
	UserID string
}

type GetSettingsOutput struct {
	Settings domain.Settings
}

type SaveSettingsInput struct {
	UserID       string
	Preference   domain.Preference
	Notification domain.Notification
	Privacy      domain.Privacy
}

type SaveSettingsOutput struct {
	Settings domain.Settings
}

type GetSettingsUseCase struct {
	settingsRepo domain.Repository
	identityRepo identityDomain.Repository
}

func NewGetSettingsUseCase(settingsRepo domain.Repository, identityRepo identityDomain.Repository) GetSettingsUseCase {
	return GetSettingsUseCase{settingsRepo: settingsRepo, identityRepo: identityRepo}
}

func (u GetSettingsUseCase) Execute(ctx context.Context, input GetSettingsInput) (GetSettingsOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return GetSettingsOutput{}, domain.ErrInvalidInput
	}
	if err := ensureUserExists(ctx, u.identityRepo, userID); err != nil {
		return GetSettingsOutput{}, err
	}

	settings, found, err := u.settingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return GetSettingsOutput{}, err
	}
	if found {
		return GetSettingsOutput{Settings: settings}, nil
	}

	defaults, err := domain.NewDefault(userID)
	if err != nil {
		return GetSettingsOutput{}, err
	}
	return GetSettingsOutput{Settings: defaults}, nil
}

type SaveSettingsUseCase struct {
	settingsRepo domain.Repository
	identityRepo identityDomain.Repository
}

func NewSaveSettingsUseCase(settingsRepo domain.Repository, identityRepo identityDomain.Repository) SaveSettingsUseCase {
	return SaveSettingsUseCase{settingsRepo: settingsRepo, identityRepo: identityRepo}
}

func (u SaveSettingsUseCase) Execute(ctx context.Context, input SaveSettingsInput) (SaveSettingsOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return SaveSettingsOutput{}, domain.ErrInvalidInput
	}
	if err := ensureUserExists(ctx, u.identityRepo, userID); err != nil {
		return SaveSettingsOutput{}, err
	}

	settings := domain.Settings{
		UserID:       userID,
		Preference:   input.Preference,
		Notification: input.Notification,
		Privacy:      input.Privacy,
	}
	if err := u.settingsRepo.Save(ctx, settings); err != nil {
		return SaveSettingsOutput{}, err
	}
	return SaveSettingsOutput{Settings: settings}, nil
}

func ensureUserExists(ctx context.Context, identityRepo identityDomain.Repository, userID string) error {
	_, found, err := identityRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !found {
		return identityDomain.ErrUserNotFound
	}
	return nil
}
