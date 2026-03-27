package ai

import (
	"context"
	"math/rand"
	"strings"
	"time"

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

	pool := []domain.MatchCandidate{
		{ProviderID: "p_001", DisplayName: "暖心倾听师-小林", ReasonText: "适合情绪陪伴和低压对话", Score: 0.93},
		{ProviderID: "p_002", DisplayName: "夜谈伙伴-阿泽", ReasonText: "夜间响应较快，偏聊天疏导", Score: 0.88},
		{ProviderID: "p_003", DisplayName: "同城散步搭子-念念", ReasonText: "适合轻社交场景和线下陪伴", Score: 0.84},
		{ProviderID: "p_004", DisplayName: "白噪音向导-阿乔", ReasonText: "关注睡眠与放松场景", Score: 0.86},
		{ProviderID: "p_005", DisplayName: "关系倾听师-木子", ReasonText: "擅长关系冲突后的情绪梳理", Score: 0.9},
		{ProviderID: "p_006", DisplayName: "成长陪伴师-周周", ReasonText: "适合目标焦虑与节奏管理话题", Score: 0.87},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})

	return domain.MatchResult{Candidates: pool[:3]}, nil
}
