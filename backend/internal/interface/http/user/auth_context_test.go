package user

import (
	"net/http/httptest"
	"testing"
)

func TestCurrentUserID_FromBearerToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer mock_at_user_123_20260328120000")

	got := currentUserID(req)
	if got != "user_123" {
		t.Fatalf("unexpected user id: %s", got)
	}
}
