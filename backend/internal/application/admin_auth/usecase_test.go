package admin_auth

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "listen/backend/internal/domain/admin_auth"
)

type fixedClock struct{}

func (fixedClock) Now() time.Time {
	return time.Date(2026, 4, 2, 9, 0, 0, 0, time.UTC)
}

func seedRepo() InMemoryRepository {
	return NewInMemoryRepository([]domain.AdminAccount{
		{
			AdminID:     "admin_001",
			Account:     "admin",
			Password:    "admin123",
			Role:        "super_admin",
			DisplayName: "平台管理员",
			Status:      "active",
		},
	})
}

func TestLoginMockUseCase_Success(t *testing.T) {
	uc := NewLoginMockUseCase(seedRepo(), fixedClock{})
	output, err := uc.Execute(context.Background(), LoginMockInput{
		Account:  "admin",
		Password: "admin123",
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if output.Identity.AdminID != "admin_001" {
		t.Fatalf("unexpected admin id: %s", output.Identity.AdminID)
	}
	if output.Identity.AccessToken == "" {
		t.Fatalf("expected token")
	}
}

func TestLoginMockUseCase_InvalidCredential(t *testing.T) {
	uc := NewLoginMockUseCase(seedRepo(), fixedClock{})
	_, err := uc.Execute(context.Background(), LoginMockInput{
		Account:  "admin",
		Password: "wrong",
	})
	if !errors.Is(err, domain.ErrInvalidCredential) {
		t.Fatalf("expected invalid credential, got: %v", err)
	}
}

func TestGetCurrentAdminUseCase_Success(t *testing.T) {
	uc := NewGetCurrentAdminUseCase(seedRepo())
	output, err := uc.Execute(context.Background(), GetCurrentAdminInput{AdminID: "admin_001"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if output.Admin.Role != "super_admin" {
		t.Fatalf("unexpected role: %s", output.Admin.Role)
	}
}
