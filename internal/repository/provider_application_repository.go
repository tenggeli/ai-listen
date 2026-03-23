package repository

import (
	"context"

	"ai-listen/internal/model"
	"gorm.io/gorm"
)

type ProviderApplicationRepository struct {
	db *gorm.DB
}

func NewProviderApplicationRepository(db *gorm.DB) *ProviderApplicationRepository {
	return &ProviderApplicationRepository{db: db}
}

func (r *ProviderApplicationRepository) Create(ctx context.Context, application *model.ProviderApplication) error {
	return r.db.WithContext(ctx).Create(application).Error
}

func (r *ProviderApplicationRepository) GetLatestByUserID(ctx context.Context, userID uint64) (*model.ProviderApplication, error) {
	var app model.ProviderApplication
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id DESC").
		First(&app).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}
