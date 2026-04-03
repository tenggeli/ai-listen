package provider

import (
	"testing"

	providerAuthApp "listen/backend/internal/application/provider_auth"
	providerAuthDomain "listen/backend/internal/domain/provider_auth"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

func TestUpdateCurrentProfileUseCase(t *testing.T) {
	repo := providerAuthApp.NewInMemoryRepository([]providerAuthDomain.ProviderAccount{
		{ProviderID: "p_pub_001", Account: "provider", Password: "provider123", DisplayName: "旧昵称", Status: "active", CityCode: "310000"},
	})
	uc := NewUpdateCurrentProfileUseCase(repo)

	output, err := uc.Execute(t.Context(), UpdateCurrentProfileInput{
		ProviderID:  "p_pub_001",
		DisplayName: "新昵称",
		CityCode:    "110000",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if output.Provider.DisplayName != "新昵称" {
		t.Fatalf("expected updated display name")
	}
	if output.Provider.CityCode != "110000" {
		t.Fatalf("expected updated city code")
	}

	fetched, found, err := repo.GetByID(t.Context(), "p_pub_001")
	if err != nil || !found {
		t.Fatalf("expected provider found, err=%v, found=%v", err, found)
	}
	if fetched.DisplayName != "新昵称" || fetched.CityCode != "110000" {
		t.Fatalf("expected repository data updated")
	}
}

func TestUpdateCurrentProfileUseCase_ValidateInput(t *testing.T) {
	repo := providerAuthApp.NewInMemoryRepository([]providerAuthDomain.ProviderAccount{})
	uc := NewUpdateCurrentProfileUseCase(repo)

	if _, err := uc.Execute(t.Context(), UpdateCurrentProfileInput{ProviderID: "", DisplayName: "昵称", CityCode: "310000"}); err == nil {
		t.Fatalf("expected unauthorized error")
	}
	if _, err := uc.Execute(t.Context(), UpdateCurrentProfileInput{ProviderID: "p_pub_001", DisplayName: "", CityCode: "310000"}); err == nil {
		t.Fatalf("expected invalid input error")
	}
	if _, err := uc.Execute(t.Context(), UpdateCurrentProfileInput{ProviderID: "missing", DisplayName: "昵称", CityCode: "310000"}); err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestListCurrentProviderServicesUseCase(t *testing.T) {
	repo := memory.NewServiceDiscoveryRepository()
	uc := NewListCurrentProviderServicesUseCase(repo)

	output, err := uc.Execute(t.Context(), ListCurrentProviderServicesInput{ProviderID: "p_pub_001"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(output.Items) == 0 {
		t.Fatalf("expected provider services")
	}

	if _, err := uc.Execute(t.Context(), ListCurrentProviderServicesInput{ProviderID: ""}); err == nil {
		t.Fatalf("expected unauthorized error")
	}
	if _, err := uc.Execute(t.Context(), ListCurrentProviderServicesInput{ProviderID: "not_exist"}); err == nil {
		t.Fatalf("expected not found error")
	}
}
