package service

import (
	"context"

	"ai-listen/internal/dto"
	"ai-listen/internal/model"
	"ai-listen/internal/pkg/apperror"
	"ai-listen/internal/repository"
)

type ProviderService struct {
	providerRepo    *repository.ProviderRepository
	userRepo        *repository.UserRepository
	serviceItemRepo *repository.ServiceItemRepository
}

func NewProviderService(
	providerRepo *repository.ProviderRepository,
	userRepo *repository.UserRepository,
	serviceItemRepo *repository.ServiceItemRepository,
) *ProviderService {
	return &ProviderService{
		providerRepo:    providerRepo,
		userRepo:        userRepo,
		serviceItemRepo: serviceItemRepo,
	}
}

func (s *ProviderService) List(ctx context.Context, cityID uint64, page, pageSize int) (*dto.ProviderListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	providers, total, err := s.providerRepo.ListByCity(ctx, cityID, offset, pageSize)
	if err != nil {
		return nil, apperror.Internal("获取服务方列表失败", err)
	}
	if len(providers) == 0 {
		return &dto.ProviderListResponse{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			List:     []dto.ProviderListItem{},
		}, nil
	}

	userMap, itemMap, err := s.batchLoadRelatedData(ctx, providers)
	if err != nil {
		return nil, err
	}

	list := make([]dto.ProviderListItem, 0, len(providers))
	for _, provider := range providers {
		tags, err := parseJSONStrings(provider.Tags)
		if err != nil {
			return nil, apperror.Internal("获取服务方列表失败", err)
		}

		list = append(list, dto.ProviderListItem{
			ProviderID:    provider.ID,
			ProviderNo:    provider.ProviderNo,
			Score:         provider.Score,
			TotalOrders:   provider.TotalOrders,
			ServiceStatus: provider.ServiceStatus,
			OnlineStatus:  provider.OnlineStatus,
			CityID:        provider.CityID,
			Intro:         provider.Intro,
			Tags:          tags,
			User:          buildUserBrief(userMap[provider.UserID]),
			ServiceItems:  itemMap[provider.ID],
		})
	}

	return &dto.ProviderListResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		List:     list,
	}, nil
}

func (s *ProviderService) Detail(ctx context.Context, id uint64) (*dto.ProviderDetailResponse, error) {
	if id == 0 {
		return nil, apperror.BadRequest("id 参数错误")
	}

	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("获取服务方详情失败", err)
	}
	if provider == nil {
		return nil, apperror.NotFound("服务方不存在")
	}

	user, err := s.userRepo.GetByID(ctx, provider.UserID)
	if err != nil {
		return nil, apperror.Internal("获取服务方详情失败", err)
	}
	if user == nil {
		return nil, apperror.NotFound("关联用户不存在")
	}

	items, err := s.serviceItemRepo.ListActiveByProviderID(ctx, provider.ID)
	if err != nil {
		return nil, apperror.Internal("获取服务方详情失败", err)
	}

	tags, err := parseJSONStrings(provider.Tags)
	if err != nil {
		return nil, apperror.Internal("获取服务方详情失败", err)
	}

	serviceItems := make([]dto.ServiceItemInfo, 0, len(items))
	for _, item := range items {
		serviceItems = append(serviceItems, buildServiceItemInfo(item))
	}

	return &dto.ProviderDetailResponse{
		Provider: dto.ProviderBaseInfo{
			ID:             provider.ID,
			UserID:         provider.UserID,
			ProviderNo:     provider.ProviderNo,
			Zodiac:         provider.Zodiac,
			Constellation:  provider.Constellation,
			Level:          provider.Level,
			Score:          provider.Score,
			TotalOrders:    provider.TotalOrders,
			TotalIncome:    provider.TotalIncome,
			ServiceStatus:  provider.ServiceStatus,
			OnlineStatus:   provider.OnlineStatus,
			CityID:         provider.CityID,
			Intro:          provider.Intro,
			Tags:           tags,
			CommissionRate: provider.CommissionRate,
			ComplaintCount: provider.ComplaintCount,
			CancelCount:    provider.CancelCount,
		},
		User:         buildUserBrief(*user),
		ServiceItems: serviceItems,
	}, nil
}

func (s *ProviderService) batchLoadRelatedData(
	ctx context.Context,
	providers []model.Provider,
) (map[uint64]model.User, map[uint64][]dto.ServiceItemInfo, error) {
	userIDs := make([]uint64, 0, len(providers))
	providerIDs := make([]uint64, 0, len(providers))
	for _, provider := range providers {
		userIDs = append(userIDs, provider.UserID)
		providerIDs = append(providerIDs, provider.ID)
	}

	users, err := s.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, nil, apperror.Internal("获取服务方列表失败", err)
	}
	userMap := make(map[uint64]model.User, len(users))
	for _, user := range users {
		userMap[user.ID] = user
	}

	items, err := s.serviceItemRepo.ListActiveByProviderIDs(ctx, providerIDs)
	if err != nil {
		return nil, nil, apperror.Internal("获取服务方列表失败", err)
	}
	itemMap := make(map[uint64][]dto.ServiceItemInfo)
	for _, item := range items {
		itemMap[item.ProviderID] = append(itemMap[item.ProviderID], buildServiceItemInfo(item))
	}

	return userMap, itemMap, nil
}

func buildUserBrief(user model.User) dto.ProviderUserBrief {
	return dto.ProviderUserBrief{
		ID:       user.ID,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		CityID:   user.CityID,
		Bio:      user.Bio,
		VipLevel: user.VipLevel,
	}
}

func buildServiceItemInfo(item model.ServiceItem) dto.ServiceItemInfo {
	return dto.ServiceItemInfo{
		ID:          item.ID,
		ProviderID:  item.ProviderID,
		CategoryID:  item.CategoryID,
		Title:       item.Title,
		Description: item.Description,
		UnitPrice:   item.UnitPrice,
		BillingType: item.BillingType,
		MinHours:    item.MinHours,
		MaxHours:    item.MaxHours,
		UnitName:    item.UnitName,
		Status:      item.Status,
	}
}
