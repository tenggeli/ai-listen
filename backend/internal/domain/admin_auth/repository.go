package admin_auth

import "context"

type Repository interface {
	GetByAccount(ctx context.Context, account string) (AdminAccount, bool, error)
	GetByID(ctx context.Context, adminID string) (AdminAccount, bool, error)
}
