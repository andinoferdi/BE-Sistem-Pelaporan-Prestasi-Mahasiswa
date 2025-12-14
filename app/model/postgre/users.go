package model

// #1 proses: import library time untuk handle timestamp
import "time"

// #2 proses: struct utama untuk menyimpan data user di database PostgreSQL
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	RoleID       string    `json:"role_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// #3 proses: struct untuk request create user baru, termasuk data opsional untuk student atau lecturer
type CreateUserRequest struct {
	Username     string `json:"username" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
	FullName     string `json:"full_name" validate:"required"`
	RoleID       string `json:"role_id" validate:"required"`
	IsActive     *bool  `json:"is_active,omitempty"`
	StudentID    string `json:"student_id,omitempty"`
	ProgramStudy string `json:"program_study,omitempty"`
	AcademicYear string `json:"academic_year,omitempty"`
	AdvisorID    string `json:"advisor_id,omitempty"`
	LecturerID   string `json:"lecturer_id,omitempty"`
	Department   string `json:"department,omitempty"`
}

// #4 proses: struct untuk request update data user yang sudah ada
type UpdateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
	RoleID   string `json:"role_id" validate:"required"`
	IsActive *bool  `json:"is_active"`
}

// #5 proses: struct response untuk get all users, return list semua user
type GetAllUsersResponse struct {
	Status string `json:"status"`
	Data   []User `json:"data"`
}

// #6 proses: struct response untuk get user by ID, return satu user
type GetUserByIDResponse struct {
	Status string `json:"status"`
	Data   User   `json:"data"`
}

// #7 proses: struct response untuk create user, return user yang baru dibuat
type CreateUserResponse struct {
	Status string `json:"status"`
	Data   User   `json:"data"`
}

// #8 proses: struct response untuk update user, return user yang sudah diupdate
type UpdateUserResponse struct {
	Status string `json:"status"`
	Data   User   `json:"data"`
}

// #9 proses: struct response untuk delete user, hanya return status
type DeleteUserResponse struct {
	Status string `json:"status"`
}
