package admin

type ReviewActionRequestDTO struct {
	Reason string `json:"reason"`
}

type ServiceItemStatusActionResponseDTO struct {
	ID            string `json:"id"`
	ServiceStatus string `json:"service_status"`
}
