package ai

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
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

func TestGetAiHomeUseCase_BuildOverview(t *testing.T) {
	repo := memory.NewMatchQuotaRepository()
	uc := NewGetAiHomeUseCase(repo, infraAi.NewMockHomeOverviewService(), fixedClock{})

	output, err := uc.Execute(context.Background(), GetHomeOverviewInput{UserID: "u1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Overview.GreetingPeriod != "早上好 · 周六" {
		t.Fatalf("unexpected greeting period: %s", output.Overview.GreetingPeriod)
	}
	if output.Overview.Remaining != domain.DailyMatchLimit {
		t.Fatalf("unexpected remaining: %d", output.Overview.Remaining)
	}
	if len(output.Overview.QuickActions) != 4 {
		t.Fatalf("unexpected quick action count: %d", len(output.Overview.QuickActions))
	}
}

func TestAiSessionUseCases_CreateAppendAndGet(t *testing.T) {
	repo := memory.NewSessionRepository()
	create := NewCreateAiSessionUseCase(repo, fixedIDGenerator{})
	appendMsg := NewAppendAiMessageUseCase(repo, NewMockReplyService(), fixedClock{})
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

	if len(sessionOutput.Session.Messages) != 2 {
		t.Fatalf("unexpected message count: %d", len(sessionOutput.Session.Messages))
	}
	if sessionOutput.Session.Messages[0].SenderType != "user" {
		t.Fatalf("first message should be user")
	}
	if sessionOutput.Session.Messages[1].SenderType != "assistant" {
		t.Fatalf("second message should be assistant")
	}
	if sessionOutput.Session.Messages[0].Content != "今晚有点焦虑" {
		t.Fatalf("unexpected message content: %s", sessionOutput.Session.Messages[0].Content)
	}
}

func TestAppendAiMessageUseCase_EmptyMessageValidation(t *testing.T) {
	repo := memory.NewSessionRepository()
	create := NewCreateAiSessionUseCase(repo, fixedIDGenerator{})
	appendMsg := NewAppendAiMessageUseCase(repo, NewMockReplyService(), fixedClock{})

	created, err := create.Execute(context.Background(), CreateSessionInput{UserID: "u1", SceneType: "listen"})
	if err != nil {
		t.Fatalf("create session failed: %v", err)
	}

	_, err = appendMsg.Execute(context.Background(), AppendMessageInput{
		SessionID:   created.SessionID,
		SenderType:  "user",
		ContentText: "   ",
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got: %v", err)
	}
}

func TestAppendAiMessageUseCase_ReplyTimeoutAndRetryFallback(t *testing.T) {
	repo := memory.NewSessionRepository()
	create := NewCreateAiSessionUseCase(repo, fixedIDGenerator{})
	replySvc := &slowReplyService{delay: 40 * time.Millisecond}
	appendMsg := NewAppendAiMessageUseCase(repo, replySvc, fixedClock{})
	appendMsg.replyCfg = replyConfig{
		timeout:    5 * time.Millisecond,
		maxRetries: 2,
	}

	created, err := create.Execute(context.Background(), CreateSessionInput{UserID: "u1", SceneType: "listen"})
	if err != nil {
		t.Fatalf("create session failed: %v", err)
	}

	output, err := appendMsg.Execute(context.Background(), AppendMessageInput{
		SessionID:   created.SessionID,
		SenderType:  "user",
		ContentText: "最近压力很大",
	})
	if err != nil {
		t.Fatalf("append should not fail on reply timeout, got: %v", err)
	}
	if got := replySvc.calls.Load(); got != 3 {
		t.Fatalf("expected 3 attempts, got: %d", got)
	}
	if len(output.Session.Messages) != 2 {
		t.Fatalf("unexpected message count: %d", len(output.Session.Messages))
	}
	if output.Session.Messages[1].Content != buildFallbackReply("最近压力很大") {
		t.Fatalf("expected fallback reply, got: %s", output.Session.Messages[1].Content)
	}
}

func TestAppendAiMessageUseCase_ConcurrentSend(t *testing.T) {
	repo := memory.NewSessionRepository()
	create := NewCreateAiSessionUseCase(repo, fixedIDGenerator{})
	appendMsg := NewAppendAiMessageUseCase(repo, NewMockReplyService(), fixedClock{})
	get := NewGetAiSessionUseCase(repo)

	created, err := create.Execute(context.Background(), CreateSessionInput{UserID: "u1", SceneType: "listen"})
	if err != nil {
		t.Fatalf("create session failed: %v", err)
	}

	messages := []string{"第一条", "第二条", "第三条"}
	var wg sync.WaitGroup
	wg.Add(len(messages))
	for _, content := range messages {
		content := content
		go func() {
			defer wg.Done()
			if _, err := appendMsg.Execute(context.Background(), AppendMessageInput{
				SessionID:   created.SessionID,
				SenderType:  "user",
				ContentText: content,
			}); err != nil {
				t.Errorf("append failed: %v", err)
			}
		}()
	}
	wg.Wait()

	sessionOutput, err := get.Execute(context.Background(), GetSessionInput{SessionID: created.SessionID})
	if err != nil {
		t.Fatalf("get session failed: %v", err)
	}
	if len(sessionOutput.Session.Messages) != len(messages)*2 {
		t.Fatalf("unexpected message count: %d", len(sessionOutput.Session.Messages))
	}

	for i := 0; i < len(sessionOutput.Session.Messages); i += 2 {
		if sessionOutput.Session.Messages[i].SenderType != "user" {
			t.Fatalf("message %d should be user, got: %s", i, sessionOutput.Session.Messages[i].SenderType)
		}
		if sessionOutput.Session.Messages[i+1].SenderType != "assistant" {
			t.Fatalf("message %d should be assistant, got: %s", i+1, sessionOutput.Session.Messages[i+1].SenderType)
		}
	}
}

type slowReplyService struct {
	delay time.Duration
	calls atomic.Int32
}

func (s *slowReplyService) GenerateReply(ctx context.Context, _ domain.Session, _ string) (string, error) {
	s.calls.Add(1)
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(s.delay):
		return "ok", nil
	}
}
