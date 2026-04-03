package ai

import (
	"context"
	"time"
)

type HomeOverviewServiceInput struct {
	UserID    string
	Now       time.Time
	Remaining int
}

type HomeOverviewService interface {
	BuildOverview(ctx context.Context, input HomeOverviewServiceInput) (HomeOverview, error)
}

type MatchService interface {
	GenerateCandidates(ctx context.Context, inputText string) (MatchResult, error)
}

type ReplyService interface {
	GenerateReply(ctx context.Context, session Session, userMessage string) (string, error)
}
