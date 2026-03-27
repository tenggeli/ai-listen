package ai

import "context"

type MatchQuotaRepository interface {
	GetByUserAndDate(ctx context.Context, userID string, date string) (DailyQuota, error)
	Save(ctx context.Context, quota DailyQuota) error
}

type SessionRepository interface {
	Create(ctx context.Context, session Session) error
	GetByID(ctx context.Context, id string) (Session, error)
	Save(ctx context.Context, session Session) error
}
