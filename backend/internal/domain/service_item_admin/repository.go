package service_item_admin

import "context"

type Repository interface {
	List(ctx context.Context, query Query) ([]ServiceItem, int, error)
	GetByID(ctx context.Context, serviceItemID string) (ServiceItem, error)
	UpdateStatus(ctx context.Context, serviceItemID string, status string, operator string) error
}
