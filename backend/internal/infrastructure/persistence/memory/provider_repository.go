package memory

import (
	"context"
	"errors"
	"sort"
	"sync"

	domain "listen/backend/internal/domain/provider"
)

type ProviderRepository struct {
	mu   sync.RWMutex
	data map[string]domain.Provider
}

func NewProviderRepository() *ProviderRepository {
	seed := map[string]domain.Provider{
		"p_001": {ID: "p_001", DisplayName: "暖心倾听师-小林", CityCode: "310100", Bio: "擅长情绪陪伴", ReviewStatus: domain.ReviewStatusSubmitted},
		"p_002": {ID: "p_002", DisplayName: "夜谈伙伴-阿泽", CityCode: "110100", Bio: "擅长夜间聊天", ReviewStatus: domain.ReviewStatusUnderReview},
		"p_003": {ID: "p_003", DisplayName: "电影散步搭子-念念", CityCode: "440100", Bio: "擅长线下陪伴", ReviewStatus: domain.ReviewStatusSupplementRequired},
	}
	return &ProviderRepository{data: seed}
}

func (r *ProviderRepository) List(_ context.Context, query domain.Query) ([]domain.Provider, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.Provider, 0, len(r.data))
	for _, item := range r.data {
		if query.ReviewStatus != "" && item.ReviewStatus != query.ReviewStatus {
			continue
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })

	total := len(items)
	start := (query.Page - 1) * query.PageSize
	if start >= total {
		return []domain.Provider{}, total, nil
	}
	end := start + query.PageSize
	if end > total {
		end = total
	}
	return items[start:end], total, nil
}

func (r *ProviderRepository) GetByID(_ context.Context, providerID string) (domain.Provider, error) {
	r.mu.RLock()
	item, ok := r.data[providerID]
	r.mu.RUnlock()
	if !ok {
		return domain.Provider{}, errors.New("provider not found")
	}
	return item, nil
}

func (r *ProviderRepository) Save(_ context.Context, provider domain.Provider, _ string, _ string, _ string) error {
	r.mu.Lock()
	r.data[provider.ID] = provider
	r.mu.Unlock()
	return nil
}
