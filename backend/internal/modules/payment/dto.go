package payment

type CreatePaymentRequest struct {
	OrderID    uint64 `json:"orderId"`
	PayChannel int    `json:"payChannel"`
}
