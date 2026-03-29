package audio

import (
	"context"
	"testing"

	infraAudio "listen/backend/internal/infrastructure/audio"
)

func TestGetSoundPageUseCase_Success(t *testing.T) {
	uc := NewGetSoundPageUseCase(infraAudio.NewMockSoundPageService())

	output, err := uc.Execute(context.Background(), GetSoundPageInput{
		Page:   "home",
		UserID: "u_test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Page.Title != "声音" {
		t.Fatalf("unexpected title: %s", output.Page.Title)
	}
	if len(output.Page.Categories) == 0 {
		t.Fatal("expected categories")
	}
	if !output.Page.HasTracks() {
		t.Fatal("expected tracks")
	}
}

func TestGetSoundPageUseCase_InvalidPage(t *testing.T) {
	uc := NewGetSoundPageUseCase(infraAudio.NewMockSoundPageService())

	_, err := uc.Execute(context.Background(), GetSoundPageInput{Page: "other"})
	if err == nil {
		t.Fatal("expected invalid page error")
	}
}
