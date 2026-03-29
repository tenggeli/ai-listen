package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	app "listen/backend/internal/application/service_discovery"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestServiceDiscoveryRoutes_ListCategories(t *testing.T) {
	mux := newServiceDiscoveryTestMux()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/services/categories", nil)

	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var envelope Envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("unmarshal envelope failed: %v", err)
	}
	dataBytes, _ := json.Marshal(envelope.Data)
	var body struct {
		Items []ServiceCategoryResponseDTO `json:"items"`
	}
	if err := json.Unmarshal(dataBytes, &body); err != nil {
		t.Fatalf("unmarshal data failed: %v", err)
	}
	if len(body.Items) == 0 {
		t.Fatal("expected seeded categories")
	}
}

func TestServiceDiscoveryRoutes_ListProviders_FilterByKeyword(t *testing.T) {
	mux := newServiceDiscoveryTestMux()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/public?keyword=睡前&page=1&page_size=10", nil)

	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var envelope Envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("unmarshal envelope failed: %v", err)
	}
	dataBytes, _ := json.Marshal(envelope.Data)
	var body struct {
		Items []ProviderPublicResponseDTO `json:"items"`
		Total int                         `json:"total"`
	}
	if err := json.Unmarshal(dataBytes, &body); err != nil {
		t.Fatalf("unmarshal data failed: %v", err)
	}
	if body.Total != 1 {
		t.Fatalf("unexpected total: %d", body.Total)
	}
	if body.Items[0].ID != "p_pub_003" {
		t.Fatalf("unexpected provider: %s", body.Items[0].ID)
	}
}

func TestServiceDiscoveryRoutes_GetProviderDetailAndItems(t *testing.T) {
	mux := newServiceDiscoveryTestMux()

	detailRec := httptest.NewRecorder()
	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/providers/public/p_pub_001", nil)
	mux.ServeHTTP(detailRec, detailReq)
	if detailRec.Code != http.StatusOK {
		t.Fatalf("unexpected detail status: %d", detailRec.Code)
	}

	itemsRec := httptest.NewRecorder()
	itemsReq := httptest.NewRequest(http.MethodGet, "/api/v1/providers/public/p_pub_001/service-items", nil)
	mux.ServeHTTP(itemsRec, itemsReq)
	if itemsRec.Code != http.StatusOK {
		t.Fatalf("unexpected items status: %d", itemsRec.Code)
	}

	var envelope Envelope
	if err := json.Unmarshal(itemsRec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("unmarshal items envelope failed: %v", err)
	}
	dataBytes, _ := json.Marshal(envelope.Data)
	var body struct {
		Items []ServiceItemResponseDTO `json:"items"`
	}
	if err := json.Unmarshal(dataBytes, &body); err != nil {
		t.Fatalf("unmarshal items data failed: %v", err)
	}
	if len(body.Items) != 3 {
		t.Fatalf("unexpected items len: %d", len(body.Items))
	}
}

func TestServiceDiscoveryRoutes_NotFound(t *testing.T) {
	mux := newServiceDiscoveryTestMux()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/public/not_found", nil)

	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func newServiceDiscoveryTestMux() *http.ServeMux {
	repo := memory.NewServiceDiscoveryRepository()
	controller := NewServiceDiscoveryController(
		app.NewListServiceCategoriesUseCase(repo),
		app.NewListPublicProvidersUseCase(repo),
		app.NewGetPublicProviderUseCase(repo),
		app.NewListProviderServiceItemsUseCase(repo),
	)

	mux := http.NewServeMux()
	RegisterServiceDiscoveryRoutes(mux, controller)
	return mux
}
