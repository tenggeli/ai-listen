package service_item_admin

import (
	"context"
	"strings"

	domain "listen/backend/internal/domain/service_item_admin"
)

type ListServiceItemsInput struct {
	ProviderID string
	CategoryID string
	Status     string
	Keyword    string
	Page       int
	PageSize   int
}

type ListServiceItemsOutput struct {
	Items []domain.ServiceItem
	Total int
}

type GetServiceItemDetailInput struct {
	ServiceItemID string
}

type GetServiceItemDetailOutput struct {
	Item domain.ServiceItem
}

type UpdateServiceItemStatusInput struct {
	ServiceItemID string
	Action        string
	Operator      string
}

type UpdateServiceItemStatusOutput struct {
	Item domain.ServiceItem
}

type ListServiceItemsUseCase struct {
	repo domain.Repository
}

func NewListServiceItemsUseCase(repo domain.Repository) ListServiceItemsUseCase {
	return ListServiceItemsUseCase{repo: repo}
}

func (u ListServiceItemsUseCase) Execute(ctx context.Context, input ListServiceItemsInput) (ListServiceItemsOutput, error) {
	page := input.Page
	if page <= 0 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}

	status, err := domain.NormalizeStatus(input.Status)
	if err != nil {
		return ListServiceItemsOutput{}, err
	}

	items, total, err := u.repo.List(ctx, domain.Query{
		ProviderID: strings.TrimSpace(input.ProviderID),
		CategoryID: strings.TrimSpace(input.CategoryID),
		Status:     status,
		Keyword:    strings.TrimSpace(input.Keyword),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		return ListServiceItemsOutput{}, err
	}
	return ListServiceItemsOutput{Items: items, Total: total}, nil
}

type GetServiceItemDetailUseCase struct {
	repo domain.Repository
}

func NewGetServiceItemDetailUseCase(repo domain.Repository) GetServiceItemDetailUseCase {
	return GetServiceItemDetailUseCase{repo: repo}
}

func (u GetServiceItemDetailUseCase) Execute(ctx context.Context, input GetServiceItemDetailInput) (GetServiceItemDetailOutput, error) {
	serviceItemID := strings.TrimSpace(input.ServiceItemID)
	if serviceItemID == "" {
		return GetServiceItemDetailOutput{}, domain.ErrInvalidInput
	}
	item, err := u.repo.GetByID(ctx, serviceItemID)
	if err != nil {
		return GetServiceItemDetailOutput{}, err
	}
	return GetServiceItemDetailOutput{Item: item}, nil
}

type UpdateServiceItemStatusUseCase struct {
	repo domain.Repository
}

func NewUpdateServiceItemStatusUseCase(repo domain.Repository) UpdateServiceItemStatusUseCase {
	return UpdateServiceItemStatusUseCase{repo: repo}
}

func (u UpdateServiceItemStatusUseCase) Execute(ctx context.Context, input UpdateServiceItemStatusInput) (UpdateServiceItemStatusOutput, error) {
	serviceItemID := strings.TrimSpace(input.ServiceItemID)
	if serviceItemID == "" {
		return UpdateServiceItemStatusOutput{}, domain.ErrInvalidInput
	}
	nextStatus, err := statusByAction(input.Action)
	if err != nil {
		return UpdateServiceItemStatusOutput{}, err
	}
	if err := u.repo.UpdateStatus(ctx, serviceItemID, nextStatus, strings.TrimSpace(input.Operator)); err != nil {
		return UpdateServiceItemStatusOutput{}, err
	}
	item, err := u.repo.GetByID(ctx, serviceItemID)
	if err != nil {
		return UpdateServiceItemStatusOutput{}, err
	}
	return UpdateServiceItemStatusOutput{Item: item}, nil
}

func statusByAction(action string) (string, error) {
	switch strings.TrimSpace(strings.ToLower(action)) {
	case "activate":
		return domain.StatusActive, nil
	case "deactivate":
		return domain.StatusInactive, nil
	default:
		return "", domain.ErrInvalidInput
	}
}
