package ai

import (
	"context"
	"fmt"
	"time"

	domain "listen/backend/internal/domain/ai"
)

type Clock interface {
	Now() time.Time
}

type IDGenerator interface {
	NewID(prefix string) string
}

type GetRemainingMatchInput struct {
	UserID string
}

type GetHomeOverviewInput struct {
	UserID string
}

type GetHomeOverviewOutput struct {
	Overview domain.HomeOverview
}

type GetRemainingMatchOutput struct {
	Remaining int
}

type SubmitMatchInput struct {
	UserID    string
	InputText string
}

type SubmitMatchOutput struct {
	Remaining  int
	Candidates []domain.MatchCandidate
}

type CreateSessionInput struct {
	UserID    string
	SceneType string
}

type CreateSessionOutput struct {
	SessionID string
}

type GetSessionInput struct {
	SessionID string
}

type GetSessionOutput struct {
	Session domain.Session
}

type AppendMessageInput struct {
	SessionID   string
	SenderType  string
	ContentText string
}

type AppendMessageOutput struct {
	Session domain.Session
}

type GetRemainingMatchUseCase struct {
	quotaRepo domain.MatchQuotaRepository
	clock     Clock
}

type GetAiHomeUseCase struct {
	quotaRepo    domain.MatchQuotaRepository
	homeOverview domain.HomeOverviewService
	clock        Clock
}

func NewGetRemainingMatchUseCase(quotaRepo domain.MatchQuotaRepository, clock Clock) GetRemainingMatchUseCase {
	return GetRemainingMatchUseCase{quotaRepo: quotaRepo, clock: clock}
}

func NewGetAiHomeUseCase(quotaRepo domain.MatchQuotaRepository, homeOverview domain.HomeOverviewService, clock Clock) GetAiHomeUseCase {
	return GetAiHomeUseCase{quotaRepo: quotaRepo, homeOverview: homeOverview, clock: clock}
}

func (u GetRemainingMatchUseCase) Execute(ctx context.Context, input GetRemainingMatchInput) (GetRemainingMatchOutput, error) {
	if input.UserID == "" {
		return GetRemainingMatchOutput{}, domain.ErrInvalidInput
	}
	quota, err := u.quotaRepo.GetByUserAndDate(ctx, input.UserID, dateKey(u.clock.Now()))
	if err != nil {
		return GetRemainingMatchOutput{}, err
	}
	return GetRemainingMatchOutput{Remaining: quota.Remaining()}, nil
}

func (u GetAiHomeUseCase) Execute(ctx context.Context, input GetHomeOverviewInput) (GetHomeOverviewOutput, error) {
	if input.UserID == "" {
		return GetHomeOverviewOutput{}, domain.ErrInvalidInput
	}

	now := u.clock.Now()
	quota, err := u.quotaRepo.GetByUserAndDate(ctx, input.UserID, dateKey(now))
	if err != nil {
		return GetHomeOverviewOutput{}, err
	}

	overview, err := u.homeOverview.BuildOverview(ctx, domain.HomeOverviewServiceInput{
		UserID:    input.UserID,
		Now:       now,
		Remaining: quota.Remaining(),
	})
	if err != nil {
		return GetHomeOverviewOutput{}, err
	}
	return GetHomeOverviewOutput{Overview: overview}, nil
}

type SubmitMatchUseCase struct {
	quotaRepo    domain.MatchQuotaRepository
	matchService domain.MatchService
	clock        Clock
}

func NewSubmitMatchUseCase(quotaRepo domain.MatchQuotaRepository, matchService domain.MatchService, clock Clock) SubmitMatchUseCase {
	return SubmitMatchUseCase{quotaRepo: quotaRepo, matchService: matchService, clock: clock}
}

func (u SubmitMatchUseCase) Execute(ctx context.Context, input SubmitMatchInput) (SubmitMatchOutput, error) {
	if input.UserID == "" || input.InputText == "" {
		return SubmitMatchOutput{}, domain.ErrInvalidInput
	}

	date := dateKey(u.clock.Now())
	quota, err := u.quotaRepo.GetByUserAndDate(ctx, input.UserID, date)
	if err != nil {
		return SubmitMatchOutput{}, err
	}

	if quota.Remaining() == 0 {
		return SubmitMatchOutput{}, domain.ErrDailyLimitReached
	}

	result, err := u.matchService.GenerateCandidates(ctx, input.InputText)
	if err != nil {
		return SubmitMatchOutput{}, err
	}

	if !result.IsEmpty() {
		if err := quota.Consume(); err != nil {
			return SubmitMatchOutput{}, err
		}
		if err := u.quotaRepo.Save(ctx, quota); err != nil {
			return SubmitMatchOutput{}, err
		}
	}

	return SubmitMatchOutput{Remaining: quota.Remaining(), Candidates: result.Candidates}, nil
}

type CreateAiSessionUseCase struct {
	sessionRepo domain.SessionRepository
	idGenerator IDGenerator
}

func NewCreateAiSessionUseCase(sessionRepo domain.SessionRepository, idGenerator IDGenerator) CreateAiSessionUseCase {
	return CreateAiSessionUseCase{sessionRepo: sessionRepo, idGenerator: idGenerator}
}

func (u CreateAiSessionUseCase) Execute(ctx context.Context, input CreateSessionInput) (CreateSessionOutput, error) {
	sessionID := u.idGenerator.NewID("sess")
	session, err := domain.NewSession(sessionID, input.UserID, input.SceneType)
	if err != nil {
		return CreateSessionOutput{}, err
	}
	if err := u.sessionRepo.Create(ctx, session); err != nil {
		return CreateSessionOutput{}, err
	}
	return CreateSessionOutput{SessionID: sessionID}, nil
}

type GetAiSessionUseCase struct {
	sessionRepo domain.SessionRepository
}

func NewGetAiSessionUseCase(sessionRepo domain.SessionRepository) GetAiSessionUseCase {
	return GetAiSessionUseCase{sessionRepo: sessionRepo}
}

func (u GetAiSessionUseCase) Execute(ctx context.Context, input GetSessionInput) (GetSessionOutput, error) {
	session, err := u.sessionRepo.GetByID(ctx, input.SessionID)
	if err != nil {
		return GetSessionOutput{}, err
	}
	return GetSessionOutput{Session: session}, nil
}

type AppendAiMessageUseCase struct {
	sessionRepo domain.SessionRepository
	clock       Clock
}

func NewAppendAiMessageUseCase(sessionRepo domain.SessionRepository, clock Clock) AppendAiMessageUseCase {
	return AppendAiMessageUseCase{sessionRepo: sessionRepo, clock: clock}
}

func (u AppendAiMessageUseCase) Execute(ctx context.Context, input AppendMessageInput) (AppendMessageOutput, error) {
	session, err := u.sessionRepo.GetByID(ctx, input.SessionID)
	if err != nil {
		return AppendMessageOutput{}, err
	}
	if err := session.AppendMessage(input.SenderType, input.ContentText, u.clock.Now()); err != nil {
		return AppendMessageOutput{}, err
	}
	if err := u.sessionRepo.Save(ctx, session); err != nil {
		return AppendMessageOutput{}, err
	}
	return AppendMessageOutput{Session: session}, nil
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now()
}

type TimestampIDGenerator struct {
	clock Clock
}

func NewTimestampIDGenerator(clock Clock) TimestampIDGenerator {
	return TimestampIDGenerator{clock: clock}
}

func (g TimestampIDGenerator) NewID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, g.clock.Now().UnixNano())
}

func dateKey(now time.Time) string {
	return now.Format("2006-01-02")
}

func greetingPeriod(now time.Time) string {
	period := "晚上好"
	switch hour := now.Hour(); {
	case hour < 11:
		period = "早上好"
	case hour < 14:
		period = "中午好"
	case hour < 18:
		period = "下午好"
	}

	weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	return period + " · " + weekdays[int(now.Weekday())]
}
