package feedback

import "context"

type ComplaintListQuery struct {
	Page     int
	PageSize int
}

type Repository interface {
	GetByOrderID(ctx context.Context, orderID string) (OrderFeedback, error)
	ListComplaints(ctx context.Context, query ComplaintListQuery) ([]OrderFeedback, int, error)
	Create(ctx context.Context, item OrderFeedback) error
}
