package ecode

type Code struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	Success              = Code{Code: 0, Message: "ok"}
	BadRequest           = Code{Code: 40000, Message: "request params invalid"}
	Unauthorized         = Code{Code: 40100, Message: "unauthorized"}
	Forbidden            = Code{Code: 40300, Message: "forbidden"}
	NotFound             = Code{Code: 40400, Message: "resource not found"}
	Conflict             = Code{Code: 40900, Message: "state conflict"}
	TooManyRequests      = Code{Code: 42900, Message: "too many requests"}
	InternalServerError  = Code{Code: 50000, Message: "internal server error"}
	SMSCodeInvalid       = Code{Code: 100101, Message: "sms code invalid or expired"}
	UserAlreadyExists    = Code{Code: 100102, Message: "user already exists"}
	ProviderAuditFailed  = Code{Code: 100201, Message: "provider audit not passed"}
	ProviderUnavailable  = Code{Code: 100202, Message: "provider unavailable"}
	OrderStateInvalid    = Code{Code: 100301, Message: "order state invalid"}
	OrderNotFound        = Code{Code: 100302, Message: "order not found"}
	PaymentStateInvalid  = Code{Code: 100303, Message: "payment state invalid"}
	BalanceInsufficient  = Code{Code: 100304, Message: "balance insufficient"}
	WithdrawAuditing     = Code{Code: 100401, Message: "withdraw under review"}
	ContentViolation     = Code{Code: 100501, Message: "content violation"}
	AIMatchLimitExceeded = Code{Code: 100601, Message: "ai match limit exceeded"}
)
