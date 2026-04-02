package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	app "listen/backend/internal/application/service_item_admin"
	domain "listen/backend/internal/domain/service_item_admin"
)

type testServiceItemRepo struct {
	items map[string]domain.ServiceItem
}

func newTestServiceItemRepo() *testServiceItemRepo {
	return &testServiceItemRepo{
		items: map[string]domain.ServiceItem{
			"si_001": {
				ID:           "si_001",
				ProviderID:   "p_001",
				ProviderName: "provider_1",
				CategoryID:   "cat_chat",
				Title:        "chat",
				Status:       domain.StatusActive,
			},
		},
	}
}

func (r *testServiceItemRepo) List(_ context.Context, query domain.Query) ([]domain.ServiceItem, int, error) {
	items := make([]domain.ServiceItem, 0)
	for _, item := range r.items {
		if query.Status != "" && item.Status != query.Status {
			continue
		}
		items = append(items, item)
	}
	return items, len(items), nil
}

func (r *testServiceItemRepo) GetByID(_ context.Context, serviceItemID string) (domain.ServiceItem, error) {
	item, found := r.items[serviceItemID]
	if !found {
		return domain.ServiceItem{}, domain.ErrServiceItemNotFound
	}
	return item, nil
}

func (r *testServiceItemRepo) UpdateStatus(_ context.Context, serviceItemID string, status string, _ string) error {
	item, found := r.items[serviceItemID]
	if !found {
		return domain.ErrServiceItemNotFound
	}
	item.Status = status
	r.items[serviceItemID] = item
	return nil
}

func TestServiceItemRoutes_ListAndAction(t *testing.T) {
	repo := newTestServiceItemRepo()
	controller := NewServiceItemController(
		app.NewListServiceItemsUseCase(repo),
		app.NewGetServiceItemDetailUseCase(repo),
		app.NewUpdateServiceItemStatusUseCase(repo),
	)

	mux := http.NewServeMux()
	RegisterServiceItemRoutes(mux, controller)

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/service-items?status=active", nil)
	listReq.Header.Set("Authorization", "Bearer mock_admin_at_admin_001_1712013723")
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("unexpected list status: %d", listRec.Code)
	}

	actionReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/service-items/si_001/deactivate", nil)
	actionReq.Header.Set("Authorization", "Bearer mock_admin_at_admin_001_1712013723")
	actionRec := httptest.NewRecorder()
	mux.ServeHTTP(actionRec, actionReq)
	if actionRec.Code != http.StatusOK {
		t.Fatalf("unexpected action status: %d", actionRec.Code)
	}
}

func TestServiceItemRoutes_DetailNotFound(t *testing.T) {
	repo := newTestServiceItemRepo()
	controller := NewServiceItemController(
		app.NewListServiceItemsUseCase(repo),
		app.NewGetServiceItemDetailUseCase(repo),
		app.NewUpdateServiceItemStatusUseCase(repo),
	)

	mux := http.NewServeMux()
	RegisterServiceItemRoutes(mux, controller)

	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/service-items/not_found", nil)
	detailReq.Header.Set("Authorization", "Bearer mock_admin_at_admin_001_1712013723")
	detailRec := httptest.NewRecorder()
	mux.ServeHTTP(detailRec, detailReq)
	if detailRec.Code != http.StatusNotFound {
		t.Fatalf("unexpected detail status: %d", detailRec.Code)
	}
}
