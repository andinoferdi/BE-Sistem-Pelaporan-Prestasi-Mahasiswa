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

// #2 proses: definisikan interface untuk operasi lecturer
type ILecturerService interface {
	GetAllLecturers(ctx context.Context) ([]modelpostgre.Lecturer, error)
	GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error)
	GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error)
	CreateLecturer(ctx context.Context, req modelpostgre.CreateLecturerRequest) (*modelpostgre.Lecturer, error)
}

// #3 proses: struct service untuk lecturer dengan dependency user dan lecturer repository
type LecturerService struct {
	userRepo     repositorypostgre.IUserRepository
	lecturerRepo repositorypostgre.ILecturerRepository
}

// #4 proses: constructor untuk membuat instance LecturerService baru
func NewLecturerService(userRepo repositorypostgre.IUserRepository, lecturerRepo repositorypostgre.ILecturerRepository) ILecturerService {
	return &LecturerService{
		userRepo:     userRepo,
		lecturerRepo: lecturerRepo,
	}
}

// #5 proses: ambil semua lecturer dari database
func (s *LecturerService) GetAllLecturers(ctx context.Context) ([]modelpostgre.Lecturer, error) {
	// #5a proses: cek repository tersedia, lalu ambil semua lecturer
	if s.lecturerRepo == nil {
		return nil, errors.New("lecturer repository tidak tersedia")
	}
	return s.lecturerRepo.GetAllLecturers(ctx)
}

// #6 proses: ambil lecturer berdasarkan user ID
func (s *LecturerService) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	// #6a proses: validasi user ID tidak kosong, lalu ambil lecturer
	if userID == "" {
		return nil, errors.New("user ID wajib diisi")
	}
	return s.userRepo.GetLecturerByUserID(ctx, userID)
}

// #7 proses: ambil lecturer berdasarkan lecturer ID
func (s *LecturerService) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	// #7a proses: validasi lecturer ID tidak kosong, lalu ambil lecturer
	if id == "" {
		return nil, errors.New("lecturer ID wajib diisi")
	}
	return s.userRepo.GetLecturerByID(ctx, id)
}

// #8 proses: buat lecturer profile baru dengan validasi role dan duplikasi
func (s *LecturerService) CreateLecturer(ctx context.Context, req modelpostgre.CreateLecturerRequest) (*modelpostgre.Lecturer, error) {
	// #8a proses: validasi user ID dan lecturer ID tidak kosong
	if req.UserID == "" {
		return nil, errors.New("user ID wajib diisi")
	}
	if req.LecturerID == "" {
		return nil, errors.New("lecturer ID wajib diisi")
	}

	// #8b proses: cari user berdasarkan user ID
	user, err := s.userRepo.FindUserByID(ctx, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, errors.New("error mengambil data user: " + err.Error())
	}

	// #8c proses: ambil role name dan validasi user harus memiliki role Dosen Wali
	roleName, err := s.userRepo.GetRoleName(ctx, user.RoleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName != "Dosen Wali" {
		return nil, errors.New("user harus memiliki role Dosen Wali untuk membuat lecturer profile")
	}

	// #8d proses: cek apakah user sudah memiliki lecturer profile
	existingLecturer, err := s.userRepo.GetLecturerByUserID(ctx, req.UserID)
	if err == nil && existingLecturer != nil {
		return nil, errors.New("user sudah memiliki lecturer profile")
	}

	// #8e proses: cek repository tersedia, lalu buat lecturer profile
	if s.lecturerRepo == nil {
		return nil, errors.New("lecturer repository tidak tersedia")
	}

	lecturer, err := s.lecturerRepo.CreateLecturer(ctx, req)
	if err != nil {
		// #8f proses: handle error duplikasi dengan pesan yang lebih jelas
		errStr := err.Error()
		if strings.Contains(errStr, "duplicate key value violates unique constraint") {
			if strings.Contains(errStr, "lecturers_lecturer_id_key") {
				return nil, errors.New("lecturer ID sudah digunakan")
			}
			if strings.Contains(errStr, "lecturers_user_id_key") {
				return nil, errors.New("user sudah memiliki lecturer profile")
			}
			return nil, errors.New("data duplikat terdeteksi")
		}
		return nil, errors.New("error membuat lecturer profile: " + errStr)
	}

	return lecturer, nil
}
