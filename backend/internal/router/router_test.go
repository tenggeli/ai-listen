package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-listen/backend/internal/router"
	"ai-listen/backend/internal/store"
	"go.uber.org/zap"
)

func TestProviderOrderFlow(t *testing.T) {
	store.ResetDefaultForTest()
	engine := router.New(zap.NewNop())

	userToken := loginBySMS(t, engine, "13800000000")
	providerToken := loginBySMS(t, engine, "13900000000")
	adminToken := adminLogin(t, engine, "admin", "admin123456")

	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/apply", map[string]any{
		"realName": "张三",
		"idCardNo": "310101199001011234",
	}, withAuth(providerToken), http.StatusOK)

	mustRequest(t, engine, http.MethodPost, "/api/v1/admin/providers/1/approve", nil, withAuth(adminToken), http.StatusOK)

	mustRequest(t, engine, http.MethodPut, "/api/v1/provider-center/profile", map[string]any{
		"displayName": "晚餐搭子",
		"intro":       "擅长吃饭聊天",
		"tags":        []string{"美食", "电影"},
	}, withAuth(providerToken), http.StatusOK)

	mustRequest(t, engine, http.MethodPut, "/api/v1/provider-center/service-items", map[string]any{
		"items": []map[string]any{
			{"serviceItemId": 1, "priceAmount": 12000, "priceUnit": "小时"},
		},
	}, withAuth(providerToken), http.StatusOK)

	mustRequest(t, engine, http.MethodPost, "/api/v1/orders", map[string]any{
		"providerId":      1,
		"serviceItemId":   1,
		"sceneText":       "一起吃晚饭",
		"cityCode":        "310100",
		"addressText":     "静安寺",
		"plannedStartAt":  "2026-03-24 19:00:00",
		"plannedDuration": 120,
	}, withAuth(userToken), http.StatusOK)

	mustRequest(t, engine, http.MethodPost, "/api/v1/payments", map[string]any{
		"orderId":    1,
		"payChannel": 1,
	}, withAuth(userToken), http.StatusOK)

	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/orders/1/accept", nil, withAuth(providerToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/orders/1/depart", nil, withAuth(providerToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/orders/1/arrive", nil, withAuth(providerToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/orders/1/start", map[string]any{
		"confirmType": 1,
		"confirmCode": "1234",
	}, withAuth(userToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/orders/1/finish", nil, withAuth(providerToken), http.StatusOK)
	resp := mustRequest(t, engine, http.MethodPost, "/api/v1/orders/1/confirm-finish", nil, withAuth(userToken), http.StatusOK)

	var body struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	orderData := body.Data["order"].(map[string]any)
	if int(orderData["status"].(float64)) != store.OrderStatusCompleted {
		t.Fatalf("expected order status %d, got %v", store.OrderStatusCompleted, orderData["status"])
	}
}

func TestAdminRouteRequiresAuth(t *testing.T) {
	store.ResetDefaultForTest()
	engine := router.New(zap.NewNop())

	mustRequest(t, engine, http.MethodGet, "/api/v1/admin/providers", nil, nil, http.StatusUnauthorized)
}

func TestAdminRoutePermissionDenied(t *testing.T) {
	store.ResetDefaultForTest()
	engine := router.New(zap.NewNop())

	contentToken := adminLogin(t, engine, "content_admin", "admin123456")

	mustRequest(t, engine, http.MethodGet, "/api/v1/admin/posts", nil, withAuth(contentToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/admin/providers/1/approve", nil, withAuth(contentToken), http.StatusForbidden)
}

func TestAdminAssignRoles(t *testing.T) {
	store.ResetDefaultForTest()
	engine := router.New(zap.NewNop())

	superToken := adminLogin(t, engine, "admin", "admin123456")
	contentToken := adminLogin(t, engine, "content_admin", "admin123456")

	mustRequest(t, engine, http.MethodGet, "/api/v1/admin/finance/reports", nil, withAuth(contentToken), http.StatusForbidden)
	mustRequest(t, engine, http.MethodGet, "/api/v1/admin/posts", nil, withAuth(contentToken), http.StatusOK)

	mustRequest(t, engine, http.MethodPut, "/api/v1/admin/rbac/users/2/roles", map[string]any{
		"roleKeys": []string{"finance_admin"},
	}, withAuth(superToken), http.StatusOK)

	mustRequest(t, engine, http.MethodGet, "/api/v1/admin/finance/reports", nil, withAuth(contentToken), http.StatusOK)
	mustRequest(t, engine, http.MethodGet, "/api/v1/admin/posts", nil, withAuth(contentToken), http.StatusForbidden)
}

func loginBySMS(t *testing.T, engine http.Handler, mobile string) string {
	t.Helper()
	mustRequest(t, engine, http.MethodPost, "/api/v1/auth/sms/send", map[string]any{
		"mobile": mobile,
		"scene":  "login",
	}, nil, http.StatusOK)
	resp := mustRequest(t, engine, http.MethodPost, "/api/v1/auth/login/sms", map[string]any{
		"mobile": mobile,
		"code":   "123456",
	}, nil, http.StatusOK)

	var body struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	return body.Data["accessToken"].(string)
}

func adminLogin(t *testing.T, engine http.Handler, username, password string) string {
	t.Helper()
	resp := mustRequest(t, engine, http.MethodPost, "/api/v1/admin/auth/login", map[string]any{
		"username": username,
		"password": password,
	}, nil, http.StatusOK)

	var body struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode admin login response: %v", err)
	}
	return body.Data["accessToken"].(string)
}

func mustRequest(t *testing.T, engine http.Handler, method, path string, payload any, headers map[string]string, wantStatus int) *httptest.ResponseRecorder {
	t.Helper()
	var body []byte
	if payload != nil {
		var err error
		body, err = json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal payload: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != wantStatus {
		t.Fatalf("%s %s expected status %d got %d body=%s", method, path, wantStatus, rec.Code, rec.Body.String())
	}
	return rec
}

func withAuth(token string) map[string]string {
	return map[string]string{"Authorization": "Bearer " + token}
}
