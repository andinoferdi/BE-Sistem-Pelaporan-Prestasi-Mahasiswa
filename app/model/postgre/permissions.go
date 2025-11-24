package model

type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type CreatePermissionRequest struct {
	Name        string `json:"name" validate:"required"`
	Resource    string `json:"resource" validate:"required"`
	Action      string `json:"action" validate:"required"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name" validate:"required"`
	Resource    string `json:"resource" validate:"required"`
	Action      string `json:"action" validate:"required"`
	Description string `json:"description"`
}

type GetAllPermissionsResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []Permission `json:"data"`
}

type GetPermissionByIDResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    Permission `json:"data"`
}

type CreatePermissionResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    Permission `json:"data"`
}

type UpdatePermissionResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    Permission `json:"data"`
}

type DeletePermissionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
