package service_discovery

import (
	"context"

	domain "listen/backend/internal/domain/service_discovery"
)

var ErrProviderNotFound = domain.ErrProviderNotFound

type ListServiceCategoriesOutput struct {
	Items []domain.ServiceCategory
}

type ListPublicProvidersInput struct {
	CategoryID string
	Keyword    string
	CityCode   string
	Page       int
	PageSize   int
}

type ListPublicProvidersOutput struct {
	Items []domain.ProviderPublicProfile
	Total int
}

type GetPublicProviderInput struct {
	ProviderID string
}

type GetPublicProviderOutput struct {
	Provider domain.ProviderPublicProfile
}

type ListProviderServiceItemsInput struct {
	ProviderID string
}

type ListProviderServiceItemsOutput struct {
	Items []domain.ServiceItem
}

type ListServiceCategoriesUseCase struct {
	repo domain.Repository
}

func NewListServiceCategoriesUseCase(repo domain.Repository) ListServiceCategoriesUseCase {
	return ListServiceCategoriesUseCase{repo: repo}
}

func (u ListServiceCategoriesUseCase) Execute(ctx context.Context) (ListServiceCategoriesOutput, error) {
	items, err := u.repo.ListCategories(ctx)
	if err != nil {
		return ListServiceCategoriesOutput{}, err
	}
	return ListServiceCategoriesOutput{Items: items}, nil
}

type ListPublicProvidersUseCase struct {
	repo domain.Repository
}

func NewListPublicProvidersUseCase(repo domain.Repository) ListPublicProvidersUseCase {
	return ListPublicProvidersUseCase{repo: repo}
}

func (u ListPublicProvidersUseCase) Execute(ctx context.Context, input ListPublicProvidersInput) (ListPublicProvidersOutput, error) {
	page := input.Page
	if page <= 0 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}

	items, total, err := u.repo.ListPublicProviders(ctx, domain.PublicProviderQuery{
		CategoryID: input.CategoryID,
		Keyword:    input.Keyword,
		CityCode:   input.CityCode,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		return ListPublicProvidersOutput{}, err
	}
	return ListPublicProvidersOutput{Items: items, Total: total}, nil
}

type GetPublicProviderUseCase struct {
	repo domain.Repository
}

func NewGetPublicProviderUseCase(repo domain.Repository) GetPublicProviderUseCase {
	return GetPublicProviderUseCase{repo: repo}
}

func (u GetPublicProviderUseCase) Execute(ctx context.Context, input GetPublicProviderInput) (GetPublicProviderOutput, error) {
	item, err := u.repo.GetPublicProviderByID(ctx, input.ProviderID)
	if err != nil {
		return GetPublicProviderOutput{}, err
	}
	return GetPublicProviderOutput{Provider: item}, nil
}

type ListProviderServiceItemsUseCase struct {
	repo domain.Repository
}

func NewListProviderServiceItemsUseCase(repo domain.Repository) ListProviderServiceItemsUseCase {
	return ListProviderServiceItemsUseCase{repo: repo}
}

func (u ListProviderServiceItemsUseCase) Execute(ctx context.Context, input ListProviderServiceItemsInput) (ListProviderServiceItemsOutput, error) {
	items, err := u.repo.ListProviderServiceItems(ctx, input.ProviderID)
	if err != nil {
		return ListProviderServiceItemsOutput{}, err
	}
	return ListProviderServiceItemsOutput{Items: items}, nil
}
