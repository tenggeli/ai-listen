package mysql

import (
	"context"
	"database/sql"

	app "listen/backend/internal/application/admin_order"
)

type OrderAdminActionRepository struct {
	db *sql.DB
}

func NewOrderAdminActionRepository(db *sql.DB) OrderAdminActionRepository {
	return OrderAdminActionRepository{db: db}
}

func (r OrderAdminActionRepository) Create(ctx context.Context, item app.ActionAudit) error {
	const insertSQL = `
INSERT INTO order_admin_action_records(
  action_id, order_id, scope, action, operator_id, reason, status_before, status_after, created_at, updated_at
) VALUES(
  ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW()
)`

	_, err := r.db.ExecContext(
		ctx,
		insertSQL,
		item.ActionID,
		item.OrderID,
		item.Scope,
		item.Action,
		item.Operator,
		item.Reason,
		item.StatusBefore,
		item.StatusAfter,
		item.UpdatedAt,
	)
	return err
}

func (r OrderAdminActionRepository) ListByOrderID(ctx context.Context, orderID string) ([]app.ActionAudit, error) {
	const query = `
SELECT action_id, order_id, scope, action, operator_id, reason, COALESCE(status_before, ''), COALESCE(status_after, ''), created_at
FROM order_admin_action_records
WHERE order_id = ?
ORDER BY created_at DESC, id DESC`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]app.ActionAudit, 0)
	for rows.Next() {
		var item app.ActionAudit
		if err := rows.Scan(
			&item.ActionID,
			&item.OrderID,
			&item.Scope,
			&item.Action,
			&item.Operator,
			&item.Reason,
			&item.StatusBefore,
			&item.StatusAfter,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

var _ app.ActionLogRepository = OrderAdminActionRepository{}
