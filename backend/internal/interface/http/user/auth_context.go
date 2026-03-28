package user

import (
	"net/http"
	"strings"
)

func currentUserID(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if userID := userIDFromMockToken(token); userID != "" {
			return userID
		}
	}

	if userID := strings.TrimSpace(r.Header.Get("X-User-ID")); userID != "" {
		return userID
	}
	if userID := strings.TrimSpace(r.URL.Query().Get("user_id")); userID != "" {
		return userID
	}
	return ""
}

func userIDFromMockToken(token string) string {
	// token pattern: mock_at_<userID>_<yyyymmddhhmmss>
	if !strings.HasPrefix(token, "mock_at_") {
		return ""
	}
	parts := strings.Split(token, "_")
	if len(parts) < 5 {
		return ""
	}
	return strings.Join(parts[2:len(parts)-1], "_")
}
