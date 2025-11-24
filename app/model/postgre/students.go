package model

import "time"

type Student struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	StudentID   string    `json:"student_id"`
	ProgramStudy string   `json:"program_study"`
	AcademicYear string   `json:"academic_year"`
	AdvisorID   string    `json:"advisor_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateStudentRequest struct {
	UserID       string `json:"user_id" validate:"required"`
	StudentID    string `json:"student_id" validate:"required"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`
	AdvisorID    string `json:"advisor_id"`
}

type UpdateStudentRequest struct {
	StudentID    string `json:"student_id" validate:"required"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`
	AdvisorID    string `json:"advisor_id"`
}

type GetAllStudentsResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    []Student `json:"data"`
}

type GetStudentByIDResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    Student `json:"data"`
}

type CreateStudentResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    Student `json:"data"`
}

type UpdateStudentResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    Student `json:"data"`
}

type DeleteStudentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

