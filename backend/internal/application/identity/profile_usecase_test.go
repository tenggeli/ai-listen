package identity

import (
	"context"
	"testing"

	infraIdentity "listen/backend/internal/infrastructure/identity"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestSaveUserProfileUseCase_CompleteProfile(t *testing.T) {
	repo := memory.NewIdentityRepository()
	auth := infraIdentity.NewMockAuthService()

	loginUC := NewLoginBySMSUseCase(repo, auth, fixedClock{}, fixedIDGenerator{})
	saveUC := NewSaveUserProfileUseCase(repo)
	getUC := NewGetUserProfileUseCase(repo)

	loginOutput, err := loginUC.Execute(context.Background(), LoginBySMSInput{
		Phone:             "13800000001",
		VerifyCode:        "123",
		AgreementAccepted: true,
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	_, err = saveUC.Execute(context.Background(), SaveUserProfileInput{
		UserID:    loginOutput.Identity.UserID,
		Nickname:  "alice",
		Gender:    "female",
		AgeRange:  "25-29",
		City:      "shanghai",
		Bio:       "hello",
		AvatarURL: "https://example.com/a.png",
	})
	if err != nil {
		t.Fatalf("save profile failed: %v", err)
	}

	got, err := getUC.Execute(context.Background(), GetUserProfileInput{UserID: loginOutput.Identity.UserID})
	if err != nil {
		t.Fatalf("get profile failed: %v", err)
	}
	if got.User.Nickname != "alice" {
		t.Fatalf("unexpected nickname: %s", got.User.Nickname)
	}
	if got.Profile.Gender != "female" {
		t.Fatalf("unexpected gender: %s", got.Profile.Gender)
	}
	if !got.User.ProfileCompleted {
		t.Fatalf("expected profile completed")
	}
}

func TestSaveUserProfileUseCase_GenderChangeNeedsConfirmation(t *testing.T) {
	repo := memory.NewIdentityRepository()
	auth := infraIdentity.NewMockAuthService()

	loginUC := NewLoginBySMSUseCase(repo, auth, fixedClock{}, fixedIDGenerator{})
	saveUC := NewSaveUserProfileUseCase(repo)

	loginOutput, err := loginUC.Execute(context.Background(), LoginBySMSInput{
		Phone:             "13800000002",
		VerifyCode:        "123",
		AgreementAccepted: true,
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	_, err = saveUC.Execute(context.Background(), SaveUserProfileInput{
		UserID:   loginOutput.Identity.UserID,
		Nickname: "bob",
		Gender:   "male",
		AgeRange: "30-34",
		City:     "beijing",
	})
	if err != nil {
		t.Fatalf("first save failed: %v", err)
	}

	_, err = saveUC.Execute(context.Background(), SaveUserProfileInput{
		UserID:   loginOutput.Identity.UserID,
		Nickname: "bob",
		Gender:   "female",
		AgeRange: "30-34",
		City:     "beijing",
	})
	if err == nil {
		t.Fatalf("expected gender change confirmation error")
	}

	_, err = saveUC.Execute(context.Background(), SaveUserProfileInput{
		UserID:                loginOutput.Identity.UserID,
		Nickname:              "bob",
		Gender:                "female",
		AgeRange:              "30-34",
		City:                  "beijing",
		GenderChangeConfirmed: true,
	})
	if err != nil {
		t.Fatalf("confirmed gender change should pass: %v", err)
	}
}
