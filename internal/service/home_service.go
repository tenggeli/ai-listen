package service

import (
	"context"
	"fmt"
	"strings"

	"ai-listen/internal/dto"
	"ai-listen/internal/model"
	"ai-listen/internal/pkg/apperror"
	"ai-listen/internal/repository"
)

type HomeService struct {
	providerRepo *repository.ProviderRepository
	userRepo     *repository.UserRepository
}

func NewHomeService(providerRepo *repository.ProviderRepository, userRepo *repository.UserRepository) *HomeService {
	return &HomeService{
		providerRepo: providerRepo,
		userRepo:     userRepo,
	}
}

func (s *HomeService) AIMatch(ctx context.Context, req dto.AIMatchRequest) (*dto.AIMatchResponse, error) {
	if req.CityID == 0 {
		return nil, apperror.BadRequest("cityId 不能为空")
	}

	providers, err := s.providerRepo.ListMatchCandidates(ctx, req.CityID, 3)
	if err != nil {
		return nil, apperror.Internal("AI 匹配失败", err)
	}
	if len(providers) == 0 {
		return &dto.AIMatchResponse{List: []dto.AIMatchItem{}}, nil
	}

	userIDs := make([]uint64, 0, len(providers))
	for _, provider := range providers {
		userIDs = append(userIDs, provider.UserID)
	}

	users, err := s.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, apperror.Internal("AI 匹配失败", err)
	}

	userMap := make(map[uint64]model.User, len(users))
	for _, user := range users {
		userMap[user.ID] = user
	}

	items := make([]dto.AIMatchItem, 0, len(providers))
	for _, provider := range providers {
		nickname := fmt.Sprintf("服务方%v", provider.ID)
		if user, ok := userMap[provider.UserID]; ok && strings.TrimSpace(user.Nickname) != "" {
			nickname = user.Nickname
		}

		items = append(items, dto.AIMatchItem{
			ProviderID:  provider.ID,
			Nickname:    nickname,
			Score:       provider.Score,
			CityID:      provider.CityID,
			MatchReason: buildMatchReason(provider),
		})
	}

	return &dto.AIMatchResponse{List: items}, nil
}

func buildMatchReason(provider model.Provider) string {
	reasons := []string{"同城且当前在线"}

	switch {
	case provider.Score >= 4.8:
		reasons = append(reasons, "评分较高")
	case provider.TotalOrders >= 50:
		reasons = append(reasons, "服务经验较丰富")
	default:
		reasons = append(reasons, "可快速响应")
	}

	return strings.Join(reasons, "，")
}
