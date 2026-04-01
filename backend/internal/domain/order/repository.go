package order

import "context"

type ListQuery struct {
	UserID   string
	Page     int
	PageSize int
}

type Repository interface {
	Create(ctx context.Context, order Order) error
	GetByID(ctx context.Context, orderID string) (Order, error)
	ListByUser(ctx context.Context, query ListQuery) ([]Order, int, error)
	Save(ctx context.Context, order Order) error
}
