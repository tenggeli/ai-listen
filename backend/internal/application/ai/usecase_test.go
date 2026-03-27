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

type fixedIDGenerator struct{}

func (fixedIDGenerator) NewID(prefix string) string {
	return prefix + "_001"
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
	if len(output.Candidates) != 3 {
		t.Fatalf("unexpected candidate count: %d", len(output.Candidates))
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

func TestAiSessionUseCases_CreateAppendAndGet(t *testing.T) {
	repo := memory.NewSessionRepository()
	create := NewCreateAiSessionUseCase(repo, fixedIDGenerator{})
	appendMsg := NewAppendAiMessageUseCase(repo, fixedClock{})
	get := NewGetAiSessionUseCase(repo)

	created, err := create.Execute(context.Background(), CreateSessionInput{UserID: "u1", SceneType: "listen"})
	if err != nil {
		t.Fatalf("create session failed: %v", err)
	}
	if created.SessionID != "sess_001" {
		t.Fatalf("unexpected session id: %s", created.SessionID)
	}

	if _, err := appendMsg.Execute(context.Background(), AppendMessageInput{
		SessionID:   created.SessionID,
		SenderType:  "user",
		ContentText: "今晚有点焦虑",
	}); err != nil {
		t.Fatalf("append message failed: %v", err)
	}

	sessionOutput, err := get.Execute(context.Background(), GetSessionInput{SessionID: created.SessionID})
	if err != nil {
		t.Fatalf("get session failed: %v", err)
	}

	if len(sessionOutput.Session.Messages) != 1 {
		t.Fatalf("unexpected message count: %d", len(sessionOutput.Session.Messages))
	}
	if sessionOutput.Session.Messages[0].Content != "今晚有点焦虑" {
		t.Fatalf("unexpected message content: %s", sessionOutput.Session.Messages[0].Content)
	}
}
