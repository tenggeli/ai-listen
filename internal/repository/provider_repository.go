package repository

import (
	"context"

	"ai-listen/internal/model"
	"gorm.io/gorm"
)

type ProviderRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{db: db}
}

func (r *ProviderRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Provider, error) {
	var provider model.Provider
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&provider).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &provider, nil
}

func (r *ProviderRepository) GetByID(ctx context.Context, id uint64) (*model.Provider, error) {
	var provider model.Provider
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&provider).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &provider, nil
}

func (r *ProviderRepository) ListByCity(ctx context.Context, cityID uint64, offset, limit int) ([]model.Provider, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Provider{})
	if cityID > 0 {
		query = query.Where("city_id = ?", cityID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []model.Provider{}, 0, nil
	}

	var providers []model.Provider
	if err := query.
		Order("score DESC").
		Order("total_orders DESC").
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&providers).Error; err != nil {
		return nil, 0, err
	}
	return providers, total, nil
}

func (r *ProviderRepository) ListMatchCandidates(ctx context.Context, cityID uint64, limit int) ([]model.Provider, error) {
	var providers []model.Provider
	err := r.db.WithContext(ctx).
		Where("city_id = ? AND service_status = ? AND online_status = ?", cityID, 1, 1).
		Order("score DESC").
		Order("total_orders DESC").
		Order("id DESC").
		Limit(limit).
		Find(&providers).Error
	if err != nil {
		return nil, err
	}
	return providers, nil
}

func (r *ProviderRepository) Create(ctx context.Context, provider *model.Provider) error {
	return r.db.WithContext(ctx).Create(provider).Error
}

func (r *ProviderRepository) UpdateByUserID(ctx context.Context, userID uint64, updates map[string]any) error {
	return r.db.WithContext(ctx).
		Model(&model.Provider{}).
		Where("user_id = ?", userID).
		Updates(updates).Error
}
