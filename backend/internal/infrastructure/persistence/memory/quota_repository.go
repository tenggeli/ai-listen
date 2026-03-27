package memory

import (
	"context"
	"sync"

	domain "listen/backend/internal/domain/ai"
)

type MatchQuotaRepository struct {
	mu   sync.RWMutex
	data map[string]domain.DailyQuota
}

func NewMatchQuotaRepository() *MatchQuotaRepository {
	return &MatchQuotaRepository{data: make(map[string]domain.DailyQuota)}
}

func (r *MatchQuotaRepository) GetByUserAndDate(_ context.Context, userID string, date string) (domain.DailyQuota, error) {
	key := userID + "#" + date
	r.mu.RLock()
	quota, ok := r.data[key]
	r.mu.RUnlock()
	if ok {
		return quota, nil
	}
	return domain.NewDailyQuota(userID, date), nil
}

func (r *MatchQuotaRepository) Save(_ context.Context, quota domain.DailyQuota) error {
	key := quota.UserID + "#" + quota.Date
	r.mu.Lock()
	r.data[key] = quota
	r.mu.Unlock()
	return nil
}
