package provider

import (
	"net/http"
	"strings"

	domain "listen/backend/internal/domain/provider_auth"
)

func currentProviderID(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if providerID := domain.ProviderIDFromAccessToken(token); providerID != "" {
			return providerID
		}
	}

	if providerID := strings.TrimSpace(r.Header.Get("X-Provider-ID")); providerID != "" {
		return providerID
	}
	return ""
}
