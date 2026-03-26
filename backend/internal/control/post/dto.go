package post

type CreatePostRequest struct {
	Content      string `json:"content"`
	Topic        string `json:"topic"`
	CityCode     string `json:"cityCode"`
	VisibleScope int    `json:"visibleScope"`
}

type CreateCommentRequest struct {
	Content        string `json:"content"`
	ReplyCommentID uint64 `json:"replyCommentId"`
}
