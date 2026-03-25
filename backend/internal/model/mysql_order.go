package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (s *MySQLStore) CreateOrder(userID uint64, req CreateOrderInput) (*Order, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var providerAuditStatus int
	err = tx.QueryRowContext(ctx, `SELECT audit_status FROM providers WHERE id = ? AND deleted_at IS NULL`, req.ProviderID).Scan(&providerAuditStatus)
	if err == sql.ErrNoRows {
		return nil, ErrProviderNotFound
	}
	if err != nil {
		return nil, err
	}
	if providerAuditStatus != 2 {
		return nil, ErrProviderAuditRequired
	}

	var item ServiceItem
	err = tx.QueryRowContext(ctx, `
		SELECT id, name, category, unit, min_price, max_price, status
		FROM service_items
		WHERE id = ?
	`, req.ServiceItemID).Scan(&item.ID, &item.Name, &item.Category, &item.Unit, &item.MinPrice, &item.MaxPrice, &item.Status)
	if err == sql.ErrNoRows {
		return nil, ErrServiceItemNotFound
	}
	if err != nil {
		return nil, err
	}

	payAmount := item.MinPrice
	var providerPrice sql.NullInt64
	err = tx.QueryRowContext(ctx, `
		SELECT price_amount
		FROM provider_service_items
		WHERE provider_id = ? AND service_item_id = ? AND status = 1
		LIMIT 1
	`, req.ProviderID, req.ServiceItemID).Scan(&providerPrice)
	if err == nil && providerPrice.Valid {
		payAmount = providerPrice.Int64
	}

	orderNo := fmt.Sprintf("L%014d", time.Now().UnixNano()%1e14)
	var plannedStartAt any
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", req.PlannedStartAt, time.Local); err == nil {
		plannedStartAt = t
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO orders(
			order_no, user_id, provider_id, service_item_id, scene_text, city_code, address_text,
			planned_start_at, unit_price, planned_duration, pay_amount, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, orderNo, userID, req.ProviderID, req.ServiceItemID, req.SceneText, req.CityCode, req.AddressText, plannedStartAt, payAmount, req.PlannedDuration, payAmount, OrderStatusPendingPayment)
	if err != nil {
		return nil, err
	}
	orderID64, _ := res.LastInsertId()
	orderID := uint64(orderID64)

	if err := s.insertOrderLogTx(ctx, tx, orderID, 0, OrderStatusPendingPayment, operatorRoleUser, userID, "order created"); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetOrder(orderID)
}

func (s *MySQLStore) OrdersByUser(userID uint64) []*Order {
	return s.listOrders(`SELECT id FROM orders WHERE user_id = ? AND deleted_at IS NULL ORDER BY id DESC`, userID)
}

func (s *MySQLStore) OrdersByProvider(userID uint64) ([]*Order, error) {
	provider, err := s.ProviderByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.listOrders(`SELECT id FROM orders WHERE provider_id = ? AND deleted_at IS NULL ORDER BY id DESC`, provider.ID), nil
}

func (s *MySQLStore) GetOrder(orderID uint64) (*Order, error) {
	return s.loadOrder(orderID)
}

func (s *MySQLStore) CreatePayment(orderID uint64) (*Payment, *Order, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	var order Order
	err = tx.QueryRowContext(ctx, `
		SELECT id, order_no, user_id, provider_id, service_item_id, status, pay_amount
		FROM orders
		WHERE id = ? FOR UPDATE
	`, orderID).Scan(&order.ID, &order.OrderNo, &order.UserID, &order.ProviderID, &order.ServiceItemID, &order.Status, &order.PayAmount)
	if err == sql.ErrNoRows {
		return nil, nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, nil, err
	}
	if order.Status != OrderStatusPendingPayment {
		return nil, nil, ErrInvalidOrderStatus
	}

	now := time.Now()
	paymentNo := fmt.Sprintf("PAY%014d", now.UnixNano()%1e14)
	res, err := tx.ExecContext(ctx, `
		INSERT INTO payments(payment_no, order_id, user_id, pay_channel, pay_amount, status, paid_at)
		VALUES (?, ?, ?, 1, ?, 20, ?)
	`, paymentNo, orderID, order.UserID, order.PayAmount, now)
	if err != nil {
		return nil, nil, err
	}
	paymentID, _ := res.LastInsertId()

	if _, err := tx.ExecContext(ctx, `UPDATE orders SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, OrderStatusPendingAccept, orderID); err != nil {
		return nil, nil, err
	}
	if err := s.insertOrderLogTx(ctx, tx, orderID, OrderStatusPendingPayment, OrderStatusPendingAccept, operatorRoleUser, order.UserID, "payment completed"); err != nil {
		return nil, nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	payment := &Payment{ID: uint64(paymentID), OrderID: orderID, PayAmount: order.PayAmount, Status: 20, CreatedAt: now, PaidAt: now}
	loadedOrder, err := s.GetOrder(orderID)
	return payment, loadedOrder, err
}

func (s *MySQLStore) CancelOrder(orderID, userID uint64, reason string) (*Order, error) {
	return s.transitionUserOrder(orderID, userID, []int{OrderStatusPendingPayment, OrderStatusPendingAccept, OrderStatusAccepted}, OrderStatusCanceled, reason, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE orders SET cancel_reason = ? WHERE id = ?`, reason, orderID)
		return err
	})
}

func (s *MySQLStore) ProviderAcceptOrder(userID, orderID uint64) (*Order, error) {
	return s.transitionProviderOrder(userID, orderID, OrderStatusPendingAccept, OrderStatusAccepted, "provider accepted", `accepted_at = CURRENT_TIMESTAMP`)
}

func (s *MySQLStore) ProviderDepartOrder(userID, orderID uint64) (*Order, error) {
	return s.transitionProviderOrder(userID, orderID, OrderStatusAccepted, OrderStatusDeparted, "provider departed", "")
}

func (s *MySQLStore) ProviderArriveOrder(userID, orderID uint64) (*Order, error) {
	return s.transitionProviderOrder(userID, orderID, OrderStatusDeparted, OrderStatusArrived, "provider arrived", `arrived_at = CURRENT_TIMESTAMP`)
}

func (s *MySQLStore) StartOrder(userID, orderID uint64, remark string) (*Order, error) {
	return s.transitionUserOrder(orderID, userID, []int{OrderStatusArrived}, OrderStatusServing, remark, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE orders SET actual_start_at = CURRENT_TIMESTAMP WHERE id = ?`, orderID)
		return err
	})
}

func (s *MySQLStore) ProviderFinishOrder(userID, orderID uint64) (*Order, error) {
	return s.transitionProviderOrder(userID, orderID, OrderStatusServing, OrderStatusPendingFinish, "provider finished service", `actual_end_at = CURRENT_TIMESTAMP, finished_at = CURRENT_TIMESTAMP`)
}

func (s *MySQLStore) ConfirmFinishOrder(userID, orderID uint64) (*Order, error) {
	return s.transitionUserOrder(orderID, userID, []int{OrderStatusPendingFinish}, OrderStatusCompleted, "user confirmed finish", nil)
}

func (s *MySQLStore) listOrders(query string, arg uint64) []*Order {
	rows, err := s.db.Query(query, arg)
	if err != nil {
		return []*Order{}
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			continue
		}
		order, err := s.loadOrder(id)
		if err == nil {
			orders = append(orders, order)
		}
	}
	return orders
}

func (s *MySQLStore) loadOrder(orderID uint64) (*Order, error) {
	order := &Order{}
	var plannedStartAt, acceptedAt, arrivedAt, startedAt, finishedAt sql.NullTime
	err := s.db.QueryRow(`
		SELECT id, order_no, user_id, provider_id, service_item_id, scene_text, city_code, address_text,
		       planned_start_at, planned_duration, status, pay_amount, cancel_reason, created_at, updated_at,
		       accepted_at, arrived_at, actual_start_at, finished_at
		FROM orders
		WHERE id = ? AND deleted_at IS NULL
	`, orderID).Scan(
		&order.ID, &order.OrderNo, &order.UserID, &order.ProviderID, &order.ServiceItemID, &order.SceneText, &order.CityCode,
		&order.AddressText, &plannedStartAt, &order.PlannedDuration, &order.Status, &order.PayAmount, &order.CancelReason,
		&order.CreatedAt, &order.UpdatedAt, &acceptedAt, &arrivedAt, &startedAt, &finishedAt,
	)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	order.PlannedStartAt = formatNullDateTime(plannedStartAt)
	order.AcceptedAt = ptrTime(acceptedAt)
	order.ArrivedAt = ptrTime(arrivedAt)
	order.StartedAt = ptrTime(startedAt)
	order.FinishedAt = ptrTime(finishedAt)
	order.Logs = s.loadOrderLogs(orderID)
	return order, nil
}

func (s *MySQLStore) loadOrderLogs(orderID uint64) []OrderLog {
	rows, err := s.db.Query(`
		SELECT from_status, to_status, operator_role, operator_id, remark, created_at
		FROM order_status_logs
		WHERE order_id = ?
		ORDER BY id ASC
	`, orderID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var logs []OrderLog
	for rows.Next() {
		var log OrderLog
		var roleCode int
		if err := rows.Scan(&log.FromStatus, &log.ToStatus, &roleCode, &log.OperatorID, &log.Remark, &log.CreatedAt); err == nil {
			log.OperatorRole = operatorRoleName(roleCode)
			logs = append(logs, log)
		}
	}
	return logs
}

func (s *MySQLStore) insertOrderLogTx(ctx context.Context, tx *sql.Tx, orderID uint64, fromStatus, toStatus, roleCode int, operatorID uint64, remark string) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO order_status_logs(order_id, from_status, to_status, operator_role, operator_id, remark)
		VALUES (?, ?, ?, ?, ?, ?)
	`, orderID, fromStatus, toStatus, roleCode, operatorID, remark)
	return err
}

func (s *MySQLStore) transitionProviderOrder(userID, orderID uint64, expectedStatus, targetStatus int, remark, extraSet string) (*Order, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var providerID uint64
	var auditStatus int
	err = tx.QueryRowContext(ctx, `SELECT id, audit_status FROM providers WHERE user_id = ?`, userID).Scan(&providerID, &auditStatus)
	if err == sql.ErrNoRows {
		return nil, ErrProviderNotFound
	}
	if err != nil {
		return nil, err
	}
	if auditStatus != 2 {
		return nil, ErrProviderAuditRequired
	}

	var currentStatus int
	var orderProviderID uint64
	err = tx.QueryRowContext(ctx, `SELECT provider_id, status FROM orders WHERE id = ? FOR UPDATE`, orderID).Scan(&orderProviderID, &currentStatus)
	if err == sql.ErrNoRows {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	if orderProviderID != providerID {
		return nil, ErrUnauthorized
	}
	if currentStatus != expectedStatus {
		return nil, ErrInvalidOrderStatus
	}

	updateSQL := `UPDATE orders SET status = ?, updated_at = CURRENT_TIMESTAMP`
	if extraSet != "" {
		updateSQL += `, ` + extraSet
	}
	updateSQL += ` WHERE id = ?`
	if _, err := tx.ExecContext(ctx, updateSQL, targetStatus, orderID); err != nil {
		return nil, err
	}
	if err := s.insertOrderLogTx(ctx, tx, orderID, expectedStatus, targetStatus, operatorRoleProvider, userID, remark); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetOrder(orderID)
}

func (s *MySQLStore) transitionUserOrder(orderID, userID uint64, allowedFrom []int, targetStatus int, remark string, beforeUpdate func(context.Context, *sql.Tx) error) (*Order, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var currentStatus int
	var orderUserID uint64
	err = tx.QueryRowContext(ctx, `SELECT user_id, status FROM orders WHERE id = ? FOR UPDATE`, orderID).Scan(&orderUserID, &currentStatus)
	if err == sql.ErrNoRows {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	if orderUserID != userID {
		return nil, ErrUnauthorized
	}
	if !containsInt(allowedFrom, currentStatus) {
		return nil, ErrInvalidOrderStatus
	}

	if beforeUpdate != nil {
		if err := beforeUpdate(ctx, tx); err != nil {
			return nil, err
		}
	}
	if _, err := tx.ExecContext(ctx, `UPDATE orders SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, targetStatus, orderID); err != nil {
		return nil, err
	}
	if err := s.insertOrderLogTx(ctx, tx, orderID, currentStatus, targetStatus, operatorRoleUser, userID, remark); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetOrder(orderID)
}
