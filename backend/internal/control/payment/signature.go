package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func buildCallbackSignText(req PaymentCallbackRequest) string {
	return fmt.Sprintf("channel=%s&notifyId=%s&payStatus=%s&paymentId=%d&thirdTradeNo=%s&timestamp=%d",
		req.Channel,
		req.NotifyID,
		strings.ToUpper(req.PayStatus),
		req.PaymentID,
		req.ThirdTradeNo,
		req.Timestamp,
	)
}

func signCallback(secret string, req PaymentCallbackRequest) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(buildCallbackSignText(req)))
	return hex.EncodeToString(mac.Sum(nil))
}

func verifyCallbackSign(secret string, req PaymentCallbackRequest) bool {
	expected := signCallback(secret, req)
	return hmac.Equal([]byte(strings.ToLower(req.Sign)), []byte(expected))
}
