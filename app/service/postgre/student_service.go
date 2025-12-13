package service

import (
	"context"
	"database/sql"
	"errors"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	"strings"
)

type IStudentService interface {
	GetAllStudents(ctx context.Context) ([]modelpostgre.Student, error)
	GetStudentByID(ctx context.Context, id string) (*modelpostgre.Student, error)
	GetStudentIDByUserID(ctx context.Context, userID string) (string, error)
	GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error)
	GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error)
	CreateStudent(ctx context.Context, req modelpostgre.CreateStudentRequest) (*modelpostgre.Student, error)
	UpdateStudentAdvisor(ctx context.Context, studentID string, advisorID string) error
}

type StudentService struct {
	studentRepo  repositorypostgre.IStudentRepository
	userRepo     repositorypostgre.IUserRepository
	lecturerRepo repositorypostgre.ILecturerRepository
}

func NewStudentService(studentRepo repositorypostgre.IStudentRepository, userRepo repositorypostgre.IUserRepository, lecturerRepo repositorypostgre.ILecturerRepository) IStudentService {
	return &StudentService{
		studentRepo:  studentRepo,
		userRepo:     userRepo,
		lecturerRepo: lecturerRepo,
	}
}

func (s *StudentService) GetAllStudents(ctx context.Context) ([]modelpostgre.Student, error) {
	return s.studentRepo.GetAllStudents(ctx)
}

func (s *StudentService) GetStudentByID(ctx context.Context, id string) (*modelpostgre.Student, error) {
	if id == "" {
		return nil, errors.New("student ID wajib diisi")
	}
	return s.studentRepo.GetStudentByID(ctx, id)
}

func (s *StudentService) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	return s.studentRepo.GetStudentIDByUserID(ctx, userID)
}

func (s *StudentService) GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error) {
	return s.studentRepo.GetStudentByUserID(ctx, userID)
}

func (s *StudentService) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error) {
	if advisorID == "" {
		return nil, errors.New("advisor ID wajib diisi")
	}
	return s.studentRepo.GetStudentsByAdvisorID(ctx, advisorID)
}

func (s *StudentService) CreateStudent(ctx context.Context, req modelpostgre.CreateStudentRequest) (*modelpostgre.Student, error) {
	if req.UserID == "" {
		return nil, errors.New("user ID wajib diisi")
	}
	if req.StudentID == "" {
		return nil, errors.New("student ID wajib diisi")
	}

	if s.userRepo == nil {
		return nil, errors.New("user repository tidak tersedia")
	}

	user, err := s.userRepo.FindUserByID(ctx, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, errors.New("error mengambil data user: " + err.Error())
	}

	roleName, err := s.userRepo.GetRoleName(ctx, user.RoleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Mahasiswa" {
		return nil, errors.New("user harus memiliki role Mahasiswa untuk membuat student profile")
	}

	existingStudent, err := s.studentRepo.GetStudentByUserID(ctx, req.UserID)
	if err == nil && existingStudent != nil {
		return nil, errors.New("user sudah memiliki student profile")
	}

	if req.AdvisorID != "" {
		if s.lecturerRepo == nil {
			return nil, errors.New("lecturer repository tidak tersedia")
		}
		_, err := s.lecturerRepo.GetLecturerByID(ctx, req.AdvisorID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("advisor ID tidak ditemukan")
			}
			return nil, errors.New("error memvalidasi advisor ID: " + err.Error())
		}
	}

	student, err := s.studentRepo.CreateStudent(ctx, req)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "duplicate key value violates unique constraint") {
			if strings.Contains(errStr, "students_student_id_key") {
				return nil, errors.New("student ID sudah digunakan")
			}
			if strings.Contains(errStr, "students_user_id_key") {
				return nil, errors.New("user sudah memiliki student profile")
			}
			return nil, errors.New("data duplikat terdeteksi")
		}
		return nil, errors.New("error membuat student profile: " + errStr)
	}

	return student, nil
}

func (s *StudentService) UpdateStudentAdvisor(ctx context.Context, studentID string, advisorID string) error {
	if studentID == "" {
		return errors.New("student ID wajib diisi")
	}

	student, err := s.studentRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("student tidak ditemukan")
		}
		return errors.New("error mengambil data student: " + err.Error())
	}

	if advisorID != "" {
		if s.lecturerRepo == nil {
			return errors.New("lecturer repository tidak tersedia")
		}
		_, err := s.lecturerRepo.GetLecturerByID(ctx, advisorID)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("advisor ID tidak ditemukan")
			}
			return errors.New("error memvalidasi advisor ID: " + err.Error())
		}
	}

	err = s.studentRepo.UpdateStudentAdvisor(ctx, student.ID, advisorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("student tidak ditemukan")
		}
		return errors.New("error mengupdate advisor: " + err.Error())
	}

	return nil
}
