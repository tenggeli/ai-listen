package audio

import (
	"context"
	"errors"
	"testing"

	domain "listen/backend/internal/domain/audio"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestGetSoundPageUseCase_Success(t *testing.T) {
	uc := NewGetSoundPageUseCase(memory.NewSoundContentRepository())

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
	uc := NewGetSoundPageUseCase(memory.NewSoundContentRepository())

	_, err := uc.Execute(context.Background(), GetSoundPageInput{Page: "other"})
	if err == nil {
		t.Fatal("expected invalid page error")
	}
}

func TestGetSoundPageUseCase_EmptyData(t *testing.T) {
	uc := NewGetSoundPageUseCase(emptySoundRepo{})

	output, err := uc.Execute(context.Background(), GetSoundPageInput{Page: "home"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Page.HasTracks() {
		t.Fatal("expected no tracks")
	}
	if output.Page.CurrentTrackID != "" {
		t.Fatalf("expected empty current track id, got: %s", output.Page.CurrentTrackID)
	}
	if output.Page.IsPlaying {
		t.Fatal("expected not playing on empty data")
	}
}

func TestGetSoundPageUseCase_InvalidCategory(t *testing.T) {
	uc := NewGetSoundPageUseCase(memory.NewSoundContentRepository())

	_, err := uc.Execute(context.Background(), GetSoundPageInput{
		Page:        "home",
		CategoryKey: "unknown",
	})
	if !errors.Is(err, domain.ErrInvalidCategory) {
		t.Fatalf("expected invalid category, got: %v", err)
	}
}

func TestGetSoundPageUseCase_PaginationBoundary(t *testing.T) {
	uc := NewGetSoundPageUseCase(memory.NewSoundContentRepository())

	output, err := uc.Execute(context.Background(), GetSoundPageInput{
		Page:     "home",
		PageNo:   999,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(output.Page.RecommendedTracks) != 0 {
		t.Fatalf("expected empty tracks on out of range page, got: %d", len(output.Page.RecommendedTracks))
	}
}

type emptySoundRepo struct{}

func (emptySoundRepo) ListCategories(context.Context) ([]domain.Category, error) {
	return []domain.Category{
		{Key: "all", Label: "全部"},
		{Key: "nature", Label: "自然白噪音"},
	}, nil
}

func (emptySoundRepo) ListTracks(context.Context, domain.TrackQuery) ([]domain.Track, int, error) {
	return []domain.Track{}, 0, nil
}

func (emptySoundRepo) ListAdminTracks(context.Context, domain.AdminTrackQuery) ([]domain.AdminTrack, int, error) {
	return []domain.AdminTrack{}, 0, nil
}

func (emptySoundRepo) GetAdminTrackByID(context.Context, string) (domain.AdminTrack, error) {
	return domain.AdminTrack{}, domain.ErrSoundNotFound
}

func (emptySoundRepo) CreateAdminTrack(context.Context, domain.AdminTrack) error {
	return nil
}

func (emptySoundRepo) UpdateAdminTrack(context.Context, domain.AdminTrack) error {
	return nil
}

func (emptySoundRepo) UpdateAdminTrackStatus(context.Context, string, string) error {
	return nil
}
