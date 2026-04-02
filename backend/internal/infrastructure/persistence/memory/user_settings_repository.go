package memory

import (
	"context"

	domain "listen/backend/internal/domain/user_settings"
)

type UserSettingsRepository struct{}

func NewUserSettingsRepository() UserSettingsRepository {
	return UserSettingsRepository{}
}

func (r UserSettingsRepository) GetByUserID(_ context.Context, _ string) (domain.Settings, bool, error) {
	return domain.Settings{}, false, domain.ErrSettingsPersistenceUnavailable
}

func (r UserSettingsRepository) Save(_ context.Context, _ domain.Settings) error {
	return domain.ErrSettingsPersistenceUnavailable
}

var _ domain.Repository = UserSettingsRepository{}
