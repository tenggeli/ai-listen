package admin

type LoginMockRequestDTO struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginMockResponseDTO struct {
	AccessToken string `json:"access_token"`
	AdminID     string `json:"admin_id"`
	Role        string `json:"role"`
}

type AdminMeResponseDTO struct {
	AdminID     string `json:"admin_id"`
	Account     string `json:"account"`
	Role        string `json:"role"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
}
