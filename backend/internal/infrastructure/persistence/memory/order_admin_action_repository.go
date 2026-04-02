package memory

import (
	"context"
	"sort"
	"sync"

	app "listen/backend/internal/application/admin_order"
)

type OrderAdminActionRepository struct {
	mu      sync.RWMutex
	items   []app.ActionAudit
	byOrder map[string][]int
}

func NewOrderAdminActionRepository() *OrderAdminActionRepository {
	return &OrderAdminActionRepository{
		items:   make([]app.ActionAudit, 0),
		byOrder: make(map[string][]int),
	}
}

func (r *OrderAdminActionRepository) Create(_ context.Context, item app.ActionAudit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cloned := item
	r.items = append(r.items, cloned)
	idx := len(r.items) - 1
	r.byOrder[item.OrderID] = append(r.byOrder[item.OrderID], idx)
	return nil
}

func (r *OrderAdminActionRepository) ListByOrderID(_ context.Context, orderID string) ([]app.ActionAudit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	indices := r.byOrder[orderID]
	result := make([]app.ActionAudit, 0, len(indices))
	for _, idx := range indices {
		if idx < 0 || idx >= len(r.items) {
			continue
		}
		result = append(result, r.items[idx])
	}
	// newest first for console readability
	sort.Slice(result, func(i, j int) bool {
		return result[i].UpdatedAt.After(result[j].UpdatedAt)
	})
	return result, nil
}

var _ app.ActionLogRepository = (*OrderAdminActionRepository)(nil)
