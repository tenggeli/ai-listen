package provider_auth

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
	ErrProviderNotFound  = errors.New("provider not found")
)

type ProviderAccount struct {
	ProviderID  string
	Account     string
	Password    string
	DisplayName string
	Status      string
	CityCode    string
}

type ProviderIdentity struct {
	ProviderID       string
	AccessToken      string
	ExpiresInSeconds int64
}

type ProviderProfile struct {
	ProviderID  string
	Account     string
	DisplayName string
	Status      string
	CityCode    string
}

func BuildAccessToken(providerID string, now time.Time) string {
	return "mock_provider_at_" + providerID + "_" + strconv.FormatInt(now.Unix(), 10)
}

func ProviderIDFromAccessToken(token string) string {
	clean := strings.TrimSpace(token)
	if !strings.HasPrefix(clean, "mock_provider_at_") {
		return ""
	}
	parts := strings.Split(clean, "_")
	if len(parts) < 6 {
		return ""
	}
	return strings.Join(parts[3:len(parts)-1], "_")
}
