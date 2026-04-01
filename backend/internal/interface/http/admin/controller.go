package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	app "listen/backend/internal/application/provider"
)

type ProviderController struct {
	listUC   app.ListReviewProvidersUseCase
	detailUC app.GetProviderDetailUseCase
	reviewUC app.ReviewProviderUseCase
}

func NewProviderController(
	listUC app.ListReviewProvidersUseCase,
	detailUC app.GetProviderDetailUseCase,
	reviewUC app.ReviewProviderUseCase,
) ProviderController {
	return ProviderController{listUC: listUC, detailUC: detailUC, reviewUC: reviewUC}
}

func (c ProviderController) HandleList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	output, err := c.listUC.Execute(r.Context(), app.ListReviewProvidersInput{
		ReviewStatus: r.URL.Query().Get("review_status"),
		Page:         page,
		PageSize:     pageSize,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	items := make([]map[string]any, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, map[string]any{
			"id":            item.ID,
			"display_name":  item.DisplayName,
			"city_code":     item.CityCode,
			"review_status": item.ReviewStatus,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": output.Total})
}

func (c ProviderController) HandleDetail(w http.ResponseWriter, r *http.Request) {
	providerID := r.PathValue("id")
	if providerID == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid provider id")
		return
	}
	output, err := c.detailUC.Execute(r.Context(), app.GetProviderDetailInput{ProviderID: providerID})
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}
	item := output.Provider
	writeJSON(w, http.StatusOK, map[string]any{
		"id":            item.ID,
		"display_name":  item.DisplayName,
		"city_code":     item.CityCode,
		"bio":           item.Bio,
		"review_status": item.ReviewStatus,
	})
}

func (c ProviderController) HandleApprove(w http.ResponseWriter, r *http.Request) {
	c.handleAction(w, r, "approve")
}

func (c ProviderController) HandleReject(w http.ResponseWriter, r *http.Request) {
	c.handleAction(w, r, "reject")
}

func (c ProviderController) HandleRequireSupplement(w http.ResponseWriter, r *http.Request) {
	c.handleAction(w, r, "require_supplement")
}

func (c ProviderController) handleAction(w http.ResponseWriter, r *http.Request, action string) {
	providerID := r.PathValue("id")
	if providerID == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid provider id")
		return
	}

	var body ReviewActionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil && err.Error() != "EOF" {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.reviewUC.Execute(r.Context(), app.ReviewActionInput{
		ProviderID: providerID,
		Action:     action,
		Operator:   currentAdminID(r),
		Reason:     body.Reason,
	})
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": output.Provider.ID, "review_status": output.Provider.ReviewStatus})
}
