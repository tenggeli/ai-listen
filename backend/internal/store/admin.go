package store

import "time"

type AdminUser struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Roles     []string  `json:"roles"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AdminRole struct {
	ID             uint64    `json:"id"`
	Name           string    `json:"name"`
	RoleKey        string    `json:"roleKey"`
	Status         int       `json:"status"`
	PermissionKeys []string  `json:"permissionKeys"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type AdminPermission struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	PermissionKey string    `json:"permissionKey"`
	Category      string    `json:"category"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func cloneAdmin(admin *AdminUser) *AdminUser {
	if admin == nil {
		return nil
	}
	copyAdmin := *admin
	copyAdmin.Roles = append([]string(nil), admin.Roles...)
	return &copyAdmin
}

func cloneAdminRole(role *AdminRole) *AdminRole {
	if role == nil {
		return nil
	}
	copyRole := *role
	copyRole.PermissionKeys = append([]string(nil), role.PermissionKeys...)
	return &copyRole
}
