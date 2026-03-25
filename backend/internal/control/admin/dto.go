package admin

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateAdminUserRolesRequest struct {
	RoleKeys []string `json:"roleKeys"`
}
