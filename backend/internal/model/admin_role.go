package model

import "time"

type AdminRole struct {
	ID             uint64    `json:"id"`
	Name           string    `json:"name"`
	RoleKey        string    `json:"roleKey"`
	Status         int       `json:"status"`
	PermissionKeys []string  `json:"permissionKeys"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func cloneAdminRole(role *AdminRole) *AdminRole {
	if role == nil {
		return nil
	}
	copyRole := *role
	copyRole.PermissionKeys = append([]string(nil), role.PermissionKeys...)
	return &copyRole
}
