package audio

type PlayLogRequest struct {
	ProgressSec int `json:"progressSec"`
	PositionSec int `json:"positionSec"`
	DurationSec int `json:"durationSec"`
}

type FavoriteRequest struct {
	Favorite bool `json:"favorite"`
}
