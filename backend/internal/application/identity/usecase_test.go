package identity

import (
	"context"
	"testing"
	"time"

	infraIdentity "listen/backend/internal/infrastructure/identity"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

type fixedClock struct{}

func (fixedClock) Now() time.Time {
	return time.Date(2026, 3, 28, 12, 0, 0, 0, time.UTC)
}

type fixedIDGenerator struct{}

func (fixedIDGenerator) NewID(prefix string) string {
	return prefix + "_001"
}

func TestLoginBySMSUseCase_NewUser(t *testing.T) {
	repo := memory.NewIdentityRepository()
	auth := infraIdentity.NewMockAuthService()
	uc := NewLoginBySMSUseCase(repo, auth, fixedClock{}, fixedIDGenerator{})

	output, err := uc.Execute(context.Background(), LoginBySMSInput{
		Phone:             "13800000000",
		VerifyCode:        "123",
		AgreementAccepted: true,
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !output.Identity.IsNewUser {
		t.Fatalf("expected new user")
	}
	if output.Identity.UserID != "user_001" {
		t.Fatalf("unexpected user id: %s", output.Identity.UserID)
	}
	if output.Identity.LoginChannel != "sms" {
		t.Fatalf("unexpected channel: %s", output.Identity.LoginChannel)
	}
}

func TestLoginBySMSUseCase_ExistingUser(t *testing.T) {
	repo := memory.NewIdentityRepository()
	auth := infraIdentity.NewMockAuthService()
	uc := NewLoginBySMSUseCase(repo, auth, fixedClock{}, fixedIDGenerator{})

	_, err := uc.Execute(context.Background(), LoginBySMSInput{
		Phone:             "13800000000",
		VerifyCode:        "123",
		AgreementAccepted: true,
	})
	if err != nil {
		t.Fatalf("setup login failed: %v", err)
	}

	output, err := uc.Execute(context.Background(), LoginBySMSInput{
		Phone:             "13800000000",
		VerifyCode:        "123",
		AgreementAccepted: true,
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if output.Identity.IsNewUser {
		t.Fatalf("expected existing user")
	}
}

func TestLoginBySMSUseCase_InvalidCode(t *testing.T) {
	repo := memory.NewIdentityRepository()
	auth := infraIdentity.NewMockAuthService()
	uc := NewLoginBySMSUseCase(repo, auth, fixedClock{}, fixedIDGenerator{})

	_, err := uc.Execute(context.Background(), LoginBySMSInput{
		Phone:             "13800000000",
		VerifyCode:        "000000",
		AgreementAccepted: true,
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestLoginByWechatUseCase_NewUser(t *testing.T) {
	repo := memory.NewIdentityRepository()
	auth := infraIdentity.NewMockAuthService()
	uc := NewLoginByWechatUseCase(repo, auth, fixedClock{}, fixedIDGenerator{})

	output, err := uc.Execute(context.Background(), LoginByWechatInput{
		AuthCode:          "wechat_code_001",
		AgreementAccepted: true,
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !output.Identity.IsNewUser {
		t.Fatalf("expected new user")
	}
	if output.Identity.LoginChannel != "wechat" {
		t.Fatalf("unexpected channel: %s", output.Identity.LoginChannel)
	}
}
