package model

// #1 proses: import library time untuk handle timestamp
import "time"

// #2 proses: struct utama untuk menyimpan data dosen di database
type Lecturer struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	LecturerID string    `json:"lecturer_id"`
	Department string    `json:"department"`
	FullName   string    `json:"full_name,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// #3 proses: struct untuk request create profil dosen baru
type CreateLecturerRequest struct {
	UserID     string `json:"user_id" validate:"required"`
	LecturerID string `json:"lecturer_id" validate:"required"`
	Department string `json:"department"`
}

// #4 proses: struct untuk request update data dosen
type UpdateLecturerRequest struct {
	LecturerID string `json:"lecturer_id" validate:"required"`
	Department string `json:"department"`
}

// #5 proses: struct response untuk get all lecturers, return list semua dosen
type GetAllLecturersResponse struct {
	Status string     `json:"status"`
	Data   []Lecturer `json:"data"`
}

// #6 proses: struct response untuk get lecturer by ID, return satu dosen
type GetLecturerByIDResponse struct {
	Status string   `json:"status"`
	Data   Lecturer `json:"data"`
}

// #7 proses: struct response untuk create lecturer, return dosen yang baru dibuat
type CreateLecturerResponse struct {
	Status string   `json:"status"`
	Data   Lecturer `json:"data"`
}

// #8 proses: struct response untuk update lecturer, return dosen yang sudah diupdate
type UpdateLecturerResponse struct {
	Status string   `json:"status"`
	Data   Lecturer `json:"data"`
}

// #9 proses: struct response untuk delete lecturer, hanya return status
type DeleteLecturerResponse struct {
	Status string `json:"status"`
}
