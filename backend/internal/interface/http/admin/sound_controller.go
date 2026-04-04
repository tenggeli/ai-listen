package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	app "listen/backend/internal/application/audio"
	domain "listen/backend/internal/domain/audio"
)

type SoundController struct {
	listUC   app.ListAdminSoundsUseCase
	createUC app.CreateAdminSoundUseCase
	updateUC app.UpdateAdminSoundUseCase
	statusUC app.UpdateAdminSoundStatusUseCase
}

func NewSoundController(
	listUC app.ListAdminSoundsUseCase,
	createUC app.CreateAdminSoundUseCase,
	updateUC app.UpdateAdminSoundUseCase,
	statusUC app.UpdateAdminSoundStatusUseCase,
) SoundController {
	return SoundController{
		listUC:   listUC,
		createUC: createUC,
		updateUC: updateUC,
		statusUC: statusUC,
	}
}

func (c SoundController) HandleList(w http.ResponseWriter, r *http.Request) {
	pageNo, _ := strconv.Atoi(r.URL.Query().Get("page_no"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	output, err := c.listUC.Execute(r.Context(), app.ListAdminSoundsInput{
		CategoryKey: r.URL.Query().Get("category_key"),
		Status:      r.URL.Query().Get("status"),
		Keyword:     r.URL.Query().Get("keyword"),
		PageNo:      pageNo,
		PageSize:    pageSize,
	})
	if err != nil {
		writeSoundError(w, err)
		return
	}
	items := make([]map[string]any, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, map[string]any{
			"id":              item.ID,
			"category_key":    item.CategoryKey,
			"title":           item.Title,
			"play_count_text": item.PlayCountText,
			"duration_text":   item.DurationText,
			"emoji":           item.Emoji,
			"author":          item.Author,
			"sort_order":      item.SortOrder,
			"status":          item.Status,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": output.Total})
}

func (c SoundController) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var body AdminSoundUpsertRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	output, err := c.createUC.Execute(r.Context(), app.CreateAdminSoundInput{
		TrackID:       body.TrackID,
		CategoryKey:   body.CategoryKey,
		Title:         body.Title,
		PlayCountText: body.PlayCountText,
		DurationText:  body.DurationText,
		Emoji:         body.Emoji,
		Author:        body.Author,
		SortOrder:     body.SortOrder,
		Status:        body.Status,
	})
	if err != nil {
		writeSoundError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildSoundItemResponse(output.Item))
}

func (c SoundController) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	trackID := r.PathValue("id")
	if trackID == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid sound id")
		return
	}
	var body AdminSoundUpsertRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	output, err := c.updateUC.Execute(r.Context(), app.UpdateAdminSoundInput{
		TrackID:       trackID,
		CategoryKey:   body.CategoryKey,
		Title:         body.Title,
		PlayCountText: body.PlayCountText,
		DurationText:  body.DurationText,
		Emoji:         body.Emoji,
		Author:        body.Author,
		SortOrder:     body.SortOrder,
	})
	if err != nil {
		writeSoundError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildSoundItemResponse(output.Item))
}

func (c SoundController) HandleActivate(w http.ResponseWriter, r *http.Request) {
	c.handleStatusAction(w, r, "activate")
}

func (c SoundController) HandleDeactivate(w http.ResponseWriter, r *http.Request) {
	c.handleStatusAction(w, r, "deactivate")
}

func (c SoundController) handleStatusAction(w http.ResponseWriter, r *http.Request, action string) {
	trackID := r.PathValue("id")
	if trackID == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid sound id")
		return
	}
	output, err := c.statusUC.Execute(r.Context(), app.UpdateAdminSoundStatusInput{
		TrackID: trackID,
		Action:  action,
	})
	if err != nil {
		writeSoundError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":     output.Item.ID,
		"status": output.Item.Status,
	})
}

func buildSoundItemResponse(item domain.AdminTrack) map[string]any {
	return map[string]any{
		"id":              item.ID,
		"category_key":    item.CategoryKey,
		"title":           item.Title,
		"play_count_text": item.PlayCountText,
		"duration_text":   item.DurationText,
		"emoji":           item.Emoji,
		"author":          item.Author,
		"sort_order":      item.SortOrder,
		"status":          item.Status,
	}
}

func writeSoundError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidCategory):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrSoundNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
