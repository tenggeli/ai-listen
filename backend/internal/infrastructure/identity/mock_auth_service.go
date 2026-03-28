package identity

import (
	"context"
	"strings"
)

type MockAuthService struct{}

func NewMockAuthService() MockAuthService {
	return MockAuthService{}
}

func (MockAuthService) VerifySMSCode(_ context.Context, _ string, code string) (bool, error) {
	return strings.TrimSpace(code) == "123456", nil
}

func (MockAuthService) ResolveWechatOpenID(_ context.Context, authCode string) (string, error) {
	clean := strings.TrimSpace(authCode)
	if clean == "" {
		return "", nil
	}
	return "wx_open_" + clean, nil
}
