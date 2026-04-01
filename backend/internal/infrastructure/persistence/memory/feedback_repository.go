package memory

import (
	"context"
	"sync"

	domain "listen/backend/internal/domain/feedback"
)

type FeedbackRepository struct {
	mu      sync.RWMutex
	byID    map[string]domain.OrderFeedback
	byOrder map[string]string
}

func NewFeedbackRepository() *FeedbackRepository {
	return &FeedbackRepository{
		byID:    make(map[string]domain.OrderFeedback),
		byOrder: make(map[string]string),
	}
}

func (r *FeedbackRepository) GetByOrderID(_ context.Context, orderID string) (domain.OrderFeedback, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	feedbackID, ok := r.byOrder[orderID]
	if !ok {
		return domain.OrderFeedback{}, domain.ErrFeedbackNotFound
	}
	item, ok := r.byID[feedbackID]
	if !ok {
		return domain.OrderFeedback{}, domain.ErrFeedbackNotFound
	}
	return item, nil
}

func (r *FeedbackRepository) Create(_ context.Context, item domain.OrderFeedback) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byOrder[item.OrderID]; ok {
		return domain.ErrFeedbackSubmitted
	}
	r.byID[item.ID] = item
	r.byOrder[item.OrderID] = item.ID
	return nil
}

var _ domain.Repository = (*FeedbackRepository)(nil)
