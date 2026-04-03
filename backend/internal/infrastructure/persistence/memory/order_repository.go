package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	domain "listen/backend/internal/domain/order"
)

type OrderRepository struct {
	mu     sync.RWMutex
	byID   map[string]domain.Order
	byUser map[string][]string
	byProv map[string][]string
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		byID:   make(map[string]domain.Order),
		byUser: make(map[string][]string),
		byProv: make(map[string][]string),
	}
}

func (r *OrderRepository) Create(_ context.Context, order domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[order.ID] = cloneOrder(order)
	r.byUser[order.UserID] = append([]string{order.ID}, r.byUser[order.UserID]...)
	r.byProv[order.ProviderID] = append([]string{order.ID}, r.byProv[order.ProviderID]...)
	return nil
}

func (r *OrderRepository) GetByID(_ context.Context, orderID string) (domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.byID[orderID]
	if !ok {
		return domain.Order{}, domain.ErrOrderNotFound
	}
	return cloneOrder(item), nil
}

func (r *OrderRepository) ListByUser(_ context.Context, query domain.ListQuery) ([]domain.Order, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byUser[query.UserID]
	all := make([]domain.Order, 0, len(ids))
	for _, orderID := range ids {
		item, ok := r.byID[orderID]
		if !ok {
			continue
		}
		all = append(all, cloneOrder(item))
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	total := len(all)
	start := (query.Page - 1) * query.PageSize
	if start >= total {
		return []domain.Order{}, total, nil
	}
	end := start + query.PageSize
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

func (r *OrderRepository) Save(_ context.Context, order domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[order.ID]; !ok {
		return domain.ErrOrderNotFound
	}
	r.byID[order.ID] = cloneOrder(order)
	return nil
}

func (r *OrderRepository) ListByProvider(_ context.Context, query domain.ProviderListQuery) ([]domain.Order, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byProv[query.ProviderID]
	all := make([]domain.Order, 0, len(ids))
	for _, orderID := range ids {
		item, ok := r.byID[orderID]
		if !ok {
			continue
		}
		all = append(all, cloneOrder(item))
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	total := len(all)
	start := (query.Page - 1) * query.PageSize
	if start >= total {
		return []domain.Order{}, total, nil
	}
	end := start + query.PageSize
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

func (r *OrderRepository) ListAll(_ context.Context, query domain.AdminListQuery) ([]domain.Order, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	all := make([]domain.Order, 0, len(r.byID))
	keyword := strings.ToLower(strings.TrimSpace(query.Keyword))
	for _, item := range r.byID {
		if query.Status != "" && item.Status != query.Status {
			continue
		}
		if keyword != "" {
			target := strings.ToLower(item.ProviderName + " " + item.ServiceItemTitle + " " + item.UserID + " " + item.ID)
			if !strings.Contains(target, keyword) {
				continue
			}
		}
		all = append(all, cloneOrder(item))
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	total := len(all)
	start := (query.Page - 1) * query.PageSize
	if start >= total {
		return []domain.Order{}, total, nil
	}
	end := start + query.PageSize
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

func cloneOrder(item domain.Order) domain.Order {
	cloned := item
	if item.PaidAt != nil {
		t := *item.PaidAt
		cloned.PaidAt = &t
	}
	if item.StatusUpdatedAt != nil {
		t := *item.StatusUpdatedAt
		cloned.StatusUpdatedAt = &t
	}
	return cloned
}

var _ domain.Repository = (*OrderRepository)(nil)
