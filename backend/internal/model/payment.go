package model

import "time"

const (
	PaymentStatusPending = 10
	PaymentStatusPaid    = 20
	PaymentStatusFailed  = 30
)

type Payment struct {
	ID        uint64    `json:"id"`
	OrderID   uint64    `json:"orderId"`
	PayAmount int64     `json:"payAmount"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	PaidAt    time.Time `json:"paidAt"`
}

func clonePayment(payment *Payment) *Payment {
	if payment == nil {
		return nil
	}
	copyPayment := *payment
	return &copyPayment
}
