package service

import (
	"context"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"ai-listen/internal/dto"
	"ai-listen/internal/model"
	"ai-listen/internal/pkg/apperror"
	"ai-listen/internal/repository"
)

type ProviderCenterService struct {
	userRepo        *repository.UserRepository
	applicationRepo *repository.ProviderApplicationRepository
	providerRepo    *repository.ProviderRepository
}

func NewProviderCenterService(
	userRepo *repository.UserRepository,
	applicationRepo *repository.ProviderApplicationRepository,
	providerRepo *repository.ProviderRepository,
) *ProviderCenterService {
	return &ProviderCenterService{
		userRepo:        userRepo,
		applicationRepo: applicationRepo,
		providerRepo:    providerRepo,
	}
}

func (s *ProviderCenterService) Apply(ctx context.Context, userID uint64, req dto.ApplyProviderRequest) (*dto.ApplyProviderResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("提交入驻申请失败", err)
	}
	if user == nil {
		return nil, apperror.NotFound("用户不存在")
	}
	if !req.AgreementSigned {
		return nil, apperror.BadRequest("请先签署平台协议")
	}

	latest, err := s.applicationRepo.GetLatestByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("提交入驻申请失败", err)
	}
	if latest != nil && latest.AuditStatus == 0 {
		return nil, apperror.BadRequest("已有待审核申请，请勿重复提交")
	}

	photosJSON, err := json.Marshal(req.Photos)
	if err != nil {
		return nil, apperror.BadRequest("照片格式错误")
	}

	now := time.Now()
	app := &model.ProviderApplication{
		UserID:           userID,
		RealName:         req.RealName,
		IDCardNo:         req.IDCardNo,
		IDCardFront:      req.IDCardFront,
		IDCardBack:       req.IDCardBack,
		FaceVerifyStatus: req.FaceVerifyStatus,
		AgreementSigned:  1,
		CityID:           req.CityID,
		Intro:            req.Intro,
		Photos:           string(photosJSON),
		ServiceDesc:      req.ServiceDesc,
		AuditStatus:      0,
		SubmittedAt:      &now,
	}
	if err := s.applicationRepo.Create(ctx, app); err != nil {
		return nil, apperror.Internal("提交入驻申请失败", err)
	}

	return &dto.ApplyProviderResponse{
		ApplicationID: app.ID,
		AuditStatus:   app.AuditStatus,
	}, nil
}

func (s *ProviderCenterService) AuditStatus(ctx context.Context, userID uint64) (*dto.AuditStatusResponse, error) {
	latest, err := s.applicationRepo.GetLatestByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("获取审核状态失败", err)
	}
	provider, err := s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("获取审核状态失败", err)
	}

	resp := &dto.AuditStatusResponse{
		HasApplication: latest != nil,
		AuditStatus:    -1,
		ProviderExists: provider != nil,
	}
	if latest == nil {
		return resp, nil
	}

	resp.AuditStatus = int(latest.AuditStatus)
	resp.RejectReason = latest.RejectReason
	if latest.SubmittedAt != nil {
		resp.SubmittedAt = latest.SubmittedAt.Format(time.RFC3339)
	}
	if latest.AuditedAt != nil {
		resp.AuditedAt = latest.AuditedAt.Format(time.RFC3339)
	}
	return resp, nil
}

func (s *ProviderCenterService) UpdateProfile(ctx context.Context, userID uint64, req dto.UpdateProviderProfileRequest) (*dto.ProviderProfileResponse, error) {
	if req.ServiceStatus != nil && (*req.ServiceStatus < 1 || *req.ServiceStatus > 3) {
		return nil, apperror.BadRequest("serviceStatus 仅支持 1-3")
	}
	if req.OnlineStatus != nil && *req.OnlineStatus > 1 {
		return nil, apperror.BadRequest("onlineStatus 仅支持 0-1")
	}

	provider, err := s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("更新服务方资料失败", err)
	}

	if provider == nil {
		latest, err := s.applicationRepo.GetLatestByUserID(ctx, userID)
		if err != nil {
			return nil, apperror.Internal("更新服务方资料失败", err)
		}
		if latest == nil || latest.AuditStatus != 1 {
			return nil, apperror.Forbidden("未通过服务方审核，无法更新资料")
		}

		providerNo, err := generateProviderNo()
		if err != nil {
			return nil, apperror.Internal("更新服务方资料失败", err)
		}

		provider = &model.Provider{
			UserID:        userID,
			ProviderNo:    providerNo,
			ServiceStatus: 3,
			OnlineStatus:  0,
		}
		if err := s.providerRepo.Create(ctx, provider); err != nil {
			return nil, apperror.Internal("更新服务方资料失败", err)
		}
	}

	updates := make(map[string]any)
	if req.Zodiac != nil {
		updates["zodiac"] = *req.Zodiac
	}
	if req.Constellation != nil {
		updates["constellation"] = *req.Constellation
	}
	if req.Level != nil {
		updates["level"] = *req.Level
	}
	if req.ServiceStatus != nil {
		updates["service_status"] = *req.ServiceStatus
	}
	if req.OnlineStatus != nil {
		updates["online_status"] = *req.OnlineStatus
	}
	if req.CityID != nil {
		updates["city_id"] = *req.CityID
	}
	if req.Intro != nil {
		updates["intro"] = *req.Intro
	}
	if req.Tags != nil {
		tagBytes, err := json.Marshal(*req.Tags)
		if err != nil {
			return nil, apperror.BadRequest("tags 格式错误")
		}
		updates["tags"] = string(tagBytes)
	}
	if len(updates) == 0 {
		return nil, apperror.BadRequest("没有可更新字段")
	}

	updates["updated_at"] = time.Now()
	if err := s.providerRepo.UpdateByUserID(ctx, userID, updates); err != nil {
		return nil, apperror.Internal("更新服务方资料失败", err)
	}

	provider, err = s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("更新服务方资料失败", err)
	}
	if provider == nil {
		return nil, apperror.NotFound("服务方不存在")
	}

	tags, err := parseJSONStrings(provider.Tags)
	if err != nil {
		return nil, apperror.Internal("更新服务方资料失败", err)
	}

	return &dto.ProviderProfileResponse{
		ID:            provider.ID,
		UserID:        provider.UserID,
		ProviderNo:    provider.ProviderNo,
		Zodiac:        provider.Zodiac,
		Constellation: provider.Constellation,
		Level:         provider.Level,
		ServiceStatus: provider.ServiceStatus,
		OnlineStatus:  provider.OnlineStatus,
		CityID:        provider.CityID,
		Intro:         provider.Intro,
		Tags:          tags,
	}, nil
}

func parseJSONStrings(raw string) ([]string, error) {
	if raw == "" {
		return []string{}, nil
	}
	var values []string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return nil, err
	}
	return values, nil
}

func generateProviderNo() (string, error) {
	n, err := crand.Int(crand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("P%s%04d", time.Now().Format("20060102"), n.Int64()), nil
}
