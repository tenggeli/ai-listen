package ai

import "context"

type MatchService interface {
	GenerateCandidates(ctx context.Context, inputText string) (MatchResult, error)
}
