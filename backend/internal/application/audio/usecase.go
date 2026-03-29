package audio

import (
	"context"

	domain "listen/backend/internal/domain/audio"
)

type GetSoundPageInput struct {
	Page   string
	UserID string
}

type GetSoundPageOutput struct {
	Page domain.HomePage
}

type GetSoundPageUseCase struct {
	service domain.HomePageService
}

func NewGetSoundPageUseCase(service domain.HomePageService) GetSoundPageUseCase {
	return GetSoundPageUseCase{service: service}
}

func (u GetSoundPageUseCase) Execute(ctx context.Context, input GetSoundPageInput) (GetSoundPageOutput, error) {
	if input.Page != "home" {
		return GetSoundPageOutput{}, domain.ErrInvalidPage
	}
	page, err := u.service.GetHomePage(ctx, input.UserID)
	if err != nil {
		return GetSoundPageOutput{}, err
	}
	return GetSoundPageOutput{Page: page}, nil
}
