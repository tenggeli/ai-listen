package order

type CreateOrderRequest struct {
	ProviderID      uint64 `json:"providerId"`
	ServiceItemID   uint64 `json:"serviceItemId"`
	SceneText       string `json:"sceneText"`
	CityCode        string `json:"cityCode"`
	AddressText     string `json:"addressText"`
	PlannedStartAt  string `json:"plannedStartAt"`
	PlannedDuration int    `json:"plannedDuration"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason"`
}

type StartOrderRequest struct {
	ConfirmType int    `json:"confirmType"`
	ConfirmCode string `json:"confirmCode"`
}
