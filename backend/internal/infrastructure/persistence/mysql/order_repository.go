package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	domain "listen/backend/internal/domain/order"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return OrderRepository{db: db}
}

func (r OrderRepository) Create(ctx context.Context, order domain.Order) error {
	const insertSQL = `
INSERT INTO order_records(
  order_id, user_id, provider_id, provider_name, service_item_id, service_item_title,
  amount, currency, status, created_at, paid_at, updated_at
) VALUES(
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW()
)`

	_, err := r.db.ExecContext(
		ctx,
		insertSQL,
		order.ID,
		order.UserID,
		order.ProviderID,
		order.ProviderName,
		order.ServiceItemID,
		order.ServiceItemTitle,
		order.Amount,
		order.Currency,
		order.Status,
		order.CreatedAt,
		nullablePaidAt(order.PaidAt),
	)
	return err
}

func (r OrderRepository) GetByID(ctx context.Context, orderID string) (domain.Order, error) {
	const query = `
SELECT order_id, user_id, provider_id, provider_name, service_item_id, service_item_title,
       amount, currency, status, created_at, paid_at
FROM order_records
WHERE order_id = ?
LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, strings.TrimSpace(orderID))
	item, err := scanOrder(row)
	if err == sql.ErrNoRows {
		return domain.Order{}, domain.ErrOrderNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	return item, nil
}

func (r OrderRepository) ListByUser(ctx context.Context, query domain.ListQuery) ([]domain.Order, int, error) {
	const countSQL = `SELECT COUNT(1) FROM order_records WHERE user_id = ?`

	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, query.UserID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const listSQL = `
SELECT order_id, user_id, provider_id, provider_name, service_item_id, service_item_title,
       amount, currency, status, created_at, paid_at
FROM order_records
WHERE user_id = ?
ORDER BY created_at DESC, id DESC
LIMIT ? OFFSET ?`

	offset := (query.Page - 1) * query.PageSize
	rows, err := r.db.QueryContext(ctx, listSQL, query.UserID, query.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.Order, 0)
	for rows.Next() {
		item, err := scanOrder(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r OrderRepository) ListByProvider(ctx context.Context, query domain.ProviderListQuery) ([]domain.Order, int, error) {
	const countSQL = `SELECT COUNT(1) FROM order_records WHERE provider_id = ?`

	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, query.ProviderID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const listSQL = `
SELECT order_id, user_id, provider_id, provider_name, service_item_id, service_item_title,
       amount, currency, status, created_at, paid_at
FROM order_records
WHERE provider_id = ?
ORDER BY created_at DESC, id DESC
LIMIT ? OFFSET ?`

	offset := (query.Page - 1) * query.PageSize
	rows, err := r.db.QueryContext(ctx, listSQL, query.ProviderID, query.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.Order, 0)
	for rows.Next() {
		item, err := scanOrder(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r OrderRepository) ListAll(ctx context.Context, query domain.AdminListQuery) ([]domain.Order, int, error) {
	where := "WHERE 1=1"
	args := make([]any, 0)

	if status := strings.TrimSpace(query.Status); status != "" {
		where += " AND status = ?"
		args = append(args, status)
	}
	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		where += " AND (order_id LIKE ? OR user_id LIKE ? OR provider_name LIKE ? OR service_item_title LIKE ?)"
		pattern := "%" + keyword + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}

	countSQL := fmt.Sprintf("SELECT COUNT(1) FROM order_records %s", where)
	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf(`
SELECT order_id, user_id, provider_id, provider_name, service_item_id, service_item_title,
       amount, currency, status, created_at, paid_at
FROM order_records
%s
ORDER BY created_at DESC, id DESC
LIMIT ? OFFSET ?`, where)

	offset := (query.Page - 1) * query.PageSize
	listArgs := append(append([]any{}, args...), query.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.Order, 0)
	for rows.Next() {
		item, err := scanOrder(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r OrderRepository) Save(ctx context.Context, order domain.Order) error {
	const updateSQL = `
UPDATE order_records
SET status = ?, paid_at = ?, updated_at = NOW()
WHERE order_id = ?`

	result, err := r.db.ExecContext(ctx, updateSQL, order.Status, nullablePaidAt(order.PaidAt), order.ID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanOrder(s scanner) (domain.Order, error) {
	var item domain.Order
	var paidAt sql.NullTime
	err := s.Scan(
		&item.ID,
		&item.UserID,
		&item.ProviderID,
		&item.ProviderName,
		&item.ServiceItemID,
		&item.ServiceItemTitle,
		&item.Amount,
		&item.Currency,
		&item.Status,
		&item.CreatedAt,
		&paidAt,
	)
	if err != nil {
		return domain.Order{}, err
	}
	item.PaidAt = nullTimePtr(paidAt)
	return item, nil
}

func nullablePaidAt(value *time.Time) any {
	if value == nil {
		return nil
	}
	return *value
}

func nullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}

var _ domain.Repository = OrderRepository{}
