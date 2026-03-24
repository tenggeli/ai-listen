package review

type CreateReviewRequest struct {
	Score       int      `json:"score"`
	Content     string   `json:"content"`
	Images      []string `json:"images"`
	IsAnonymous bool     `json:"isAnonymous"`
}
