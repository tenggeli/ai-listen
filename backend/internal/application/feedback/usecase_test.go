package feedback

import (
	"context"
	"testing"
	"time"

	aiApp "listen/backend/internal/application/ai"
	orderApp "listen/backend/internal/application/order"
	feedbackDomain "listen/backend/internal/domain/feedback"
	orderDomain "listen/backend/internal/domain/order"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestSubmitOrderFeedbackUseCase_Success(t *testing.T) {
	orderRepo := memory.NewOrderRepository()
	feedbackRepo := memory.NewFeedbackRepository()
	clock := aiApp.SystemClock{}
	idGenerator := aiApp.NewTimestampIDGenerator(clock)

	createOrderUC := orderApp.NewCreateOrderUseCase(orderRepo, idGenerator, clock)
	payOrderUC := orderApp.NewPayOrderMockSuccessUseCase(orderRepo, clock)
	submitUC := NewSubmitOrderFeedbackUseCase(feedbackRepo, orderRepo, idGenerator, clock)
	getUC := NewGetOrderFeedbackUseCase(feedbackRepo, orderRepo)

	createdOrder, err := createOrderUC.Execute(context.Background(), orderApp.CreateOrderInput{
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
	if _, err := payOrderUC.Execute(context.Background(), orderApp.PayOrderMockSuccessInput{
		UserID:  "u_001",
		OrderID: createdOrder.Order.ID,
	}); err != nil {
		t.Fatalf("pay order failed: %v", err)
	}

	output, err := submitUC.Execute(context.Background(), SubmitFeedbackInput{
		UserID:           "u_001",
		OrderID:          createdOrder.Order.ID,
		RatingScore:      9,
		ReviewTags:       []string{"很会接情绪", "不尴尬"},
		ReviewContent:    "体验很好，沟通自然",
		ComplaintReason:  "",
		ComplaintContent: "",
	})
	if err != nil {
		t.Fatalf("submit feedback failed: %v", err)
	}
	if output.Item.ID == "" {
		t.Fatal("expected feedback id")
	}

	got, err := getUC.Execute(context.Background(), GetFeedbackInput{
		UserID:  "u_001",
		OrderID: createdOrder.Order.ID,
	})
	if err != nil {
		t.Fatalf("get feedback failed: %v", err)
	}
	if got.Item.RatingScore != 9 {
		t.Fatalf("unexpected rating score: %d", got.Item.RatingScore)
	}
}

func TestSubmitOrderFeedbackUseCase_RejectDuplicate(t *testing.T) {
	orderRepo := memory.NewOrderRepository()
	feedbackRepo := memory.NewFeedbackRepository()
	idGenerator := fixedIDGenerator{id: "fixed_id"}
	clock := fixedClock{now: time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)}

	createOrderUC := orderApp.NewCreateOrderUseCase(orderRepo, idGenerator, clock)
	payOrderUC := orderApp.NewPayOrderMockSuccessUseCase(orderRepo, clock)
	submitUC := NewSubmitOrderFeedbackUseCase(feedbackRepo, orderRepo, idGenerator, clock)

	createdOrder, err := createOrderUC.Execute(context.Background(), orderApp.CreateOrderInput{
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
	if _, err := payOrderUC.Execute(context.Background(), orderApp.PayOrderMockSuccessInput{
		UserID:  "u_001",
		OrderID: createdOrder.Order.ID,
	}); err != nil {
		t.Fatalf("pay order failed: %v", err)
	}

	_, err = submitUC.Execute(context.Background(), SubmitFeedbackInput{
		UserID:          "u_001",
		OrderID:         createdOrder.Order.ID,
		RatingScore:     10,
		ReviewTags:      []string{"稳定"},
		ReviewContent:   "good",
		ComplaintReason: "",
	})
	if err != nil {
		t.Fatalf("first submit failed: %v", err)
	}

	_, err = submitUC.Execute(context.Background(), SubmitFeedbackInput{
		UserID:          "u_001",
		OrderID:         createdOrder.Order.ID,
		RatingScore:     8,
		ReviewTags:      []string{"再次提交"},
		ReviewContent:   "should fail",
		ComplaintReason: "",
	})
	if err == nil {
		t.Fatal("expected duplicate error")
	}
	if err != feedbackDomain.ErrFeedbackSubmitted {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSubmitOrderFeedbackUseCase_ComplaintMovesOrderToAfterSale(t *testing.T) {
	orderRepo := memory.NewOrderRepository()
	feedbackRepo := memory.NewFeedbackRepository()
	idGenerator := fixedIDGenerator{id: "fixed_id"}
	clock := fixedClock{now: time.Date(2026, 4, 4, 10, 0, 0, 0, time.UTC)}

	createOrderUC := orderApp.NewCreateOrderUseCase(orderRepo, idGenerator, clock)
	payOrderUC := orderApp.NewPayOrderMockSuccessUseCase(orderRepo, clock)
	submitUC := NewSubmitOrderFeedbackUseCase(feedbackRepo, orderRepo, idGenerator, clock)

	createdOrder, err := createOrderUC.Execute(context.Background(), orderApp.CreateOrderInput{
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
	if _, err := payOrderUC.Execute(context.Background(), orderApp.PayOrderMockSuccessInput{
		UserID:  "u_001",
		OrderID: createdOrder.Order.ID,
	}); err != nil {
		t.Fatalf("pay order failed: %v", err)
	}

	if _, err := submitUC.Execute(context.Background(), SubmitFeedbackInput{
		UserID:           "u_001",
		OrderID:          createdOrder.Order.ID,
		RatingScore:      8,
		ReviewTags:       []string{"需介入"},
		ReviewContent:    "有争议",
		ComplaintReason:  "服务内容与描述不符",
		ComplaintContent: "希望平台处理",
	}); err != nil {
		t.Fatalf("submit complaint feedback failed: %v", err)
	}

	orderItem, err := orderRepo.GetByID(context.Background(), createdOrder.Order.ID)
	if err != nil {
		t.Fatalf("get order failed: %v", err)
	}
	if orderItem.Status != orderDomain.StatusAfterSale {
		t.Fatalf("expected after sale status, got %s", orderItem.Status)
	}
	if orderItem.StatusActionReason != "用户发起投诉" {
		t.Fatalf("expected complaint reason mark, got %s", orderItem.StatusActionReason)
	}
	if orderItem.StatusUpdatedAt == nil {
		t.Fatalf("expected status updated at")
	}
}

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
