package service_item_admin

import (
	"context"
	"testing"

	domain "listen/backend/internal/domain/service_item_admin"
)

type fakeServiceItemRepo struct {
	items map[string]domain.ServiceItem
}

func newFakeServiceItemRepo() *fakeServiceItemRepo {
	return &fakeServiceItemRepo{
		items: map[string]domain.ServiceItem{
			"si_001": {
				ID:           "si_001",
				ProviderID:   "p_001",
				ProviderName: "provider_1",
				CategoryID:   "cat_chat",
				Title:        "chat 30 min",
				Status:       domain.StatusActive,
			},
			"si_002": {
				ID:           "si_002",
				ProviderID:   "p_002",
				ProviderName: "provider_2",
				CategoryID:   "cat_movie",
				Title:        "movie",
				Status:       domain.StatusInactive,
			},
		},
	}
}

func (r *fakeServiceItemRepo) List(_ context.Context, query domain.Query) ([]domain.ServiceItem, int, error) {
	items := make([]domain.ServiceItem, 0)
	for _, item := range r.items {
		if query.Status != "" && item.Status != query.Status {
			continue
		}
		items = append(items, item)
	}
	return items, len(items), nil
}

func (r *fakeServiceItemRepo) GetByID(_ context.Context, serviceItemID string) (domain.ServiceItem, error) {
	item, found := r.items[serviceItemID]
	if !found {
		return domain.ServiceItem{}, domain.ErrServiceItemNotFound
	}
	return item, nil
}

func (r *fakeServiceItemRepo) UpdateStatus(_ context.Context, serviceItemID string, status string, _ string) error {
	item, found := r.items[serviceItemID]
	if !found {
		return domain.ErrServiceItemNotFound
	}
	item.Status = status
	r.items[serviceItemID] = item
	return nil
}

func TestListServiceItemsUseCase_FilterByStatus(t *testing.T) {
	repo := newFakeServiceItemRepo()
	uc := NewListServiceItemsUseCase(repo)

	output, err := uc.Execute(context.Background(), ListServiceItemsInput{Status: domain.StatusActive})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if output.Total != 1 {
		t.Fatalf("unexpected total: %d", output.Total)
	}
}

func TestGetServiceItemDetailUseCase_NotFound(t *testing.T) {
	repo := newFakeServiceItemRepo()
	uc := NewGetServiceItemDetailUseCase(repo)
	if _, err := uc.Execute(context.Background(), GetServiceItemDetailInput{ServiceItemID: "missing"}); err == nil {
		t.Fatalf("expected not found")
	}
}

func TestUpdateServiceItemStatusUseCase(t *testing.T) {
	repo := newFakeServiceItemRepo()
	uc := NewUpdateServiceItemStatusUseCase(repo)

	output, err := uc.Execute(context.Background(), UpdateServiceItemStatusInput{
		ServiceItemID: "si_001",
		Action:        "deactivate",
		Operator:      "admin_1",
	})
	if err != nil {
		t.Fatalf("update status failed: %v", err)
	}
	if output.Item.Status != domain.StatusInactive {
		t.Fatalf("unexpected status: %s", output.Item.Status)
	}
}
