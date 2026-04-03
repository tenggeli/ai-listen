package user

type CreateOrderRequestDTO struct {
	ProviderID       string `json:"provider_id"`
	ProviderName     string `json:"provider_name"`
	ServiceItemID    string `json:"service_item_id"`
	ServiceItemTitle string `json:"service_item_title"`
	Amount           int    `json:"amount"`
	Currency         string `json:"currency"`
}

type OrderResponseDTO struct {
	ID               string  `json:"id"`
	UserID           string  `json:"user_id"`
	ProviderID       string  `json:"provider_id"`
	ProviderName     string  `json:"provider_name"`
	ServiceItemID    string  `json:"service_item_id"`
	ServiceItemTitle string  `json:"service_item_title"`
	Amount           int     `json:"amount"`
	Currency         string  `json:"currency"`
	Status           string  `json:"status"`
	StatusReason     string  `json:"status_reason"`
	CreatedAt        string  `json:"created_at"`
	PaidAt           *string `json:"paid_at,omitempty"`
}
