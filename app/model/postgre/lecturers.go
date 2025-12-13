package model

import "time"

type Lecturer struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	LecturerID string    `json:"lecturer_id"`
	Department string    `json:"department"`
	FullName   string    `json:"full_name,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateLecturerRequest struct {
	UserID     string `json:"user_id" validate:"required"`
	LecturerID string `json:"lecturer_id" validate:"required"`
	Department string `json:"department"`
}

type UpdateLecturerRequest struct {
	LecturerID string `json:"lecturer_id" validate:"required"`
	Department string `json:"department"`
}

type GetAllLecturersResponse struct {
	Status string     `json:"status"`
	Data   []Lecturer `json:"data"`
}

type GetLecturerByIDResponse struct {
	Status string   `json:"status"`
	Data   Lecturer `json:"data"`
}

type CreateLecturerResponse struct {
	Status string   `json:"status"`
	Data   Lecturer `json:"data"`
}

type UpdateLecturerResponse struct {
	Status string   `json:"status"`
	Data   Lecturer `json:"data"`
}

type DeleteLecturerResponse struct {
	Status string `json:"status"`
}
