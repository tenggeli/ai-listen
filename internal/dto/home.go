package dto

type AIMatchRequest struct {
	InputType string  `json:"inputType"`
	Content   string  `json:"content"`
	CityID    uint64  `json:"cityId"`
	Lng       float64 `json:"lng"`
	Lat       float64 `json:"lat"`
}

type AIMatchItem struct {
	ProviderID  uint64  `json:"providerId"`
	Nickname    string  `json:"nickname"`
	Score       float64 `json:"score"`
	CityID      uint64  `json:"cityId"`
	MatchReason string  `json:"matchReason"`
}

type AIMatchResponse struct {
	List []AIMatchItem `json:"list"`
}
