package user

import (
	"errors"
	"net/http"
	"strconv"

	app "listen/backend/internal/application/audio"
	domain "listen/backend/internal/domain/audio"
)

type SoundController struct {
	getSoundPageUC app.GetSoundPageUseCase
}

func NewSoundController(getSoundPageUC app.GetSoundPageUseCase) SoundController {
	return SoundController{getSoundPageUC: getSoundPageUC}
}

func (c SoundController) HandleGetSounds(w http.ResponseWriter, r *http.Request) {
	pageNo, _ := strconv.Atoi(r.URL.Query().Get("page_no"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	output, err := c.getSoundPageUC.Execute(r.Context(), app.GetSoundPageInput{
		Page:        r.URL.Query().Get("page"),
		UserID:      r.URL.Query().Get("user_id"),
		CategoryKey: r.URL.Query().Get("category_key"),
		PageNo:      pageNo,
		PageSize:    pageSize,
	})
	if err != nil {
		writeSoundError(w, err)
		return
	}

	categories := make([]map[string]any, 0, len(output.Page.Categories))
	for _, item := range output.Page.Categories {
		categories = append(categories, map[string]any{
			"key":   item.Key,
			"label": item.Label,
		})
	}

	tracks := make([]map[string]any, 0, len(output.Page.RecommendedTracks))
	for _, item := range output.Page.RecommendedTracks {
		tracks = append(tracks, map[string]any{
			"id":              item.ID,
			"title":           item.Title,
			"category":        item.Category,
			"play_count_text": item.PlayCountText,
			"duration_text":   item.DurationText,
			"emoji":           item.Emoji,
			"author":          item.Author,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"title":                 output.Page.Title,
		"subtitle":              output.Page.Subtitle,
		"categories":            categories,
		"recommended_tracks":    tracks,
		"current_track_id":      output.Page.CurrentTrackID,
		"current_progress_text": output.Page.CurrentProgressText,
		"total_duration_text":   output.Page.TotalDurationText,
		"is_playing":            output.Page.IsPlaying,
	})
}

func writeSoundError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidPage):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidCategory):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
