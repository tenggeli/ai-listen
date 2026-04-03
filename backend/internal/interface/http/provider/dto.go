package provider

type LoginMockRequestDTO struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginMockResponseDTO struct {
	AccessToken string `json:"access_token"`
	ProviderID  string `json:"provider_id"`
}

type ProviderMeResponseDTO struct {
	ProviderID  string `json:"provider_id"`
	Account     string `json:"account"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	CityCode    string `json:"city_code"`
}

type ProviderProfileUpdateRequestDTO struct {
	DisplayName string `json:"display_name"`
	CityCode    string `json:"city_code"`
}

type ProviderProfileResponseDTO struct {
	ProviderID  string `json:"provider_id"`
	Account     string `json:"account"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	CityCode    string `json:"city_code"`
}

type ProviderServiceItemDTO struct {
	ItemID        string `json:"item_id"`
	ProviderID    string `json:"provider_id"`
	CategoryID    string `json:"category_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	PriceAmount   int    `json:"price_amount"`
	PriceUnit     string `json:"price_unit"`
	SupportOnline bool   `json:"support_online"`
	SortOrder     int    `json:"sort_order"`
}

type OrderStatusActionResponseDTO struct {
	ID                 string  `json:"id"`
	Status             string  `json:"status"`
	StatusReason       string  `json:"status_reason"`
	StatusActionReason string  `json:"status_action_reason"`
	StatusUpdatedAt    *string `json:"status_updated_at,omitempty"`
}
