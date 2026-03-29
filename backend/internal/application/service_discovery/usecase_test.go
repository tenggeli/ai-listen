package service_discovery

import (
	"context"
	"testing"
	"time"

	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestListServiceCategoriesUseCase_Success(t *testing.T) {
	repo := memory.NewServiceDiscoveryRepository()
	uc := NewListServiceCategoriesUseCase(repo)

	output, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(output.Items) < 3 {
		t.Fatalf("expected seeded categories, got %d", len(output.Items))
	}
	if output.Items[0].ID != "cat_all" {
		t.Fatalf("unexpected first category: %s", output.Items[0].ID)
	}
}

func TestListPublicProvidersUseCase_FilterByCategory(t *testing.T) {
	repo := memory.NewServiceDiscoveryRepository()
	uc := NewListPublicProvidersUseCase(repo)

	output, err := uc.Execute(context.Background(), ListPublicProvidersInput{
		CategoryID: "cat_relax",
		Page:       1,
		PageSize:   10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Total != 1 {
		t.Fatalf("expected one provider in cat_relax, got %d", output.Total)
	}
	if output.Items[0].ID != "p_pub_003" {
		t.Fatalf("unexpected provider id: %s", output.Items[0].ID)
	}
}

func TestGetPublicProviderUseCase_NotFound(t *testing.T) {
	repo := memory.NewServiceDiscoveryRepository()
	uc := NewGetPublicProviderUseCase(repo)

	_, err := uc.Execute(context.Background(), GetPublicProviderInput{ProviderID: "not_exist"})
	if err == nil {
		t.Fatal("expected not found error")
	}
	if err != ErrProviderNotFound {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListProviderServiceItemsUseCase_Success(t *testing.T) {
	repo := memory.NewServiceDiscoveryRepository()
	uc := NewListProviderServiceItemsUseCase(repo)

	output, err := uc.Execute(context.Background(), ListProviderServiceItemsInput{ProviderID: "p_pub_001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(output.Items) != 3 {
		t.Fatalf("expected 3 service items, got %d", len(output.Items))
	}
	if output.Items[0].PriceAmount != 99 {
		t.Fatalf("unexpected first item price: %d", output.Items[0].PriceAmount)
	}
}

func TestGetPublicProviderUseCase_OnlineStatusChangesWithHeartbeat(t *testing.T) {
	repo := memory.NewServiceDiscoveryRepository()
	uc := NewGetPublicProviderUseCase(repo)

	before, err := uc.Execute(context.Background(), GetPublicProviderInput{ProviderID: "p_pub_003"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if before.Provider.Online {
		t.Fatal("expected provider p_pub_003 to be offline before heartbeat")
	}

	if err := repo.TouchProviderHeartbeat("p_pub_003", time.Now()); err != nil {
		t.Fatalf("touch heartbeat failed: %v", err)
	}

	after, err := uc.Execute(context.Background(), GetPublicProviderInput{ProviderID: "p_pub_003"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !after.Provider.Online {
		t.Fatal("expected provider p_pub_003 to be online after heartbeat")
	}
}
