package provider

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	aiApp "listen/backend/internal/application/ai"
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
