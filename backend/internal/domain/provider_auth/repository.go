package provider_auth

import "context"

type Repository interface {
	GetByAccount(ctx context.Context, account string) (ProviderAccount, bool, error)
	GetByID(ctx context.Context, providerID string) (ProviderAccount, bool, error)
}
