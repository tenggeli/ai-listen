package model

import "time"

type Post struct {
	ID            uint64    `json:"id"`
	AuthorRole    int       `json:"authorRole"`
	AuthorID      uint64    `json:"authorId"`
	Content       string    `json:"content"`
	Topic         string    `json:"topic"`
	CityCode      string    `json:"cityCode"`
	VisibleScope  int       `json:"visibleScope"`
	LikeCount     int       `json:"likeCount"`
	CommentCount  int       `json:"commentCount"`
	FavoriteCount int       `json:"favoriteCount"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type PostComment struct {
	ID             uint64    `json:"id"`
	PostID         uint64    `json:"postId"`
	UserRole       int       `json:"userRole"`
	UserID         uint64    `json:"userId"`
	ReplyCommentID uint64    `json:"replyCommentId"`
	Content        string    `json:"content"`
	Status         int       `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Review struct {
	ID           uint64    `json:"id"`
	OrderID      uint64    `json:"orderId"`
	ReviewerRole int       `json:"reviewerRole"`
	ReviewerID   uint64    `json:"reviewerId"`
	TargetUserID uint64    `json:"targetUserId"`
	Score        int       `json:"score"`
	Content      string    `json:"content"`
	Images       []string  `json:"images"`
	IsAnonymous  bool      `json:"isAnonymous"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Complaint struct {
	ID              uint64    `json:"id"`
	OrderID         uint64    `json:"orderId"`
	ComplainantRole int       `json:"complainantRole"`
	ComplainantID   uint64    `json:"complainantId"`
	RespondentID    uint64    `json:"respondentId"`
	ComplaintType   string    `json:"complaintType"`
	Content         string    `json:"content"`
	EvidenceImages  []string  `json:"evidenceImages"`
	Status          int       `json:"status"`
	ResultRemark    string    `json:"resultRemark"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type CreatePostInput struct {
	Content      string
	Topic        string
	CityCode     string
	VisibleScope int
}

type CreatePostCommentInput struct {
	Content        string
	ReplyCommentID uint64
}

type CreateReviewInput struct {
	Score       int
	Content     string
	Images      []string
	IsAnonymous bool
}

type CreateComplaintInput struct {
	ComplaintType  string
	Content        string
	EvidenceImages []string
}

func clonePost(post *Post) *Post {
	if post == nil {
		return nil
	}
	copied := *post
	return &copied
}

func clonePostComment(comment *PostComment) *PostComment {
	if comment == nil {
		return nil
	}
	copied := *comment
	return &copied
}

func cloneReview(review *Review) *Review {
	if review == nil {
		return nil
	}
	copied := *review
	copied.Images = append([]string(nil), review.Images...)
	return &copied
}

func cloneComplaint(complaint *Complaint) *Complaint {
	if complaint == nil {
		return nil
	}
	copied := *complaint
	copied.EvidenceImages = append([]string(nil), complaint.EvidenceImages...)
	return &copied
}
