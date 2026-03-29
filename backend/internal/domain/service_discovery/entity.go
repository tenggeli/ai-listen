package service_discovery

type ServiceCategory struct {
	ID        string
	Name      string
	Icon      string
	SortOrder int
}

type ProviderPublicProfile struct {
	ID               string
	DisplayName      string
	AvatarURL        string
	CityCode         string
	Bio              string
	RatingAvg        float64
	CompletedOrders  int
	Online           bool
	VerificationText string
	Tags             []string
	PriceFrom        int
	PriceUnit        string
}

type ServiceItem struct {
	ID            string
	ProviderID    string
	CategoryID    string
	Title         string
	Description   string
	PriceAmount   int
	PriceUnit     string
	SupportOnline bool
	SortOrder     int
}
