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

func (r *ProviderRepository) Create(ctx context.Context, provider *model.Provider) error {
	return r.db.WithContext(ctx).Create(provider).Error
}

func (r *ProviderRepository) UpdateByUserID(ctx context.Context, userID uint64, updates map[string]any) error {
	return r.db.WithContext(ctx).
		Model(&model.Provider{}).
		Where("user_id = ?", userID).
		Updates(updates).Error
}
