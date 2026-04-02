package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	domain "listen/backend/internal/domain/service_item_admin"
)

type ServiceItemAdminRepository struct {
	db *sql.DB
}

func NewServiceItemAdminRepository(db *sql.DB) ServiceItemAdminRepository {
	return ServiceItemAdminRepository{db: db}
}

func (r ServiceItemAdminRepository) List(ctx context.Context, query domain.Query) ([]domain.ServiceItem, int, error) {
	where := "WHERE 1=1"
	args := make([]any, 0, 8)
	if query.ProviderID != "" {
		where += " AND si.provider_id = ?"
		args = append(args, query.ProviderID)
	}
	if query.CategoryID != "" {
		where += " AND si.category_id = ?"
		args = append(args, query.CategoryID)
	}
	if query.Status != "" {
		where += " AND si.service_status = ?"
		args = append(args, query.Status)
	}
	if query.Keyword != "" {
		where += " AND (si.title LIKE ? OR si.description LIKE ? OR pp.display_name LIKE ?)"
		keyword := "%" + query.Keyword + "%"
		args = append(args, keyword, keyword, keyword)
	}

	countSQL := fmt.Sprintf(`
SELECT COUNT(1)
FROM provider_service_items si
LEFT JOIN provider_profiles pp ON pp.provider_id = si.provider_id
%s`, where)

	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf(`
SELECT
  si.item_id,
  si.provider_id,
  COALESCE(pp.display_name, ''),
  si.category_id,
  si.title,
  si.description,
  si.price_amount,
  si.price_unit,
  si.support_online,
  si.sort_order,
  si.service_status
FROM provider_service_items si
LEFT JOIN provider_profiles pp ON pp.provider_id = si.provider_id
%s
ORDER BY si.updated_at DESC, si.id DESC
LIMIT ? OFFSET ?`, where)

	offset := (query.Page - 1) * query.PageSize
	args = append(args, query.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.ServiceItem, 0)
	for rows.Next() {
		var item domain.ServiceItem
		var supportOnline int
		if err := rows.Scan(
			&item.ID,
			&item.ProviderID,
			&item.ProviderName,
			&item.CategoryID,
			&item.Title,
			&item.Description,
			&item.PriceAmount,
			&item.PriceUnit,
			&supportOnline,
			&item.SortOrder,
			&item.Status,
		); err != nil {
			return nil, 0, err
		}
		item.SupportOnline = supportOnline == 1
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r ServiceItemAdminRepository) GetByID(ctx context.Context, serviceItemID string) (domain.ServiceItem, error) {
	serviceItemID = strings.TrimSpace(serviceItemID)
	if serviceItemID == "" {
		return domain.ServiceItem{}, domain.ErrInvalidInput
	}

	const query = `
SELECT
  si.item_id,
  si.provider_id,
  COALESCE(pp.display_name, ''),
  si.category_id,
  si.title,
  si.description,
  si.price_amount,
  si.price_unit,
  si.support_online,
  si.sort_order,
  si.service_status
FROM provider_service_items si
LEFT JOIN provider_profiles pp ON pp.provider_id = si.provider_id
WHERE si.item_id = ?
LIMIT 1`

	var item domain.ServiceItem
	var supportOnline int
	err := r.db.QueryRowContext(ctx, query, serviceItemID).Scan(
		&item.ID,
		&item.ProviderID,
		&item.ProviderName,
		&item.CategoryID,
		&item.Title,
		&item.Description,
		&item.PriceAmount,
		&item.PriceUnit,
		&supportOnline,
		&item.SortOrder,
		&item.Status,
	)
	if err == sql.ErrNoRows {
		return domain.ServiceItem{}, domain.ErrServiceItemNotFound
	}
	if err != nil {
		return domain.ServiceItem{}, err
	}
	item.SupportOnline = supportOnline == 1
	return item, nil
}

func (r ServiceItemAdminRepository) UpdateStatus(ctx context.Context, serviceItemID string, status string, _ string) error {
	serviceItemID = strings.TrimSpace(serviceItemID)
	status, err := domain.NormalizeStatus(status)
	if err != nil {
		return err
	}
	if serviceItemID == "" || status == "" {
		return domain.ErrInvalidInput
	}

	result, err := r.db.ExecContext(ctx, `
UPDATE provider_service_items
SET service_status = ?, updated_at = NOW()
WHERE item_id = ?`, status, serviceItemID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrServiceItemNotFound
	}
	return nil
}

var _ domain.Repository = ServiceItemAdminRepository{}
