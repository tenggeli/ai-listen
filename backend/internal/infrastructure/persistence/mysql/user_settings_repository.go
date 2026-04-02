package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	domain "listen/backend/internal/domain/user_settings"
)

type UserSettingsRepository struct {
	db *sql.DB
}

func NewUserSettingsRepository(db *sql.DB) UserSettingsRepository {
	return UserSettingsRepository{db: db}
}

func (r UserSettingsRepository) GetByUserID(ctx context.Context, userID string) (domain.Settings, bool, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return domain.Settings{}, false, domain.ErrInvalidInput
	}

	const query = `
SELECT user_id, preference_json, notification_json, privacy_json
FROM user_settings
WHERE user_id = ? LIMIT 1`

	var settings domain.Settings
	var preferenceJSON string
	var notificationJSON string
	var privacyJSON string

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&settings.UserID, &preferenceJSON, &notificationJSON, &privacyJSON)
	if err == sql.ErrNoRows {
		return domain.Settings{}, false, nil
	}
	if err != nil {
		return domain.Settings{}, false, err
	}

	if err := json.Unmarshal([]byte(preferenceJSON), &settings.Preference); err != nil {
		return domain.Settings{}, false, err
	}
	if err := json.Unmarshal([]byte(notificationJSON), &settings.Notification); err != nil {
		return domain.Settings{}, false, err
	}
	if err := json.Unmarshal([]byte(privacyJSON), &settings.Privacy); err != nil {
		return domain.Settings{}, false, err
	}
	return settings, true, nil
}

func (r UserSettingsRepository) Save(ctx context.Context, settings domain.Settings) error {
	if strings.TrimSpace(settings.UserID) == "" {
		return domain.ErrInvalidInput
	}

	preferenceJSON, err := json.Marshal(settings.Preference)
	if err != nil {
		return err
	}
	notificationJSON, err := json.Marshal(settings.Notification)
	if err != nil {
		return err
	}
	privacyJSON, err := json.Marshal(settings.Privacy)
	if err != nil {
		return err
	}

	const upsert = `
INSERT INTO user_settings(
  user_id, preference_json, notification_json, privacy_json, created_at, updated_at
) VALUES(
  ?, ?, ?, ?, NOW(), NOW()
)
ON DUPLICATE KEY UPDATE
  preference_json = VALUES(preference_json),
  notification_json = VALUES(notification_json),
  privacy_json = VALUES(privacy_json),
  updated_at = NOW()`

	_, err = r.db.ExecContext(ctx, upsert, settings.UserID, string(preferenceJSON), string(notificationJSON), string(privacyJSON))
	return err
}

var _ domain.Repository = UserSettingsRepository{}
