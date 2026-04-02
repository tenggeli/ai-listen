package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	settingsApp "listen/backend/internal/application/user_settings"
	identityDomain "listen/backend/internal/domain/identity"
	settingsDomain "listen/backend/internal/domain/user_settings"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

type testSettingsRepo struct {
	data map[string]settingsDomain.Settings
}

func newTestSettingsRepo() *testSettingsRepo {
	return &testSettingsRepo{data: map[string]settingsDomain.Settings{}}
}

func (r *testSettingsRepo) GetByUserID(_ context.Context, userID string) (settingsDomain.Settings, bool, error) {
	item, found := r.data[userID]
	return item, found, nil
}

func (r *testSettingsRepo) Save(_ context.Context, settings settingsDomain.Settings) error {
	r.data[settings.UserID] = settings
	return nil
}

func TestSettingsRoutes_FullFlow(t *testing.T) {
	identityRepo := memory.NewIdentityRepository()
	settingsRepo := newTestSettingsRepo()
	account, _ := identityDomain.NewUserAccount("u_settings_1", "13800138002", "")
	if err := identityRepo.Save(context.Background(), account); err != nil {
		t.Fatalf("save user failed: %v", err)
	}

	controller := NewSettingsController(
		settingsApp.NewGetSettingsUseCase(settingsRepo, identityRepo),
		settingsApp.NewSaveSettingsUseCase(settingsRepo, identityRepo),
	)

	mux := http.NewServeMux()
	RegisterSettingsRoutes(mux, controller)

	saveBody := map[string]any{
		"preference": map[string]any{
			"prefer_same_city_providers": true,
			"auto_play_sound_preview":    false,
			"hide_offline_providers":     true,
		},
		"notification": map[string]any{
			"order_status_update":     true,
			"complaint_result_notice": false,
			"marketing_activity":      false,
		},
		"privacy": map[string]any{
			"profile_public_visible":      false,
			"personalized_recommendation": true,
			"risk_control_data_sharing":   false,
		},
	}
	rawBody, _ := json.Marshal(saveBody)
	saveReq := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/settings", bytes.NewReader(rawBody))
	saveReq.Header.Set("Authorization", "Bearer mock_at_u_settings_1_20260402010101")
	saveReq.Header.Set("Content-Type", "application/json")
	saveRec := httptest.NewRecorder()
	mux.ServeHTTP(saveRec, saveReq)
	if saveRec.Code != http.StatusOK {
		t.Fatalf("unexpected save status: %d", saveRec.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/settings", nil)
	getReq.Header.Set("Authorization", "Bearer mock_at_u_settings_1_20260402010101")
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("unexpected get status: %d", getRec.Code)
	}
}

func TestSettingsRoutes_Unauthorized(t *testing.T) {
	identityRepo := memory.NewIdentityRepository()
	settingsRepo := newTestSettingsRepo()
	controller := NewSettingsController(
		settingsApp.NewGetSettingsUseCase(settingsRepo, identityRepo),
		settingsApp.NewSaveSettingsUseCase(settingsRepo, identityRepo),
	)
	mux := http.NewServeMux()
	RegisterSettingsRoutes(mux, controller)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/settings", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestSettingsRoutes_MysqlOnlyInMemoryMode(t *testing.T) {
	identityRepo := memory.NewIdentityRepository()
	unsupportedRepo := memory.NewUserSettingsRepository()
	account, _ := identityDomain.NewUserAccount("u_settings_2", "13800138003", "")
	if err := identityRepo.Save(context.Background(), account); err != nil {
		t.Fatalf("save user failed: %v", err)
	}

	controller := NewSettingsController(
		settingsApp.NewGetSettingsUseCase(unsupportedRepo, identityRepo),
		settingsApp.NewSaveSettingsUseCase(unsupportedRepo, identityRepo),
	)
	mux := http.NewServeMux()
	RegisterSettingsRoutes(mux, controller)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/settings", nil)
	req.Header.Set("Authorization", "Bearer mock_at_u_settings_2_20260402010101")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestSettingsRoutes_InvalidBody(t *testing.T) {
	identityRepo := memory.NewIdentityRepository()
	settingsRepo := newTestSettingsRepo()
	account, _ := identityDomain.NewUserAccount("u_settings_3", "13800138004", "")
	if err := identityRepo.Save(context.Background(), account); err != nil {
		t.Fatalf("save user failed: %v", err)
	}

	controller := NewSettingsController(
		settingsApp.NewGetSettingsUseCase(settingsRepo, identityRepo),
		settingsApp.NewSaveSettingsUseCase(settingsRepo, identityRepo),
	)
	mux := http.NewServeMux()
	RegisterSettingsRoutes(mux, controller)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/settings", bytes.NewBufferString("invalid"))
	req.Header.Set("Authorization", "Bearer mock_at_u_settings_3_20260402010101")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}
