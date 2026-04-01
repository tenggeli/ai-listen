package admin_auth

import (
	"context"
	"strings"
	"time"

	domain "listen/backend/internal/domain/admin_auth"
)

type Clock interface {
	Now() time.Time
}

type LoginMockInput struct {
	Account  string
	Password string
}

type LoginMockOutput struct {
	Identity domain.AdminIdentity
}

type LoginMockUseCase struct {
	repo  domain.Repository
	clock Clock
}

func NewLoginMockUseCase(repo domain.Repository, clock Clock) LoginMockUseCase {
	return LoginMockUseCase{repo: repo, clock: clock}
}

func (u LoginMockUseCase) Execute(ctx context.Context, input LoginMockInput) (LoginMockOutput, error) {
	account := strings.TrimSpace(input.Account)
	password := strings.TrimSpace(input.Password)
	if account == "" || password == "" {
		return LoginMockOutput{}, domain.ErrInvalidInput
	}

	admin, found, err := u.repo.GetByAccount(ctx, account)
	if err != nil {
		return LoginMockOutput{}, err
	}
	if !found || admin.Password != password {
		return LoginMockOutput{}, domain.ErrInvalidCredential
	}
	if admin.Status != "active" {
		return LoginMockOutput{}, domain.ErrUnauthorized
	}

	now := u.clock.Now()
	return LoginMockOutput{Identity: domain.AdminIdentity{
		AdminID:          admin.AdminID,
		Role:             admin.Role,
		AccessToken:      domain.BuildAccessToken(admin.AdminID, now),
		ExpiresInSeconds: 7200,
	}}, nil
}

type GetCurrentAdminInput struct {
	AdminID string
}

type GetCurrentAdminOutput struct {
	Admin domain.AdminProfile
}

type GetCurrentAdminUseCase struct {
	repo domain.Repository
}

func NewGetCurrentAdminUseCase(repo domain.Repository) GetCurrentAdminUseCase {
	return GetCurrentAdminUseCase{repo: repo}
}

func (u GetCurrentAdminUseCase) Execute(ctx context.Context, input GetCurrentAdminInput) (GetCurrentAdminOutput, error) {
	adminID := strings.TrimSpace(input.AdminID)
	if adminID == "" {
		return GetCurrentAdminOutput{}, domain.ErrUnauthorized
	}

	admin, found, err := u.repo.GetByID(ctx, adminID)
	if err != nil {
		return GetCurrentAdminOutput{}, err
	}
	if !found {
		return GetCurrentAdminOutput{}, domain.ErrAdminNotFound
	}
	if admin.Status != "active" {
		return GetCurrentAdminOutput{}, domain.ErrUnauthorized
	}

	return GetCurrentAdminOutput{Admin: domain.AdminProfile{
		AdminID:     admin.AdminID,
		Account:     admin.Account,
		Role:        admin.Role,
		DisplayName: admin.DisplayName,
		Status:      admin.Status,
	}}, nil
}

type InMemoryRepository struct {
	byID      map[string]domain.AdminAccount
	byAccount map[string]string
}

func NewInMemoryRepository(accounts []domain.AdminAccount) InMemoryRepository {
	byID := make(map[string]domain.AdminAccount, len(accounts))
	byAccount := make(map[string]string, len(accounts))
	for _, item := range accounts {
		cleanID := strings.TrimSpace(item.AdminID)
		cleanAccount := strings.TrimSpace(item.Account)
		if cleanID == "" || cleanAccount == "" {
			continue
		}
		byID[cleanID] = item
		byAccount[cleanAccount] = cleanID
	}
	return InMemoryRepository{byID: byID, byAccount: byAccount}
}

func (r InMemoryRepository) GetByAccount(_ context.Context, account string) (domain.AdminAccount, bool, error) {
	adminID, ok := r.byAccount[strings.TrimSpace(account)]
	if !ok {
		return domain.AdminAccount{}, false, nil
	}
	item, ok := r.byID[adminID]
	if !ok {
		return domain.AdminAccount{}, false, nil
	}
	return item, true, nil
}

func (r InMemoryRepository) GetByID(_ context.Context, adminID string) (domain.AdminAccount, bool, error) {
	item, ok := r.byID[strings.TrimSpace(adminID)]
	if !ok {
		return domain.AdminAccount{}, false, nil
	}
	return item, true, nil
}
