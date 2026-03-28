package identity

import (
	"context"
	"testing"

	infraIdentity "listen/backend/internal/infrastructure/identity"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestSaveUserPersonalityUseCase_Update(t *testing.T) {
	repo := memory.NewIdentityRepository()
	auth := infraIdentity.NewMockAuthService()

	loginUC := NewLoginBySMSUseCase(repo, auth, fixedClock{}, fixedIDGenerator{})
	saveUC := NewSaveUserPersonalityUseCase(repo)
	getUC := NewGetUserProfileUseCase(repo)

	loginOutput, err := loginUC.Execute(context.Background(), LoginBySMSInput{
		Phone:             "13800000003",
		VerifyCode:        "123456",
		AgreementAccepted: true,
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	_, err = saveUC.Execute(context.Background(), SaveUserPersonalityInput{
		UserID:       loginOutput.Identity.UserID,
		MBTI:         "infj",
		InterestTags: []string{"movie", "night_run", "movie"},
	})
	if err != nil {
		t.Fatalf("save personality failed: %v", err)
	}

	got, err := getUC.Execute(context.Background(), GetUserProfileInput{UserID: loginOutput.Identity.UserID})
	if err != nil {
		t.Fatalf("get profile failed: %v", err)
	}
	if got.Personality.MBTI != "INFJ" {
		t.Fatalf("unexpected mbti: %s", got.Personality.MBTI)
	}
	if len(got.Personality.InterestTags) != 2 {
		t.Fatalf("unexpected interest tag size: %d", len(got.Personality.InterestTags))
	}
	if !got.User.PersonalityCompleted {
		t.Fatalf("expected personality completed")
	}
}

func TestSkipUserPersonalityUseCase(t *testing.T) {
	repo := memory.NewIdentityRepository()
	auth := infraIdentity.NewMockAuthService()

	loginUC := NewLoginBySMSUseCase(repo, auth, fixedClock{}, fixedIDGenerator{})
	skipUC := NewSkipUserPersonalityUseCase(repo)

	loginOutput, err := loginUC.Execute(context.Background(), LoginBySMSInput{
		Phone:             "13800000004",
		VerifyCode:        "123456",
		AgreementAccepted: true,
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	output, err := skipUC.Execute(context.Background(), SkipUserPersonalityInput{
		UserID: loginOutput.Identity.UserID,
	})
	if err != nil {
		t.Fatalf("skip failed: %v", err)
	}
	if !output.Personality.Skipped {
		t.Fatalf("expected skipped")
	}
	if !output.User.PersonalityCompleted {
		t.Fatalf("expected personality completed after skip")
	}
}
