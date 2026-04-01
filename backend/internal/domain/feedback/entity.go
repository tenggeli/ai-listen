package feedback

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrFeedbackNotFound   = errors.New("feedback not found")
	ErrFeedbackForbidden  = errors.New("feedback forbidden")
	ErrFeedbackSubmitted  = errors.New("feedback already submitted")
	ErrOrderNotFinishable = errors.New("order is not ready for feedback")
)

type OrderFeedback struct {
	ID               string
	OrderID          string
	UserID           string
	RatingScore      int
	ReviewTags       []string
	ReviewContent    string
	HasComplaint     bool
	ComplaintReason  string
	ComplaintContent string
	CreatedAt        time.Time
}

func NewOrderFeedback(
	id string,
	orderID string,
	userID string,
	ratingScore int,
	reviewTags []string,
	reviewContent string,
	complaintReason string,
	complaintContent string,
	createdAt time.Time,
) (OrderFeedback, error) {
	id = strings.TrimSpace(id)
	orderID = strings.TrimSpace(orderID)
	userID = strings.TrimSpace(userID)
	reviewContent = strings.TrimSpace(reviewContent)
	complaintReason = strings.TrimSpace(complaintReason)
	complaintContent = strings.TrimSpace(complaintContent)
	tags := normalizeTags(reviewTags)

	hasReview := ratingScore > 0 || reviewContent != "" || len(tags) > 0
	hasComplaint := complaintReason != "" || complaintContent != ""
	if id == "" || orderID == "" || userID == "" {
		return OrderFeedback{}, ErrInvalidInput
	}
	if !hasReview && !hasComplaint {
		return OrderFeedback{}, ErrInvalidInput
	}
	if hasReview && (ratingScore < 1 || ratingScore > 10) {
		return OrderFeedback{}, ErrInvalidInput
	}
	if hasComplaint && complaintReason == "" {
		return OrderFeedback{}, ErrInvalidInput
	}

	return OrderFeedback{
		ID:               id,
		OrderID:          orderID,
		UserID:           userID,
		RatingScore:      ratingScore,
		ReviewTags:       tags,
		ReviewContent:    reviewContent,
		HasComplaint:     hasComplaint,
		ComplaintReason:  complaintReason,
		ComplaintContent: complaintContent,
		CreatedAt:        createdAt,
	}, nil
}

func normalizeTags(tags []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(tags))
	for _, item := range tags {
		clean := strings.TrimSpace(item)
		if clean == "" {
			continue
		}
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		result = append(result, clean)
	}
	return result
}
