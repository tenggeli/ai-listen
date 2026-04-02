package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	adminOrderApp "listen/backend/internal/application/admin_order"
	aiApp "listen/backend/internal/application/ai"
	feedbackApp "listen/backend/internal/application/feedback"
	orderApp "listen/backend/internal/application/order"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

func TestOrderAndComplaintRoutes(t *testing.T) {
	orderRepo := memory.NewOrderRepository()
	feedbackRepo := memory.NewFeedbackRepository()
	actionRepo := memory.NewOrderAdminActionRepository()
	idGenerator := aiApp.NewTimestampIDGenerator(fixedClock{now: time.Date(2026, 4, 2, 10, 0, 0, 0, time.UTC)})
	clock := fixedClock{now: time.Date(2026, 4, 2, 10, 0, 0, 0, time.UTC)}

	createOrderUC := orderApp.NewCreateOrderUseCase(orderRepo, idGenerator, clock)
	payOrderUC := orderApp.NewPayOrderMockSuccessUseCase(orderRepo, clock)
	submitFeedbackUC := feedbackApp.NewSubmitOrderFeedbackUseCase(feedbackRepo, orderRepo, idGenerator, clock)

	created, err := createOrderUC.Execute(context.Background(), orderApp.CreateOrderInput{
		UserID:           "u_001",
		ProviderID:       "p_pub_001",
		ProviderName:     "provider",
		ServiceItemID:    "si_001",
		ServiceItemTitle: "service",
		Amount:           99,
		Currency:         "CNY",
	})
	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}
	if _, err = payOrderUC.Execute(context.Background(), orderApp.PayOrderMockSuccessInput{UserID: "u_001", OrderID: created.Order.ID}); err != nil {
		t.Fatalf("pay order failed: %v", err)
	}
	if _, err = submitFeedbackUC.Execute(context.Background(), feedbackApp.SubmitFeedbackInput{
		UserID:          "u_001",
		OrderID:         created.Order.ID,
		ComplaintReason: "服务未按约开始",
	}); err != nil {
		t.Fatalf("submit feedback failed: %v", err)
	}

	controller := NewOrderController(adminOrderApp.NewUseCase(orderRepo, feedbackRepo, actionRepo, idGenerator, clock))
	mux := http.NewServeMux()
	RegisterOrderRoutes(mux, controller)
	RegisterComplaintRoutes(mux, controller)

	authHeader := "Bearer mock_admin_at_admin_001_1712013723"

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/orders?page=1&page_size=10", nil)
	listReq.Header.Set("Authorization", authHeader)
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("unexpected order list status: %d", listRec.Code)
	}

	interveneReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/orders/"+created.Order.ID+"/intervene", nil)
	interveneReq.Header.Set("Authorization", authHeader)
	interveneRec := httptest.NewRecorder()
	mux.ServeHTTP(interveneRec, interveneReq)
	if interveneRec.Code != http.StatusOK {
		t.Fatalf("unexpected intervene status: %d", interveneRec.Code)
	}

	complaintListReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/complaints?page=1&page_size=10", nil)
	complaintListReq.Header.Set("Authorization", authHeader)
	complaintListRec := httptest.NewRecorder()
	mux.ServeHTTP(complaintListRec, complaintListReq)
	if complaintListRec.Code != http.StatusOK {
		t.Fatalf("unexpected complaint list status: %d", complaintListRec.Code)
	}

	resolveReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/complaints/"+created.Order.ID+"/resolve", nil)
	resolveReq.Header.Set("Authorization", authHeader)
	resolveRec := httptest.NewRecorder()
	mux.ServeHTTP(resolveRec, resolveReq)
	if resolveRec.Code != http.StatusOK {
		t.Fatalf("unexpected resolve status: %d", resolveRec.Code)
	}

	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/complaints/"+created.Order.ID, nil)
	detailReq.Header.Set("Authorization", authHeader)
	detailRec := httptest.NewRecorder()
	mux.ServeHTTP(detailRec, detailReq)
	if detailRec.Code != http.StatusOK {
		t.Fatalf("unexpected complaint detail status: %d", detailRec.Code)
	}
	var envelope Envelope
	if err := json.Unmarshal(detailRec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode detail envelope failed: %v", err)
	}
	data, ok := envelope.Data.(map[string]any)
	if !ok {
		t.Fatalf("unexpected data type")
	}
	logs, ok := data["action_logs"].([]any)
	if !ok || len(logs) == 0 {
		t.Fatalf("expected action logs")
	}
}
