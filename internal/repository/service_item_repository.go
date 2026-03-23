package repository

import (
	"context"

	"ai-listen/internal/model"
	"gorm.io/gorm"
)

type ServiceItemRepository struct {
	db *gorm.DB
}

func NewServiceItemRepository(db *gorm.DB) *ServiceItemRepository {
	return &ServiceItemRepository{db: db}
}

func (r *ServiceItemRepository) ListActiveByProviderIDs(ctx context.Context, providerIDs []uint64) ([]model.ServiceItem, error) {
	if len(providerIDs) == 0 {
		return []model.ServiceItem{}, nil
	}

	var items []model.ServiceItem
	if err := r.db.WithContext(ctx).
		Where("provider_id IN ? AND status = ?", providerIDs, 1).
		Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ServiceItemRepository) ListActiveByProviderID(ctx context.Context, providerID uint64) ([]model.ServiceItem, error) {
	return r.ListActiveByProviderIDs(ctx, []uint64{providerID})
}
