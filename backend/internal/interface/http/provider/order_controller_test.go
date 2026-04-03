package provider

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	adminOrderApp "listen/backend/internal/application/admin_order"
	aiApp "listen/backend/internal/application/ai"
	feedbackApp "listen/backend/internal/application/feedback"
	orderApp "listen/backend/internal/application/order"
	providerAuthApp "listen/backend/internal/application/provider_auth"
	providerAuthDomain "listen/backend/internal/domain/provider_auth"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestProviderRoutes_OrderFlow(t *testing.T) {
	orderRepo := memory.NewOrderRepository()
	clock := aiApp.SystemClock{}
	idGenerator := aiApp.NewTimestampIDGenerator(clock)

	providerRepo := providerAuthApp.NewInMemoryRepository([]providerAuthDomain.ProviderAccount{
		{ProviderID: "p_pub_001", Account: "provider", Password: "provider123", DisplayName: "provider", Status: "active", CityCode: "310000"},
	})

	userCreateOrderUC := orderApp.NewCreateOrderUseCase(orderRepo, idGenerator, clock)
	userPayOrderUC := orderApp.NewPayOrderMockSuccessUseCase(orderRepo, clock)
	providerAuthController := NewAuthController(
		providerAuthApp.NewLoginMockUseCase(providerRepo, clock),
		providerAuthApp.NewGetCurrentProviderUseCase(providerRepo),
	)
	providerOrderController := NewOrderController(
		orderApp.NewProviderListOrdersUseCase(orderRepo),
		orderApp.NewProviderGetOrderUseCase(orderRepo),
		orderApp.NewProviderOperateOrderUseCase(orderRepo),
	)

	mux := http.NewServeMux()
	RegisterAuthRoutes(mux, providerAuthController)
	RegisterOrderRoutes(mux, providerOrderController)

	created, err := userCreateOrderUC.Execute(t.Context(), orderApp.CreateOrderInput{
		UserID:           "u_1",
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
	if _, err = userPayOrderUC.Execute(t.Context(), orderApp.PayOrderMockSuccessInput{UserID: "u_1", OrderID: created.Order.ID}); err != nil {
		t.Fatalf("pay order failed: %v", err)
	}

	loginRaw, _ := json.Marshal(map[string]any{"account": "provider", "password": "provider123"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/provider/auth/login/mock", bytes.NewReader(loginRaw))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	mux.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("unexpected login status: %d", loginRec.Code)
	}
	var loginEnv Envelope
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginEnv); err != nil {
		t.Fatalf("decode login failed: %v", err)
	}
	loginDataRaw, _ := json.Marshal(loginEnv.Data)
	var loginData LoginMockResponseDTO
	if err := json.Unmarshal(loginDataRaw, &loginData); err != nil {
		t.Fatalf("decode login data failed: %v", err)
	}
	if loginData.AccessToken == "" {
		t.Fatalf("expected access token")
	}

	header := "Bearer " + loginData.AccessToken
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/provider/orders?page=1&page_size=10", nil)
	listReq.Header.Set("Authorization", header)
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("unexpected list status: %d", listRec.Code)
	}

	actions := []string{"accept", "depart", "arrive", "start", "complete"}
	for _, action := range actions {
		actionReq := httptest.NewRequest(http.MethodPost, "/api/v1/provider/orders/"+created.Order.ID+"/"+action, nil)
		actionReq.Header.Set("Authorization", header)
		actionRec := httptest.NewRecorder()
		mux.ServeHTTP(actionRec, actionReq)
		if actionRec.Code != http.StatusOK {
			t.Fatalf("unexpected %s status: %d", action, actionRec.Code)
		}
	}
}

func TestProviderRoutes_OrderActionConflict(t *testing.T) {
	orderRepo := memory.NewOrderRepository()
	clock := aiApp.SystemClock{}
	idGenerator := aiApp.NewTimestampIDGenerator(clock)

	providerRepo := providerAuthApp.NewInMemoryRepository([]providerAuthDomain.ProviderAccount{
		{ProviderID: "p_pub_001", Account: "provider", Password: "provider123", DisplayName: "provider", Status: "active", CityCode: "310000"},
	})

	userCreateOrderUC := orderApp.NewCreateOrderUseCase(orderRepo, idGenerator, clock)
	userPayOrderUC := orderApp.NewPayOrderMockSuccessUseCase(orderRepo, clock)
	providerAuthController := NewAuthController(
		providerAuthApp.NewLoginMockUseCase(providerRepo, clock),
		providerAuthApp.NewGetCurrentProviderUseCase(providerRepo),
	)
	providerOrderController := NewOrderController(
		orderApp.NewProviderListOrdersUseCase(orderRepo),
		orderApp.NewProviderGetOrderUseCase(orderRepo),
		orderApp.NewProviderOperateOrderUseCase(orderRepo),
	)

	mux := http.NewServeMux()
	RegisterAuthRoutes(mux, providerAuthController)
	RegisterOrderRoutes(mux, providerOrderController)

	created, err := userCreateOrderUC.Execute(t.Context(), orderApp.CreateOrderInput{
		UserID:           "u_1",
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
	if _, err = userPayOrderUC.Execute(t.Context(), orderApp.PayOrderMockSuccessInput{UserID: "u_1", OrderID: created.Order.ID}); err != nil {
		t.Fatalf("pay order failed: %v", err)
	}

	loginRaw, _ := json.Marshal(map[string]any{"account": "provider", "password": "provider123"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/provider/auth/login/mock", bytes.NewReader(loginRaw))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	mux.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("unexpected login status: %d", loginRec.Code)
	}
	var loginEnv Envelope
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginEnv); err != nil {
		t.Fatalf("decode login failed: %v", err)
	}
	loginDataRaw, _ := json.Marshal(loginEnv.Data)
	var loginData LoginMockResponseDTO
	if err := json.Unmarshal(loginDataRaw, &loginData); err != nil {
		t.Fatalf("decode login data failed: %v", err)
	}

	header := "Bearer " + loginData.AccessToken
	firstReq := httptest.NewRequest(http.MethodPost, "/api/v1/provider/orders/"+created.Order.ID+"/accept", nil)
	firstReq.Header.Set("Authorization", header)
	firstRec := httptest.NewRecorder()
	mux.ServeHTTP(firstRec, firstReq)
	if firstRec.Code != http.StatusOK {
		t.Fatalf("unexpected first accept status: %d", firstRec.Code)
	}

	secondReq := httptest.NewRequest(http.MethodPost, "/api/v1/provider/orders/"+created.Order.ID+"/accept", nil)
	secondReq.Header.Set("Authorization", header)
	secondRec := httptest.NewRecorder()
	mux.ServeHTTP(secondRec, secondReq)
	if secondRec.Code != http.StatusConflict {
		t.Fatalf("unexpected second accept status: %d", secondRec.Code)
	}
}

func TestProviderRoutes_AfterSaleBackflowVisibleToProviderDetail(t *testing.T) {
	orderRepo := memory.NewOrderRepository()
	feedbackRepo := memory.NewFeedbackRepository()
	actionRepo := memory.NewOrderAdminActionRepository()
	clock := aiApp.SystemClock{}
	idGenerator := aiApp.NewTimestampIDGenerator(clock)

	providerRepo := providerAuthApp.NewInMemoryRepository([]providerAuthDomain.ProviderAccount{
		{ProviderID: "p_pub_001", Account: "provider", Password: "provider123", DisplayName: "provider", Status: "active", CityCode: "310000"},
	})

	userCreateOrderUC := orderApp.NewCreateOrderUseCase(orderRepo, idGenerator, clock)
	userPayOrderUC := orderApp.NewPayOrderMockSuccessUseCase(orderRepo, clock)
	submitFeedbackUC := feedbackApp.NewSubmitOrderFeedbackUseCase(feedbackRepo, orderRepo, idGenerator, clock)
	adminUC := adminOrderApp.NewUseCase(orderRepo, feedbackRepo, actionRepo, idGenerator, clock)
	providerAuthController := NewAuthController(
		providerAuthApp.NewLoginMockUseCase(providerRepo, clock),
		providerAuthApp.NewGetCurrentProviderUseCase(providerRepo),
	)
	providerOrderController := NewOrderController(
		orderApp.NewProviderListOrdersUseCase(orderRepo),
		orderApp.NewProviderGetOrderUseCase(orderRepo),
		orderApp.NewProviderOperateOrderUseCase(orderRepo),
	)

	mux := http.NewServeMux()
	RegisterAuthRoutes(mux, providerAuthController)
	RegisterOrderRoutes(mux, providerOrderController)

	created, err := userCreateOrderUC.Execute(t.Context(), orderApp.CreateOrderInput{
		UserID:           "u_1",
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
	if _, err = userPayOrderUC.Execute(t.Context(), orderApp.PayOrderMockSuccessInput{UserID: "u_1", OrderID: created.Order.ID}); err != nil {
		t.Fatalf("pay order failed: %v", err)
	}
	if _, err = submitFeedbackUC.Execute(t.Context(), feedbackApp.SubmitFeedbackInput{
		UserID:           "u_1",
		OrderID:          created.Order.ID,
		RatingScore:      8,
		ReviewTags:       []string{"投诉"},
		ReviewContent:    "需处理",
		ComplaintReason:  "服务内容与描述不符",
		ComplaintContent: "请平台协助",
	}); err != nil {
		t.Fatalf("submit feedback failed: %v", err)
	}
	if _, err = adminUC.ActionOrder(t.Context(), adminOrderApp.ActionOrderInput{
		OrderID:  created.Order.ID,
		Action:   "intervene",
		Operator: "admin_001",
		Reason:   "平台人工介入",
	}); err != nil {
		t.Fatalf("admin intervene failed: %v", err)
	}

	loginRaw, _ := json.Marshal(map[string]any{"account": "provider", "password": "provider123"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/provider/auth/login/mock", bytes.NewReader(loginRaw))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	mux.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("unexpected login status: %d", loginRec.Code)
	}
	var loginEnv Envelope
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginEnv); err != nil {
		t.Fatalf("decode login failed: %v", err)
	}
	loginDataRaw, _ := json.Marshal(loginEnv.Data)
	var loginData LoginMockResponseDTO
	if err := json.Unmarshal(loginDataRaw, &loginData); err != nil {
		t.Fatalf("decode login data failed: %v", err)
	}

	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/provider/orders/"+created.Order.ID, nil)
	detailReq.Header.Set("Authorization", "Bearer "+loginData.AccessToken)
	detailRec := httptest.NewRecorder()
	mux.ServeHTTP(detailRec, detailReq)
	if detailRec.Code != http.StatusOK {
		t.Fatalf("unexpected detail status: %d", detailRec.Code)
	}
	var detailEnv Envelope
	if err := json.Unmarshal(detailRec.Body.Bytes(), &detailEnv); err != nil {
		t.Fatalf("decode detail envelope failed: %v", err)
	}
	detailDataRaw, _ := json.Marshal(detailEnv.Data)
	var detail struct {
		Status             string  `json:"status"`
		StatusReason       string  `json:"status_reason"`
		StatusActionReason string  `json:"status_action_reason"`
		StatusUpdatedAt    *string `json:"status_updated_at"`
	}
	if err := json.Unmarshal(detailDataRaw, &detail); err != nil {
		t.Fatalf("decode detail data failed: %v", err)
	}
	if detail.Status != "after_sale_processing" {
		t.Fatalf("expected after sale status, got %s", detail.Status)
	}
	if detail.StatusReason != "介入中" {
		t.Fatalf("expected intervene reason, got %s", detail.StatusReason)
	}
	if detail.StatusActionReason != "平台人工介入" {
		t.Fatalf("expected action reason, got %s", detail.StatusActionReason)
	}
	if detail.StatusUpdatedAt == nil {
		t.Fatalf("expected status updated at")
	}
	if _, err := time.Parse(time.RFC3339, *detail.StatusUpdatedAt); err != nil {
		t.Fatalf("invalid status updated at: %v", err)
	}
}
