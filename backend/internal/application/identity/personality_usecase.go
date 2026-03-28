package identity

import (
	"context"
	"strings"

	domain "listen/backend/internal/domain/identity"
)

type SaveUserPersonalityInput struct {
	UserID       string
	MBTI         string
	InterestTags []string
}

type SaveUserPersonalityOutput struct {
	User        domain.User
	Profile     domain.UserProfile
	Personality domain.UserPersonalityProfile
}

type SkipUserPersonalityInput struct {
	UserID string
}

type SkipUserPersonalityOutput struct {
	User        domain.User
	Profile     domain.UserProfile
	Personality domain.UserPersonalityProfile
}

type SaveUserPersonalityUseCase struct {
	repo domain.Repository
}

func NewSaveUserPersonalityUseCase(repo domain.Repository) SaveUserPersonalityUseCase {
	return SaveUserPersonalityUseCase{repo: repo}
}

func (u SaveUserPersonalityUseCase) Execute(ctx context.Context, input SaveUserPersonalityInput) (SaveUserPersonalityOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return SaveUserPersonalityOutput{}, domain.ErrInvalidInput
	}

	account, found, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return SaveUserPersonalityOutput{}, err
	}
	if !found {
		return SaveUserPersonalityOutput{}, domain.ErrUserNotFound
	}

	if err := account.UpdatePersonality(input.MBTI, input.InterestTags); err != nil {
		return SaveUserPersonalityOutput{}, err
	}
	if err := u.repo.Save(ctx, account); err != nil {
		return SaveUserPersonalityOutput{}, err
	}

	return SaveUserPersonalityOutput{
		User:        account.ToUser(),
		Profile:     account.ToProfile(),
		Personality: account.ToPersonalityProfile(),
	}, nil
}

type SkipUserPersonalityUseCase struct {
	repo domain.Repository
}

func NewSkipUserPersonalityUseCase(repo domain.Repository) SkipUserPersonalityUseCase {
	return SkipUserPersonalityUseCase{repo: repo}
}

func (u SkipUserPersonalityUseCase) Execute(ctx context.Context, input SkipUserPersonalityInput) (SkipUserPersonalityOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return SkipUserPersonalityOutput{}, domain.ErrInvalidInput
	}

	account, found, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return SkipUserPersonalityOutput{}, err
	}
	if !found {
		return SkipUserPersonalityOutput{}, domain.ErrUserNotFound
	}

	account.SkipForNow()
	if err := u.repo.Save(ctx, account); err != nil {
		return SkipUserPersonalityOutput{}, err
	}

	return SkipUserPersonalityOutput{
		User:        account.ToUser(),
		Profile:     account.ToProfile(),
		Personality: account.ToPersonalityProfile(),
	}, nil
}
