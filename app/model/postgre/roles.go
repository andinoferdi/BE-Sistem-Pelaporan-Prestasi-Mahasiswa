package model

import "time"

type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type GetAllRolesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []Role `json:"data"`
}

type GetRoleByIDResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Role   `json:"data"`
}

type CreateRoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Role   `json:"data"`
}

type UpdateRoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Role   `json:"data"`
}

type DeleteRoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

