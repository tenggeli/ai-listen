package order

import (
	"context"
	"strings"
	"time"

	domain "listen/backend/internal/domain/order"
)

type Clock interface {
	Now() time.Time
}

type IDGenerator interface {
	NewID(prefix string) string
}

type CreateOrderInput struct {
	UserID           string
	ProviderID       string
	ProviderName     string
	ServiceItemID    string
	ServiceItemTitle string
	Amount           int
	Currency         string
}

type CreateOrderOutput struct {
	Order domain.Order
}

type ListOrdersInput struct {
	UserID   string
	Page     int
	PageSize int
}

type ListOrdersOutput struct {
	Items []domain.Order
	Total int
}

type GetOrderInput struct {
	UserID  string
	OrderID string
}

type GetOrderOutput struct {
	Order domain.Order
}

type PayOrderMockSuccessInput struct {
	UserID  string
	OrderID string
}

type PayOrderMockSuccessOutput struct {
	Order domain.Order
}

type CreateOrderUseCase struct {
	repo        domain.Repository
	idGenerator IDGenerator
	clock       Clock
}

func NewCreateOrderUseCase(repo domain.Repository, idGenerator IDGenerator, clock Clock) CreateOrderUseCase {
	return CreateOrderUseCase{repo: repo, idGenerator: idGenerator, clock: clock}
}

func (u CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (CreateOrderOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return CreateOrderOutput{}, domain.ErrInvalidInput
	}

	order, err := domain.NewOrder(
		u.idGenerator.NewID("ord"),
		userID,
		input.ProviderID,
		input.ProviderName,
		input.ServiceItemID,
		input.ServiceItemTitle,
		input.Amount,
		input.Currency,
		u.clock.Now(),
	)
	if err != nil {
		return CreateOrderOutput{}, err
	}

	if err := u.repo.Create(ctx, order); err != nil {
		return CreateOrderOutput{}, err
	}
	return CreateOrderOutput{Order: order}, nil
}

type ListOrdersUseCase struct {
	repo domain.Repository
}

func NewListOrdersUseCase(repo domain.Repository) ListOrdersUseCase {
	return ListOrdersUseCase{repo: repo}
}

func (u ListOrdersUseCase) Execute(ctx context.Context, input ListOrdersInput) (ListOrdersOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return ListOrdersOutput{}, domain.ErrInvalidInput
	}

	page := input.Page
	if page <= 0 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}

	items, total, err := u.repo.ListByUser(ctx, domain.ListQuery{
		UserID:   userID,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return ListOrdersOutput{}, err
	}
	return ListOrdersOutput{Items: items, Total: total}, nil
}

type GetOrderUseCase struct {
	repo domain.Repository
}

func NewGetOrderUseCase(repo domain.Repository) GetOrderUseCase {
	return GetOrderUseCase{repo: repo}
}

func (u GetOrderUseCase) Execute(ctx context.Context, input GetOrderInput) (GetOrderOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	orderID := strings.TrimSpace(input.OrderID)
	if userID == "" || orderID == "" {
		return GetOrderOutput{}, domain.ErrInvalidInput
	}

	item, err := u.repo.GetByID(ctx, orderID)
	if err != nil {
		return GetOrderOutput{}, err
	}
	if item.UserID != userID {
		return GetOrderOutput{}, domain.ErrOrderForbidden
	}
	return GetOrderOutput{Order: item}, nil
}

type PayOrderMockSuccessUseCase struct {
	repo  domain.Repository
	clock Clock
}

func NewPayOrderMockSuccessUseCase(repo domain.Repository, clock Clock) PayOrderMockSuccessUseCase {
	return PayOrderMockSuccessUseCase{repo: repo, clock: clock}
}

func (u PayOrderMockSuccessUseCase) Execute(ctx context.Context, input PayOrderMockSuccessInput) (PayOrderMockSuccessOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	orderID := strings.TrimSpace(input.OrderID)
	if userID == "" || orderID == "" {
		return PayOrderMockSuccessOutput{}, domain.ErrInvalidInput
	}

	item, err := u.repo.GetByID(ctx, orderID)
	if err != nil {
		return PayOrderMockSuccessOutput{}, err
	}
	if item.UserID != userID {
		return PayOrderMockSuccessOutput{}, domain.ErrOrderForbidden
	}
	if err := item.MarkPaid(u.clock.Now()); err != nil {
		return PayOrderMockSuccessOutput{}, err
	}
	if err := u.repo.Save(ctx, item); err != nil {
		return PayOrderMockSuccessOutput{}, err
	}
	return PayOrderMockSuccessOutput{Order: item}, nil
}
