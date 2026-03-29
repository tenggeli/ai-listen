package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	domain "listen/backend/internal/domain/service_discovery"
)

type ServiceDiscoveryRepository struct {
	db *sql.DB
}

func NewServiceDiscoveryRepository(db *sql.DB) ServiceDiscoveryRepository {
	return ServiceDiscoveryRepository{db: db}
}

func (r ServiceDiscoveryRepository) ListCategories(ctx context.Context) ([]domain.ServiceCategory, error) {
	const query = `
SELECT category_id, name, icon, sort_order
FROM service_categories
WHERE status = 'active'
ORDER BY sort_order ASC, category_id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.ServiceCategory, 0)
	for rows.Next() {
		var item domain.ServiceCategory
		if err := rows.Scan(&item.ID, &item.Name, &item.Icon, &item.SortOrder); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r ServiceDiscoveryRepository) ListPublicProviders(ctx context.Context, query domain.PublicProviderQuery) ([]domain.ProviderPublicProfile, int, error) {
	where := "WHERE p.review_status = 'approved'"
	args := make([]any, 0, 8)
	if query.CityCode != "" {
		where += " AND pp.city_code = ?"
		args = append(args, query.CityCode)
	}

	if strings.TrimSpace(query.Keyword) != "" {
		where += " AND (pp.display_name LIKE ? OR pp.bio LIKE ? OR EXISTS (SELECT 1 FROM provider_public_tags ppt WHERE ppt.provider_id = p.id AND ppt.tag_name LIKE ?))"
		k := "%" + strings.TrimSpace(query.Keyword) + "%"
		args = append(args, k, k, k)
	}
	if query.CategoryID != "" && query.CategoryID != "cat_all" {
		where += " AND EXISTS (SELECT 1 FROM provider_service_items psi WHERE psi.provider_id = p.id AND psi.category_id = ? AND psi.service_status = 'active')"
		args = append(args, query.CategoryID)
	}

	countSQL := fmt.Sprintf(`SELECT COUNT(1) FROM providers p LEFT JOIN provider_profiles pp ON pp.provider_id = p.id %s`, where)
	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf(`
SELECT
  p.id,
  COALESCE(pp.display_name, ''),
  COALESCE(pp.avatar_url, ''),
  COALESCE(pp.city_code, ''),
  COALESCE(pp.bio, ''),
  p.rating_avg,
  p.order_completed_count,
  CASE
    WHEN p.provider_status = 'active' AND pr.last_seen_at IS NOT NULL AND pr.last_seen_at >= DATE_SUB(NOW(), INTERVAL 5 MINUTE) THEN 1
    ELSE 0
  END AS is_online,
  COALESCE(pp.verification_text, ''),
  COALESCE(min_price.price_amount, 0),
  COALESCE(min_price.price_unit, '')
FROM providers p
LEFT JOIN provider_profiles pp ON pp.provider_id = p.id
LEFT JOIN (
  SELECT provider_id, MIN(price_amount) AS price_amount, SUBSTRING_INDEX(GROUP_CONCAT(price_unit ORDER BY price_amount ASC, sort_order ASC, item_id ASC), ',', 1) AS price_unit
  FROM provider_service_items
  WHERE service_status = 'active'
  GROUP BY provider_id
) min_price ON min_price.provider_id = p.id
LEFT JOIN provider_presence pr ON pr.provider_id = p.id
%s
ORDER BY p.rating_avg DESC, p.order_completed_count DESC
LIMIT ? OFFSET ?`, where)

	offset := (query.Page - 1) * query.PageSize
	args = append(args, query.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.ProviderPublicProfile, 0)
	for rows.Next() {
		var item domain.ProviderPublicProfile
		var isOnline int
		if err := rows.Scan(
			&item.ID,
			&item.DisplayName,
			&item.AvatarURL,
			&item.CityCode,
			&item.Bio,
			&item.RatingAvg,
			&item.CompletedOrders,
			&isOnline,
			&item.VerificationText,
			&item.PriceFrom,
			&item.PriceUnit,
		); err != nil {
			return nil, 0, err
		}
		item.Online = isOnline == 1
		tags, err := r.loadTags(ctx, item.ID)
		if err != nil {
			return nil, 0, err
		}
		item.Tags = tags
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r ServiceDiscoveryRepository) GetPublicProviderByID(ctx context.Context, providerID string) (domain.ProviderPublicProfile, error) {
	const query = `
SELECT
  p.id,
  COALESCE(pp.display_name, ''),
  COALESCE(pp.avatar_url, ''),
  COALESCE(pp.city_code, ''),
  COALESCE(pp.bio, ''),
  p.rating_avg,
  p.order_completed_count,
  CASE
    WHEN p.provider_status = 'active' AND pr.last_seen_at IS NOT NULL AND pr.last_seen_at >= DATE_SUB(NOW(), INTERVAL 5 MINUTE) THEN 1
    ELSE 0
  END AS is_online,
  COALESCE(pp.verification_text, ''),
  COALESCE(min_price.price_amount, 0),
  COALESCE(min_price.price_unit, '')
FROM providers p
LEFT JOIN provider_profiles pp ON pp.provider_id = p.id
LEFT JOIN (
  SELECT provider_id, MIN(price_amount) AS price_amount, SUBSTRING_INDEX(GROUP_CONCAT(price_unit ORDER BY price_amount ASC, sort_order ASC, item_id ASC), ',', 1) AS price_unit
  FROM provider_service_items
  WHERE service_status = 'active'
  GROUP BY provider_id
) min_price ON min_price.provider_id = p.id
LEFT JOIN provider_presence pr ON pr.provider_id = p.id
WHERE p.id = ? AND p.review_status = 'approved'
LIMIT 1`

	var item domain.ProviderPublicProfile
	var isOnline int
	err := r.db.QueryRowContext(ctx, query, providerID).Scan(
		&item.ID,
		&item.DisplayName,
		&item.AvatarURL,
		&item.CityCode,
		&item.Bio,
		&item.RatingAvg,
		&item.CompletedOrders,
		&isOnline,
		&item.VerificationText,
		&item.PriceFrom,
		&item.PriceUnit,
	)
	if err == sql.ErrNoRows {
		return domain.ProviderPublicProfile{}, domain.ErrProviderNotFound
	}
	if err != nil {
		return domain.ProviderPublicProfile{}, err
	}
	item.Online = isOnline == 1
	tags, err := r.loadTags(ctx, item.ID)
	if err != nil {
		return domain.ProviderPublicProfile{}, err
	}
	item.Tags = tags
	return item, nil
}

func (r ServiceDiscoveryRepository) ListProviderServiceItems(ctx context.Context, providerID string) ([]domain.ServiceItem, error) {
	if _, err := r.GetPublicProviderByID(ctx, providerID); err != nil {
		return nil, err
	}

	const query = `
SELECT item_id, provider_id, category_id, title, description, price_amount, price_unit, support_online, sort_order
FROM provider_service_items
WHERE provider_id = ? AND service_status = 'active'
ORDER BY sort_order ASC, item_id ASC`

	rows, err := r.db.QueryContext(ctx, query, providerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.ServiceItem, 0)
	for rows.Next() {
		var item domain.ServiceItem
		var supportOnline int
		if err := rows.Scan(
			&item.ID,
			&item.ProviderID,
			&item.CategoryID,
			&item.Title,
			&item.Description,
			&item.PriceAmount,
			&item.PriceUnit,
			&supportOnline,
			&item.SortOrder,
		); err != nil {
			return nil, err
		}
		item.SupportOnline = supportOnline == 1
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r ServiceDiscoveryRepository) loadTags(ctx context.Context, providerID string) ([]string, error) {
	const query = `
SELECT tag_name
FROM provider_public_tags
WHERE provider_id = ?
ORDER BY sort_order ASC, id ASC`

	rows, err := r.db.QueryContext(ctx, query, providerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

var _ domain.Repository = ServiceDiscoveryRepository{}
