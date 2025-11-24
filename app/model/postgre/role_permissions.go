package model

type RolePermission struct {
	RoleID       string `json:"role_id"`
	PermissionID string `json:"permission_id"`
}

type CreateRolePermissionRequest struct {
	RoleID       string `json:"role_id" validate:"required"`
	PermissionID string `json:"permission_id" validate:"required"`
}

type GetAllRolePermissionsResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    []RolePermission `json:"data"`
}

type CreateRolePermissionResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    RolePermission `json:"data"`
}

type DeleteRolePermissionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
