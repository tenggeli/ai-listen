package service

import (
	"context"

	"ai-listen/internal/dto"
	"ai-listen/internal/pkg/apperror"
	"ai-listen/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(ctx context.Context, userID uint64) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("获取用户资料失败", err)
	}
	if user == nil {
		return nil, apperror.NotFound("用户不存在")
	}

	return &dto.UserProfileResponse{
		ID:         user.ID,
		Mobile:     user.Mobile,
		Nickname:   user.Nickname,
		Avatar:     user.Avatar,
		Gender:     user.Gender,
		CityID:     user.CityID,
		Bio:        user.Bio,
		VipLevel:   user.VipLevel,
		UserStatus: user.UserStatus,
	}, nil
}
