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
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrInvalidPaymentStatus  = errors.New("invalid payment status")
	ErrPostNotFound          = errors.New("post not found")
	ErrReviewNotFound        = errors.New("review not found")
	ErrComplaintNotFound     = errors.New("complaint not found")
	ErrReviewAlreadyExists   = errors.New("review already exists")
	ErrPaymentCallbackKeyUse = errors.New("payment callback idempotency key already used by another payment")
	ErrInvalidSMSCode        = errors.New("invalid sms code")
	ErrServiceItemNotFound   = errors.New("service item not found")
)
