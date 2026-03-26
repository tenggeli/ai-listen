package payment

type CreatePaymentRequest struct {
	OrderID    uint64 `json:"orderId"`
	PayChannel int    `json:"payChannel"`
}

type PaymentCallbackRequest struct {
	PaymentID    uint64 `json:"paymentId"`
	PayStatus    string `json:"payStatus"`
	ThirdTradeNo string `json:"thirdTradeNo"`
	Channel      string `json:"channel"`
	NotifyID     string `json:"notifyId"`
	Timestamp    int64  `json:"timestamp"`
	Sign         string `json:"sign"`
	NotifyRaw    string `json:"notifyRaw"`
}
