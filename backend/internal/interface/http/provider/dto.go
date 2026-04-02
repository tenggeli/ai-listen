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

type OrderStatusActionResponseDTO struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
