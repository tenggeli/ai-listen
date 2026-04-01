package admin

import (
	"net/http"
	"strings"

	domain "listen/backend/internal/domain/admin_auth"
)

func currentAdminID(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if adminID := domain.AdminIDFromAccessToken(token); adminID != "" {
			return adminID
		}
	}

	if adminID := strings.TrimSpace(r.Header.Get("X-Admin-ID")); adminID != "" {
		return adminID
	}

	return ""
}
