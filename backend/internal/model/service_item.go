package model

type ServiceItem struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Unit     string `json:"unit"`
	MinPrice int64  `json:"minPrice"`
	MaxPrice int64  `json:"maxPrice"`
	Status   int    `json:"status"`
}
