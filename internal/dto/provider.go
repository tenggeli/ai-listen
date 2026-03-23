package dto

type ProviderUserBrief struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   uint8  `json:"gender"`
	CityID   uint64 `json:"cityId"`
	Bio      string `json:"bio"`
	VipLevel uint8  `json:"vipLevel"`
}

type ServiceItemInfo struct {
	ID          uint64  `json:"id"`
	ProviderID  uint64  `json:"providerId"`
	CategoryID  uint64  `json:"categoryId"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unitPrice"`
	BillingType uint8   `json:"billingType"`
	MinHours    int     `json:"minHours"`
	MaxHours    int     `json:"maxHours"`
	UnitName    string  `json:"unitName"`
	Status      uint8   `json:"status"`
}

type ProviderListItem struct {
	ProviderID    uint64            `json:"providerId"`
	ProviderNo    string            `json:"providerNo"`
	Score         float64           `json:"score"`
	TotalOrders   int               `json:"totalOrders"`
	ServiceStatus uint8             `json:"serviceStatus"`
	OnlineStatus  uint8             `json:"onlineStatus"`
	CityID        uint64            `json:"cityId"`
	Intro         string            `json:"intro"`
	Tags          []string          `json:"tags"`
	User          ProviderUserBrief `json:"user"`
	ServiceItems  []ServiceItemInfo `json:"serviceItems"`
}

type ProviderListResponse struct {
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
	Total    int64              `json:"total"`
	List     []ProviderListItem `json:"list"`
}

type ProviderBaseInfo struct {
	ID             uint64   `json:"id"`
	UserID         uint64   `json:"userId"`
	ProviderNo     string   `json:"providerNo"`
	Zodiac         string   `json:"zodiac"`
	Constellation  string   `json:"constellation"`
	Level          string   `json:"level"`
	Score          float64  `json:"score"`
	TotalOrders    int      `json:"totalOrders"`
	TotalIncome    float64  `json:"totalIncome"`
	ServiceStatus  uint8    `json:"serviceStatus"`
	OnlineStatus   uint8    `json:"onlineStatus"`
	CityID         uint64   `json:"cityId"`
	Intro          string   `json:"intro"`
	Tags           []string `json:"tags"`
	CommissionRate float64  `json:"commissionRate"`
	ComplaintCount int      `json:"complaintCount"`
	CancelCount    int      `json:"cancelCount"`
}

type ProviderDetailResponse struct {
	Provider     ProviderBaseInfo  `json:"provider"`
	User         ProviderUserBrief `json:"user"`
	ServiceItems []ServiceItemInfo `json:"serviceItems"`
}
