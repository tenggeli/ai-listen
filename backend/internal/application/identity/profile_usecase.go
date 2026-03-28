package identity

import (
	"context"
	"strings"

	domain "listen/backend/internal/domain/identity"
)

type GetUserProfileInput struct {
	UserID string
}

type GetUserProfileOutput struct {
	User        domain.User
	Profile     domain.UserProfile
	Personality domain.UserPersonalityProfile
}

type SaveUserProfileInput struct {
	UserID                string
	Nickname              string
	AvatarURL             string
	Gender                string
	AgeRange              string
	City                  string
	Bio                   string
	GenderChangeConfirmed bool
}

type SaveUserProfileOutput struct {
	User        domain.User
	Profile     domain.UserProfile
	Personality domain.UserPersonalityProfile
}

type GetUserProfileUseCase struct {
	repo domain.Repository
}

func NewGetUserProfileUseCase(repo domain.Repository) GetUserProfileUseCase {
	return GetUserProfileUseCase{repo: repo}
}

func (u GetUserProfileUseCase) Execute(ctx context.Context, input GetUserProfileInput) (GetUserProfileOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return GetUserProfileOutput{}, domain.ErrInvalidInput
	}

	account, found, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return GetUserProfileOutput{}, err
	}
	if !found {
		return GetUserProfileOutput{}, domain.ErrUserNotFound
	}
	return GetUserProfileOutput{
		User:        account.ToUser(),
		Profile:     account.ToProfile(),
		Personality: account.ToPersonalityProfile(),
	}, nil
}

type SaveUserProfileUseCase struct {
	repo domain.Repository
}

func NewSaveUserProfileUseCase(repo domain.Repository) SaveUserProfileUseCase {
	return SaveUserProfileUseCase{repo: repo}
}

func (u SaveUserProfileUseCase) Execute(ctx context.Context, input SaveUserProfileInput) (SaveUserProfileOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return SaveUserProfileOutput{}, domain.ErrInvalidInput
	}

	account, found, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return SaveUserProfileOutput{}, err
	}
	if !found {
		return SaveUserProfileOutput{}, domain.ErrUserNotFound
	}

	if err := account.CompleteBasicProfile(
		input.Nickname,
		input.AvatarURL,
		input.Gender,
		input.AgeRange,
		input.City,
		input.Bio,
		input.GenderChangeConfirmed,
	); err != nil {
		return SaveUserProfileOutput{}, err
	}
	if err := u.repo.Save(ctx, account); err != nil {
		return SaveUserProfileOutput{}, err
	}

	return SaveUserProfileOutput{
		User:        account.ToUser(),
		Profile:     account.ToProfile(),
		Personality: account.ToPersonalityProfile(),
	}, nil
}
