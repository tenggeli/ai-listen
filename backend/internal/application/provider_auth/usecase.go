package provider_auth

import (
	"context"
	"strings"
	"time"

	domain "listen/backend/internal/domain/provider_auth"
)

type Clock interface {
	Now() time.Time
}

type LoginMockInput struct {
	Account  string
	Password string
}

type LoginMockOutput struct {
	Identity domain.ProviderIdentity
}

type LoginMockUseCase struct {
	repo  domain.Repository
	clock Clock
}

func NewLoginMockUseCase(repo domain.Repository, clock Clock) LoginMockUseCase {
	return LoginMockUseCase{repo: repo, clock: clock}
}

func (u LoginMockUseCase) Execute(ctx context.Context, input LoginMockInput) (LoginMockOutput, error) {
	account := strings.TrimSpace(input.Account)
	password := strings.TrimSpace(input.Password)
	if account == "" || password == "" {
		return LoginMockOutput{}, domain.ErrInvalidInput
	}

	provider, found, err := u.repo.GetByAccount(ctx, account)
	if err != nil {
		return LoginMockOutput{}, err
	}
	if !found || provider.Password != password {
		return LoginMockOutput{}, domain.ErrInvalidCredential
	}
	if provider.Status != "active" {
		return LoginMockOutput{}, domain.ErrUnauthorized
	}

	now := u.clock.Now()
	return LoginMockOutput{Identity: domain.ProviderIdentity{
		ProviderID:       provider.ProviderID,
		AccessToken:      domain.BuildAccessToken(provider.ProviderID, now),
		ExpiresInSeconds: 7200,
	}}, nil
}

type GetCurrentProviderInput struct {
	ProviderID string
}

type GetCurrentProviderOutput struct {
	Provider domain.ProviderProfile
}

type GetCurrentProviderUseCase struct {
	repo domain.Repository
}

func NewGetCurrentProviderUseCase(repo domain.Repository) GetCurrentProviderUseCase {
	return GetCurrentProviderUseCase{repo: repo}
}

func (u GetCurrentProviderUseCase) Execute(ctx context.Context, input GetCurrentProviderInput) (GetCurrentProviderOutput, error) {
	providerID := strings.TrimSpace(input.ProviderID)
	if providerID == "" {
		return GetCurrentProviderOutput{}, domain.ErrUnauthorized
	}

	provider, found, err := u.repo.GetByID(ctx, providerID)
	if err != nil {
		return GetCurrentProviderOutput{}, err
	}
	if !found {
		return GetCurrentProviderOutput{}, domain.ErrProviderNotFound
	}
	if provider.Status != "active" {
		return GetCurrentProviderOutput{}, domain.ErrUnauthorized
	}

	return GetCurrentProviderOutput{Provider: domain.ProviderProfile{
		ProviderID:  provider.ProviderID,
		Account:     provider.Account,
		DisplayName: provider.DisplayName,
		Status:      provider.Status,
		CityCode:    provider.CityCode,
	}}, nil
}

type InMemoryRepository struct {
	byID      map[string]domain.ProviderAccount
	byAccount map[string]string
}

func NewInMemoryRepository(accounts []domain.ProviderAccount) InMemoryRepository {
	byID := make(map[string]domain.ProviderAccount, len(accounts))
	byAccount := make(map[string]string, len(accounts))
	for _, item := range accounts {
		cleanID := strings.TrimSpace(item.ProviderID)
		cleanAccount := strings.TrimSpace(item.Account)
		if cleanID == "" || cleanAccount == "" {
			continue
		}
		byID[cleanID] = item
		byAccount[cleanAccount] = cleanID
	}
	return InMemoryRepository{byID: byID, byAccount: byAccount}
}

func (r InMemoryRepository) GetByAccount(_ context.Context, account string) (domain.ProviderAccount, bool, error) {
	providerID, ok := r.byAccount[strings.TrimSpace(account)]
	if !ok {
		return domain.ProviderAccount{}, false, nil
	}
	item, ok := r.byID[providerID]
	if !ok {
		return domain.ProviderAccount{}, false, nil
	}
	return item, true, nil
}

func (r InMemoryRepository) GetByID(_ context.Context, providerID string) (domain.ProviderAccount, bool, error) {
	item, ok := r.byID[strings.TrimSpace(providerID)]
	if !ok {
		return domain.ProviderAccount{}, false, nil
	}
	return item, true, nil
}

func (r InMemoryRepository) UpdateProfile(_ context.Context, providerID string, displayName string, cityCode string) (domain.ProviderAccount, bool, error) {
	cleanID := strings.TrimSpace(providerID)
	item, ok := r.byID[cleanID]
	if !ok {
		return domain.ProviderAccount{}, false, nil
	}
	item.DisplayName = strings.TrimSpace(displayName)
	item.CityCode = strings.TrimSpace(cityCode)
	r.byID[cleanID] = item
	return item, true, nil
}
