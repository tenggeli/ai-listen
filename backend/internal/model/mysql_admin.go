package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

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
