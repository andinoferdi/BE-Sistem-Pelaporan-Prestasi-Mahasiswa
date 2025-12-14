package model

// #1 proses: import library time untuk handle timestamp
import "time"

// #2 proses: struct utama untuk menyimpan data mahasiswa di database
type Student struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	StudentID    string    `json:"student_id"`
	ProgramStudy string    `json:"program_study"`
	AcademicYear string    `json:"academic_year"`
	AdvisorID    string    `json:"advisor_id"`
	FullName     string    `json:"full_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// #3 proses: struct untuk request create profil mahasiswa baru
type CreateStudentRequest struct {
	UserID       string `json:"user_id" validate:"required"`
	StudentID    string `json:"student_id" validate:"required"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`
	AdvisorID    string `json:"advisor_id"`
}

// #4 proses: struct untuk request update data mahasiswa
type UpdateStudentRequest struct {
	StudentID    string `json:"student_id" validate:"required"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`
	AdvisorID    string `json:"advisor_id"`
}

// #5 proses: struct response untuk get all students, return list semua mahasiswa
type GetAllStudentsResponse struct {
	Status string    `json:"status"`
	Data   []Student `json:"data"`
}

// #6 proses: struct response untuk get student by ID, return satu mahasiswa
type GetStudentByIDResponse struct {
	Status string  `json:"status"`
	Data   Student `json:"data"`
}

// #7 proses: struct response untuk create student, return mahasiswa yang baru dibuat
type CreateStudentResponse struct {
	Status string  `json:"status"`
	Data   Student `json:"data"`
}

// #8 proses: struct response untuk update student, return mahasiswa yang sudah diupdate
type UpdateStudentResponse struct {
	Status string  `json:"status"`
	Data   Student `json:"data"`
}

// #9 proses: struct response untuk delete student, hanya return status
type DeleteStudentResponse struct {
	Status string `json:"status"`
}
