package model

// #1 proses: import library time untuk handle timestamp
import "time"

// #2 proses: struct utama untuk menyimpan data role di database
type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// #3 proses: struct untuk request create role baru
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// #4 proses: struct untuk request update data role
type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// #5 proses: struct response untuk get all roles, return list semua role
type GetAllRolesResponse struct {
	Status string `json:"status"`
	Data   []Role `json:"data"`
}

// #6 proses: struct response untuk get role by ID, return satu role
type GetRoleByIDResponse struct {
	Status string `json:"status"`
	Data   Role   `json:"data"`
}

// #7 proses: struct response untuk create role, return role yang baru dibuat
type CreateRoleResponse struct {
	Status string `json:"status"`
	Data   Role   `json:"data"`
}

// #8 proses: struct response untuk update role, return role yang sudah diupdate
type UpdateRoleResponse struct {
	Status string `json:"status"`
	Data   Role   `json:"data"`
}

// #9 proses: struct response untuk delete role, hanya return status
type DeleteRoleResponse struct {
	Status string `json:"status"`
}
