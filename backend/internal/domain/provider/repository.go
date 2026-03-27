package provider

import "context"

type Query struct {
	ReviewStatus string
	Page         int
	PageSize     int
}

type Repository interface {
	List(ctx context.Context, query Query) ([]Provider, int, error)
	GetByID(ctx context.Context, providerID string) (Provider, error)
	Save(ctx context.Context, provider Provider, operator string, action string, reason string) error
}
