package order

import (
	"context"
	"strings"
	"sync"
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

type ProviderListOrdersInput struct {
	ProviderID string
	Page       int
	PageSize   int
}

type ProviderListOrdersOutput struct {
	Items []domain.Order
	Total int
}

type ProviderOperateOrderInput struct {
	ProviderID string
	OrderID    string
	Action     string
}

type ProviderOperateOrderOutput struct {
	Order domain.Order
}

type ProviderGetOrderInput struct {
	ProviderID string
	OrderID    string
}

type ProviderGetOrderOutput struct {
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

type ProviderListOrdersUseCase struct {
	repo domain.Repository
}

func NewProviderListOrdersUseCase(repo domain.Repository) ProviderListOrdersUseCase {
	return ProviderListOrdersUseCase{repo: repo}
}

func (u ProviderListOrdersUseCase) Execute(ctx context.Context, input ProviderListOrdersInput) (ProviderListOrdersOutput, error) {
	providerID := strings.TrimSpace(input.ProviderID)
	if providerID == "" {
		return ProviderListOrdersOutput{}, domain.ErrInvalidInput
	}
	page := input.Page
	if page <= 0 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}

	items, total, err := u.repo.ListByProvider(ctx, domain.ProviderListQuery{
		ProviderID: providerID,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		return ProviderListOrdersOutput{}, err
	}
	return ProviderListOrdersOutput{Items: items, Total: total}, nil
}

type ProviderGetOrderUseCase struct {
	repo domain.Repository
}

func NewProviderGetOrderUseCase(repo domain.Repository) ProviderGetOrderUseCase {
	return ProviderGetOrderUseCase{repo: repo}
}

func (u ProviderGetOrderUseCase) Execute(ctx context.Context, input ProviderGetOrderInput) (ProviderGetOrderOutput, error) {
	providerID := strings.TrimSpace(input.ProviderID)
	orderID := strings.TrimSpace(input.OrderID)
	if providerID == "" || orderID == "" {
		return ProviderGetOrderOutput{}, domain.ErrInvalidInput
	}

	item, err := u.repo.GetByID(ctx, orderID)
	if err != nil {
		return ProviderGetOrderOutput{}, err
	}
	if item.ProviderID != providerID {
		return ProviderGetOrderOutput{}, domain.ErrOrderForbidden
	}
	return ProviderGetOrderOutput{Order: item}, nil
}

type ProviderOperateOrderUseCase struct {
	repo      domain.Repository
	orderLock *sync.Map
	clock     Clock
}

type systemClock struct{}

func (systemClock) Now() time.Time {
	return time.Now()
}

func NewProviderOperateOrderUseCase(repo domain.Repository, clocks ...Clock) ProviderOperateOrderUseCase {
	clock := Clock(systemClock{})
	if len(clocks) > 0 && clocks[0] != nil {
		clock = clocks[0]
	}
	return ProviderOperateOrderUseCase{
		repo:      repo,
		orderLock: &sync.Map{},
		clock:     clock,
	}
}

func (u ProviderOperateOrderUseCase) Execute(ctx context.Context, input ProviderOperateOrderInput) (ProviderOperateOrderOutput, error) {
	providerID := strings.TrimSpace(input.ProviderID)
	orderID := strings.TrimSpace(input.OrderID)
	action := strings.TrimSpace(input.Action)
	if providerID == "" || orderID == "" || action == "" {
		return ProviderOperateOrderOutput{}, domain.ErrInvalidInput
	}
	lock := u.getOrderLock(orderID)
	lock.Lock()
	defer lock.Unlock()

	item, err := u.repo.GetByID(ctx, orderID)
	if err != nil {
		return ProviderOperateOrderOutput{}, err
	}
	if item.ProviderID != providerID {
		return ProviderOperateOrderOutput{}, domain.ErrOrderForbidden
	}

	switch action {
	case "accept":
		err = item.MarkAccepted(u.clock.Now())
	case "depart":
		err = item.MarkOnTheWay(u.clock.Now())
	case "arrive":
		err = item.MarkArrived(u.clock.Now())
	case "start":
		err = item.MarkInService(u.clock.Now())
	case "complete":
		err = item.MarkCompleted(u.clock.Now())
	default:
		return ProviderOperateOrderOutput{}, domain.ErrInvalidInput
	}
	if err != nil {
		return ProviderOperateOrderOutput{}, err
	}
	if err := u.repo.Save(ctx, item); err != nil {
		return ProviderOperateOrderOutput{}, err
	}
	return ProviderOperateOrderOutput{Order: item}, nil
}

func (u ProviderOperateOrderUseCase) getOrderLock(orderID string) *sync.Mutex {
	if u.orderLock == nil {
		u.orderLock = &sync.Map{}
	}
	value, _ := u.orderLock.LoadOrStore(orderID, &sync.Mutex{})
	lock, ok := value.(*sync.Mutex)
	if !ok {
		return &sync.Mutex{}
	}
	return lock
}
