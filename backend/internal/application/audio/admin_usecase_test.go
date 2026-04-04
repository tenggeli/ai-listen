package audio

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "listen/backend/internal/domain/audio"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

type fixedAudioIDGenerator struct{}

func (fixedAudioIDGenerator) NewID(prefix string) string {
	return prefix + "_001"
}

func TestAdminSoundUseCases_CreateUpdateStatusAndList(t *testing.T) {
	repo := memory.NewSoundContentRepository()
	createUC := NewCreateAdminSoundUseCase(repo, fixedAudioIDGenerator{})
	updateUC := NewUpdateAdminSoundUseCase(repo)
	statusUC := NewUpdateAdminSoundStatusUseCase(repo)
	listUC := NewListAdminSoundsUseCase(repo)
	readUC := NewGetSoundPageUseCase(repo)

	created, err := createUC.Execute(context.Background(), CreateAdminSoundInput{
		CategoryKey:   "sleep",
		Title:         "入眠呼吸练习",
		PlayCountText: "0 次播放",
		DurationText:  "10:00",
		Emoji:         "🌙",
		Author:        "listen 运营",
		SortOrder:     7,
	})
	if err != nil {
		t.Fatalf("create sound failed: %v", err)
	}
	if created.Item.Status != domain.SoundStatusInactive {
		t.Fatalf("expected default inactive, got: %s", created.Item.Status)
	}

	updated, err := updateUC.Execute(context.Background(), UpdateAdminSoundInput{
		TrackID:       created.Item.ID,
		CategoryKey:   "sleep",
		Title:         "入眠呼吸练习 · 更新版",
		PlayCountText: "3 次播放",
		DurationText:  "11:00",
		Emoji:         "🌌",
		Author:        "listen 运营",
		SortOrder:     8,
	})
	if err != nil {
		t.Fatalf("update sound failed: %v", err)
	}
	if updated.Item.Title != "入眠呼吸练习 · 更新版" {
		t.Fatalf("unexpected updated title: %s", updated.Item.Title)
	}

	activated, err := statusUC.Execute(context.Background(), UpdateAdminSoundStatusInput{
		TrackID: created.Item.ID,
		Action:  "activate",
	})
	if err != nil {
		t.Fatalf("activate sound failed: %v", err)
	}
	if activated.Item.Status != domain.SoundStatusActive {
		t.Fatalf("expected active status, got: %s", activated.Item.Status)
	}

	listOutput, err := listUC.Execute(context.Background(), ListAdminSoundsInput{
		Status: domain.SoundStatusActive,
	})
	if err != nil {
		t.Fatalf("list sounds failed: %v", err)
	}
	found := false
	for _, item := range listOutput.Items {
		if item.ID == created.Item.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected new active sound in admin list")
	}

	userOutput, err := readUC.Execute(context.Background(), GetSoundPageInput{
		Page: "home",
	})
	if err != nil {
		t.Fatalf("get sound page failed: %v", err)
	}
	userFound := false
	for _, item := range userOutput.Page.RecommendedTracks {
		if item.ID == created.Item.ID {
			userFound = true
			break
		}
	}
	if !userFound {
		t.Fatal("expected active sound visible on user sound page")
	}
}

func TestCreateAdminSoundUseCase_ValidateInput(t *testing.T) {
	repo := memory.NewSoundContentRepository()
	createUC := NewCreateAdminSoundUseCase(repo, fixedAudioIDGenerator{})

	_, err := createUC.Execute(context.Background(), CreateAdminSoundInput{
		CategoryKey:  "unknown",
		Title:        "abc",
		DurationText: "01:00",
	})
	if !errors.Is(err, domain.ErrInvalidCategory) {
		t.Fatalf("expected invalid category, got: %v", err)
	}

	_, err = createUC.Execute(context.Background(), CreateAdminSoundInput{
		CategoryKey:  "sleep",
		Title:        "",
		DurationText: "01:00",
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected invalid input for empty title, got: %v", err)
	}
}

func TestListAdminSoundsUseCase_PaginationBoundary(t *testing.T) {
	repo := memory.NewSoundContentRepository()
	listUC := NewListAdminSoundsUseCase(repo)

	output, err := listUC.Execute(context.Background(), ListAdminSoundsInput{
		PageNo:   999,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(output.Items) != 0 {
		t.Fatalf("expected empty page, got: %d", len(output.Items))
	}
}

func TestUpdateAdminSoundStatusUseCase_InvalidAction(t *testing.T) {
	repo := memory.NewSoundContentRepository()
	statusUC := NewUpdateAdminSoundStatusUseCase(repo)

	_, err := statusUC.Execute(context.Background(), UpdateAdminSoundStatusInput{
		TrackID: "track-rain",
		Action:  "invalid-action",
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got: %v", err)
	}
}

func TestCreateAdminSoundUseCase_RequireGeneratorWhenTrackIDEmpty(t *testing.T) {
	repo := memory.NewSoundContentRepository()
	createUC := NewCreateAdminSoundUseCase(repo, nil)

	_, err := createUC.Execute(context.Background(), CreateAdminSoundInput{
		CategoryKey:  "sleep",
		Title:        "xx",
		DurationText: "01:00",
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected invalid input for nil generator, got: %v", err)
	}
}

func TestCreateAdminSoundUseCase_UseProvidedTrackID(t *testing.T) {
	repo := memory.NewSoundContentRepository()
	createUC := NewCreateAdminSoundUseCase(repo, fixedAudioIDGenerator{})

	const fixedID = "track_custom_id"
	output, err := createUC.Execute(context.Background(), CreateAdminSoundInput{
		TrackID:      fixedID,
		CategoryKey:  "sleep",
		Title:        "manual",
		DurationText: "01:00",
		Status:       domain.SoundStatusInactive,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Item.ID != fixedID {
		t.Fatalf("expected fixed track id, got: %s", output.Item.ID)
	}
}

func TestListAdminSoundsUseCase_InvalidStatus(t *testing.T) {
	repo := memory.NewSoundContentRepository()
	listUC := NewListAdminSoundsUseCase(repo)

	_, err := listUC.Execute(context.Background(), ListAdminSoundsInput{
		Status: "bad-status",
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected invalid input for status, got: %v", err)
	}
}

func TestUpdateAdminSoundUseCase_NotFound(t *testing.T) {
	repo := memory.NewSoundContentRepository()
	updateUC := NewUpdateAdminSoundUseCase(repo)

	_, err := updateUC.Execute(context.Background(), UpdateAdminSoundInput{
		TrackID:      "track_not_found_" + time.Now().Format("150405"),
		CategoryKey:  "sleep",
		Title:        "x",
		DurationText: "01:00",
	})
	if !errors.Is(err, domain.ErrSoundNotFound) {
		t.Fatalf("expected sound not found, got: %v", err)
	}
}
