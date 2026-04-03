package provider

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	aiApp "listen/backend/internal/application/ai"
	providerApp "listen/backend/internal/application/provider"
	providerAuthApp "listen/backend/internal/application/provider_auth"
	providerAuthDomain "listen/backend/internal/domain/provider_auth"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestProviderRoutes_ProfileUpdateAndGet(t *testing.T) {
	mux := buildProviderProfileTestMux()
	token := providerLoginToken(t, mux)
	header := "Bearer " + token

	updateBody, _ := json.Marshal(map[string]any{"display_name": "倾听师·夏", "city_code": "320100"})
	updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/provider/profile", bytes.NewReader(updateBody))
	updateReq.Header.Set("Authorization", header)
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	mux.ServeHTTP(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("unexpected update status: %d", updateRec.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/provider/profile", nil)
	getReq.Header.Set("Authorization", header)
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("unexpected get status: %d", getRec.Code)
	}

	var env Envelope
	if err := json.Unmarshal(getRec.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}
	raw, _ := json.Marshal(env.Data)
	var profile ProviderMeResponseDTO
	if err := json.Unmarshal(raw, &profile); err != nil {
		t.Fatalf("decode profile failed: %v", err)
	}
	if profile.DisplayName != "倾听师·夏" {
		t.Fatalf("expected updated display_name, got %s", profile.DisplayName)
	}
	if profile.CityCode != "320100" {
		t.Fatalf("expected updated city_code, got %s", profile.CityCode)
	}
}

func TestProviderRoutes_ProfileUpdateValidationAndErrors(t *testing.T) {
	mux := buildProviderProfileTestMux()
	token := providerLoginToken(t, mux)
	header := "Bearer " + token

	invalidReq := httptest.NewRequest(http.MethodPut, "/api/v1/provider/profile", bytes.NewBufferString("{bad-json"))
	invalidReq.Header.Set("Authorization", header)
	invalidReq.Header.Set("Content-Type", "application/json")
	invalidRec := httptest.NewRecorder()
	mux.ServeHTTP(invalidRec, invalidReq)
	if invalidRec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected invalid body status: %d", invalidRec.Code)
	}

	emptyNameBody, _ := json.Marshal(map[string]any{"display_name": "  ", "city_code": "310000"})
	emptyNameReq := httptest.NewRequest(http.MethodPut, "/api/v1/provider/profile", bytes.NewReader(emptyNameBody))
	emptyNameReq.Header.Set("Authorization", header)
	emptyNameReq.Header.Set("Content-Type", "application/json")
	emptyNameRec := httptest.NewRecorder()
	mux.ServeHTTP(emptyNameRec, emptyNameReq)
	if emptyNameRec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected empty name status: %d", emptyNameRec.Code)
	}

	unauthorizedReq := httptest.NewRequest(http.MethodPut, "/api/v1/provider/profile", bytes.NewReader(emptyNameBody))
	unauthorizedRec := httptest.NewRecorder()
	mux.ServeHTTP(unauthorizedRec, unauthorizedReq)
	if unauthorizedRec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected unauthorized status: %d", unauthorizedRec.Code)
	}

	notFoundBody, _ := json.Marshal(map[string]any{"display_name": "昵称", "city_code": "310000"})
	notFoundReq := httptest.NewRequest(http.MethodPut, "/api/v1/provider/profile", bytes.NewReader(notFoundBody))
	notFoundReq.Header.Set("X-Provider-ID", "provider_not_found")
	notFoundReq.Header.Set("Content-Type", "application/json")
	notFoundRec := httptest.NewRecorder()
	mux.ServeHTTP(notFoundRec, notFoundReq)
	if notFoundRec.Code != http.StatusNotFound {
		t.Fatalf("unexpected not found status: %d", notFoundRec.Code)
	}
}

func TestProviderRoutes_ListServices(t *testing.T) {
	mux := buildProviderProfileTestMux()
	token := providerLoginToken(t, mux)
	header := "Bearer " + token

	successReq := httptest.NewRequest(http.MethodGet, "/api/v1/provider/services", nil)
	successReq.Header.Set("Authorization", header)
	successRec := httptest.NewRecorder()
	mux.ServeHTTP(successRec, successReq)
	if successRec.Code != http.StatusOK {
		t.Fatalf("unexpected success status: %d", successRec.Code)
	}

	var env Envelope
	if err := json.Unmarshal(successRec.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}
	raw, _ := json.Marshal(env.Data)
	var payload struct {
		Items []ProviderServiceItemDTO `json:"items"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("decode services failed: %v", err)
	}
	if len(payload.Items) == 0 {
		t.Fatalf("expected non-empty services")
	}

	unauthorizedReq := httptest.NewRequest(http.MethodGet, "/api/v1/provider/services", nil)
	unauthorizedRec := httptest.NewRecorder()
	mux.ServeHTTP(unauthorizedRec, unauthorizedReq)
	if unauthorizedRec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected unauthorized status: %d", unauthorizedRec.Code)
	}

	notFoundReq := httptest.NewRequest(http.MethodGet, "/api/v1/provider/services", nil)
	notFoundReq.Header.Set("X-Provider-ID", "provider_not_found")
	notFoundRec := httptest.NewRecorder()
	mux.ServeHTTP(notFoundRec, notFoundReq)
	if notFoundRec.Code != http.StatusNotFound {
		t.Fatalf("unexpected not found status: %d", notFoundRec.Code)
	}
}

func buildProviderProfileTestMux() *http.ServeMux {
	clock := aiApp.SystemClock{}
	providerAuthRepo := providerAuthApp.NewInMemoryRepository([]providerAuthDomain.ProviderAccount{
		{ProviderID: "p_pub_001", Account: "provider", Password: "provider123", DisplayName: "暖心倾听师 · 小林", Status: "active", CityCode: "310000"},
	})
	serviceRepo := memory.NewServiceDiscoveryRepository()

	authController := NewAuthController(
		providerAuthApp.NewLoginMockUseCase(providerAuthRepo, clock),
		providerAuthApp.NewGetCurrentProviderUseCase(providerAuthRepo),
	)
	profileController := NewProfileController(
		providerApp.NewUpdateCurrentProfileUseCase(providerAuthRepo),
		providerApp.NewListCurrentProviderServicesUseCase(serviceRepo),
	)

	mux := http.NewServeMux()
	RegisterAuthRoutes(mux, authController)
	RegisterProfileRoutes(mux, profileController)
	return mux
}

func providerLoginToken(t *testing.T, mux *http.ServeMux) string {
	t.Helper()
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
		t.Fatalf("decode login response failed: %v", err)
	}
	raw, _ := json.Marshal(loginEnv.Data)
	var loginData LoginMockResponseDTO
	if err := json.Unmarshal(raw, &loginData); err != nil {
		t.Fatalf("decode login data failed: %v", err)
	}
	if loginData.AccessToken == "" {
		t.Fatalf("expected access token")
	}
	return loginData.AccessToken
}
