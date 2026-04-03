package provider

import (
	"context"
	"errors"
	"strings"

	domain "listen/backend/internal/domain/provider"
	providerAuthDomain "listen/backend/internal/domain/provider_auth"
	serviceDiscoveryDomain "listen/backend/internal/domain/service_discovery"
)

var ErrProviderNotFound = errors.New("provider not found")
var ErrInvalidInput = providerAuthDomain.ErrInvalidInput
var ErrUnauthorized = providerAuthDomain.ErrUnauthorized

type ListReviewProvidersInput struct {
	ReviewStatus string
	Page         int
	PageSize     int
}

type ListReviewProvidersOutput struct {
	Items []domain.Provider
	Total int
}

type GetProviderDetailInput struct {
	ProviderID string
}

type GetProviderDetailOutput struct {
	Provider domain.Provider
}

type ReviewActionInput struct {
	ProviderID string
	Action     string
	Operator   string
	Reason     string
}

type ReviewActionOutput struct {
	Provider domain.Provider
}

type ListReviewProvidersUseCase struct {
	repo domain.Repository
}

func NewListReviewProvidersUseCase(repo domain.Repository) ListReviewProvidersUseCase {
	return ListReviewProvidersUseCase{repo: repo}
}

func (u ListReviewProvidersUseCase) Execute(ctx context.Context, input ListReviewProvidersInput) (ListReviewProvidersOutput, error) {
	page := input.Page
	if page <= 0 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}

	items, total, err := u.repo.List(ctx, domain.Query{
		ReviewStatus: input.ReviewStatus,
		Page:         page,
		PageSize:     pageSize,
	})
	if err != nil {
		return ListReviewProvidersOutput{}, err
	}
	return ListReviewProvidersOutput{Items: items, Total: total}, nil
}

type GetProviderDetailUseCase struct {
	repo domain.Repository
}

func NewGetProviderDetailUseCase(repo domain.Repository) GetProviderDetailUseCase {
	return GetProviderDetailUseCase{repo: repo}
}

func (u GetProviderDetailUseCase) Execute(ctx context.Context, input GetProviderDetailInput) (GetProviderDetailOutput, error) {
	provider, err := u.repo.GetByID(ctx, input.ProviderID)
	if err != nil {
		return GetProviderDetailOutput{}, err
	}
	return GetProviderDetailOutput{Provider: provider}, nil
}

type ReviewProviderUseCase struct {
	repo domain.Repository
}

func NewReviewProviderUseCase(repo domain.Repository) ReviewProviderUseCase {
	return ReviewProviderUseCase{repo: repo}
}

func (u ReviewProviderUseCase) Execute(ctx context.Context, input ReviewActionInput) (ReviewActionOutput, error) {
	provider, err := u.repo.GetByID(ctx, input.ProviderID)
	if err != nil {
		return ReviewActionOutput{}, err
	}

	switch input.Action {
	case "approve":
		err = provider.Approve()
	case "reject":
		err = provider.Reject()
	case "require_supplement":
		err = provider.RequireSupplement()
	default:
		err = domain.ErrInvalidProviderTransition
	}
	if err != nil {
		return ReviewActionOutput{}, err
	}

	if err := u.repo.Save(ctx, provider, input.Operator, input.Action, input.Reason); err != nil {
		return ReviewActionOutput{}, err
	}
	return ReviewActionOutput{Provider: provider}, nil
}

type UpdateCurrentProfileInput struct {
	ProviderID  string
	DisplayName string
	CityCode    string
}

type UpdateCurrentProfileOutput struct {
	Provider providerAuthDomain.ProviderProfile
}

type UpdateCurrentProfileRepository interface {
	UpdateProfile(ctx context.Context, providerID string, displayName string, cityCode string) (providerAuthDomain.ProviderAccount, bool, error)
}

type UpdateCurrentProfileUseCase struct {
	repo UpdateCurrentProfileRepository
}

func NewUpdateCurrentProfileUseCase(repo UpdateCurrentProfileRepository) UpdateCurrentProfileUseCase {
	return UpdateCurrentProfileUseCase{repo: repo}
}

func (u UpdateCurrentProfileUseCase) Execute(ctx context.Context, input UpdateCurrentProfileInput) (UpdateCurrentProfileOutput, error) {
	providerID := strings.TrimSpace(input.ProviderID)
	displayName := strings.TrimSpace(input.DisplayName)
	cityCode := strings.TrimSpace(input.CityCode)

	if providerID == "" {
		return UpdateCurrentProfileOutput{}, providerAuthDomain.ErrUnauthorized
	}
	if displayName == "" {
		return UpdateCurrentProfileOutput{}, providerAuthDomain.ErrInvalidInput
	}

	account, found, err := u.repo.UpdateProfile(ctx, providerID, displayName, cityCode)
	if err != nil {
		return UpdateCurrentProfileOutput{}, err
	}
	if !found {
		return UpdateCurrentProfileOutput{}, providerAuthDomain.ErrProviderNotFound
	}
	return UpdateCurrentProfileOutput{
		Provider: providerAuthDomain.ProviderProfile{
			ProviderID:  account.ProviderID,
			Account:     account.Account,
			DisplayName: account.DisplayName,
			Status:      account.Status,
			CityCode:    account.CityCode,
		},
	}, nil
}

type ListCurrentProviderServicesInput struct {
	ProviderID string
}

type ListCurrentProviderServicesOutput struct {
	Items []serviceDiscoveryDomain.ServiceItem
}

type ListCurrentProviderServicesUseCase struct {
	repo serviceDiscoveryDomain.Repository
}

func NewListCurrentProviderServicesUseCase(repo serviceDiscoveryDomain.Repository) ListCurrentProviderServicesUseCase {
	return ListCurrentProviderServicesUseCase{repo: repo}
}

func (u ListCurrentProviderServicesUseCase) Execute(ctx context.Context, input ListCurrentProviderServicesInput) (ListCurrentProviderServicesOutput, error) {
	providerID := strings.TrimSpace(input.ProviderID)
	if providerID == "" {
		return ListCurrentProviderServicesOutput{}, providerAuthDomain.ErrUnauthorized
	}
	items, err := u.repo.ListProviderServiceItems(ctx, providerID)
	if err != nil {
		return ListCurrentProviderServicesOutput{}, err
	}
	return ListCurrentProviderServicesOutput{Items: items}, nil
}
