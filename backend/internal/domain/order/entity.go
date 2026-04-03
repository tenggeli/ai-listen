package order

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidInput           = errors.New("invalid input")
	ErrOrderNotFound          = errors.New("order not found")
	ErrOrderForbidden         = errors.New("order forbidden")
	ErrInvalidOrderTransition = errors.New("invalid order transition")
)

const (
	StatusCreated   = "created"
	StatusPaid      = "paid"
	StatusAccepted  = "accepted"
	StatusOnTheWay  = "on_the_way"
	StatusArrived   = "arrived"
	StatusInService = "in_service"
	StatusCompleted = "completed"
	StatusAfterSale = "after_sale_processing"
	StatusClosed    = "closed"
)

var orderStatusReasonMap = map[string]string{
	StatusCreated:   "待支付",
	StatusPaid:      "待服务方接单",
	StatusAccepted:  "服务方已接单",
	StatusOnTheWay:  "服务方出发中",
	StatusArrived:   "服务方已到达，待开始服务",
	StatusInService: "服务进行中",
	StatusCompleted: "服务已完成",
	StatusAfterSale: "订单售后处理中",
	StatusClosed:    "订单已关闭",
}

type Order struct {
	ID               string
	UserID           string
	ProviderID       string
	ProviderName     string
	ServiceItemID    string
	ServiceItemTitle string
	Amount           int
	Currency         string
	Status           string
	CreatedAt        time.Time
	PaidAt           *time.Time
}

func IsKnownStatus(status string) bool {
	_, ok := orderStatusReasonMap[status]
	return ok
}

func StatusReason(status string) string {
	if reason, ok := orderStatusReasonMap[status]; ok {
		return reason
	}
	return "状态未知"
}

func NewOrder(
	id string,
	userID string,
	providerID string,
	providerName string,
	serviceItemID string,
	serviceItemTitle string,
	amount int,
	currency string,
	createdAt time.Time,
) (Order, error) {
	if strings.TrimSpace(id) == "" ||
		strings.TrimSpace(userID) == "" ||
		strings.TrimSpace(providerID) == "" ||
		strings.TrimSpace(providerName) == "" ||
		strings.TrimSpace(serviceItemID) == "" ||
		strings.TrimSpace(serviceItemTitle) == "" ||
		amount <= 0 {
		return Order{}, ErrInvalidInput
	}

	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		currency = "CNY"
	}

	return Order{
		ID:               strings.TrimSpace(id),
		UserID:           strings.TrimSpace(userID),
		ProviderID:       strings.TrimSpace(providerID),
		ProviderName:     strings.TrimSpace(providerName),
		ServiceItemID:    strings.TrimSpace(serviceItemID),
		ServiceItemTitle: strings.TrimSpace(serviceItemTitle),
		Amount:           amount,
		Currency:         currency,
		Status:           StatusCreated,
		CreatedAt:        createdAt,
	}, nil
}

func (o *Order) MarkPaid(now time.Time) error {
	if o.Status == StatusPaid {
		return nil
	}
	if o.Status != StatusCreated {
		return ErrInvalidOrderTransition
	}
	o.Status = StatusPaid
	paidAt := now
	o.PaidAt = &paidAt
	return nil
}

func (o *Order) MarkAccepted() error {
	if o.Status != StatusPaid {
		return ErrInvalidOrderTransition
	}
	o.Status = StatusAccepted
	return nil
}

func (o *Order) MarkOnTheWay() error {
	if o.Status != StatusAccepted {
		return ErrInvalidOrderTransition
	}
	o.Status = StatusOnTheWay
	return nil
}

func (o *Order) MarkArrived() error {
	if o.Status != StatusOnTheWay {
		return ErrInvalidOrderTransition
	}
	o.Status = StatusArrived
	return nil
}

func (o *Order) MarkInService() error {
	if o.Status != StatusArrived {
		return ErrInvalidOrderTransition
	}
	o.Status = StatusInService
	return nil
}

func (o *Order) MarkCompleted() error {
	if o.Status != StatusInService {
		return ErrInvalidOrderTransition
	}
	o.Status = StatusCompleted
	return nil
}

func (o *Order) MarkAfterSaleProcessing() error {
	if o.Status == StatusAfterSale {
		return nil
	}
	switch o.Status {
	case StatusPaid, StatusAccepted, StatusOnTheWay, StatusArrived, StatusInService, StatusCompleted, StatusClosed:
		o.Status = StatusAfterSale
		return nil
	default:
		return ErrInvalidOrderTransition
	}
}

func (o *Order) MarkClosedByAdmin() error {
	if o.Status == StatusClosed {
		return nil
	}
	switch o.Status {
	case StatusAfterSale, StatusCompleted:
		o.Status = StatusClosed
		return nil
	default:
		return ErrInvalidOrderTransition
	}
}
