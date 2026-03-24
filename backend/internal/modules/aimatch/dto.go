package aimatch

type MatchRequest struct {
	InputType string `json:"inputType"`
	Content   string `json:"content"`
	CityCode  string `json:"cityCode"`
}
