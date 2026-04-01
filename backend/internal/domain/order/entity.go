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
	StatusCreated = "created"
	StatusPaid    = "paid"
)

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
