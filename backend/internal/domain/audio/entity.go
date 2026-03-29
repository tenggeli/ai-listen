package audio

type Category struct {
	Key   string
	Label string
}

type Track struct {
	ID            string
	Title         string
	Category      string
	PlayCountText string
	DurationText  string
	Emoji         string
	Author        string
}

type HomePage struct {
	Title               string
	Subtitle            string
	Categories          []Category
	RecommendedTracks   []Track
	CurrentTrackID      string
	CurrentProgressText string
	TotalDurationText   string
	IsPlaying           bool
}

func (h HomePage) HasTracks() bool {
	return len(h.RecommendedTracks) > 0
}
