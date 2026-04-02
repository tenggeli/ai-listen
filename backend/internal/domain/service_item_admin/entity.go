package service_item_admin

import (
	"errors"
	"strings"
)

var (
	ErrInvalidInput                   = errors.New("invalid input")
	ErrServiceItemNotFound            = errors.New("service item not found")
	ErrPersistenceUnavailableInMemory = errors.New("service item management only available in mysql mode")
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

type ServiceItem struct {
	ID            string
	ProviderID    string
	ProviderName  string
	CategoryID    string
	Title         string
	Description   string
	PriceAmount   int
	PriceUnit     string
	SupportOnline bool
	SortOrder     int
	Status        string
}

type Query struct {
	ProviderID string
	CategoryID string
	Status     string
	Keyword    string
	Page       int
	PageSize   int
}

func NormalizeStatus(status string) (string, error) {
	s := strings.TrimSpace(strings.ToLower(status))
	switch s {
	case "":
		return "", nil
	case StatusActive:
		return StatusActive, nil
	case StatusInactive:
		return StatusInactive, nil
	default:
		return "", ErrInvalidInput
	}
}
