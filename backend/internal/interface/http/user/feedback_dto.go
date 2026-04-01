package user

type SubmitFeedbackRequestDTO struct {
	RatingScore      int      `json:"rating_score"`
	ReviewTags       []string `json:"review_tags"`
	ReviewContent    string   `json:"review_content"`
	ComplaintReason  string   `json:"complaint_reason"`
	ComplaintContent string   `json:"complaint_content"`
}

type FeedbackResponseDTO struct {
	ID               string   `json:"id"`
	OrderID          string   `json:"order_id"`
	UserID           string   `json:"user_id"`
	RatingScore      int      `json:"rating_score"`
	ReviewTags       []string `json:"review_tags"`
	ReviewContent    string   `json:"review_content"`
	HasComplaint     bool     `json:"has_complaint"`
	ComplaintReason  string   `json:"complaint_reason"`
	ComplaintContent string   `json:"complaint_content"`
	CreatedAt        string   `json:"created_at"`
}
