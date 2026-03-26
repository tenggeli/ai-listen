package model

import (
	"fmt"
	"sort"
	"time"
)

func (s *MemoryStore) CreatePost(userID uint64, input CreatePostInput) (*Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	post := &Post{
		ID:            s.nextPostID,
		AuthorRole:    operatorRoleUser,
		AuthorID:      userID,
		Content:       input.Content,
		Topic:         input.Topic,
		CityCode:      input.CityCode,
		VisibleScope:  input.VisibleScope,
		LikeCount:     0,
		CommentCount:  0,
		FavoriteCount: 0,
		Status:        1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	s.posts[post.ID] = post
	s.nextPostID++
	return clonePost(post), nil
}

func (s *MemoryStore) Posts(page, pageSize int) ([]*Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*Post, 0, len(s.posts))
	for _, post := range s.posts {
		if post.Status == 1 {
			list = append(list, clonePost(post))
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID > list[j].ID })

	start, end := paginateRange(page, pageSize, len(list))
	if start >= len(list) {
		return []*Post{}, nil
	}
	return list[start:end], nil
}

func (s *MemoryStore) GetPost(postID uint64) (*Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[postID]
	if !ok || post.Status != 1 {
		return nil, ErrPostNotFound
	}
	return clonePost(post), nil
}

func (s *MemoryStore) CreatePostComment(userID, postID uint64, input CreatePostCommentInput) (*PostComment, *Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[postID]
	if !ok || post.Status != 1 {
		return nil, nil, ErrPostNotFound
	}

	now := time.Now()
	comment := &PostComment{
		ID:             s.nextCommentID,
		PostID:         postID,
		UserRole:       operatorRoleUser,
		UserID:         userID,
		ReplyCommentID: input.ReplyCommentID,
		Content:        input.Content,
		Status:         1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	s.postComments[comment.ID] = comment
	s.nextCommentID++

	post.CommentCount++
	post.UpdatedAt = now
	return clonePostComment(comment), clonePost(post), nil
}

func (s *MemoryStore) LikePost(userID, postID uint64) (*Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[postID]
	if !ok || post.Status != 1 {
		return nil, ErrPostNotFound
	}

	key := fmt.Sprintf("%d:%d:%d", postID, operatorRoleUser, userID)
	if _, exists := s.postLikes[key]; !exists {
		s.postLikes[key] = struct{}{}
		post.LikeCount++
		post.UpdatedAt = time.Now()
	}
	return clonePost(post), nil
}

func (s *MemoryStore) UnlikePost(userID, postID uint64) (*Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[postID]
	if !ok || post.Status != 1 {
		return nil, ErrPostNotFound
	}

	key := fmt.Sprintf("%d:%d:%d", postID, operatorRoleUser, userID)
	if _, exists := s.postLikes[key]; exists {
		delete(s.postLikes, key)
		if post.LikeCount > 0 {
			post.LikeCount--
		}
		post.UpdatedAt = time.Now()
	}
	return clonePost(post), nil
}

func (s *MemoryStore) CreateReview(userID, orderID uint64, input CreateReviewInput) (*Review, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrUnauthorized
	}
	if order.Status != OrderStatusCompleted {
		return nil, ErrInvalidOrderStatus
	}
	provider, ok := s.providers[order.ProviderID]
	if !ok {
		return nil, ErrProviderNotFound
	}

	uniqueKey := fmt.Sprintf("%d:%d", orderID, userID)
	if _, exists := s.reviewByOrderUser[uniqueKey]; exists {
		return nil, ErrReviewAlreadyExists
	}

	now := time.Now()
	review := &Review{
		ID:           s.nextReviewID,
		OrderID:      orderID,
		ReviewerRole: operatorRoleUser,
		ReviewerID:   userID,
		TargetUserID: provider.UserID,
		Score:        input.Score,
		Content:      input.Content,
		Images:       append([]string(nil), input.Images...),
		IsAnonymous:  input.IsAnonymous,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	s.reviews[review.ID] = review
	s.reviewByOrderUser[uniqueKey] = review.ID
	s.nextReviewID++
	return cloneReview(review), nil
}

func (s *MemoryStore) ReviewsByOrder(orderID uint64) ([]*Review, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*Review, 0)
	for _, review := range s.reviews {
		if review.OrderID == orderID {
			list = append(list, cloneReview(review))
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list, nil
}

func (s *MemoryStore) CreateComplaint(userID, orderID uint64, input CreateComplaintInput) (*Complaint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrUnauthorized
	}
	provider, ok := s.providers[order.ProviderID]
	if !ok {
		return nil, ErrProviderNotFound
	}

	now := time.Now()
	complaint := &Complaint{
		ID:              s.nextComplaintID,
		OrderID:         orderID,
		ComplainantRole: operatorRoleUser,
		ComplainantID:   userID,
		RespondentID:    provider.UserID,
		ComplaintType:   input.ComplaintType,
		Content:         input.Content,
		EvidenceImages:  append([]string(nil), input.EvidenceImages...),
		Status:          10,
		ResultRemark:    "",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	s.complaints[complaint.ID] = complaint
	s.nextComplaintID++
	return cloneComplaint(complaint), nil
}

func (s *MemoryStore) GetComplaint(complaintID uint64) (*Complaint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	complaint, ok := s.complaints[complaintID]
	if !ok {
		return nil, ErrComplaintNotFound
	}
	return cloneComplaint(complaint), nil
}

func paginateRange(page, pageSize, total int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	return start, end
}
