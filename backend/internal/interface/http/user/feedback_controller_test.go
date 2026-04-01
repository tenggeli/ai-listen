package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	aiApp "listen/backend/internal/application/ai"
	feedbackApp "listen/backend/internal/application/feedback"
	orderApp "listen/backend/internal/application/order"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestFeedbackRoutes_SubmitAndGet(t *testing.T) {
	orderRepo := memory.NewOrderRepository()
	feedbackRepo := memory.NewFeedbackRepository()
	clock := aiApp.SystemClock{}
	idGenerator := aiApp.NewTimestampIDGenerator(clock)

	orderController := NewOrderController(
		orderApp.NewCreateOrderUseCase(orderRepo, idGenerator, clock),
		orderApp.NewListOrdersUseCase(orderRepo),
		orderApp.NewGetOrderUseCase(orderRepo),
		orderApp.NewPayOrderMockSuccessUseCase(orderRepo, clock),
	)
	feedbackController := NewFeedbackController(
		feedbackApp.NewSubmitOrderFeedbackUseCase(feedbackRepo, orderRepo, idGenerator, clock),
		feedbackApp.NewGetOrderFeedbackUseCase(feedbackRepo, orderRepo),
	)

	mux := http.NewServeMux()
	RegisterOrderRoutes(mux, orderController)
	RegisterFeedbackRoutes(mux, feedbackController)

	createRaw, _ := json.Marshal(map[string]any{
		"provider_id":        "p_001",
		"provider_name":      "provider-a",
		"service_item_id":    "si_001",
		"service_item_title": "service-a",
		"amount":             99,
		"currency":           "CNY",
	})
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(createRaw))
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
	createDataRaw, _ := json.Marshal(createEnvelope.Data)
	var orderData OrderResponseDTO
	if err := json.Unmarshal(createDataRaw, &orderData); err != nil {
		t.Fatalf("decode order data failed: %v", err)
	}

	payReq := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderData.ID+"/pay/mock-success", nil)
	payReq.Header.Set("X-User-ID", "u_1")
	payRec := httptest.NewRecorder()
	mux.ServeHTTP(payRec, payReq)
	if payRec.Code != http.StatusOK {
		t.Fatalf("unexpected pay status: %d", payRec.Code)
	}

	submitRaw, _ := json.Marshal(map[string]any{
		"rating_score":      9,
		"review_tags":       []string{"很会接情绪", "不尴尬"},
		"review_content":    "服务体验不错",
		"complaint_reason":  "",
		"complaint_content": "",
	})
	submitReq := httptest.NewRequest(http.MethodPost, "/api/v1/orders/"+orderData.ID+"/feedback", bytes.NewReader(submitRaw))
	submitReq.Header.Set("X-User-ID", "u_1")
	submitReq.Header.Set("Content-Type", "application/json")
	submitRec := httptest.NewRecorder()
	mux.ServeHTTP(submitRec, submitReq)
	if submitRec.Code != http.StatusOK {
		t.Fatalf("unexpected submit status: %d", submitRec.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+orderData.ID+"/feedback", nil)
	getReq.Header.Set("X-User-ID", "u_1")
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("unexpected get feedback status: %d", getRec.Code)
	}
}
