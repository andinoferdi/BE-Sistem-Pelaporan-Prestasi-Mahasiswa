package service

import (
	"context"
	"errors"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
)

type ILecturerService interface {
	GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error)
	GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error)
}

type LecturerService struct {
	userRepo repositorypostgre.IUserRepository
}

func NewLecturerService(userRepo repositorypostgre.IUserRepository) ILecturerService {
	return &LecturerService{userRepo: userRepo}
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
