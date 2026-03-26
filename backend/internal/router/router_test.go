package router_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/router"
	"go.uber.org/zap"
)

func TestProviderOrderFlow(t *testing.T) {
	model.ResetDefaultForTest()
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
	mustRequest(t, engine, http.MethodPost, "/api/v1/payments/callback", paymentCallbackPayload(1, "SUCCESS", "TRADE-10001", "wechat", "N10001"), nil, http.StatusOK)

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
	if int(orderData["status"].(float64)) != model.OrderStatusCompleted {
		t.Fatalf("expected order status %d, got %v", model.OrderStatusCompleted, orderData["status"])
	}
}

func TestPaymentCallbackIdempotent(t *testing.T) {
	model.ResetDefaultForTest()
	engine := router.New(zap.NewNop())

	userToken := loginBySMS(t, engine, "13800000011")
	providerToken := loginBySMS(t, engine, "13900000011")
	adminToken := adminLogin(t, engine, "admin", "admin123456")

	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/apply", map[string]any{
		"realName": "李四",
		"idCardNo": "310101199002022345",
	}, withAuth(providerToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/admin/providers/1/approve", nil, withAuth(adminToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/orders", map[string]any{
		"providerId":      1,
		"serviceItemId":   1,
		"sceneText":       "第一次支付回调测试",
		"cityCode":        "310100",
		"addressText":     "虹桥",
		"plannedStartAt":  "2026-03-24 19:00:00",
		"plannedDuration": 60,
	}, withAuth(userToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/payments", map[string]any{
		"orderId":    1,
		"payChannel": 1,
	}, withAuth(userToken), http.StatusOK)

	firstPayload := paymentCallbackPayload(1, "SUCCESS", "TRADE-20001", "wechat", "N20001")
	mustRequest(t, engine, http.MethodPost, "/api/v1/payments/callback", firstPayload, nil, http.StatusOK)
	second := mustRequest(t, engine, http.MethodPost, "/api/v1/payments/callback", firstPayload, nil, http.StatusOK)

	var body struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(second.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode callback response: %v", err)
	}
	orderData := body.Data["order"].(map[string]any)
	if int(orderData["status"].(float64)) != model.OrderStatusPendingAccept {
		t.Fatalf("expected order status %d, got %v", model.OrderStatusPendingAccept, orderData["status"])
	}
	if idempotent, ok := body.Data["idempotent"].(bool); !ok || !idempotent {
		t.Fatalf("expected idempotent=true, got %v", body.Data["idempotent"])
	}
}

func TestPaymentCallbackSignatureInvalid(t *testing.T) {
	model.ResetDefaultForTest()
	engine := router.New(zap.NewNop())

	userToken := loginBySMS(t, engine, "13800000012")
	providerToken := loginBySMS(t, engine, "13900000012")
	adminToken := adminLogin(t, engine, "admin", "admin123456")

	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/apply", map[string]any{
		"realName": "王五",
		"idCardNo": "310101199003033456",
	}, withAuth(providerToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/admin/providers/1/approve", nil, withAuth(adminToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/orders", map[string]any{
		"providerId":      1,
		"serviceItemId":   1,
		"sceneText":       "签名校验测试",
		"cityCode":        "310100",
		"addressText":     "浦东",
		"plannedStartAt":  "2026-03-24 19:00:00",
		"plannedDuration": 60,
	}, withAuth(userToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/payments", map[string]any{
		"orderId":    1,
		"payChannel": 1,
	}, withAuth(userToken), http.StatusOK)

	payload := paymentCallbackPayload(1, "SUCCESS", "TRADE-30001", "wechat", "N30001")
	payload["sign"] = "invalid-sign"
	mustRequest(t, engine, http.MethodPost, "/api/v1/payments/callback", payload, nil, http.StatusUnauthorized)
}

func TestPostFlow(t *testing.T) {
	model.ResetDefaultForTest()
	engine := router.New(zap.NewNop())

	userToken := loginBySMS(t, engine, "13800000021")
	createResp := mustRequest(t, engine, http.MethodPost, "/api/v1/posts", map[string]any{
		"content":      "今天心情不错，想找人聊聊",
		"topic":        "日常",
		"cityCode":     "310100",
		"visibleScope": 1,
	}, withAuth(userToken), http.StatusOK)

	var createBody struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(createResp.Body.Bytes(), &createBody); err != nil {
		t.Fatalf("decode create post response: %v", err)
	}
	post := createBody.Data["post"].(map[string]any)
	postID := uint64(post["id"].(float64))
	postIDText := strconv.FormatUint(postID, 10)

	mustRequest(t, engine, http.MethodGet, "/api/v1/posts", nil, nil, http.StatusOK)
	mustRequest(t, engine, http.MethodGet, "/api/v1/posts/"+postIDText, nil, nil, http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/posts/"+postIDText+"/comments", map[string]any{
		"content": "欢迎来到广场",
	}, withAuth(userToken), http.StatusOK)
	liked := mustRequest(t, engine, http.MethodPost, "/api/v1/posts/"+postIDText+"/likes", nil, withAuth(userToken), http.StatusOK)
	mustRequest(t, engine, http.MethodDelete, "/api/v1/posts/"+postIDText+"/likes", nil, withAuth(userToken), http.StatusOK)

	var likeBody struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(liked.Body.Bytes(), &likeBody); err != nil {
		t.Fatalf("decode like response: %v", err)
	}
	likedPost := likeBody.Data["post"].(map[string]any)
	if uint64(likedPost["id"].(float64)) != postID {
		t.Fatalf("expected liked post id %d, got %v", postID, likedPost["id"])
	}
	if int(likedPost["likeCount"].(float64)) != 1 {
		t.Fatalf("expected likeCount 1, got %v", likedPost["likeCount"])
	}
}

func TestReviewAndComplaintFlow(t *testing.T) {
	model.ResetDefaultForTest()
	engine := router.New(zap.NewNop())

	userToken := loginBySMS(t, engine, "13800000022")
	providerToken := loginBySMS(t, engine, "13900000022")
	adminToken := adminLogin(t, engine, "admin", "admin123456")

	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/apply", map[string]any{
		"realName": "赵六",
		"idCardNo": "310101199004044567",
	}, withAuth(providerToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/admin/providers/1/approve", nil, withAuth(adminToken), http.StatusOK)

	mustRequest(t, engine, http.MethodPost, "/api/v1/orders", map[string]any{
		"providerId":      1,
		"serviceItemId":   1,
		"sceneText":       "review complaint test",
		"cityCode":        "310100",
		"addressText":     "徐汇",
		"plannedStartAt":  "2026-03-24 19:00:00",
		"plannedDuration": 90,
	}, withAuth(userToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/payments", map[string]any{
		"orderId":    1,
		"payChannel": 1,
	}, withAuth(userToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/payments/callback", paymentCallbackPayload(1, "SUCCESS", "TRADE-40001", "wechat", "N40001"), nil, http.StatusOK)

	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/orders/1/accept", nil, withAuth(providerToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/orders/1/depart", nil, withAuth(providerToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/orders/1/arrive", nil, withAuth(providerToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/orders/1/start", map[string]any{
		"confirmType": 1,
		"confirmCode": "5678",
	}, withAuth(userToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/provider-center/orders/1/finish", nil, withAuth(providerToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/orders/1/confirm-finish", nil, withAuth(userToken), http.StatusOK)

	mustRequest(t, engine, http.MethodPost, "/api/v1/orders/1/reviews", map[string]any{
		"score":       9,
		"content":     "服务体验很好",
		"images":      []string{"https://example.com/review1.png"},
		"isAnonymous": false,
	}, withAuth(userToken), http.StatusOK)
	mustRequest(t, engine, http.MethodGet, "/api/v1/orders/1/reviews", nil, nil, http.StatusOK)

	complaintResp := mustRequest(t, engine, http.MethodPost, "/api/v1/orders/1/complaints", map[string]any{
		"complaintType":  "late",
		"content":        "迟到20分钟",
		"evidenceImages": []string{"https://example.com/evi1.png"},
	}, withAuth(userToken), http.StatusOK)

	var complaintBody struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(complaintResp.Body.Bytes(), &complaintBody); err != nil {
		t.Fatalf("decode complaint create response: %v", err)
	}
	complaint := complaintBody.Data["complaint"].(map[string]any)
	complaintID := int(complaint["id"].(float64))
	mustRequest(t, engine, http.MethodGet, "/api/v1/complaints/"+strconv.Itoa(complaintID), nil, nil, http.StatusOK)
}

func TestAdminRouteRequiresAuth(t *testing.T) {
	model.ResetDefaultForTest()
	engine := router.New(zap.NewNop())

	mustRequest(t, engine, http.MethodGet, "/api/v1/admin/providers", nil, nil, http.StatusUnauthorized)
}

func TestAdminRoutePermissionDenied(t *testing.T) {
	model.ResetDefaultForTest()
	engine := router.New(zap.NewNop())

	contentToken := adminLogin(t, engine, "content_admin", "admin123456")

	mustRequest(t, engine, http.MethodGet, "/api/v1/admin/posts", nil, withAuth(contentToken), http.StatusOK)
	mustRequest(t, engine, http.MethodPost, "/api/v1/admin/providers/1/approve", nil, withAuth(contentToken), http.StatusForbidden)
}

func TestAdminAssignRoles(t *testing.T) {
	model.ResetDefaultForTest()
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

func paymentCallbackPayload(paymentID uint64, payStatus, thirdTradeNo, channel, notifyID string) map[string]any {
	timestamp := time.Now().Unix()
	message := fmt.Sprintf("channel=%s&notifyId=%s&payStatus=%s&paymentId=%d&thirdTradeNo=%s&timestamp=%d",
		channel, notifyID, strings.ToUpper(payStatus), paymentID, thirdTradeNo, timestamp)
	mac := hmac.New(sha256.New, []byte("listen-dev-callback-secret"))
	_, _ = mac.Write([]byte(message))
	sign := hex.EncodeToString(mac.Sum(nil))
	return map[string]any{
		"paymentId":    paymentID,
		"payStatus":    payStatus,
		"thirdTradeNo": thirdTradeNo,
		"channel":      channel,
		"notifyId":     notifyID,
		"timestamp":    timestamp,
		"sign":         sign,
	}
}
