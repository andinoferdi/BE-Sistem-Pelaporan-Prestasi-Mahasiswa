package service

import (
	"context"
	"errors"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repository "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
)

type IUserService interface {
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.User, error)
	UpdateUser(ctx context.Context, id string, req model.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUserRole(ctx context.Context, id string, roleID string) error
}

type UserService struct {
	userRepo repository.IUserRepository
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return nil, errors.New("fitur ini belum diimplementasikan")
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("user ID wajib diisi")
	}
	return s.userRepo.FindUserByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.User, error) {
	return nil, errors.New("fitur ini belum diimplementasikan")
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req model.UpdateUserRequest) (*model.User, error) {
	return nil, errors.New("fitur ini belum diimplementasikan")
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return errors.New("fitur ini belum diimplementasikan")
}

func (s *UserService) UpdateUserRole(ctx context.Context, id string, roleID string) error {
	return errors.New("fitur ini belum diimplementasikan")
}
