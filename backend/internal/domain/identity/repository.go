package identity

import "context"

type Repository interface {
	GetByID(ctx context.Context, userID string) (UserAccount, bool, error)
	GetByPhone(ctx context.Context, phone string) (UserAccount, bool, error)
	GetByWechatOpenID(ctx context.Context, openID string) (UserAccount, bool, error)
	Save(ctx context.Context, account UserAccount) error
}

type AuthGateway interface {
	VerifySMSCode(ctx context.Context, phone string, code string) (bool, error)
	ResolveWechatOpenID(ctx context.Context, authCode string) (string, error)
}
