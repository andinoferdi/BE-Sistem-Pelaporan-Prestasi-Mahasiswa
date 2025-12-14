package service

// #1 proses: import library yang diperlukan untuk context, database, errors, dan utils
import (
	"context"
	"database/sql"
	"errors"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repository "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	utilspostgre "sistem-pelaporan-prestasi-mahasiswa/utils/postgre"
)

// #2 proses: definisikan interface untuk operasi user
type IUserService interface {
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.User, error)
	UpdateUser(ctx context.Context, id string, req model.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUserRole(ctx context.Context, id string, roleID string) error
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}

// #3 proses: struct service untuk user dengan dependency user, student, lecturer repository dan database connection
type UserService struct {
	userRepo     repository.IUserRepository
	studentRepo  repository.IStudentRepository
	lecturerRepo repository.ILecturerRepository
	db           *sql.DB
}

// #4 proses: constructor untuk membuat instance UserService baru
func NewUserService(userRepo repository.IUserRepository, studentRepo repository.IStudentRepository, lecturerRepo repository.ILecturerRepository, db *sql.DB) IUserService {
	return &UserService{
		userRepo:     userRepo,
		studentRepo:  studentRepo,
		lecturerRepo: lecturerRepo,
		db:           db,
	}
}

// #5 proses: ambil semua user dari database
func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

// #6 proses: ambil user berdasarkan user ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	// #6a proses: validasi user ID tidak kosong, lalu ambil user
	if id == "" {
		return nil, errors.New("user ID wajib diisi")
	}
	return s.userRepo.FindUserByID(ctx, id)
}

// #7 proses: buat user baru dengan transaction, bisa sekaligus buat student atau lecturer profile
func (s *UserService) CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.User, error) {
	// #7a proses: validasi semua field wajib tidak kosong
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

	// #7b proses: cek email dan username belum digunakan
	existingUser, _ := s.userRepo.FindUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah digunakan")
	}

	existingUser, _ = s.userRepo.FindUserByUsernameOrEmail(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username sudah digunakan")
	}

	// #7c proses: hash password sebelum disimpan
	passwordHash, err := utilspostgre.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error hashing password: " + err.Error())
	}

	// #7d proses: ambil role name untuk menentukan apakah perlu buat profile
	roleName, err := s.userRepo.GetRoleName(ctx, req.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role tidak ditemukan")
		}
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	// #7e proses: mulai transaction untuk operasi multi-step
	tx, err := s.db.Begin()
	if err != nil {
		return nil, errors.New("error memulai transaction: " + err.Error())
	}
	defer tx.Rollback()

	// #7f proses: buat user object dengan data dari request
	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	// #7g proses: cast repository ke concrete type untuk akses method WithTx
	userRepoImpl, ok := s.userRepo.(*repository.UserRepository)
	if !ok {
		return nil, errors.New("error casting user repository")
	}

	// #7h proses: buat user dalam transaction
	createdUser, err := userRepoImpl.CreateUserWithTx(ctx, tx, user)
	if err != nil {
		return nil, errors.New("error membuat user: " + err.Error())
	}

	// #7i proses: jika role Mahasiswa dan ada student ID, buat student profile dalam transaction
	if roleName == "Mahasiswa" && req.StudentID != "" {
		studentReq := model.CreateStudentRequest{
			UserID:       createdUser.ID,
			StudentID:    req.StudentID,
			ProgramStudy: req.ProgramStudy,
			AcademicYear: req.AcademicYear,
			AdvisorID:    req.AdvisorID,
		}
		studentRepoImpl, ok := s.studentRepo.(*repository.StudentRepository)
		if !ok {
			return nil, errors.New("error casting student repository")
		}
		_, err := studentRepoImpl.CreateStudentWithTx(ctx, tx, studentReq)
		if err != nil {
			return nil, errors.New("error membuat student profile: " + err.Error())
		}
	}

	// #7j proses: jika role Dosen Wali dan ada lecturer ID, buat lecturer profile dalam transaction
	if roleName == "Dosen Wali" && req.LecturerID != "" {
		lecturerReq := model.CreateLecturerRequest{
			UserID:     createdUser.ID,
			LecturerID: req.LecturerID,
			Department: req.Department,
		}
		lecturerRepoImpl, ok := s.lecturerRepo.(*repository.LecturerRepository)
		if !ok {
			return nil, errors.New("error casting lecturer repository")
		}
		_, err := lecturerRepoImpl.CreateLecturerWithTx(ctx, tx, lecturerReq)
		if err != nil {
			return nil, errors.New("error membuat lecturer profile: " + err.Error())
		}
	}

	// #7k proses: commit transaction jika semua operasi berhasil
	if err := tx.Commit(); err != nil {
		return nil, errors.New("error commit transaction: " + err.Error())
	}

	return createdUser, nil
}

// #8 proses: update user dengan validasi email, username, dan role change
func (s *UserService) UpdateUser(ctx context.Context, id string, req model.UpdateUserRequest) (*model.User, error) {
	// #8a proses: validasi user ID tidak kosong
	if id == "" {
		return nil, errors.New("user ID wajib diisi")
	}

	// #8b proses: cari user yang akan diupdate
	existingUser, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}

	// #8c proses: validasi semua field wajib tidak kosong
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

	// #8d proses: cek email dan username tidak digunakan oleh user lain
	userByEmail, _ := s.userRepo.FindUserByEmail(ctx, req.Email)
	if userByEmail != nil && userByEmail.ID != id {
		return nil, errors.New("email sudah digunakan oleh user lain")
	}

	userByUsername, _ := s.userRepo.FindUserByUsernameOrEmail(ctx, req.Username)
	if userByUsername != nil && userByUsername.ID != id {
		return nil, errors.New("username sudah digunakan oleh user lain")
	}

	// #8e proses: ambil role name baru
	roleName, err := s.userRepo.GetRoleName(ctx, req.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role tidak ditemukan")
		}
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	// #8f proses: set isActive dari request atau gunakan nilai existing
	isActive := existingUser.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// #8g proses: buat user object dengan data update
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		RoleID:   req.RoleID,
		IsActive: isActive,
	}

	// #8h proses: update user di database
	updatedUser, err := s.userRepo.UpdateUser(ctx, id, user)
	if err != nil {
		return nil, errors.New("error mengupdate user: " + err.Error())
	}

	// #8i proses: jika role berubah, cek apakah user sudah punya profile yang harus dihapus dulu
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

// #9 proses: hapus user dengan validasi dosen wali tidak punya mahasiswa bimbingan
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	// #9a proses: validasi user ID tidak kosong
	if id == "" {
		return errors.New("user ID wajib diisi")
	}

	// #9b proses: cari user yang akan dihapus
	existingUser, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user tidak ditemukan")
		}
		return err
	}

	// #9c proses: jika user adalah dosen wali, cek apakah masih punya mahasiswa bimbingan
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

	// #9d proses: hapus user dari database
	err = s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user tidak ditemukan")
		}
		return errors.New("error menghapus user: " + err.Error())
	}

	return nil
}

// #10 proses: update role user dengan validasi profile dan mahasiswa bimbingan
func (s *UserService) UpdateUserRole(ctx context.Context, id string, roleID string) error {
	// #10a proses: validasi user ID dan role ID tidak kosong
	if id == "" {
		return errors.New("user ID wajib diisi")
	}
	if roleID == "" {
		return errors.New("role ID wajib diisi")
	}

	// #10b proses: cari user yang akan diupdate
	existingUser, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user tidak ditemukan")
		}
		return err
	}

	// #10c proses: ambil role name lama dan baru
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

	// #10d proses: jika role tidak berubah, tidak perlu update
	if oldRoleName == roleName {
		return nil
	}

	// #10e proses: jika role lama adalah Mahasiswa, cek apakah sudah punya profile
	if oldRoleName == "Mahasiswa" {
		_, err = s.studentRepo.GetStudentByUserID(ctx, id)
		if err == nil {
			return errors.New("tidak dapat mengubah role user yang sudah memiliki profil mahasiswa. Hapus profil terlebih dahulu")
		}
	} else if oldRoleName == "Dosen Wali" {
		// #10f proses: jika role lama adalah Dosen Wali, cek apakah masih punya mahasiswa bimbingan
		lecturer, err := s.lecturerRepo.GetLecturerByUserID(ctx, id)
		if err == nil {
			students, _ := s.studentRepo.GetStudentsByAdvisorID(ctx, lecturer.ID)
			if len(students) > 0 {
				return errors.New("tidak dapat mengubah role dosen wali yang masih memiliki mahasiswa bimbingan. Pindahkan mahasiswa terlebih dahulu")
			}
		}
	}

	// #10g proses: update role user di database
	err = s.userRepo.UpdateUserRole(ctx, id, roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user tidak ditemukan")
		}
		return errors.New("error mengupdate role user: " + err.Error())
	}

	return nil
}

// #11 proses: ambil semua role dari database
func (s *UserService) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return s.userRepo.GetAllRoles(ctx)
}
