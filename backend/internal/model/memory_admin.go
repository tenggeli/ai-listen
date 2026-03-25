package model

import (
	"fmt"
	"strings"
	"time"
)

func (s *MemoryStore) seedAdminRBAC() {
	now := time.Now()
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
	for idx, item := range permissionSeed {
		s.adminPerms[item.PermissionKey] = &AdminPermission{
			ID:            uint64(idx + 1),
			Name:          item.Name,
			PermissionKey: item.PermissionKey,
			Category:      item.Category,
			Status:        1,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
	}

	roleSeed := []struct {
		Name    string
		RoleKey string
	}{
		{Name: "超级管理员", RoleKey: "super_admin"},
		{Name: "运营管理员", RoleKey: "ops_admin"},
		{Name: "财务管理员", RoleKey: "finance_admin"},
		{Name: "内容管理员", RoleKey: "content_admin"},
	}
	for idx, item := range roleSeed {
		s.adminRoles[item.RoleKey] = &AdminRole{
			ID:        uint64(idx + 1),
			Name:      item.Name,
			RoleKey:   item.RoleKey,
			Status:    1,
			CreatedAt: now,
			UpdatedAt: now,
		}
		s.rolePerms[item.RoleKey] = map[string]struct{}{}
	}

	s.bindRolePermissions("super_admin", []string{"*"})
	s.bindRolePermissions("ops_admin", []string{
		"admin:dashboard:overview",
		"admin:user:list", "admin:user:detail", "admin:user:status:update",
		"admin:provider:list", "admin:provider:detail", "admin:provider:approve", "admin:provider:reject", "admin:provider:status:update",
		"admin:service_item:list", "admin:service_item:create", "admin:service_item:update", "admin:service_item:delete",
		"admin:order:list", "admin:order:detail", "admin:order:manual_complete",
		"admin:complaint:list", "admin:complaint:detail", "admin:complaint:resolve", "admin:risk_event:list",
	})
	s.bindRolePermissions("finance_admin", []string{
		"admin:dashboard:overview",
		"admin:order:list", "admin:order:detail", "admin:order:refund",
		"admin:withdraw:list", "admin:withdraw:approve", "admin:withdraw:reject", "admin:finance:reports",
	})
	s.bindRolePermissions("content_admin", []string{
		"admin:dashboard:overview",
		"admin:post:list", "admin:post:hide",
		"admin:audio:list", "admin:audio:off_shelf",
		"admin:complaint:list", "admin:complaint:detail", "admin:complaint:resolve",
	})
}

func (s *MemoryStore) bindRolePermissions(roleKey string, permissionKeys []string) {
	set, ok := s.rolePerms[roleKey]
	if !ok {
		set = map[string]struct{}{}
		s.rolePerms[roleKey] = set
	}
	for _, permissionKey := range permissionKeys {
		set[permissionKey] = struct{}{}
	}
}

func (s *MemoryStore) AdminLogin(username, password string) (*AdminUser, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	admin, ok := s.adminByName[username]
	if !ok {
		return nil, "", ErrUnauthorized
	}
	if s.adminPasswords[admin.ID] != password || admin.Status != 1 {
		return nil, "", ErrUnauthorized
	}

	token := fmt.Sprintf("admin-token-%d-%d", admin.ID, time.Now().UnixNano())
	s.adminTokens[token] = admin.ID
	return cloneAdmin(admin), token, nil
}

func (s *MemoryStore) AdminByToken(raw string) (*AdminUser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	token := strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
	adminID, ok := s.adminTokens[token]
	if !ok {
		return nil, ErrUnauthorized
	}
	admin, ok := s.adminUsers[adminID]
	if !ok {
		return nil, ErrUnauthorized
	}
	return cloneAdmin(admin), nil
}

func (s *MemoryStore) AdminUserByID(adminUserID uint64) (*AdminUser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	admin, ok := s.adminUsers[adminUserID]
	if !ok {
		return nil, ErrAdminNotFound
	}
	return cloneAdmin(admin), nil
}

func (s *MemoryStore) AdminPermissionsByRoles(roleKeys []string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	permissionSet := map[string]struct{}{}
	for _, roleKey := range roleKeys {
		for permissionKey := range s.rolePerms[roleKey] {
			if permissionKey == "*" {
				for key := range s.adminPerms {
					permissionSet[key] = struct{}{}
				}
				continue
			}
			permissionSet[permissionKey] = struct{}{}
		}
	}

	result := make([]string, 0, len(permissionSet))
	for key := range permissionSet {
		result = append(result, key)
	}
	return result
}

func (s *MemoryStore) AdminHasPermission(adminUserID uint64, permission string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	admin, ok := s.adminUsers[adminUserID]
	if !ok {
		return false
	}
	for _, roleKey := range admin.Roles {
		if _, ok := s.rolePerms[roleKey]["*"]; ok {
			return true
		}
		if _, ok := s.rolePerms[roleKey][permission]; ok {
			return true
		}
	}
	return false
}

func (s *MemoryStore) AdminRoles() []*AdminRole {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*AdminRole, 0, len(s.adminRoles))
	for roleKey, role := range s.adminRoles {
		copyRole := cloneAdminRole(role)
		for permissionKey := range s.rolePerms[roleKey] {
			copyRole.PermissionKeys = append(copyRole.PermissionKeys, permissionKey)
		}
		result = append(result, copyRole)
	}
	return result
}

func (s *MemoryStore) AdminPermissions() []*AdminPermission {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*AdminPermission, 0, len(s.adminPerms))
	for _, permission := range s.adminPerms {
		copyPermission := *permission
		result = append(result, &copyPermission)
	}
	return result
}

func (s *MemoryStore) UpdateAdminUserRoles(adminUserID uint64, roleKeys []string) (*AdminUser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	admin, ok := s.adminUsers[adminUserID]
	if !ok {
		return nil, ErrAdminNotFound
	}

	normalized := make([]string, 0, len(roleKeys))
	seen := map[string]struct{}{}
	for _, roleKey := range roleKeys {
		if _, exists := s.adminRoles[roleKey]; !exists {
			return nil, ErrRoleNotFound
		}
		if _, duplicated := seen[roleKey]; duplicated {
			continue
		}
		seen[roleKey] = struct{}{}
		normalized = append(normalized, roleKey)
	}

	admin.Roles = normalized
	admin.UpdatedAt = time.Now()
	return cloneAdmin(admin), nil
}
