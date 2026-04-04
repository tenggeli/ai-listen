package admin

type ReviewActionRequestDTO struct {
	Reason string `json:"reason"`
}

type ServiceItemStatusActionResponseDTO struct {
	ID            string `json:"id"`
	ServiceStatus string `json:"service_status"`
}

type AdminSoundUpsertRequestDTO struct {
	TrackID       string `json:"track_id"`
	CategoryKey   string `json:"category_key"`
	Title         string `json:"title"`
	PlayCountText string `json:"play_count_text"`
	DurationText  string `json:"duration_text"`
	Emoji         string `json:"emoji"`
	Author        string `json:"author"`
	SortOrder     int    `json:"sort_order"`
	Status        string `json:"status"`
}
