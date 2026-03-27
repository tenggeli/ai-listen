package mysql

import (
	"context"
	"database/sql"

	domain "listen/backend/internal/domain/ai"
)

type MatchQuotaRepository struct {
	db *sql.DB
}

func NewMatchQuotaRepository(db *sql.DB) MatchQuotaRepository {
	return MatchQuotaRepository{db: db}
}

func (r MatchQuotaRepository) GetByUserAndDate(ctx context.Context, userID string, date string) (domain.DailyQuota, error) {
	quota := domain.NewDailyQuota(userID, date)
	const query = `SELECT used_count FROM ai_daily_quotas WHERE user_id = ? AND quota_date = ? LIMIT 1`
	var used int
	err := r.db.QueryRowContext(ctx, query, userID, date).Scan(&used)
	if err == sql.ErrNoRows {
		return quota, nil
	}
	if err != nil {
		return domain.DailyQuota{}, err
	}
	quota.Used = used
	return quota, nil
}

func (r MatchQuotaRepository) Save(ctx context.Context, quota domain.DailyQuota) error {
	const upsert = `
INSERT INTO ai_daily_quotas(user_id, quota_date, used_count, created_at, updated_at)
VALUES(?, ?, ?, NOW(), NOW())
ON DUPLICATE KEY UPDATE used_count = VALUES(used_count), updated_at = NOW()`
	_, err := r.db.ExecContext(ctx, upsert, quota.UserID, quota.Date, quota.Used)
	return err
}
