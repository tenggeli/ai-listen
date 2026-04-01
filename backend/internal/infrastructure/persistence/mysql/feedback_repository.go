package mysql

import (
	"context"
	"database/sql"
	"encoding/json"

	gmysql "github.com/go-sql-driver/mysql"

	domain "listen/backend/internal/domain/feedback"
)

type FeedbackRepository struct {
	db *sql.DB
}

func NewFeedbackRepository(db *sql.DB) FeedbackRepository {
	return FeedbackRepository{db: db}
}

func (r FeedbackRepository) GetByOrderID(ctx context.Context, orderID string) (domain.OrderFeedback, error) {
	const query = `
SELECT feedback_id, order_id, user_id, rating_score, review_tags_json, review_content,
       has_complaint, complaint_reason, complaint_content, created_at
FROM order_feedback_records
WHERE order_id = ?
LIMIT 1`

	var item domain.OrderFeedback
	var tagsJSON string
	var hasComplaint int
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&item.ID,
		&item.OrderID,
		&item.UserID,
		&item.RatingScore,
		&tagsJSON,
		&item.ReviewContent,
		&hasComplaint,
		&item.ComplaintReason,
		&item.ComplaintContent,
		&item.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return domain.OrderFeedback{}, domain.ErrFeedbackNotFound
	}
	if err != nil {
		return domain.OrderFeedback{}, err
	}

	item.ReviewTags = parseFeedbackTags(tagsJSON)
	item.HasComplaint = hasComplaint == 1
	return item, nil
}

func (r FeedbackRepository) Create(ctx context.Context, item domain.OrderFeedback) error {
	tagsJSON, err := json.Marshal(item.ReviewTags)
	if err != nil {
		return err
	}

	const insertSQL = `
INSERT INTO order_feedback_records(
  feedback_id, order_id, user_id, rating_score, review_tags_json, review_content,
  has_complaint, complaint_reason, complaint_content, created_at, updated_at
) VALUES(
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW()
)`

	_, err = r.db.ExecContext(
		ctx,
		insertSQL,
		item.ID,
		item.OrderID,
		item.UserID,
		item.RatingScore,
		string(tagsJSON),
		item.ReviewContent,
		boolToInt(item.HasComplaint),
		item.ComplaintReason,
		item.ComplaintContent,
		item.CreatedAt,
	)
	if isMySQLDuplicateError(err) {
		return domain.ErrFeedbackSubmitted
	}
	return err
}

func parseFeedbackTags(raw string) []string {
	if raw == "" {
		return []string{}
	}
	items := make([]string, 0)
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return []string{}
	}
	return items
}

func isMySQLDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	mysqlErr, ok := err.(*gmysql.MySQLError)
	return ok && mysqlErr.Number == 1062
}

var _ domain.Repository = FeedbackRepository{}
