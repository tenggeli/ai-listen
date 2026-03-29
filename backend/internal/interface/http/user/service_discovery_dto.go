package user

type ListPublicProvidersRequestDTO struct {
	CategoryID string
	Keyword    string
	CityCode   string
	Page       int
	PageSize   int
}

type ServiceCategoryResponseDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type ProviderPublicResponseDTO struct {
	ID                string   `json:"id"`
	DisplayName       string   `json:"display_name"`
	AvatarURL         string   `json:"avatar_url"`
	CityCode          string   `json:"city_code"`
	Bio               string   `json:"bio"`
	RatingAvg         float64  `json:"rating_avg"`
	CompletedOrders   int      `json:"completed_orders"`
	Online            bool     `json:"online"`
	VerificationLabel string   `json:"verification_label"`
	Tags              []string `json:"tags"`
	PriceFrom         int      `json:"price_from"`
	PriceUnit         string   `json:"price_unit"`
}

type ServiceItemResponseDTO struct {
	ID            string `json:"id"`
	ProviderID    string `json:"provider_id"`
	CategoryID    string `json:"category_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	PriceAmount   int    `json:"price_amount"`
	PriceUnit     string `json:"price_unit"`
	SupportOnline bool   `json:"support_online"`
}
