package model

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

func cloneAdmin(admin *AdminUser) *AdminUser {
	if admin == nil {
		return nil
	}
	copyAdmin := *admin
	copyAdmin.Roles = append([]string(nil), admin.Roles...)
	return &copyAdmin
}
