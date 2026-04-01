package feedback

import "context"

type Repository interface {
	GetByOrderID(ctx context.Context, orderID string) (OrderFeedback, error)
	Create(ctx context.Context, item OrderFeedback) error
}
