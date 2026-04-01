package admin_auth

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidCredential = errors.New("invalid credential")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrAdminNotFound     = errors.New("admin not found")
)

type AdminAccount struct {
	AdminID     string
	Account     string
	Password    string
	Role        string
	DisplayName string
	Status      string
}

type AdminIdentity struct {
	AdminID          string
	Role             string
	AccessToken      string
	ExpiresInSeconds int64
}

type AdminProfile struct {
	AdminID     string
	Account     string
	Role        string
	DisplayName string
	Status      string
}

func BuildAccessToken(adminID string, now time.Time) string {
	return "mock_admin_at_" + adminID + "_" + strconv.FormatInt(now.Unix(), 10)
}

func AdminIDFromAccessToken(token string) string {
	clean := strings.TrimSpace(token)
	if !strings.HasPrefix(clean, "mock_admin_at_") {
		return ""
	}

	parts := strings.Split(clean, "_")
	if len(parts) < 5 {
		return ""
	}
	return strings.Join(parts[3:len(parts)-1], "_")
}
