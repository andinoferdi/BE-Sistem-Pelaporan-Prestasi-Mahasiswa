package service

// #1 proses: import library yang diperlukan untuk context, database, errors, dan strings
import (
	"context"
	"database/sql"
	"errors"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	"strings"
)

// #2 proses: definisikan interface untuk operasi student
type IStudentService interface {
	GetAllStudents(ctx context.Context) ([]modelpostgre.Student, error)
	GetStudentByID(ctx context.Context, id string) (*modelpostgre.Student, error)
	GetStudentIDByUserID(ctx context.Context, userID string) (string, error)
	GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error)
	GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error)
	CreateStudent(ctx context.Context, req modelpostgre.CreateStudentRequest) (*modelpostgre.Student, error)
	UpdateStudentAdvisor(ctx context.Context, studentID string, advisorID string) error
}

// #3 proses: struct service untuk student dengan dependency student, user, dan lecturer repository
type StudentService struct {
	studentRepo  repositorypostgre.IStudentRepository
	userRepo     repositorypostgre.IUserRepository
	lecturerRepo repositorypostgre.ILecturerRepository
}

// #4 proses: constructor untuk membuat instance StudentService baru
func NewStudentService(studentRepo repositorypostgre.IStudentRepository, userRepo repositorypostgre.IUserRepository, lecturerRepo repositorypostgre.ILecturerRepository) IStudentService {
	return &StudentService{
		studentRepo:  studentRepo,
		userRepo:     userRepo,
		lecturerRepo: lecturerRepo,
	}
}

// #5 proses: ambil semua student dari database
func (s *StudentService) GetAllStudents(ctx context.Context) ([]modelpostgre.Student, error) {
	return s.studentRepo.GetAllStudents(ctx)
}

// #6 proses: ambil student berdasarkan student ID
func (s *StudentService) GetStudentByID(ctx context.Context, id string) (*modelpostgre.Student, error) {
	// #6a proses: validasi student ID tidak kosong, lalu ambil student
	if id == "" {
		return nil, errors.New("student ID wajib diisi")
	}
	return s.studentRepo.GetStudentByID(ctx, id)
}

// #7 proses: ambil student ID berdasarkan user ID
func (s *StudentService) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	return s.studentRepo.GetStudentIDByUserID(ctx, userID)
}

// #8 proses: ambil student berdasarkan user ID
func (s *StudentService) GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error) {
	return s.studentRepo.GetStudentByUserID(ctx, userID)
}

// #9 proses: ambil semua student yang dibimbing oleh advisor tertentu
func (s *StudentService) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error) {
	// #9a proses: validasi advisor ID tidak kosong, lalu ambil students
	if advisorID == "" {
		return nil, errors.New("advisor ID wajib diisi")
	}
	return s.studentRepo.GetStudentsByAdvisorID(ctx, advisorID)
}

// #10 proses: buat student profile baru dengan validasi role dan advisor
func (s *StudentService) CreateStudent(ctx context.Context, req modelpostgre.CreateStudentRequest) (*modelpostgre.Student, error) {
	// #10a proses: validasi user ID dan student ID tidak kosong
	if req.UserID == "" {
		return nil, errors.New("user ID wajib diisi")
	}
	if req.StudentID == "" {
		return nil, errors.New("student ID wajib diisi")
	}

	// #10b proses: cek repository tersedia, lalu cari user
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

	// #10c proses: ambil role name dan validasi user harus memiliki role Mahasiswa
	roleName, err := s.userRepo.GetRoleName(ctx, user.RoleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Mahasiswa" {
		return nil, errors.New("user harus memiliki role Mahasiswa untuk membuat student profile")
	}

	// #10d proses: cek apakah user sudah memiliki student profile
	existingStudent, err := s.studentRepo.GetStudentByUserID(ctx, req.UserID)
	if err == nil && existingStudent != nil {
		return nil, errors.New("user sudah memiliki student profile")
	}

	// #10e proses: jika ada advisor ID, validasi advisor ID ada di database
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

	// #10f proses: buat student profile
	student, err := s.studentRepo.CreateStudent(ctx, req)
	if err != nil {
		// #10g proses: handle error duplikasi dengan pesan yang lebih jelas
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

// #11 proses: update advisor untuk student tertentu
func (s *StudentService) UpdateStudentAdvisor(ctx context.Context, studentID string, advisorID string) error {
	// #11a proses: validasi student ID tidak kosong
	if studentID == "" {
		return errors.New("student ID wajib diisi")
	}

	// #11b proses: cari student berdasarkan student ID
	student, err := s.studentRepo.GetStudentByID(ctx, studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("student tidak ditemukan")
		}
		return errors.New("error mengambil data student: " + err.Error())
	}

	// #11c proses: jika ada advisor ID, validasi advisor ID ada di database
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

	// #11d proses: update advisor student
	err = s.studentRepo.UpdateStudentAdvisor(ctx, student.ID, advisorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("student tidak ditemukan")
		}
		return errors.New("error mengupdate advisor: " + err.Error())
	}

	return nil
}
