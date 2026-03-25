package model

import "time"

type AdminPermission struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	PermissionKey string    `json:"permissionKey"`
	Category      string    `json:"category"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
