package identity

import (
	"context"
	"regexp"
	"strings"
	"time"

	domain "listen/backend/internal/domain/identity"
)

var chinaPhoneRegexp = regexp.MustCompile(`^1\d{10}$`)

type Clock interface {
	Now() time.Time
}

type IDGenerator interface {
	NewID(prefix string) string
}

type LoginBySMSInput struct {
	Phone             string
	VerifyCode        string
	AgreementAccepted bool
}

type LoginByWechatInput struct {
	AuthCode          string
	AgreementAccepted bool
}

type LoginOutput struct {
	Identity domain.UserIdentity
}

type LoginBySMSUseCase struct {
	repo        domain.Repository
	auth        domain.AuthGateway
	clock       Clock
	idGenerator IDGenerator
}

func NewLoginBySMSUseCase(repo domain.Repository, auth domain.AuthGateway, clock Clock, idGenerator IDGenerator) LoginBySMSUseCase {
	return LoginBySMSUseCase{repo: repo, auth: auth, clock: clock, idGenerator: idGenerator}
}

func (u LoginBySMSUseCase) Execute(ctx context.Context, input LoginBySMSInput) (LoginOutput, error) {
	phone := strings.TrimSpace(input.Phone)
	code := strings.TrimSpace(input.VerifyCode)
	if !input.AgreementAccepted || !chinaPhoneRegexp.MatchString(phone) || code == "" {
		return LoginOutput{}, domain.ErrInvalidInput
	}

	ok, err := u.auth.VerifySMSCode(ctx, phone, code)
	if err != nil {
		return LoginOutput{}, err
	}
	if !ok {
		return LoginOutput{}, domain.ErrInvalidCredential
	}

	account, found, err := u.repo.GetByPhone(ctx, phone)
	if err != nil {
		return LoginOutput{}, err
	}

	isNewUser := false
	if !found {
		isNewUser = true
		account, err = domain.NewUserAccount(u.idGenerator.NewID("user"), phone, "")
		if err != nil {
			return LoginOutput{}, err
		}
		if err := u.repo.Save(ctx, account); err != nil {
			return LoginOutput{}, err
		}
	}

	return LoginOutput{Identity: domain.BuildUserIdentity(account, domain.LoginChannelSMS, u.clock.Now(), isNewUser)}, nil
}

type LoginByWechatUseCase struct {
	repo        domain.Repository
	auth        domain.AuthGateway
	clock       Clock
	idGenerator IDGenerator
}

func NewLoginByWechatUseCase(repo domain.Repository, auth domain.AuthGateway, clock Clock, idGenerator IDGenerator) LoginByWechatUseCase {
	return LoginByWechatUseCase{repo: repo, auth: auth, clock: clock, idGenerator: idGenerator}
}

func (u LoginByWechatUseCase) Execute(ctx context.Context, input LoginByWechatInput) (LoginOutput, error) {
	if !input.AgreementAccepted || strings.TrimSpace(input.AuthCode) == "" {
		return LoginOutput{}, domain.ErrInvalidInput
	}

	openID, err := u.auth.ResolveWechatOpenID(ctx, strings.TrimSpace(input.AuthCode))
	if err != nil {
		return LoginOutput{}, err
	}
	if openID == "" {
		return LoginOutput{}, domain.ErrInvalidCredential
	}

	account, found, err := u.repo.GetByWechatOpenID(ctx, openID)
	if err != nil {
		return LoginOutput{}, err
	}

	isNewUser := false
	if !found {
		isNewUser = true
		account, err = domain.NewUserAccount(u.idGenerator.NewID("user"), "", openID)
		if err != nil {
			return LoginOutput{}, err
		}
		if err := u.repo.Save(ctx, account); err != nil {
			return LoginOutput{}, err
		}
	}

	return LoginOutput{Identity: domain.BuildUserIdentity(account, domain.LoginChannelWechat, u.clock.Now(), isNewUser)}, nil
}
