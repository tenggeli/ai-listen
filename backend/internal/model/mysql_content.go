package model

import (
	"context"
	"database/sql"
)

func (s *MySQLStore) CreatePost(userID uint64, input CreatePostInput) (*Post, error) {
	ctx := context.Background()
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO posts(author_role, author_id, content, topic, city_code, visible_scope, status)
		VALUES (?, ?, ?, ?, ?, ?, 1)
	`, operatorRoleUser, userID, input.Content, input.Topic, input.CityCode, input.VisibleScope)
	if err != nil {
		return nil, err
	}
	postID, _ := res.LastInsertId()
	return s.GetPost(uint64(postID))
}

func (s *MySQLStore) Posts(page, pageSize int) ([]*Post, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	rows, err := s.db.Query(`
		SELECT id, author_role, author_id, content, topic, city_code, visible_scope, like_count, comment_count, favorite_count, status, created_at, updated_at
		FROM posts
		WHERE status = 1 AND deleted_at IS NULL
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(
			&post.ID, &post.AuthorRole, &post.AuthorID, &post.Content, &post.Topic, &post.CityCode,
			&post.VisibleScope, &post.LikeCount, &post.CommentCount, &post.FavoriteCount, &post.Status,
			&post.CreatedAt, &post.UpdatedAt,
		); err == nil {
			list = append(list, &post)
		}
	}
	return list, nil
}

func (s *MySQLStore) GetPost(postID uint64) (*Post, error) {
	var post Post
	err := s.db.QueryRow(`
		SELECT id, author_role, author_id, content, topic, city_code, visible_scope, like_count, comment_count, favorite_count, status, created_at, updated_at
		FROM posts
		WHERE id = ? AND status = 1 AND deleted_at IS NULL
	`, postID).Scan(
		&post.ID, &post.AuthorRole, &post.AuthorID, &post.Content, &post.Topic, &post.CityCode,
		&post.VisibleScope, &post.LikeCount, &post.CommentCount, &post.FavoriteCount, &post.Status,
		&post.CreatedAt, &post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrPostNotFound
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *MySQLStore) CreatePostComment(userID, postID uint64, input CreatePostCommentInput) (*PostComment, *Post, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	var exists uint64
	err = tx.QueryRowContext(ctx, `SELECT id FROM posts WHERE id = ? AND status = 1 AND deleted_at IS NULL FOR UPDATE`, postID).Scan(&exists)
	if err == sql.ErrNoRows {
		return nil, nil, ErrPostNotFound
	}
	if err != nil {
		return nil, nil, err
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO post_comments(post_id, user_role, user_id, reply_comment_id, content, status)
		VALUES (?, ?, ?, ?, ?, 1)
	`, postID, operatorRoleUser, userID, input.ReplyCommentID, input.Content)
	if err != nil {
		return nil, nil, err
	}
	commentID, _ := res.LastInsertId()

	if _, err := tx.ExecContext(ctx, `
		UPDATE posts
		SET comment_count = comment_count + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, postID); err != nil {
		return nil, nil, err
	}

	var comment PostComment
	err = tx.QueryRowContext(ctx, `
		SELECT id, post_id, user_role, user_id, reply_comment_id, content, status, created_at, updated_at
		FROM post_comments
		WHERE id = ?
	`, commentID).Scan(
		&comment.ID, &comment.PostID, &comment.UserRole, &comment.UserID, &comment.ReplyCommentID,
		&comment.Content, &comment.Status, &comment.CreatedAt, &comment.UpdatedAt,
	)
	if err != nil {
		return nil, nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}
	post, err := s.GetPost(postID)
	return &comment, post, err
}

func (s *MySQLStore) LikePost(userID, postID uint64) (*Post, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var exists uint64
	err = tx.QueryRowContext(ctx, `SELECT id FROM posts WHERE id = ? AND status = 1 AND deleted_at IS NULL FOR UPDATE`, postID).Scan(&exists)
	if err == sql.ErrNoRows {
		return nil, ErrPostNotFound
	}
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT IGNORE INTO post_likes(post_id, user_role, user_id)
		VALUES (?, ?, ?)
	`, postID, operatorRoleUser, userID)
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE posts
		SET like_count = (SELECT COUNT(1) FROM post_likes WHERE post_id = ?), updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, postID, postID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetPost(postID)
}

func (s *MySQLStore) UnlikePost(userID, postID uint64) (*Post, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var exists uint64
	err = tx.QueryRowContext(ctx, `SELECT id FROM posts WHERE id = ? AND status = 1 AND deleted_at IS NULL FOR UPDATE`, postID).Scan(&exists)
	if err == sql.ErrNoRows {
		return nil, ErrPostNotFound
	}
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM post_likes
		WHERE post_id = ? AND user_role = ? AND user_id = ?
	`, postID, operatorRoleUser, userID); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE posts
		SET like_count = (SELECT COUNT(1) FROM post_likes WHERE post_id = ?), updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, postID, postID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetPost(postID)
}

func (s *MySQLStore) CreateReview(userID, orderID uint64, input CreateReviewInput) (*Review, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var orderUserID, providerID uint64
	var orderStatus int
	err = tx.QueryRowContext(ctx, `
		SELECT user_id, provider_id, status
		FROM orders
		WHERE id = ? FOR UPDATE
	`, orderID).Scan(&orderUserID, &providerID, &orderStatus)
	if err == sql.ErrNoRows {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	if orderUserID != userID {
		return nil, ErrUnauthorized
	}
	if orderStatus != OrderStatusCompleted {
		return nil, ErrInvalidOrderStatus
	}

	var targetUserID uint64
	err = tx.QueryRowContext(ctx, `SELECT user_id FROM providers WHERE id = ? AND deleted_at IS NULL`, providerID).Scan(&targetUserID)
	if err == sql.ErrNoRows {
		return nil, ErrProviderNotFound
	}
	if err != nil {
		return nil, err
	}

	var existed uint64
	err = tx.QueryRowContext(ctx, `
		SELECT id FROM reviews
		WHERE order_id = ? AND reviewer_role = ? AND reviewer_id = ?
		LIMIT 1
	`, orderID, operatorRoleUser, userID).Scan(&existed)
	if err == nil {
		return nil, ErrReviewAlreadyExists
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO reviews(order_id, reviewer_role, reviewer_id, target_user_id, score, content, images, is_anonymous)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, orderID, operatorRoleUser, userID, targetUserID, input.Score, input.Content, toJSON(input.Images), boolToTinyint(input.IsAnonymous))
	if err != nil {
		return nil, err
	}
	reviewID, _ := res.LastInsertId()

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.loadReview(uint64(reviewID))
}

func (s *MySQLStore) ReviewsByOrder(orderID uint64) ([]*Review, error) {
	rows, err := s.db.Query(`
		SELECT id
		FROM reviews
		WHERE order_id = ?
		ORDER BY id ASC
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Review
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			continue
		}
		review, err := s.loadReview(id)
		if err == nil {
			list = append(list, review)
		}
	}
	return list, nil
}

func (s *MySQLStore) CreateComplaint(userID, orderID uint64, input CreateComplaintInput) (*Complaint, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var orderUserID, providerID uint64
	err = tx.QueryRowContext(ctx, `
		SELECT user_id, provider_id
		FROM orders
		WHERE id = ? FOR UPDATE
	`, orderID).Scan(&orderUserID, &providerID)
	if err == sql.ErrNoRows {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	if orderUserID != userID {
		return nil, ErrUnauthorized
	}

	var respondentID uint64
	err = tx.QueryRowContext(ctx, `SELECT user_id FROM providers WHERE id = ? AND deleted_at IS NULL`, providerID).Scan(&respondentID)
	if err == sql.ErrNoRows {
		return nil, ErrProviderNotFound
	}
	if err != nil {
		return nil, err
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO complaints(order_id, complainant_role, complainant_id, respondent_id, complaint_type, content, evidence_images, status, result_remark)
		VALUES (?, ?, ?, ?, ?, ?, ?, 10, '')
	`, orderID, operatorRoleUser, userID, respondentID, input.ComplaintType, input.Content, toJSON(input.EvidenceImages))
	if err != nil {
		return nil, err
	}
	complaintID, _ := res.LastInsertId()

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetComplaint(uint64(complaintID))
}

func (s *MySQLStore) GetComplaint(complaintID uint64) (*Complaint, error) {
	complaint := &Complaint{}
	var evidenceRaw sql.NullString
	err := s.db.QueryRow(`
		SELECT id, order_id, complainant_role, complainant_id, respondent_id, complaint_type, content, evidence_images, status, result_remark, created_at, updated_at
		FROM complaints
		WHERE id = ?
	`, complaintID).Scan(
		&complaint.ID, &complaint.OrderID, &complaint.ComplainantRole, &complaint.ComplainantID,
		&complaint.RespondentID, &complaint.ComplaintType, &complaint.Content, &evidenceRaw,
		&complaint.Status, &complaint.ResultRemark, &complaint.CreatedAt, &complaint.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrComplaintNotFound
	}
	if err != nil {
		return nil, err
	}
	complaint.EvidenceImages = parseJSONStringArray(evidenceRaw.String)
	return complaint, nil
}

func (s *MySQLStore) loadReview(reviewID uint64) (*Review, error) {
	review := &Review{}
	var imagesRaw sql.NullString
	err := s.db.QueryRow(`
		SELECT id, order_id, reviewer_role, reviewer_id, target_user_id, score, content, images, is_anonymous, created_at, updated_at
		FROM reviews
		WHERE id = ?
	`, reviewID).Scan(
		&review.ID, &review.OrderID, &review.ReviewerRole, &review.ReviewerID,
		&review.TargetUserID, &review.Score, &review.Content, &imagesRaw,
		&review.IsAnonymous, &review.CreatedAt, &review.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrReviewNotFound
	}
	if err != nil {
		return nil, err
	}
	review.Images = parseJSONStringArray(imagesRaw.String)
	return review, nil
}

func boolToTinyint(value bool) int {
	if value {
		return 1
	}
	return 0
}
