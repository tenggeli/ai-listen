package model

import "time"

type OrderLog struct {
	FromStatus   int       `json:"fromStatus"`
	ToStatus     int       `json:"toStatus"`
	OperatorRole string    `json:"operatorRole"`
	OperatorID   uint64    `json:"operatorId"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Order struct {
	ID              uint64     `json:"id"`
	OrderNo         string     `json:"orderNo"`
	UserID          uint64     `json:"userId"`
	ProviderID      uint64     `json:"providerId"`
	ServiceItemID   uint64     `json:"serviceItemId"`
	SceneText       string     `json:"sceneText"`
	CityCode        string     `json:"cityCode"`
	AddressText     string     `json:"addressText"`
	PlannedStartAt  string     `json:"plannedStartAt"`
	PlannedDuration int        `json:"plannedDuration"`
	Status          int        `json:"status"`
	PayAmount       int64      `json:"payAmount"`
	CancelReason    string     `json:"cancelReason,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	AcceptedAt      *time.Time `json:"acceptedAt,omitempty"`
	ArrivedAt       *time.Time `json:"arrivedAt,omitempty"`
	StartedAt       *time.Time `json:"startedAt,omitempty"`
	FinishedAt      *time.Time `json:"finishedAt,omitempty"`
	Logs            []OrderLog `json:"logs"`
}

type CreateOrderInput struct {
	ProviderID      uint64
	ServiceItemID   uint64
	SceneText       string
	CityCode        string
	AddressText     string
	PlannedStartAt  string
	PlannedDuration int
}

func cloneOrder(order *Order) *Order {
	if order == nil {
		return nil
	}
	copyOrder := *order
	copyOrder.Logs = append([]OrderLog(nil), order.Logs...)
	return &copyOrder
}
