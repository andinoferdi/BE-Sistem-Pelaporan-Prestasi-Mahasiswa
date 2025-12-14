package model

// #1 proses: struct untuk relasi many-to-many antara role dan permission
type RolePermission struct {
	RoleID       string `json:"role_id"`
	PermissionID string `json:"permission_id"`
}

// #2 proses: struct untuk request assign permission ke role
type CreateRolePermissionRequest struct {
	RoleID       string `json:"role_id" validate:"required"`
	PermissionID string `json:"permission_id" validate:"required"`
}

// #3 proses: struct response untuk get all role permissions, return list semua relasi
type GetAllRolePermissionsResponse struct {
	Status string           `json:"status"`
	Data   []RolePermission `json:"data"`
}

// #4 proses: struct response untuk create role permission, return relasi yang baru dibuat
type CreateRolePermissionResponse struct {
	Status string         `json:"status"`
	Data   RolePermission `json:"data"`
}

// #5 proses: struct response untuk delete role permission, hanya return status
type DeleteRolePermissionResponse struct {
	Status string `json:"status"`
}
