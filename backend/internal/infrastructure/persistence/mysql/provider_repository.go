package mysql

import (
	"context"
	"database/sql"
	"fmt"

	domain "listen/backend/internal/domain/provider"
)

type ProviderRepository struct {
	db *sql.DB
}

func NewProviderRepository(db *sql.DB) ProviderRepository {
	return ProviderRepository{db: db}
}

func (r ProviderRepository) List(ctx context.Context, query domain.Query) ([]domain.Provider, int, error) {
	args := make([]any, 0, 4)
	where := "WHERE 1=1"
	if query.ReviewStatus != "" {
		where += " AND p.review_status = ?"
		args = append(args, query.ReviewStatus)
	}

	countSQL := fmt.Sprintf(`SELECT COUNT(1) FROM providers p %s`, where)
	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf(`
SELECT p.id, COALESCE(pp.display_name, ''), COALESCE(pp.city_code, ''), COALESCE(pp.bio, ''), p.review_status
FROM providers p
LEFT JOIN provider_profiles pp ON pp.provider_id = p.id
%s
ORDER BY p.updated_at DESC
LIMIT ? OFFSET ?`, where)
	offset := (query.Page - 1) * query.PageSize
	args = append(args, query.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.Provider, 0)
	for rows.Next() {
		var p domain.Provider
		if err := rows.Scan(&p.ID, &p.DisplayName, &p.CityCode, &p.Bio, &p.ReviewStatus); err != nil {
			return nil, 0, err
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r ProviderRepository) GetByID(ctx context.Context, providerID string) (domain.Provider, error) {
	const query = `
SELECT p.id, COALESCE(pp.display_name, ''), COALESCE(pp.city_code, ''), COALESCE(pp.bio, ''), p.review_status
FROM providers p
LEFT JOIN provider_profiles pp ON pp.provider_id = p.id
WHERE p.id = ? LIMIT 1`
	var p domain.Provider
	err := r.db.QueryRowContext(ctx, query, providerID).Scan(&p.ID, &p.DisplayName, &p.CityCode, &p.Bio, &p.ReviewStatus)
	if err == sql.ErrNoRows {
		return domain.Provider{}, fmt.Errorf("provider not found")
	}
	if err != nil {
		return domain.Provider{}, err
	}
	return p, nil
}

func (r ProviderRepository) Save(ctx context.Context, provider domain.Provider, operator string, action string, reason string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const update = `UPDATE providers SET review_status = ?, updated_at = NOW() WHERE id = ?`
	if _, err = tx.ExecContext(ctx, update, provider.ReviewStatus, provider.ID); err != nil {
		return err
	}

	const insertRecord = `
INSERT INTO provider_review_records(provider_id, action, reason, operator_id, created_at, updated_at)
VALUES(?, ?, ?, ?, NOW(), NOW())`
	if _, err = tx.ExecContext(ctx, insertRecord, provider.ID, action, reason, operator); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
