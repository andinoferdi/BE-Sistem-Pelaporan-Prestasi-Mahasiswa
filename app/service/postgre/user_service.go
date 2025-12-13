package service

import (
	"context"
	"database/sql"
	"errors"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repository "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	utilspostgre "sistem-pelaporan-prestasi-mahasiswa/utils/postgre"
)

type IUserService interface {
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.User, error)
	UpdateUser(ctx context.Context, id string, req model.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUserRole(ctx context.Context, id string, roleID string) error
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}

type UserService struct {
	userRepo     repository.IUserRepository
	studentRepo  repository.IStudentRepository
	lecturerRepo repository.ILecturerRepository
	db           *sql.DB
}

func NewUserService(userRepo repository.IUserRepository, studentRepo repository.IStudentRepository, lecturerRepo repository.ILecturerRepository, db *sql.DB) IUserService {
	return &UserService{
		userRepo:     userRepo,
		studentRepo:  studentRepo,
		lecturerRepo: lecturerRepo,
		db:           db,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("user ID wajib diisi")
	}
	return s.userRepo.FindUserByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.User, error) {
	if req.Username == "" {
		return nil, errors.New("username wajib diisi")
	}
	if req.Email == "" {
		return nil, errors.New("email wajib diisi")
	}
	if req.Password == "" {
		return nil, errors.New("password wajib diisi")
	}
	if req.FullName == "" {
		return nil, errors.New("full name wajib diisi")
	}
	if req.RoleID == "" {
		return nil, errors.New("role ID wajib diisi")
	}

	existingUser, _ := s.userRepo.FindUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah digunakan")
	}

	existingUser, _ = s.userRepo.FindUserByUsernameOrEmail(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username sudah digunakan")
	}

	_, err := s.userRepo.GetRoleName(ctx, req.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role tidak ditemukan")
		}
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	passwordHash, err := utilspostgre.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error hashing password: " + err.Error())
	}

	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.New("error membuat user: " + err.Error())
	}

	roleName, err := s.userRepo.GetRoleName(ctx, createdUser.RoleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	if roleName == "Mahasiswa" && req.StudentID != "" {
		studentReq := model.CreateStudentRequest{
			UserID:       createdUser.ID,
			StudentID:    req.StudentID,
			ProgramStudy: req.ProgramStudy,
			AcademicYear: req.AcademicYear,
			AdvisorID:    req.AdvisorID,
		}
		_, err := s.studentRepo.CreateStudent(ctx, studentReq)
		if err != nil {
			return nil, errors.New("error membuat student profile: " + err.Error())
		}
	}

	if roleName == "Dosen Wali" && req.LecturerID != "" {
		lecturerReq := model.CreateLecturerRequest{
			UserID:     createdUser.ID,
			LecturerID: req.LecturerID,
			Department: req.Department,
		}
		_, err := s.lecturerRepo.CreateLecturer(ctx, lecturerReq)
		if err != nil {
			return nil, errors.New("error membuat lecturer profile: " + err.Error())
		}
	}

	return createdUser, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req model.UpdateUserRequest) (*model.User, error) {
	if id == "" {
		return nil, errors.New("user ID wajib diisi")
	}

	existingUser, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}

	if req.Username == "" {
		return nil, errors.New("username wajib diisi")
	}
	if req.Email == "" {
		return nil, errors.New("email wajib diisi")
	}
	if req.FullName == "" {
		return nil, errors.New("full name wajib diisi")
	}
	if req.RoleID == "" {
		return nil, errors.New("role ID wajib diisi")
	}

	userByEmail, _ := s.userRepo.FindUserByEmail(ctx, req.Email)
	if userByEmail != nil && userByEmail.ID != id {
		return nil, errors.New("email sudah digunakan oleh user lain")
	}

	userByUsername, _ := s.userRepo.FindUserByUsernameOrEmail(ctx, req.Username)
	if userByUsername != nil && userByUsername.ID != id {
		return nil, errors.New("username sudah digunakan oleh user lain")
	}

	roleName, err := s.userRepo.GetRoleName(ctx, req.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role tidak ditemukan")
		}
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	isActive := existingUser.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		RoleID:   req.RoleID,
		IsActive: isActive,
	}

	updatedUser, err := s.userRepo.UpdateUser(ctx, id, user)
	if err != nil {
		return nil, errors.New("error mengupdate user: " + err.Error())
	}

	oldRoleName, _ := s.userRepo.GetRoleName(ctx, existingUser.RoleID)
	if oldRoleName != roleName {
		if oldRoleName == "Mahasiswa" {
			_, err = s.studentRepo.GetStudentByUserID(ctx, id)
			if err == nil {
				return nil, errors.New("tidak dapat mengubah role user yang sudah memiliki profil. Hapus profil terlebih dahulu")
			}
		} else if oldRoleName == "Dosen Wali" {
			_, err = s.lecturerRepo.GetLecturerByUserID(ctx, id)
			if err == nil {
				return nil, errors.New("tidak dapat mengubah role user yang sudah memiliki profil. Hapus profil terlebih dahulu")
			}
		}
	}

	return updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID wajib diisi")
	}

	existingUser, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user tidak ditemukan")
		}
		return err
	}

	roleName, err := s.userRepo.GetRoleName(ctx, existingUser.RoleID)
	if err == nil {
		if roleName == "Dosen Wali" {
			lecturer, err := s.lecturerRepo.GetLecturerByUserID(ctx, id)
			if err == nil {
				students, _ := s.studentRepo.GetStudentsByAdvisorID(ctx, lecturer.ID)
				if len(students) > 0 {
					return errors.New("tidak dapat menghapus dosen wali yang masih memiliki mahasiswa bimbingan. Pindahkan mahasiswa terlebih dahulu")
				}
			}
		}
	}

	err = s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user tidak ditemukan")
		}
		return errors.New("error menghapus user: " + err.Error())
	}

	return nil
}

func (s *UserService) UpdateUserRole(ctx context.Context, id string, roleID string) error {
	if id == "" {
		return errors.New("user ID wajib diisi")
	}
	if roleID == "" {
		return errors.New("role ID wajib diisi")
	}

	existingUser, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user tidak ditemukan")
		}
		return err
	}

	oldRoleName, err := s.userRepo.GetRoleName(ctx, existingUser.RoleID)
	if err != nil {
		return errors.New("error mengambil role name: " + err.Error())
	}

	roleName, err := s.userRepo.GetRoleName(ctx, roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("role tidak ditemukan")
		}
		return errors.New("error mengambil role name: " + err.Error())
	}

	if oldRoleName == roleName {
		return nil
	}

	if oldRoleName == "Mahasiswa" {
		_, err = s.studentRepo.GetStudentByUserID(ctx, id)
		if err == nil {
			return errors.New("tidak dapat mengubah role user yang sudah memiliki profil mahasiswa. Hapus profil terlebih dahulu")
		}
	} else if oldRoleName == "Dosen Wali" {
		lecturer, err := s.lecturerRepo.GetLecturerByUserID(ctx, id)
		if err == nil {
			students, _ := s.studentRepo.GetStudentsByAdvisorID(ctx, lecturer.ID)
			if len(students) > 0 {
				return errors.New("tidak dapat mengubah role dosen wali yang masih memiliki mahasiswa bimbingan. Pindahkan mahasiswa terlebih dahulu")
			}
		}
	}

	err = s.userRepo.UpdateUserRole(ctx, id, roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user tidak ditemukan")
		}
		return errors.New("error mengupdate role user: " + err.Error())
	}

	return nil
}

func (s *UserService) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return s.userRepo.GetAllRoles(ctx)
}
