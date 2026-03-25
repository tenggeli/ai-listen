package model

import "errors"

var (
	ErrUnauthorized          = errors.New("unauthorized")
	ErrUserNotFound          = errors.New("user not found")
	ErrAdminNotFound         = errors.New("admin user not found")
	ErrProviderNotFound      = errors.New("provider not found")
	ErrRoleNotFound          = errors.New("admin role not found")
	ErrProviderAuditRequired = errors.New("provider audit not passed")
	ErrOrderNotFound         = errors.New("order not found")
	ErrInvalidOrderStatus    = errors.New("invalid order status")
	ErrInvalidSMSCode        = errors.New("invalid sms code")
	ErrServiceItemNotFound   = errors.New("service item not found")
)
