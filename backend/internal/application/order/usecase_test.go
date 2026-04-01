package order

import (
	"context"
	"testing"
	"time"

	domain "listen/backend/internal/domain/order"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

type fixedIDGenerator struct {
	id string
}

func (g fixedIDGenerator) NewID(_ string) string {
	return g.id
}

func TestCreateAndPayOrderUseCase(t *testing.T) {
	repo := memory.NewOrderRepository()
	clock := fixedClock{now: time.Date(2026, 3, 30, 10, 0, 0, 0, time.UTC)}
	idGenerator := fixedIDGenerator{id: "ord_001"}

	createUC := NewCreateOrderUseCase(repo, idGenerator, clock)
	listUC := NewListOrdersUseCase(repo)
	getUC := NewGetOrderUseCase(repo)
	payUC := NewPayOrderMockSuccessUseCase(repo, fixedClock{
		now: time.Date(2026, 3, 30, 10, 5, 0, 0, time.UTC),
	})

	created, err := createUC.Execute(context.Background(), CreateOrderInput{
		UserID:           "u_001",
		ProviderID:       "p_001",
		ProviderName:     "provider-a",
		ServiceItemID:    "si_001",
		ServiceItemTitle: "service-a",
		Amount:           99,
		Currency:         "cny",
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}
	if created.Order.Status != domain.StatusCreated {
		t.Fatalf("unexpected order status: %s", created.Order.Status)
	}

	paid, err := payUC.Execute(context.Background(), PayOrderMockSuccessInput{
		UserID:  "u_001",
		OrderID: created.Order.ID,
	})
	if err != nil {
		t.Fatalf("pay order failed: %v", err)
	}
	if paid.Order.Status != domain.StatusPaid {
		t.Fatalf("unexpected paid status: %s", paid.Order.Status)
	}
	if paid.Order.PaidAt == nil {
		t.Fatalf("expected paid_at")
	}

	listOutput, err := listUC.Execute(context.Background(), ListOrdersInput{
		UserID:   "u_001",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("list orders failed: %v", err)
	}
	if listOutput.Total != 1 || len(listOutput.Items) != 1 {
		t.Fatalf("unexpected list result total=%d len=%d", listOutput.Total, len(listOutput.Items))
	}

	detail, err := getUC.Execute(context.Background(), GetOrderInput{
		UserID:  "u_001",
		OrderID: created.Order.ID,
	})
	if err != nil {
		t.Fatalf("get order failed: %v", err)
	}
	if detail.Order.ID != created.Order.ID {
		t.Fatalf("unexpected order id: %s", detail.Order.ID)
	}
}

func TestGetOrderUseCase_Forbidden(t *testing.T) {
	repo := memory.NewOrderRepository()
	clock := fixedClock{now: time.Date(2026, 3, 30, 10, 0, 0, 0, time.UTC)}
	idGenerator := fixedIDGenerator{id: "ord_001"}

	createUC := NewCreateOrderUseCase(repo, idGenerator, clock)
	getUC := NewGetOrderUseCase(repo)

	created, err := createUC.Execute(context.Background(), CreateOrderInput{
		UserID:           "u_001",
		ProviderID:       "p_001",
		ProviderName:     "provider-a",
		ServiceItemID:    "si_001",
		ServiceItemTitle: "service-a",
		Amount:           99,
		Currency:         "CNY",
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	_, err = getUC.Execute(context.Background(), GetOrderInput{
		UserID:  "u_002",
		OrderID: created.Order.ID,
	})
	if err == nil {
		t.Fatalf("expected forbidden error")
	}
	if err != domain.ErrOrderForbidden {
		t.Fatalf("unexpected error: %v", err)
	}
}
