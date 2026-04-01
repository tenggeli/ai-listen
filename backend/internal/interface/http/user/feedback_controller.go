package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	app "listen/backend/internal/application/feedback"
	domain "listen/backend/internal/domain/feedback"
)

type FeedbackController struct {
	submitUC app.SubmitOrderFeedbackUseCase
	getUC    app.GetOrderFeedbackUseCase
}

func NewFeedbackController(
	submitUC app.SubmitOrderFeedbackUseCase,
	getUC app.GetOrderFeedbackUseCase,
) FeedbackController {
	return FeedbackController{
		submitUC: submitUC,
		getUC:    getUC,
	}
}

func (c FeedbackController) HandleSubmitFeedback(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)
	orderID := strings.TrimSpace(r.PathValue("id"))

	var body SubmitFeedbackRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.submitUC.Execute(r.Context(), app.SubmitFeedbackInput{
		UserID:           userID,
		OrderID:          orderID,
		RatingScore:      body.RatingScore,
		ReviewTags:       body.ReviewTags,
		ReviewContent:    body.ReviewContent,
		ComplaintReason:  body.ComplaintReason,
		ComplaintContent: body.ComplaintContent,
	})
	if err != nil {
		writeFeedbackError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildFeedbackResponse(output.Item))
}

func (c FeedbackController) HandleGetFeedback(w http.ResponseWriter, r *http.Request) {
	userID := currentUserID(r)
	orderID := strings.TrimSpace(r.PathValue("id"))

	output, err := c.getUC.Execute(r.Context(), app.GetFeedbackInput{
		UserID:  userID,
		OrderID: orderID,
	})
	if err != nil {
		writeFeedbackError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildFeedbackResponse(output.Item))
}

func buildFeedbackResponse(item domain.OrderFeedback) FeedbackResponseDTO {
	return FeedbackResponseDTO{
		ID:               item.ID,
		OrderID:          item.OrderID,
		UserID:           item.UserID,
		RatingScore:      item.RatingScore,
		ReviewTags:       item.ReviewTags,
		ReviewContent:    item.ReviewContent,
		HasComplaint:     item.HasComplaint,
		ComplaintReason:  item.ComplaintReason,
		ComplaintContent: item.ComplaintContent,
		CreatedAt:        item.CreatedAt.Format(time.RFC3339),
	}
}

func writeFeedbackError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrFeedbackForbidden):
		writeJSONError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, domain.ErrFeedbackNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrFeedbackSubmitted):
		writeJSONError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrOrderNotFinishable):
		writeJSONError(w, http.StatusConflict, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}
