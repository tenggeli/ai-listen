package model

import "time"

type Payment struct {
	ID        uint64    `json:"id"`
	OrderID   uint64    `json:"orderId"`
	PayAmount int64     `json:"payAmount"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	PaidAt    time.Time `json:"paidAt"`
}
