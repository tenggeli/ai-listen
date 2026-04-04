package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	audioApp "listen/backend/internal/application/audio"
	memory "listen/backend/internal/infrastructure/persistence/memory"
	userHTTP "listen/backend/internal/interface/http/user"
)

type fixedSoundIDGenerator struct{}

func (fixedSoundIDGenerator) NewID(prefix string) string {
	return prefix + "_001"
}

func TestSoundRoutes_CRUDAndStatusAction(t *testing.T) {
	repo := memory.NewSoundContentRepository()
	controller := NewSoundController(
		audioApp.NewListAdminSoundsUseCase(repo),
		audioApp.NewCreateAdminSoundUseCase(repo, fixedSoundIDGenerator{}),
		audioApp.NewUpdateAdminSoundUseCase(repo),
		audioApp.NewUpdateAdminSoundStatusUseCase(repo),
	)

	mux := http.NewServeMux()
	RegisterSoundRoutes(mux, controller)

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/sounds", nil)
	listReq.Header.Set("Authorization", "Bearer mock_admin_at_admin_001_1712013723")
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("unexpected list status: %d", listRec.Code)
	}

	createBody := map[string]any{
		"category_key":    "sleep",
		"title":           "联调新增",
		"play_count_text": "0 次播放",
		"duration_text":   "10:00",
		"emoji":           "🌙",
		"author":          "运营后台",
		"sort_order":      99,
	}
	rawCreate, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/sounds", bytes.NewReader(rawCreate))
	createReq.Header.Set("Authorization", "Bearer mock_admin_at_admin_001_1712013723")
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusOK {
		t.Fatalf("unexpected create status: %d", createRec.Code)
	}

	var createdEnvelope Envelope
	if err := json.Unmarshal(createRec.Body.Bytes(), &createdEnvelope); err != nil {
		t.Fatalf("decode create response failed: %v", err)
	}
	dataBytes, _ := json.Marshal(createdEnvelope.Data)
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(dataBytes, &created); err != nil {
		t.Fatalf("decode create data failed: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected created id")
	}

	updateBody := map[string]any{
		"category_key":    "sleep",
		"title":           "联调新增-更新",
		"play_count_text": "1 次播放",
		"duration_text":   "11:00",
		"emoji":           "🌌",
		"author":          "运营后台",
		"sort_order":      88,
	}
	rawUpdate, _ := json.Marshal(updateBody)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/admin/sounds/"+created.ID, bytes.NewReader(rawUpdate))
	updateReq.Header.Set("Authorization", "Bearer mock_admin_at_admin_001_1712013723")
	updateRec := httptest.NewRecorder()
	mux.ServeHTTP(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("unexpected update status: %d", updateRec.Code)
	}

	activeReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/sounds/"+created.ID+"/activate", nil)
	activeReq.Header.Set("Authorization", "Bearer mock_admin_at_admin_001_1712013723")
	activeRec := httptest.NewRecorder()
	mux.ServeHTTP(activeRec, activeReq)
	if activeRec.Code != http.StatusOK {
		t.Fatalf("unexpected activate status: %d", activeRec.Code)
	}

	inactiveReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/sounds/"+created.ID+"/deactivate", nil)
	inactiveReq.Header.Set("Authorization", "Bearer mock_admin_at_admin_001_1712013723")
	inactiveRec := httptest.NewRecorder()
	mux.ServeHTTP(inactiveRec, inactiveReq)
	if inactiveRec.Code != http.StatusOK {
		t.Fatalf("unexpected deactivate status: %d", inactiveRec.Code)
	}
}

func TestSoundRoutes_AuthRequired(t *testing.T) {
	repo := memory.NewSoundContentRepository()
	controller := NewSoundController(
		audioApp.NewListAdminSoundsUseCase(repo),
		audioApp.NewCreateAdminSoundUseCase(repo, fixedSoundIDGenerator{}),
		audioApp.NewUpdateAdminSoundUseCase(repo),
		audioApp.NewUpdateAdminSoundStatusUseCase(repo),
	)

	mux := http.NewServeMux()
	RegisterSoundRoutes(mux, controller)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/sounds", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized, got: %d", rec.Code)
	}
}

func TestSoundRoutes_RegressUserSoundsReadActiveData(t *testing.T) {
	repo := memory.NewSoundContentRepository()
	adminController := NewSoundController(
		audioApp.NewListAdminSoundsUseCase(repo),
		audioApp.NewCreateAdminSoundUseCase(repo, fixedSoundIDGenerator{}),
		audioApp.NewUpdateAdminSoundUseCase(repo),
		audioApp.NewUpdateAdminSoundStatusUseCase(repo),
	)
	userController := userHTTP.NewSoundController(audioApp.NewGetSoundPageUseCase(repo))

	adminMux := http.NewServeMux()
	RegisterSoundRoutes(adminMux, adminController)

	userMux := http.NewServeMux()
	userHTTP.RegisterSoundRoutes(userMux, userController)

	createBody := map[string]any{
		"track_id":        "track_from_admin",
		"category_key":    "sleep",
		"title":           "用户侧可见验证",
		"play_count_text": "0 次播放",
		"duration_text":   "09:00",
		"emoji":           "🎧",
		"author":          "运营后台",
		"sort_order":      10,
		"status":          "active",
	}
	rawCreate, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/sounds", bytes.NewReader(rawCreate))
	createReq.Header.Set("Authorization", "Bearer mock_admin_at_admin_001_1712013723")
	createRec := httptest.NewRecorder()
	adminMux.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusOK {
		t.Fatalf("create status: %d", createRec.Code)
	}

	readReq := httptest.NewRequest(http.MethodGet, "/api/v1/sounds?page=home&user_id=u_1", nil)
	readRec := httptest.NewRecorder()
	userMux.ServeHTTP(readRec, readReq)
	if readRec.Code != http.StatusOK {
		t.Fatalf("unexpected user sound status: %d", readRec.Code)
	}

	var envelope userHTTP.Envelope
	if err := json.Unmarshal(readRec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}
	dataRaw, _ := json.Marshal(envelope.Data)
	var body struct {
		RecommendedTracks []struct {
			ID string `json:"id"`
		} `json:"recommended_tracks"`
	}
	if err := json.Unmarshal(dataRaw, &body); err != nil {
		t.Fatalf("decode data failed: %v", err)
	}
	found := false
	for _, item := range body.RecommendedTracks {
		if item.ID == "track_from_admin" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected active sound from admin to be visible to user endpoint")
	}
}
