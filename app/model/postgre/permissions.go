package model

// #1 proses: struct utama untuk menyimpan data permission di database
type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

// #2 proses: struct untuk request create permission baru
type CreatePermissionRequest struct {
	Name        string `json:"name" validate:"required"`
	Resource    string `json:"resource" validate:"required"`
	Action      string `json:"action" validate:"required"`
	Description string `json:"description"`
}

// #3 proses: struct untuk request update data permission
type UpdatePermissionRequest struct {
	Name        string `json:"name" validate:"required"`
	Resource    string `json:"resource" validate:"required"`
	Action      string `json:"action" validate:"required"`
	Description string `json:"description"`
}

// #4 proses: struct response untuk get all permissions, return list semua permission
type GetAllPermissionsResponse struct {
	Status string       `json:"status"`
	Data   []Permission `json:"data"`
}

// #5 proses: struct response untuk get permission by ID, return satu permission
type GetPermissionByIDResponse struct {
	Status string     `json:"status"`
	Data   Permission `json:"data"`
}

// #6 proses: struct response untuk create permission, return permission yang baru dibuat
type CreatePermissionResponse struct {
	Status string     `json:"status"`
	Data   Permission `json:"data"`
}

// #7 proses: struct response untuk update permission, return permission yang sudah diupdate
type UpdatePermissionResponse struct {
	Status string     `json:"status"`
	Data   Permission `json:"data"`
}

// #8 proses: struct response untuk delete permission, hanya return status
type DeletePermissionResponse struct {
	Status string `json:"status"`
}
