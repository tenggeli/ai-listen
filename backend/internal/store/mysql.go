package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"ai-listen/backend/migrations"
	_ "github.com/go-sql-driver/mysql"
)

const (
	operatorRoleUser     = 1
	operatorRoleProvider = 2
	operatorRoleAdmin    = 3
)

type MySQLStore struct {
	db           *sql.DB
	adminTokenMu sync.RWMutex
	adminTokens  map[string]uint64
}

func NewMySQLStore(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	store := &MySQLStore{
		db:          db,
		adminTokens: map[string]uint64{},
	}
	if err := store.initSchema(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.seedServiceItems(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.seedAdminBootstrap(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *MySQLStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *MySQLStore) initSchema(ctx context.Context) error {
	for _, stmt := range splitSQLStatements(migrations.InitSchemaSQL) {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("exec schema: %w", err)
		}
	}
	return nil
}

func (s *MySQLStore) seedServiceItems(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM service_items").Scan(&count); err != nil {
		return fmt.Errorf("count service_items: %w", err)
	}
	if count > 0 {
		return nil
	}

	items := []struct {
		name     string
		category string
		unit     string
		minPrice int64
		maxPrice int64
		sort     int
	}{
		{name: "陪吃饭", category: "陪伴", unit: "小时", minPrice: 10000, maxPrice: 50000, sort: 10},
		{name: "观影搭子", category: "娱乐", unit: "小时", minPrice: 8000, maxPrice: 40000, sort: 20},
		{name: "心理疏导", category: "情绪支持", unit: "小时", minPrice: 20000, maxPrice: 100000, sort: 30},
	}
	for _, item := range items {
		if _, err := s.db.ExecContext(ctx, `
			INSERT INTO service_items(name, category, unit, min_price, max_price, sort, status)
			VALUES (?, ?, ?, ?, ?, ?, 1)
		`, item.name, item.category, item.unit, item.minPrice, item.maxPrice, item.sort); err != nil {
			return fmt.Errorf("seed service_items: %w", err)
		}
	}
	return nil
}

func (s *MySQLStore) IssueSMSCode(mobile string) string {
	code := "123456"
	expiresAt := time.Now().Add(10 * time.Minute)
	_, _ = s.db.Exec(`
		INSERT INTO auth_sms_codes(mobile, code, expires_at)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE code = VALUES(code), expires_at = VALUES(expires_at), updated_at = CURRENT_TIMESTAMP
	`, mobile, code, expiresAt)
	_, _ = s.db.Exec(`
		INSERT INTO sms_logs(mobile, scene, template_code, status, provider_msg_id)
		VALUES (?, 'login', 'debug', 20, '')
	`, mobile)
	return code
}

func (s *MySQLStore) LoginBySMS(mobile, code string) (*User, string, string, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, "", "", err
	}
	defer tx.Rollback()

	var storedCode string
	var expiresAt time.Time
	err = tx.QueryRowContext(ctx, `
		SELECT code, expires_at
		FROM auth_sms_codes
		WHERE mobile = ?
	`, mobile).Scan(&storedCode, &expiresAt)
	if err != nil || storedCode != code || expiresAt.Before(time.Now()) {
		return nil, "", "", ErrInvalidSMSCode
	}

	user, err := s.findOrCreateUserByMobileTx(ctx, tx, mobile)
	if err != nil {
		return nil, "", "", err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO user_auths(user_id, auth_type, auth_key)
		VALUES (?, 'sms', ?)
		ON DUPLICATE KEY UPDATE user_id = VALUES(user_id), updated_at = CURRENT_TIMESTAMP
	`, user.ID, mobile); err != nil {
		return nil, "", "", err
	}

	accessToken := fmt.Sprintf("token-%d-%d", user.ID, time.Now().UnixNano())
	refreshToken := fmt.Sprintf("refresh-%d-%d", user.ID, time.Now().UnixNano())
	if err := s.insertTokenTx(ctx, tx, user.ID, accessToken, "access", 24*time.Hour); err != nil {
		return nil, "", "", err
	}
	if err := s.insertTokenTx(ctx, tx, user.ID, refreshToken, "refresh", 7*24*time.Hour); err != nil {
		return nil, "", "", err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id = ?`, user.ID); err != nil {
		return nil, "", "", err
	}
	if err := tx.Commit(); err != nil {
		return nil, "", "", err
	}
	return user, accessToken, refreshToken, nil
}

func (s *MySQLStore) RefreshToken(refreshToken string) (string, error) {
	token := strings.TrimSpace(refreshToken)
	var userID uint64
	var expiresAt time.Time
	err := s.db.QueryRow(`
		SELECT user_id, expires_at
		FROM auth_tokens
		WHERE token = ? AND token_type = 'refresh' AND status = 1
	`, token).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		return "", ErrUnauthorized
	}

	accessToken := fmt.Sprintf("token-%d-%d", userID, time.Now().UnixNano())
	if _, err := s.db.Exec(`
		INSERT INTO auth_tokens(user_id, token, token_type, expires_at, status)
		VALUES (?, ?, 'access', ?, 1)
	`, userID, accessToken, time.Now().Add(24*time.Hour)); err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *MySQLStore) AdminLogin(username, password string) (*AdminUser, string, error) {
	var adminID uint64
	var passwordHash string
	err := s.db.QueryRow(`
		SELECT id, password_hash
		FROM admin_users
		WHERE username = ? AND status = 1
	`, username).Scan(&adminID, &passwordHash)
	if err != nil || passwordHash != password {
		return nil, "", ErrUnauthorized
	}

	token := fmt.Sprintf("admin-token-%d-%d", adminID, time.Now().UnixNano())
	s.adminTokenMu.Lock()
	s.adminTokens[token] = adminID
	s.adminTokenMu.Unlock()

	_, _ = s.db.Exec(`UPDATE admin_users SET last_login_at = CURRENT_TIMESTAMP WHERE id = ?`, adminID)

	admin, err := s.AdminUserByID(adminID)
	if err != nil {
		return nil, "", err
	}
	return admin, token, nil
}

func (s *MySQLStore) AdminByToken(raw string) (*AdminUser, error) {
	token := strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
	s.adminTokenMu.RLock()
	adminUserID, ok := s.adminTokens[token]
	s.adminTokenMu.RUnlock()
	if !ok {
		return nil, ErrUnauthorized
	}
	return s.AdminUserByID(adminUserID)
}

func (s *MySQLStore) AdminUserByID(adminUserID uint64) (*AdminUser, error) {
	admin := &AdminUser{}
	err := s.db.QueryRow(`
		SELECT id, username, nickname, status, created_at, updated_at
		FROM admin_users
		WHERE id = ? AND status = 1
	`, adminUserID).Scan(&admin.ID, &admin.Username, &admin.Nickname, &admin.Status, &admin.CreatedAt, &admin.UpdatedAt)
	if err != nil {
		return nil, ErrAdminNotFound
	}

	rows, err := s.db.Query(`
		SELECT r.role_key
		FROM admin_user_roles ur
		JOIN admin_roles r ON r.id = ur.admin_role_id
		WHERE ur.admin_user_id = ? AND r.status = 1
		ORDER BY r.id ASC
	`, adminUserID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var roleKey string
			if scanErr := rows.Scan(&roleKey); scanErr == nil {
				admin.Roles = append(admin.Roles, roleKey)
			}
		}
	}
	return admin, nil
}

func (s *MySQLStore) AdminPermissionsByRoles(roleKeys []string) []string {
	if len(roleKeys) == 0 {
		return []string{}
	}

	for _, roleKey := range roleKeys {
		if roleKey == "super_admin" {
			permissions := s.AdminPermissions()
			result := make([]string, 0, len(permissions))
			for _, item := range permissions {
				result = append(result, item.PermissionKey)
			}
			return result
		}
	}

	placeholders := strings.TrimRight(strings.Repeat("?,", len(roleKeys)), ",")
	args := make([]any, 0, len(roleKeys))
	for _, roleKey := range roleKeys {
		args = append(args, roleKey)
	}
	rows, err := s.db.Query(`
		SELECT DISTINCT p.permission_key
		FROM admin_roles r
		JOIN admin_role_permissions rp ON rp.admin_role_id = r.id
		JOIN admin_permissions p ON p.id = rp.admin_permission_id
		WHERE r.role_key IN (`+placeholders+`) AND r.status = 1 AND p.status = 1
		ORDER BY p.id ASC
	`, args...)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var permissionKey string
		if scanErr := rows.Scan(&permissionKey); scanErr == nil {
			result = append(result, permissionKey)
		}
	}
	return result
}

func (s *MySQLStore) AdminHasPermission(adminUserID uint64, permission string) bool {
	var exists int
	err := s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM admin_user_roles ur
			JOIN admin_roles r ON r.id = ur.admin_role_id
			LEFT JOIN admin_role_permissions rp ON rp.admin_role_id = r.id
			LEFT JOIN admin_permissions p ON p.id = rp.admin_permission_id
			WHERE ur.admin_user_id = ? AND r.status = 1
			  AND (r.role_key = 'super_admin' OR (p.permission_key = ? AND p.status = 1))
		)
	`, adminUserID, permission).Scan(&exists)
	return err == nil && exists == 1
}

func (s *MySQLStore) AdminRoles() []*AdminRole {
	rows, err := s.db.Query(`
		SELECT id, name, role_key, status, created_at, updated_at
		FROM admin_roles
		WHERE status = 1
		ORDER BY id ASC
	`)
	if err != nil {
		return []*AdminRole{}
	}
	defer rows.Close()

	var result []*AdminRole
	for rows.Next() {
		role := &AdminRole{}
		if scanErr := rows.Scan(&role.ID, &role.Name, &role.RoleKey, &role.Status, &role.CreatedAt, &role.UpdatedAt); scanErr != nil {
			continue
		}
		result = append(result, role)
	}

	for _, role := range result {
		permissionRows, queryErr := s.db.Query(`
			SELECT p.permission_key
			FROM admin_role_permissions rp
			JOIN admin_permissions p ON p.id = rp.admin_permission_id
			WHERE rp.admin_role_id = ? AND p.status = 1
			ORDER BY p.id ASC
		`, role.ID)
		if queryErr != nil {
			continue
		}
		for permissionRows.Next() {
			var permissionKey string
			if scanErr := permissionRows.Scan(&permissionKey); scanErr == nil {
				role.PermissionKeys = append(role.PermissionKeys, permissionKey)
			}
		}
		permissionRows.Close()
	}
	return result
}

func (s *MySQLStore) AdminPermissions() []*AdminPermission {
	rows, err := s.db.Query(`
		SELECT id, name, permission_key, category, status, created_at, updated_at
		FROM admin_permissions
		WHERE status = 1
		ORDER BY id ASC
	`)
	if err != nil {
		return []*AdminPermission{}
	}
	defer rows.Close()

	var result []*AdminPermission
	for rows.Next() {
		permission := &AdminPermission{}
		if scanErr := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.PermissionKey,
			&permission.Category,
			&permission.Status,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		); scanErr == nil {
			result = append(result, permission)
		}
	}
	return result
}

func (s *MySQLStore) UpdateAdminUserRoles(adminUserID uint64, roleKeys []string) (*AdminUser, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var exists int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM admin_users WHERE id = ? AND status = 1`, adminUserID).Scan(&exists); err != nil || exists == 0 {
		return nil, ErrAdminNotFound
	}

	normalized := make([]string, 0, len(roleKeys))
	seen := map[string]struct{}{}
	for _, roleKey := range roleKeys {
		if _, duplicated := seen[roleKey]; duplicated {
			continue
		}
		seen[roleKey] = struct{}{}
		normalized = append(normalized, roleKey)
	}

	roleIDs := make([]uint64, 0, len(normalized))
	for _, roleKey := range normalized {
		var roleID uint64
		if err := tx.QueryRowContext(ctx, `SELECT id FROM admin_roles WHERE role_key = ? AND status = 1`, roleKey).Scan(&roleID); err != nil {
			return nil, ErrRoleNotFound
		}
		roleIDs = append(roleIDs, roleID)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_user_roles WHERE admin_user_id = ?`, adminUserID); err != nil {
		return nil, err
	}
	for _, roleID := range roleIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO admin_user_roles(admin_user_id, admin_role_id)
			VALUES (?, ?)
		`, adminUserID, roleID); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.AdminUserByID(adminUserID)
}

func (s *MySQLStore) UserByToken(raw string) (*User, error) {
	token := strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
	user := &User{}
	var birthday sql.NullTime
	err := s.db.QueryRow(`
		SELECT u.id, u.mobile, u.nickname, u.avatar, u.gender, u.birthday, u.city_code, u.created_at, u.updated_at
		FROM auth_tokens t
		JOIN users u ON u.id = t.user_id
		WHERE t.token = ? AND t.token_type = 'access' AND t.status = 1 AND t.expires_at > CURRENT_TIMESTAMP
	`, token).Scan(
		&user.ID, &user.Mobile, &user.Nickname, &user.Avatar, &user.Gender, &birthday, &user.CityCode, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, ErrUnauthorized
	}
	user.Birthday = formatNullDate(birthday)
	return user, nil
}

func (s *MySQLStore) GetUser(userID uint64) (*User, error) {
	return s.getUser(userID)
}

func (s *MySQLStore) UpdateUser(userID uint64, updater func(*User)) (*User, error) {
	user, err := s.getUser(userID)
	if err != nil {
		return nil, err
	}
	updater(user)

	var birthday any
	if user.Birthday != "" {
		birthday = user.Birthday
	}

	if _, err := s.db.Exec(`
		UPDATE users
		SET nickname = ?, avatar = ?, gender = ?, birthday = ?, city_code = ?
		WHERE id = ?
	`, user.Nickname, user.Avatar, user.Gender, birthday, user.CityCode, userID); err != nil {
		return nil, err
	}
	return s.getUser(userID)
}

func (s *MySQLStore) ServiceItems() []*ServiceItem {
	rows, err := s.db.Query(`
		SELECT id, name, category, unit, min_price, max_price, status
		FROM service_items
		WHERE status = 1
		ORDER BY sort ASC, id ASC
	`)
	if err != nil {
		return []*ServiceItem{}
	}
	defer rows.Close()

	var items []*ServiceItem
	for rows.Next() {
		item := &ServiceItem{}
		if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.Unit, &item.MinPrice, &item.MaxPrice, &item.Status); err == nil {
			items = append(items, item)
		}
	}
	return items
}

func (s *MySQLStore) GetServiceItem(id uint64) (*ServiceItem, error) {
	item := &ServiceItem{}
	err := s.db.QueryRow(`
		SELECT id, name, category, unit, min_price, max_price, status
		FROM service_items
		WHERE id = ?
	`, id).Scan(&item.ID, &item.Name, &item.Category, &item.Unit, &item.MinPrice, &item.MaxPrice, &item.Status)
	if err != nil {
		return nil, ErrServiceItemNotFound
	}
	return item, nil
}

func (s *MySQLStore) ApplyProvider(userID uint64, realName, idCardNo string) (*Provider, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var providerID uint64
	err = tx.QueryRowContext(ctx, `SELECT id FROM providers WHERE user_id = ?`, userID).Scan(&providerID)
	switch {
	case err == nil:
		if _, err := tx.ExecContext(ctx, `
			UPDATE providers
			SET real_name = ?, id_card_no = ?, audit_status = 1
			WHERE id = ?
		`, realName, idCardNo, providerID); err != nil {
			return nil, err
		}
	case err == sql.ErrNoRows:
		res, err := tx.ExecContext(ctx, `
			INSERT INTO providers(user_id, provider_no, real_name, id_card_no, audit_status, work_status, status)
			VALUES (?, ?, ?, ?, 1, 1, 1)
		`, userID, fmt.Sprintf("P%08d", userID), realName, idCardNo)
		if err != nil {
			return nil, err
		}
		id, _ := res.LastInsertId()
		providerID = uint64(id)
	default:
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO audit_records(biz_type, biz_id, audit_status, auditor_id, audit_remark)
		VALUES ('provider', ?, 10, 0, '')
	`, providerID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.ProviderByUserID(userID)
}

func (s *MySQLStore) ProviderByUserID(userID uint64) (*Provider, error) {
	var providerID uint64
	if err := s.db.QueryRow(`SELECT id FROM providers WHERE user_id = ? AND deleted_at IS NULL`, userID).Scan(&providerID); err != nil {
		return nil, ErrProviderNotFound
	}
	return s.ProviderByID(providerID)
}

func (s *MySQLStore) ProviderByID(providerID uint64) (*Provider, error) {
	provider, err := s.loadProvider(providerID)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (s *MySQLStore) Providers() []*Provider {
	rows, err := s.db.Query(`SELECT id FROM providers WHERE deleted_at IS NULL ORDER BY id ASC`)
	if err != nil {
		return []*Provider{}
	}
	defer rows.Close()

	var list []*Provider
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			continue
		}
		provider, err := s.loadProvider(id)
		if err == nil {
			list = append(list, provider)
		}
	}
	return list
}

func (s *MySQLStore) UpdateProviderProfile(userID uint64, displayName, intro string, tags []string) (*Provider, error) {
	provider, err := s.ProviderByUserID(userID)
	if err != nil {
		return nil, err
	}
	tagsJSON := toJSON(tags)
	if _, err := s.db.Exec(`
		INSERT INTO provider_profiles(provider_id, display_name, intro, tags)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE display_name = VALUES(display_name), intro = VALUES(intro), tags = VALUES(tags), updated_at = CURRENT_TIMESTAMP
	`, provider.ID, displayName, intro, tagsJSON); err != nil {
		return nil, err
	}
	return s.ProviderByID(provider.ID)
}

func (s *MySQLStore) UpdateProviderServiceItems(userID uint64, items []ProviderServiceItem) (*Provider, error) {
	provider, err := s.ProviderByUserID(userID)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM provider_service_items WHERE provider_id = ?`, provider.ID); err != nil {
		return nil, err
	}
	for _, item := range items {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO provider_service_items(provider_id, service_item_id, price_amount, price_unit, status)
			VALUES (?, ?, ?, ?, 1)
		`, provider.ID, item.ServiceItemID, item.PriceAmount, item.PriceUnit); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.ProviderByID(provider.ID)
}

func (s *MySQLStore) UpdateProviderWorkStatus(userID uint64, workStatus int) (*Provider, error) {
	provider, err := s.ProviderByUserID(userID)
	if err != nil {
		return nil, err
	}
	if _, err := s.db.Exec(`UPDATE providers SET work_status = ? WHERE id = ?`, workStatus, provider.ID); err != nil {
		return nil, err
	}
	return s.ProviderByID(provider.ID)
}

func (s *MySQLStore) ApproveProvider(providerID uint64, remark string) (*Provider, error) {
	return s.updateProviderAuditStatus(providerID, 2, 20, remark)
}

func (s *MySQLStore) RejectProvider(providerID uint64, remark string) (*Provider, error) {
	return s.updateProviderAuditStatus(providerID, 3, 30, remark)
}

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

func (s *MySQLStore) getUser(userID uint64) (*User, error) {
	user := &User{}
	var birthday sql.NullTime
	err := s.db.QueryRow(`
		SELECT id, mobile, nickname, avatar, gender, birthday, city_code, created_at, updated_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`, userID).Scan(&user.ID, &user.Mobile, &user.Nickname, &user.Avatar, &user.Gender, &birthday, &user.CityCode, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, ErrUserNotFound
	}
	user.Birthday = formatNullDate(birthday)
	return user, nil
}

func (s *MySQLStore) findOrCreateUserByMobileTx(ctx context.Context, tx *sql.Tx, mobile string) (*User, error) {
	user := &User{}
	var birthday sql.NullTime
	err := tx.QueryRowContext(ctx, `
		SELECT id, mobile, nickname, avatar, gender, birthday, city_code, created_at, updated_at
		FROM users
		WHERE mobile = ? AND deleted_at IS NULL
	`, mobile).Scan(&user.ID, &user.Mobile, &user.Nickname, &user.Avatar, &user.Gender, &birthday, &user.CityCode, &user.CreatedAt, &user.UpdatedAt)
	switch {
	case err == nil:
		user.Birthday = formatNullDate(birthday)
		return user, nil
	case err != sql.ErrNoRows:
		return nil, err
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO users(mobile, nickname, city_code, status)
		VALUES (?, ?, '310100', 1)
	`, mobile, fmt.Sprintf("listen用户%s", mobile[len(mobile)-4:]))
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	user.ID = uint64(id)
	user.Mobile = mobile
	user.Nickname = fmt.Sprintf("listen用户%s", mobile[len(mobile)-4:])
	user.CityCode = "310100"
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	return user, nil
}

func (s *MySQLStore) insertTokenTx(ctx context.Context, tx *sql.Tx, userID uint64, token, tokenType string, ttl time.Duration) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO auth_tokens(user_id, token, token_type, expires_at, status)
		VALUES (?, ?, ?, ?, 1)
	`, userID, token, tokenType, time.Now().Add(ttl))
	return err
}

func (s *MySQLStore) updateProviderAuditStatus(providerID uint64, providerStatus, auditRecordStatus int, remark string) (*Provider, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `UPDATE providers SET audit_status = ? WHERE id = ?`, providerStatus, providerID)
	if err != nil {
		return nil, err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return nil, ErrProviderNotFound
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO audit_records(biz_type, biz_id, audit_status, auditor_id, audit_remark)
		VALUES ('provider', ?, ?, 1, ?)
	`, providerID, auditRecordStatus, remark); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.ProviderByID(providerID)
}

func (s *MySQLStore) loadProvider(providerID uint64) (*Provider, error) {
	provider := &Provider{}
	err := s.db.QueryRow(`
		SELECT id, user_id, provider_no, real_name, id_card_no, audit_status, work_status, created_at, updated_at
		FROM providers
		WHERE id = ? AND deleted_at IS NULL
	`, providerID).Scan(
		&provider.ID, &provider.UserID, &provider.ProviderNo, &provider.RealName, &provider.IDCardNo,
		&provider.AuditStatus, &provider.WorkStatus, &provider.CreatedAt, &provider.UpdatedAt,
	)
	if err != nil {
		return nil, ErrProviderNotFound
	}

	var tagsRaw sql.NullString
	_ = s.db.QueryRow(`SELECT display_name, intro, tags FROM provider_profiles WHERE provider_id = ?`, providerID).Scan(&provider.DisplayName, &provider.Intro, &tagsRaw)
	provider.Tags = parseJSONStringArray(tagsRaw.String)
	provider.AuditRemark = s.latestAuditRemark("provider", providerID)
	provider.ServiceItems = s.loadProviderServiceItems(providerID)
	return provider, nil
}

func (s *MySQLStore) latestAuditRemark(bizType string, bizID uint64) string {
	var remark string
	_ = s.db.QueryRow(`
		SELECT audit_remark
		FROM audit_records
		WHERE biz_type = ? AND biz_id = ?
		ORDER BY id DESC
		LIMIT 1
	`, bizType, bizID).Scan(&remark)
	return remark
}

func (s *MySQLStore) loadProviderServiceItems(providerID uint64) []ProviderServiceItem {
	rows, err := s.db.Query(`
		SELECT service_item_id, price_amount, price_unit
		FROM provider_service_items
		WHERE provider_id = ? AND status = 1
		ORDER BY id ASC
	`, providerID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var items []ProviderServiceItem
	for rows.Next() {
		var item ProviderServiceItem
		if err := rows.Scan(&item.ServiceItemID, &item.PriceAmount, &item.PriceUnit); err == nil {
			items = append(items, item)
		}
	}
	return items
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

func splitSQLStatements(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	stmts := make([]string, 0, len(parts))
	for _, part := range parts {
		stmt := strings.TrimSpace(part)
		if stmt != "" {
			stmts = append(stmts, stmt)
		}
	}
	return stmts
}

func parseJSONStringArray(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil
	}
	return items
}

func (s *MySQLStore) seedAdminBootstrap(ctx context.Context) error {
	roleSeed := []struct {
		Name    string
		RoleKey string
	}{
		{Name: "超级管理员", RoleKey: "super_admin"},
		{Name: "运营管理员", RoleKey: "ops_admin"},
		{Name: "财务管理员", RoleKey: "finance_admin"},
		{Name: "内容管理员", RoleKey: "content_admin"},
	}
	roleIDMap := map[string]uint64{}
	for _, item := range roleSeed {
		var roleID uint64
		err := s.db.QueryRowContext(ctx, `SELECT id FROM admin_roles WHERE role_key = ? LIMIT 1`, item.RoleKey).Scan(&roleID)
		if err == sql.ErrNoRows {
			res, execErr := s.db.ExecContext(ctx, `
				INSERT INTO admin_roles(name, role_key, status)
				VALUES (?, ?, 1)
			`, item.Name, item.RoleKey)
			if execErr != nil {
				return fmt.Errorf("seed admin role %s: %w", item.RoleKey, execErr)
			}
			id, _ := res.LastInsertId()
			roleID = uint64(id)
		} else if err != nil {
			return fmt.Errorf("query admin role %s: %w", item.RoleKey, err)
		}
		roleIDMap[item.RoleKey] = roleID
	}

	permissionSeed := []struct {
		Name          string
		PermissionKey string
		Category      string
	}{
		{Name: "看板总览", PermissionKey: "admin:dashboard:overview", Category: "dashboard"},
		{Name: "用户列表", PermissionKey: "admin:user:list", Category: "user"},
		{Name: "用户详情", PermissionKey: "admin:user:detail", Category: "user"},
		{Name: "用户状态更新", PermissionKey: "admin:user:status:update", Category: "user"},
		{Name: "服务方列表", PermissionKey: "admin:provider:list", Category: "provider"},
		{Name: "服务方详情", PermissionKey: "admin:provider:detail", Category: "provider"},
		{Name: "服务方审核通过", PermissionKey: "admin:provider:approve", Category: "provider"},
		{Name: "服务方审核拒绝", PermissionKey: "admin:provider:reject", Category: "provider"},
		{Name: "服务方状态更新", PermissionKey: "admin:provider:status:update", Category: "provider"},
		{Name: "服务项目列表", PermissionKey: "admin:service_item:list", Category: "service_item"},
		{Name: "服务项目创建", PermissionKey: "admin:service_item:create", Category: "service_item"},
		{Name: "服务项目更新", PermissionKey: "admin:service_item:update", Category: "service_item"},
		{Name: "服务项目删除", PermissionKey: "admin:service_item:delete", Category: "service_item"},
		{Name: "订单列表", PermissionKey: "admin:order:list", Category: "order"},
		{Name: "订单详情", PermissionKey: "admin:order:detail", Category: "order"},
		{Name: "订单人工完结", PermissionKey: "admin:order:manual_complete", Category: "order"},
		{Name: "订单退款", PermissionKey: "admin:order:refund", Category: "order"},
		{Name: "提现列表", PermissionKey: "admin:withdraw:list", Category: "withdraw"},
		{Name: "提现审核通过", PermissionKey: "admin:withdraw:approve", Category: "withdraw"},
		{Name: "提现审核拒绝", PermissionKey: "admin:withdraw:reject", Category: "withdraw"},
		{Name: "财务报表", PermissionKey: "admin:finance:reports", Category: "finance"},
		{Name: "帖子列表", PermissionKey: "admin:post:list", Category: "content"},
		{Name: "帖子隐藏", PermissionKey: "admin:post:hide", Category: "content"},
		{Name: "声音列表", PermissionKey: "admin:audio:list", Category: "content"},
		{Name: "声音下架", PermissionKey: "admin:audio:off_shelf", Category: "content"},
		{Name: "投诉列表", PermissionKey: "admin:complaint:list", Category: "complaint"},
		{Name: "投诉详情", PermissionKey: "admin:complaint:detail", Category: "complaint"},
		{Name: "投诉处理", PermissionKey: "admin:complaint:resolve", Category: "complaint"},
		{Name: "风控事件列表", PermissionKey: "admin:risk_event:list", Category: "risk"},
		{Name: "配置列表", PermissionKey: "admin:config:list", Category: "config"},
		{Name: "配置更新", PermissionKey: "admin:config:update", Category: "config"},
		{Name: "RBAC角色列表", PermissionKey: "admin:rbac:role:list", Category: "rbac"},
		{Name: "RBAC权限列表", PermissionKey: "admin:rbac:permission:list", Category: "rbac"},
		{Name: "管理员角色分配", PermissionKey: "admin:rbac:user_role:assign", Category: "rbac"},
	}
	permissionIDMap := map[string]uint64{}
	for _, item := range permissionSeed {
		var permissionID uint64
		err := s.db.QueryRowContext(ctx, `SELECT id FROM admin_permissions WHERE permission_key = ? LIMIT 1`, item.PermissionKey).Scan(&permissionID)
		if err == sql.ErrNoRows {
			res, execErr := s.db.ExecContext(ctx, `
				INSERT INTO admin_permissions(name, permission_key, category, status)
				VALUES (?, ?, ?, 1)
			`, item.Name, item.PermissionKey, item.Category)
			if execErr != nil {
				return fmt.Errorf("seed admin permission %s: %w", item.PermissionKey, execErr)
			}
			id, _ := res.LastInsertId()
			permissionID = uint64(id)
		} else if err != nil {
			return fmt.Errorf("query admin permission %s: %w", item.PermissionKey, err)
		}
		permissionIDMap[item.PermissionKey] = permissionID
	}

	rolePermissionSeed := map[string][]string{
		"super_admin": {
			"admin:dashboard:overview",
			"admin:user:list", "admin:user:detail", "admin:user:status:update",
			"admin:provider:list", "admin:provider:detail", "admin:provider:approve", "admin:provider:reject", "admin:provider:status:update",
			"admin:service_item:list", "admin:service_item:create", "admin:service_item:update", "admin:service_item:delete",
			"admin:order:list", "admin:order:detail", "admin:order:manual_complete", "admin:order:refund",
			"admin:withdraw:list", "admin:withdraw:approve", "admin:withdraw:reject",
			"admin:finance:reports",
			"admin:post:list", "admin:post:hide", "admin:audio:list", "admin:audio:off_shelf",
			"admin:complaint:list", "admin:complaint:detail", "admin:complaint:resolve",
			"admin:risk_event:list", "admin:config:list", "admin:config:update",
			"admin:rbac:role:list", "admin:rbac:permission:list", "admin:rbac:user_role:assign",
		},
		"ops_admin": {
			"admin:dashboard:overview",
			"admin:user:list", "admin:user:detail", "admin:user:status:update",
			"admin:provider:list", "admin:provider:detail", "admin:provider:approve", "admin:provider:reject", "admin:provider:status:update",
			"admin:service_item:list", "admin:service_item:create", "admin:service_item:update", "admin:service_item:delete",
			"admin:order:list", "admin:order:detail", "admin:order:manual_complete",
			"admin:complaint:list", "admin:complaint:detail", "admin:complaint:resolve", "admin:risk_event:list",
		},
		"finance_admin": {
			"admin:dashboard:overview",
			"admin:order:list", "admin:order:detail", "admin:order:refund",
			"admin:withdraw:list", "admin:withdraw:approve", "admin:withdraw:reject", "admin:finance:reports",
		},
		"content_admin": {
			"admin:dashboard:overview",
			"admin:post:list", "admin:post:hide", "admin:audio:list", "admin:audio:off_shelf",
			"admin:complaint:list", "admin:complaint:detail", "admin:complaint:resolve",
		},
	}
	for roleKey, permissionKeys := range rolePermissionSeed {
		roleID := roleIDMap[roleKey]
		for _, permissionKey := range permissionKeys {
			permissionID := permissionIDMap[permissionKey]
			if _, err := s.db.ExecContext(ctx, `
				INSERT IGNORE INTO admin_role_permissions(admin_role_id, admin_permission_id)
				VALUES (?, ?)
			`, roleID, permissionID); err != nil {
				return fmt.Errorf("seed role permissions %s: %w", roleKey, err)
			}
		}
	}

	adminSeed := []struct {
		Username string
		Password string
		Nickname string
		RoleKey  string
	}{
		{Username: "admin", Password: "admin123456", Nickname: "系统管理员", RoleKey: "super_admin"},
		{Username: "content_admin", Password: "admin123456", Nickname: "内容管理员", RoleKey: "content_admin"},
	}
	for _, item := range adminSeed {
		var adminID uint64
		err := s.db.QueryRowContext(ctx, `SELECT id FROM admin_users WHERE username = ? LIMIT 1`, item.Username).Scan(&adminID)
		if err == sql.ErrNoRows {
			res, execErr := s.db.ExecContext(ctx, `
				INSERT INTO admin_users(username, password_hash, nickname, status)
				VALUES (?, ?, ?, 1)
			`, item.Username, item.Password, item.Nickname)
			if execErr != nil {
				return fmt.Errorf("seed admin user %s: %w", item.Username, execErr)
			}
			id, _ := res.LastInsertId()
			adminID = uint64(id)
		} else if err != nil {
			return fmt.Errorf("query admin user %s: %w", item.Username, err)
		}

		if _, err := s.db.ExecContext(ctx, `
			INSERT IGNORE INTO admin_user_roles(admin_user_id, admin_role_id)
			VALUES (?, ?)
		`, adminID, roleIDMap[item.RoleKey]); err != nil {
			return fmt.Errorf("seed admin user role %s: %w", item.Username, err)
		}
	}
	return nil
}

func toJSON(v any) string {
	if v == nil {
		return "[]"
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func formatNullDate(v sql.NullTime) string {
	if !v.Valid {
		return ""
	}
	return v.Time.Format("2006-01-02")
}

func formatNullDateTime(v sql.NullTime) string {
	if !v.Valid {
		return ""
	}
	return v.Time.Format("2006-01-02 15:04:05")
}

func ptrTime(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}

func operatorRoleName(roleCode int) string {
	switch roleCode {
	case operatorRoleProvider:
		return "provider"
	case operatorRoleAdmin:
		return "admin"
	default:
		return "user"
	}
}

func containsInt(list []int, target int) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
