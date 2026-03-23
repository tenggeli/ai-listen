package dto

type SendCodeRequest struct {
	Mobile string `json:"mobile" binding:"required"`
}

type SendCodeResponse struct {
	DebugCode string `json:"debugCode,omitempty"`
}

type MobileLoginRequest struct {
	Mobile string `json:"mobile" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type MobileLoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
	UserID    uint64 `json:"userId"`
}
