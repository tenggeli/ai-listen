package ai

import (
	"context"
	"testing"
	"time"

	domain "listen/backend/internal/domain/ai"
	infraAi "listen/backend/internal/infrastructure/ai"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

type fixedClock struct{}

func (fixedClock) Now() time.Time {
	return time.Date(2026, 3, 28, 10, 0, 0, 0, time.UTC)
}

func TestSubmitMatchUseCase_ConsumeQuotaOnSuccess(t *testing.T) {
	repo := memory.NewMatchQuotaRepository()
	uc := NewSubmitMatchUseCase(repo, infraAi.NewMockMatchService(), fixedClock{})

	output, err := uc.Execute(context.Background(), SubmitMatchInput{UserID: "u1", InputText: "我有点焦虑"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Remaining != domain.DailyMatchLimit-1 {
		t.Fatalf("unexpected remaining: %d", output.Remaining)
	}
}

func TestSubmitMatchUseCase_NoConsumeOnEmptyResult(t *testing.T) {
	repo := memory.NewMatchQuotaRepository()
	uc := NewSubmitMatchUseCase(repo, infraAi.NewMockMatchService(), fixedClock{})

	output, err := uc.Execute(context.Background(), SubmitMatchInput{UserID: "u1", InputText: "   "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Remaining != domain.DailyMatchLimit {
		t.Fatalf("quota should not be consumed when result is empty, got: %d", output.Remaining)
	}
}

func TestGetRemainingMatchUseCase_DefaultFive(t *testing.T) {
	repo := memory.NewMatchQuotaRepository()
	uc := NewGetRemainingMatchUseCase(repo, fixedClock{})

	output, err := uc.Execute(context.Background(), GetRemainingMatchInput{UserID: "u1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Remaining != domain.DailyMatchLimit {
		t.Fatalf("unexpected remaining: %d", output.Remaining)
	}
}
