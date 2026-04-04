package audio

import (
	"context"
	"strings"

	domain "listen/backend/internal/domain/audio"
)

type IDGenerator interface {
	NewID(prefix string) string
}

type ListAdminSoundsInput struct {
	CategoryKey string
	Status      string
	Keyword     string
	PageNo      int
	PageSize    int
}

type ListAdminSoundsOutput struct {
	Items []domain.AdminTrack
	Total int
}

type CreateAdminSoundInput struct {
	TrackID       string
	CategoryKey   string
	Title         string
	PlayCountText string
	DurationText  string
	Emoji         string
	Author        string
	SortOrder     int
	Status        string
}

type CreateAdminSoundOutput struct {
	Item domain.AdminTrack
}

type UpdateAdminSoundInput struct {
	TrackID       string
	CategoryKey   string
	Title         string
	PlayCountText string
	DurationText  string
	Emoji         string
	Author        string
	SortOrder     int
}

type UpdateAdminSoundOutput struct {
	Item domain.AdminTrack
}

type UpdateAdminSoundStatusInput struct {
	TrackID string
	Action  string
}

type UpdateAdminSoundStatusOutput struct {
	Item domain.AdminTrack
}

type ListAdminSoundsUseCase struct {
	repo domain.Repository
}

func NewListAdminSoundsUseCase(repo domain.Repository) ListAdminSoundsUseCase {
	return ListAdminSoundsUseCase{repo: repo}
}

func (u ListAdminSoundsUseCase) Execute(ctx context.Context, input ListAdminSoundsInput) (ListAdminSoundsOutput, error) {
	status, err := normalizeStatus(input.Status)
	if err != nil {
		return ListAdminSoundsOutput{}, err
	}
	pageNo := input.PageNo
	if pageNo <= 0 {
		pageNo = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	items, total, err := u.repo.ListAdminTracks(ctx, domain.AdminTrackQuery{
		CategoryKey: strings.TrimSpace(input.CategoryKey),
		Status:      status,
		Keyword:     strings.TrimSpace(input.Keyword),
		PageNo:      pageNo,
		PageSize:    pageSize,
	})
	if err != nil {
		return ListAdminSoundsOutput{}, err
	}
	return ListAdminSoundsOutput{Items: items, Total: total}, nil
}

type CreateAdminSoundUseCase struct {
	repo        domain.Repository
	idGenerator IDGenerator
}

func NewCreateAdminSoundUseCase(repo domain.Repository, idGenerator IDGenerator) CreateAdminSoundUseCase {
	return CreateAdminSoundUseCase{repo: repo, idGenerator: idGenerator}
}

func (u CreateAdminSoundUseCase) Execute(ctx context.Context, input CreateAdminSoundInput) (CreateAdminSoundOutput, error) {
	trackID := strings.TrimSpace(input.TrackID)
	if trackID == "" {
		if u.idGenerator == nil {
			return CreateAdminSoundOutput{}, domain.ErrInvalidInput
		}
		trackID = u.idGenerator.NewID("track")
	}
	status, err := normalizeStatus(input.Status)
	if err != nil {
		return CreateAdminSoundOutput{}, err
	}
	if status == "" {
		status = domain.SoundStatusInactive
	}
	categoryKey := strings.TrimSpace(input.CategoryKey)
	if !u.isValidCategory(ctx, categoryKey) {
		return CreateAdminSoundOutput{}, domain.ErrInvalidCategory
	}
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.DurationText) == "" {
		return CreateAdminSoundOutput{}, domain.ErrInvalidInput
	}

	item := domain.AdminTrack{
		ID:            trackID,
		CategoryKey:   categoryKey,
		Title:         strings.TrimSpace(input.Title),
		PlayCountText: strings.TrimSpace(input.PlayCountText),
		DurationText:  strings.TrimSpace(input.DurationText),
		Emoji:         strings.TrimSpace(input.Emoji),
		Author:        strings.TrimSpace(input.Author),
		SortOrder:     input.SortOrder,
		Status:        status,
	}
	if err := u.repo.CreateAdminTrack(ctx, item); err != nil {
		return CreateAdminSoundOutput{}, err
	}
	return CreateAdminSoundOutput{Item: item}, nil
}

func (u CreateAdminSoundUseCase) isValidCategory(ctx context.Context, categoryKey string) bool {
	if categoryKey == "" || categoryKey == "all" {
		return false
	}
	categories, err := u.repo.ListCategories(ctx)
	if err != nil {
		return false
	}
	for _, category := range categories {
		if category.Key == categoryKey {
			return true
		}
	}
	return false
}

type UpdateAdminSoundUseCase struct {
	repo domain.Repository
}

func NewUpdateAdminSoundUseCase(repo domain.Repository) UpdateAdminSoundUseCase {
	return UpdateAdminSoundUseCase{repo: repo}
}

func (u UpdateAdminSoundUseCase) Execute(ctx context.Context, input UpdateAdminSoundInput) (UpdateAdminSoundOutput, error) {
	trackID := strings.TrimSpace(input.TrackID)
	categoryKey := strings.TrimSpace(input.CategoryKey)
	if trackID == "" || categoryKey == "" {
		return UpdateAdminSoundOutput{}, domain.ErrInvalidInput
	}
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.DurationText) == "" {
		return UpdateAdminSoundOutput{}, domain.ErrInvalidInput
	}
	categories, err := u.repo.ListCategories(ctx)
	if err != nil {
		return UpdateAdminSoundOutput{}, err
	}
	validCategory := false
	for _, category := range categories {
		if category.Key == categoryKey {
			validCategory = true
			break
		}
	}
	if !validCategory {
		return UpdateAdminSoundOutput{}, domain.ErrInvalidCategory
	}

	item, err := u.repo.GetAdminTrackByID(ctx, trackID)
	if err != nil {
		return UpdateAdminSoundOutput{}, err
	}
	item.CategoryKey = categoryKey
	item.Title = strings.TrimSpace(input.Title)
	item.PlayCountText = strings.TrimSpace(input.PlayCountText)
	item.DurationText = strings.TrimSpace(input.DurationText)
	item.Emoji = strings.TrimSpace(input.Emoji)
	item.Author = strings.TrimSpace(input.Author)
	item.SortOrder = input.SortOrder

	if err := u.repo.UpdateAdminTrack(ctx, item); err != nil {
		return UpdateAdminSoundOutput{}, err
	}
	return UpdateAdminSoundOutput{Item: item}, nil
}

type UpdateAdminSoundStatusUseCase struct {
	repo domain.Repository
}

func NewUpdateAdminSoundStatusUseCase(repo domain.Repository) UpdateAdminSoundStatusUseCase {
	return UpdateAdminSoundStatusUseCase{repo: repo}
}

func (u UpdateAdminSoundStatusUseCase) Execute(ctx context.Context, input UpdateAdminSoundStatusInput) (UpdateAdminSoundStatusOutput, error) {
	trackID := strings.TrimSpace(input.TrackID)
	if trackID == "" {
		return UpdateAdminSoundStatusOutput{}, domain.ErrInvalidInput
	}
	status, err := statusByAction(input.Action)
	if err != nil {
		return UpdateAdminSoundStatusOutput{}, err
	}
	if err := u.repo.UpdateAdminTrackStatus(ctx, trackID, status); err != nil {
		return UpdateAdminSoundStatusOutput{}, err
	}
	item, err := u.repo.GetAdminTrackByID(ctx, trackID)
	if err != nil {
		return UpdateAdminSoundStatusOutput{}, err
	}
	return UpdateAdminSoundStatusOutput{Item: item}, nil
}

func normalizeStatus(status string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "":
		return "", nil
	case domain.SoundStatusActive:
		return domain.SoundStatusActive, nil
	case domain.SoundStatusInactive:
		return domain.SoundStatusInactive, nil
	default:
		return "", domain.ErrInvalidInput
	}
}

func statusByAction(action string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "activate":
		return domain.SoundStatusActive, nil
	case "deactivate":
		return domain.SoundStatusInactive, nil
	default:
		return "", domain.ErrInvalidInput
	}
}
