package memory

import (
	"context"
	"sort"
	"sync"

	domain "listen/backend/internal/domain/order"
)

type OrderRepository struct {
	mu     sync.RWMutex
	byID   map[string]domain.Order
	byUser map[string][]string
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		byID:   make(map[string]domain.Order),
		byUser: make(map[string][]string),
	}
}

func (r *OrderRepository) Create(_ context.Context, order domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[order.ID] = cloneOrder(order)
	r.byUser[order.UserID] = append([]string{order.ID}, r.byUser[order.UserID]...)
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

func cloneOrder(item domain.Order) domain.Order {
	cloned := item
	if item.PaidAt != nil {
		t := *item.PaidAt
		cloned.PaidAt = &t
	}
	return cloned
}

var _ domain.Repository = (*OrderRepository)(nil)
