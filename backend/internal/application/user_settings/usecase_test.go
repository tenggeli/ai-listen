package user_settings

import (
	"context"
	"testing"

	identityDomain "listen/backend/internal/domain/identity"
	domain "listen/backend/internal/domain/user_settings"
	memory "listen/backend/internal/infrastructure/persistence/memory"
)

type stubSettingsRepo struct {
	data map[string]domain.Settings
}

func newStubSettingsRepo() *stubSettingsRepo {
	return &stubSettingsRepo{data: map[string]domain.Settings{}}
}

func (r *stubSettingsRepo) GetByUserID(_ context.Context, userID string) (domain.Settings, bool, error) {
	item, found := r.data[userID]
	return item, found, nil
}

func (r *stubSettingsRepo) Save(_ context.Context, settings domain.Settings) error {
	r.data[settings.UserID] = settings
	return nil
}

func TestSettingsUseCase_SaveAndGetSuccess(t *testing.T) {
	identityRepo := memory.NewIdentityRepository()
	settingsRepo := newStubSettingsRepo()
	account, err := identityDomain.NewUserAccount("user_settings_001", "13800138000", "")
	if err != nil {
		t.Fatalf("create user failed: %v", err)
	}
	if err := identityRepo.Save(context.Background(), account); err != nil {
		t.Fatalf("save user failed: %v", err)
	}

	saveUC := NewSaveSettingsUseCase(settingsRepo, identityRepo)
	getUC := NewGetSettingsUseCase(settingsRepo, identityRepo)

	_, err = saveUC.Execute(context.Background(), SaveSettingsInput{
		UserID: "user_settings_001",
		Preference: domain.Preference{
			PreferSameCityProviders: true,
			AutoPlaySoundPreview:    false,
			HideOfflineProviders:    true,
		},
		Notification: domain.Notification{
			OrderStatusUpdate:     true,
			ComplaintResultNotice: false,
			MarketingActivity:     true,
		},
		Privacy: domain.Privacy{
			ProfilePublicVisible:       false,
			PersonalizedRecommendation: true,
			RiskControlDataSharing:     false,
		},
	})
	if err != nil {
		t.Fatalf("save settings failed: %v", err)
	}

	got, err := getUC.Execute(context.Background(), GetSettingsInput{UserID: "user_settings_001"})
	if err != nil {
		t.Fatalf("get settings failed: %v", err)
	}
	if got.Settings.Preference.AutoPlaySoundPreview {
		t.Fatalf("expected auto play false")
	}
	if got.Settings.Privacy.ProfilePublicVisible {
		t.Fatalf("expected profile public false")
	}
}

func TestSettingsUseCase_MysqlOnlyInMemoryMode(t *testing.T) {
	identityRepo := memory.NewIdentityRepository()
	account, err := identityDomain.NewUserAccount("user_settings_002", "13800138001", "")
	if err != nil {
		t.Fatalf("create user failed: %v", err)
	}
	if err := identityRepo.Save(context.Background(), account); err != nil {
		t.Fatalf("save user failed: %v", err)
	}

	repo := memory.NewUserSettingsRepository()
	getUC := NewGetSettingsUseCase(repo, identityRepo)

	_, err = getUC.Execute(context.Background(), GetSettingsInput{UserID: "user_settings_002"})
	if err == nil {
		t.Fatalf("expected mysql only error")
	}
	if err != domain.ErrSettingsPersistenceUnavailable {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSettingsUseCase_InvalidInput(t *testing.T) {
	identityRepo := memory.NewIdentityRepository()
	settingsRepo := newStubSettingsRepo()
	getUC := NewGetSettingsUseCase(settingsRepo, identityRepo)
	if _, err := getUC.Execute(context.Background(), GetSettingsInput{UserID: ""}); err == nil {
		t.Fatalf("expected invalid input error")
	}
}
