package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	adminOrderApp "listen/backend/internal/application/admin_order"
	aiApp "listen/backend/internal/application/ai"
	userFeedbackApp "listen/backend/internal/application/feedback"
	orderApp "listen/backend/internal/application/order"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestOrderRoutes_FullFlow(t *testing.T) {
	repo := memory.NewOrderRepository()
	clock := aiApp.SystemClock{}
	idGenerator := aiApp.NewTimestampIDGenerator(clock)

	controller := NewOrderController(
		orderApp.NewCreateOrderUseCase(repo, idGenerator, clock),
		orderApp.NewListOrdersUseCase(repo),
		orderApp.NewGetOrderUseCase(repo),
		orderApp.NewPayOrderMockSuccessUseCase(repo, clock),
	)

	mux := http.NewServeMux()
	RegisterOrderRoutes(mux, controller)

	createBody := map[string]any{
		"provider_id":        "p_pub_001",
		"provider_name":      "暖心倾听师 · 小林",
		"service_item_id":    "si_001",
		"service_item_title": "线上倾听 30 分钟",
		"amount":             99,
		"currency":           "CNY",
	}
	rawBody, _ := json.Marshal(createBody)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(rawBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-ID", "u_1")
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusOK {
		t.Fatalf("unexpected create status: %d", createRec.Code)
	}

	var createEnvelope Envelope
	if err := json.Unmarshal(createRec.Body.Bytes(), &createEnvelope); err != nil {
		t.Fatalf("decode create envelope failed: %v", err)
	}
	createData, _ := json.Marshal(createEnvelope.Data)
	var created OrderResponseDTO
	if err := json.Unmarshal(createData, &created); err != nil {
		t.Fatalf("decode create data failed: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected order id")
	}

	payReq := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+created.ID+"/pay/mock-success", nil)
	payReq.Header.Set("X-User-ID", "u_1")
	payRec := httptest.NewRecorder()
	mux.ServeHTTP(payRec, payReq)
	if payRec.Code != http.StatusOK {
		t.Fatalf("unexpected pay status: %d", payRec.Code)
	}

	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+created.ID, nil)
	detailReq.Header.Set("X-User-ID", "u_1")
	detailRec := httptest.NewRecorder()
	mux.ServeHTTP(detailRec, detailReq)
	if detailRec.Code != http.StatusOK {
		t.Fatalf("unexpected detail status: %d", detailRec.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders?page=1&page_size=10", nil)
	listReq.Header.Set("X-User-ID", "u_1")
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("unexpected list status: %d", listRec.Code)
	}
}

func TestOrderRoutes_Forbidden(t *testing.T) {
	repo := memory.NewOrderRepository()
	clock := aiApp.SystemClock{}
	idGenerator := aiApp.NewTimestampIDGenerator(clock)

	controller := NewOrderController(
		orderApp.NewCreateOrderUseCase(repo, idGenerator, clock),
		orderApp.NewListOrdersUseCase(repo),
		orderApp.NewGetOrderUseCase(repo),
		orderApp.NewPayOrderMockSuccessUseCase(repo, clock),
	)

	mux := http.NewServeMux()
	RegisterOrderRoutes(mux, controller)

	createBody := map[string]any{
		"provider_id":        "p_pub_001",
		"provider_name":      "provider",
		"service_item_id":    "si_001",
		"service_item_title": "service",
		"amount":             99,
		"currency":           "CNY",
	}
	rawBody, _ := json.Marshal(createBody)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(rawBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-ID", "u_1")
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusOK {
		t.Fatalf("unexpected create status: %d", createRec.Code)
	}

	var createEnvelope Envelope
	if err := json.Unmarshal(createRec.Body.Bytes(), &createEnvelope); err != nil {
		t.Fatalf("decode create envelope failed: %v", err)
	}
	createData, _ := json.Marshal(createEnvelope.Data)
	var created OrderResponseDTO
	if err := json.Unmarshal(createData, &created); err != nil {
		t.Fatalf("decode create data failed: %v", err)
	}

	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+created.ID, nil)
	detailReq.Header.Set("X-User-ID", "u_2")
	detailRec := httptest.NewRecorder()
	mux.ServeHTTP(detailRec, detailReq)
	if detailRec.Code != http.StatusForbidden {
		t.Fatalf("unexpected forbidden status: %d", detailRec.Code)
	}
}

func TestOrderRoutes_AfterSaleBackflowVisibleToUserDetail(t *testing.T) {
	orderRepo := memory.NewOrderRepository()
	feedbackRepo := memory.NewFeedbackRepository()
	actionRepo := memory.NewOrderAdminActionRepository()
	clock := aiApp.SystemClock{}
	idGenerator := aiApp.NewTimestampIDGenerator(clock)

	controller := NewOrderController(
		orderApp.NewCreateOrderUseCase(orderRepo, idGenerator, clock),
		orderApp.NewListOrdersUseCase(orderRepo),
		orderApp.NewGetOrderUseCase(orderRepo),
		orderApp.NewPayOrderMockSuccessUseCase(orderRepo, clock),
	)
	adminUC := adminOrderApp.NewUseCase(orderRepo, feedbackRepo, actionRepo, idGenerator, clock)
	feedbackUC := userFeedbackApp.NewSubmitOrderFeedbackUseCase(feedbackRepo, orderRepo, idGenerator, clock)

	mux := http.NewServeMux()
	RegisterOrderRoutes(mux, controller)

	createBody := map[string]any{
		"provider_id":        "p_pub_001",
		"provider_name":      "暖心倾听师 · 小林",
		"service_item_id":    "si_001",
		"service_item_title": "线上倾听 30 分钟",
		"amount":             99,
		"currency":           "CNY",
	}
	rawBody, _ := json.Marshal(createBody)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(rawBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-ID", "u_1")
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusOK {
		t.Fatalf("unexpected create status: %d", createRec.Code)
	}
	var createEnvelope Envelope
	if err := json.Unmarshal(createRec.Body.Bytes(), &createEnvelope); err != nil {
		t.Fatalf("decode create envelope failed: %v", err)
	}
	createData, _ := json.Marshal(createEnvelope.Data)
	var created OrderResponseDTO
	if err := json.Unmarshal(createData, &created); err != nil {
		t.Fatalf("decode create data failed: %v", err)
	}

	payReq := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+created.ID+"/pay/mock-success", nil)
	payReq.Header.Set("X-User-ID", "u_1")
	payRec := httptest.NewRecorder()
	mux.ServeHTTP(payRec, payReq)
	if payRec.Code != http.StatusOK {
		t.Fatalf("unexpected pay status: %d", payRec.Code)
	}

	if _, err := feedbackUC.Execute(t.Context(), userFeedbackApp.SubmitFeedbackInput{
		UserID:           "u_1",
		OrderID:          created.ID,
		RatingScore:      8,
		ReviewTags:       []string{"投诉"},
		ReviewContent:    "需处理",
		ComplaintReason:  "服务未按约开始",
		ComplaintContent: "请平台介入",
	}); err != nil {
		t.Fatalf("submit complaint failed: %v", err)
	}

	if _, err := adminUC.ActionOrder(t.Context(), adminOrderApp.ActionOrderInput{
		OrderID:  created.ID,
		Action:   "close",
		Operator: "admin_001",
		Reason:   "投诉处理完成",
	}); err != nil {
		t.Fatalf("admin close failed: %v", err)
	}

	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+created.ID, nil)
	detailReq.Header.Set("X-User-ID", "u_1")
	detailRec := httptest.NewRecorder()
	mux.ServeHTTP(detailRec, detailReq)
	if detailRec.Code != http.StatusOK {
		t.Fatalf("unexpected detail status: %d", detailRec.Code)
	}
	var detailEnvelope Envelope
	if err := json.Unmarshal(detailRec.Body.Bytes(), &detailEnvelope); err != nil {
		t.Fatalf("decode detail envelope failed: %v", err)
	}
	detailData, _ := json.Marshal(detailEnvelope.Data)
	var detail OrderResponseDTO
	if err := json.Unmarshal(detailData, &detail); err != nil {
		t.Fatalf("decode detail data failed: %v", err)
	}
	if detail.Status != "closed" {
		t.Fatalf("expected closed status, got %s", detail.Status)
	}
	if detail.StatusReason != "已关闭" {
		t.Fatalf("expected closed reason, got %s", detail.StatusReason)
	}
	if detail.StatusActionReason != "投诉处理完成" {
		t.Fatalf("expected action reason backflow, got %s", detail.StatusActionReason)
	}
	if detail.StatusUpdatedAt == nil {
		t.Fatalf("expected status updated at")
	}
	if _, err := time.Parse(time.RFC3339, *detail.StatusUpdatedAt); err != nil {
		t.Fatalf("invalid status updated at: %v", err)
	}
}
