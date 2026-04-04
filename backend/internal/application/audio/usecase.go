package audio

import (
	"context"
	"strings"

	domain "listen/backend/internal/domain/audio"
)

type GetSoundPageInput struct {
	Page        string
	UserID      string
	CategoryKey string
	PageNo      int
	PageSize    int
}

type GetSoundPageOutput struct {
	Page domain.HomePage
}

type GetSoundPageUseCase struct {
	repo domain.Repository
}

func NewGetSoundPageUseCase(repo domain.Repository) GetSoundPageUseCase {
	return GetSoundPageUseCase{repo: repo}
}

func (u GetSoundPageUseCase) Execute(ctx context.Context, input GetSoundPageInput) (GetSoundPageOutput, error) {
	if input.Page != "home" {
		return GetSoundPageOutput{}, domain.ErrInvalidPage
	}
	categories, err := u.repo.ListCategories(ctx)
	if err != nil {
		return GetSoundPageOutput{}, err
	}

	categoryKey := normalizeCategoryKey(input.CategoryKey)
	if !isValidCategory(categoryKey, categories) {
		return GetSoundPageOutput{}, domain.ErrInvalidCategory
	}

	pageNo := normalizePageNo(input.PageNo)
	pageSize := normalizePageSize(input.PageSize)
	tracks, _, err := u.repo.ListTracks(ctx, domain.TrackQuery{
		CategoryKey: categoryKey,
		PageNo:      pageNo,
		PageSize:    pageSize,
	})
	if err != nil {
		return GetSoundPageOutput{}, err
	}

	currentTrackID := ""
	currentProgressText := "00:00"
	totalDurationText := "00:00"
	isPlaying := false
	if len(tracks) > 0 {
		currentTrackID = tracks[0].ID
		totalDurationText = tracks[0].DurationText
		isPlaying = true
	}

	page := domain.HomePage{
		Title:               "声音",
		Subtitle:            "用声音抚慰此刻的你",
		Categories:          categories,
		RecommendedTracks:   tracks,
		CurrentTrackID:      currentTrackID,
		CurrentProgressText: currentProgressText,
		TotalDurationText:   totalDurationText,
		IsPlaying:           isPlaying,
	}
	return GetSoundPageOutput{Page: page}, nil
}

func normalizeCategoryKey(raw string) string {
	key := strings.TrimSpace(raw)
	if key == "" {
		return "all"
	}
	return key
}

func isValidCategory(categoryKey string, categories []domain.Category) bool {
	if categoryKey == "all" {
		return true
	}
	for _, item := range categories {
		if item.Key == categoryKey {
			return true
		}
	}
	return false
}

func normalizePageNo(pageNo int) int {
	if pageNo <= 0 {
		return 1
	}
	return pageNo
}

func normalizePageSize(pageSize int) int {
	if pageSize <= 0 {
		return 20
	}
	if pageSize > 100 {
		return 100
	}
	return pageSize
}
