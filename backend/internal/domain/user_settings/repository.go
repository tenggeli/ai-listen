package user_settings

import "context"

type Repository interface {
	GetByUserID(ctx context.Context, userID string) (Settings, bool, error)
	Save(ctx context.Context, settings Settings) error
}
