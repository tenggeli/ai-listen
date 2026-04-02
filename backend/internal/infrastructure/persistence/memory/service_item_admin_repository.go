package memory

import (
	"context"

	domain "listen/backend/internal/domain/service_item_admin"
)

type ServiceItemAdminRepository struct{}

func NewServiceItemAdminRepository() ServiceItemAdminRepository {
	return ServiceItemAdminRepository{}
}

func (r ServiceItemAdminRepository) List(_ context.Context, _ domain.Query) ([]domain.ServiceItem, int, error) {
	return nil, 0, domain.ErrPersistenceUnavailableInMemory
}

func (r ServiceItemAdminRepository) GetByID(_ context.Context, _ string) (domain.ServiceItem, error) {
	return domain.ServiceItem{}, domain.ErrPersistenceUnavailableInMemory
}

func (r ServiceItemAdminRepository) UpdateStatus(_ context.Context, _ string, _ string, _ string) error {
	return domain.ErrPersistenceUnavailableInMemory
}

var _ domain.Repository = ServiceItemAdminRepository{}
