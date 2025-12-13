package service

import (
	"context"
	"database/sql"
	"errors"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	"strings"
)

type ILecturerService interface {
	GetAllLecturers(ctx context.Context) ([]modelpostgre.Lecturer, error)
	GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error)
	GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error)
	CreateLecturer(ctx context.Context, req modelpostgre.CreateLecturerRequest) (*modelpostgre.Lecturer, error)
}

type LecturerService struct {
	userRepo     repositorypostgre.IUserRepository
	lecturerRepo repositorypostgre.ILecturerRepository
}

func NewLecturerService(userRepo repositorypostgre.IUserRepository) ILecturerService {
	return &LecturerService{userRepo: userRepo}
}

func NewLecturerServiceWithDeps(userRepo repositorypostgre.IUserRepository, lecturerRepo repositorypostgre.ILecturerRepository) ILecturerService {
	return &LecturerService{
		userRepo:     userRepo,
		lecturerRepo: lecturerRepo,
	}
}

func (s *LecturerService) GetAllLecturers(ctx context.Context) ([]modelpostgre.Lecturer, error) {
	if s.lecturerRepo == nil {
		return nil, errors.New("lecturer repository tidak tersedia")
	}
	return s.lecturerRepo.GetAllLecturers(ctx)
}

func (s *LecturerService) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	if userID == "" {
		return nil, errors.New("user ID wajib diisi")
	}
	return s.userRepo.GetLecturerByUserID(ctx, userID)
}

func (s *LecturerService) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	if id == "" {
		return nil, errors.New("lecturer ID wajib diisi")
	}
	return s.userRepo.GetLecturerByID(ctx, id)
}

func (s *LecturerService) CreateLecturer(ctx context.Context, req modelpostgre.CreateLecturerRequest) (*modelpostgre.Lecturer, error) {
	if req.UserID == "" {
		return nil, errors.New("user ID wajib diisi")
	}
	if req.LecturerID == "" {
		return nil, errors.New("lecturer ID wajib diisi")
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

	if roleName != "Dosen Wali" {
		return nil, errors.New("user harus memiliki role Dosen Wali untuk membuat lecturer profile")
	}

	existingLecturer, err := s.userRepo.GetLecturerByUserID(ctx, req.UserID)
	if err == nil && existingLecturer != nil {
		return nil, errors.New("user sudah memiliki lecturer profile")
	}

	if s.lecturerRepo == nil {
		return nil, errors.New("lecturer repository tidak tersedia")
	}

	lecturer, err := s.lecturerRepo.CreateLecturer(ctx, req)
	if err != nil {
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
