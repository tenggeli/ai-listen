package service_discovery

import "context"

type PublicProviderQuery struct {
	CategoryID string
	Keyword    string
	CityCode   string
	Page       int
	PageSize   int
}

type Repository interface {
	ListCategories(ctx context.Context) ([]ServiceCategory, error)
	ListPublicProviders(ctx context.Context, query PublicProviderQuery) ([]ProviderPublicProfile, int, error)
	GetPublicProviderByID(ctx context.Context, providerID string) (ProviderPublicProfile, error)
	ListProviderServiceItems(ctx context.Context, providerID string) ([]ServiceItem, error)
}
