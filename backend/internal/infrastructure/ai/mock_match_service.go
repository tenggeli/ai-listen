package ai

import (
	"context"
	"strings"

	domain "listen/backend/internal/domain/ai"
)

type MockMatchService struct{}

func NewMockMatchService() MockMatchService {
	return MockMatchService{}
}

func (MockMatchService) GenerateCandidates(_ context.Context, inputText string) (domain.MatchResult, error) {
	text := strings.TrimSpace(inputText)
	if text == "" {
		return domain.MatchResult{Candidates: []domain.MatchCandidate{}}, nil
	}

	result := []domain.MatchCandidate{
		{ProviderID: "p_001", DisplayName: "暖心倾听师-小林", ReasonText: "适合情绪陪伴和低压对话", Score: 0.93},
		{ProviderID: "p_002", DisplayName: "夜谈伙伴-阿泽", ReasonText: "夜间响应较快，偏聊天疏导", Score: 0.88},
		{ProviderID: "p_003", DisplayName: "同城散步搭子-念念", ReasonText: "适合轻社交场景和线下陪伴", Score: 0.84},
	}

	if strings.Contains(text, "音乐") || strings.Contains(text, "睡眠") {
		result[2] = domain.MatchCandidate{ProviderID: "p_004", DisplayName: "白噪音向导-阿乔", ReasonText: "关注睡眠与放松场景", Score: 0.86}
	}

	return domain.MatchResult{Candidates: result}, nil
}
