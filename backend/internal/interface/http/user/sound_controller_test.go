package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	app "listen/backend/internal/application/audio"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestSoundRoutes_GetSounds(t *testing.T) {
	controller := NewSoundController(app.NewGetSoundPageUseCase(memory.NewSoundContentRepository()))
	mux := http.NewServeMux()
	RegisterSoundRoutes(mux, controller)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sounds?page=home&user_id=u_1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var envelope Envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode envelope failed: %v", err)
	}
	dataBytes, _ := json.Marshal(envelope.Data)
	var body struct {
		Title             string `json:"title"`
		RecommendedTracks []any  `json:"recommended_tracks"`
	}
	if err := json.Unmarshal(dataBytes, &body); err != nil {
		t.Fatalf("decode data failed: %v", err)
	}
	if body.Title != "声音" {
		t.Fatalf("unexpected title: %s", body.Title)
	}
	if len(body.RecommendedTracks) == 0 {
		t.Fatal("expected recommended tracks")
	}
}

func TestSoundRoutes_InvalidPage(t *testing.T) {
	controller := NewSoundController(app.NewGetSoundPageUseCase(memory.NewSoundContentRepository()))
	mux := http.NewServeMux()
	RegisterSoundRoutes(mux, controller)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sounds?page=detail&user_id=u_1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestSoundRoutes_InvalidCategory(t *testing.T) {
	controller := NewSoundController(app.NewGetSoundPageUseCase(memory.NewSoundContentRepository()))
	mux := http.NewServeMux()
	RegisterSoundRoutes(mux, controller)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sounds?page=home&category_key=invalid&user_id=u_1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestSoundRoutes_PaginationBoundary(t *testing.T) {
	controller := NewSoundController(app.NewGetSoundPageUseCase(memory.NewSoundContentRepository()))
	mux := http.NewServeMux()
	RegisterSoundRoutes(mux, controller)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sounds?page=home&page_no=999&page_size=10&user_id=u_1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var envelope Envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode envelope failed: %v", err)
	}
	dataBytes, _ := json.Marshal(envelope.Data)
	var body struct {
		RecommendedTracks []any `json:"recommended_tracks"`
	}
	if err := json.Unmarshal(dataBytes, &body); err != nil {
		t.Fatalf("decode data failed: %v", err)
	}
	if len(body.RecommendedTracks) != 0 {
		t.Fatalf("expected empty tracks on boundary page, got: %d", len(body.RecommendedTracks))
	}
}
