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

// #2 proses: definisikan interface untuk operasi autentikasi
type IAuthService interface {
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.RefreshTokenResponse, error)
	Logout(ctx context.Context, userID string) error
	GetProfile(ctx context.Context, userID string) (*model.GetProfileResponse, error)
}

// #3 proses: struct service untuk autentikasi dengan dependency user repository
type AuthService struct {
	userRepo repository.IUserRepository
}

// #4 proses: constructor untuk membuat instance AuthService baru
func NewAuthService(userRepo repository.IUserRepository) IAuthService {
	return &AuthService{userRepo: userRepo}
}

// #5 proses: proses login user dengan validasi kredensial dan generate token
func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	// #5a proses: validasi input username dan password tidak kosong
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username dan password wajib diisi")
	}

	// #5b proses: cari user berdasarkan username atau email
	user, err := s.userRepo.FindUserByUsernameOrEmail(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("username atau password tidak valid")
		}
		return nil, err
	}

	// #5c proses: cek apakah user aktif
	if !user.IsActive {
		return nil, errors.New("akun Anda tidak aktif. Silakan hubungi administrator")
	}

	// #5d proses: verifikasi password dengan hash yang tersimpan
	if !utilspostgre.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("username atau password tidak valid")
	}

	// #5e proses: generate access token dan refresh token
	token, err := utilspostgre.GenerateToken(*user)
	if err != nil {
		return nil, errors.New("error generating token: " + err.Error())
	}

	refreshToken, err := utilspostgre.GenerateRefreshToken(*user)
	if err != nil {
		return nil, errors.New("error generating refresh token: " + err.Error())
	}

	// #5f proses: ambil permissions dan role name user
	permissions, err := s.userRepo.GetUserPermissions(ctx, user.ID)
	if err != nil {
		return nil, errors.New("error mengambil permissions: " + err.Error())
	}

	roleName, err := s.userRepo.GetRoleName(ctx, user.RoleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	// #5g proses: build response dengan token dan data user
	response := &model.LoginResponse{
		Status: "success",
		Data: struct {
			Token        string                  `json:"token"`
			RefreshToken string                  `json:"refreshToken"`
			User         model.LoginUserResponse `json:"user"`
		}{
			Token:        token,
			RefreshToken: refreshToken,
			User: model.LoginUserResponse{
				ID:          user.ID,
				Username:    user.Username,
				FullName:    user.FullName,
				Role:        roleName,
				Permissions: permissions,
			},
		},
	}

	return response, nil
}

// #6 proses: refresh access token menggunakan refresh token yang valid
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.RefreshTokenResponse, error) {
	// #6a proses: validasi refresh token tidak kosong
	if refreshToken == "" {
		return nil, errors.New("refresh token wajib diisi")
	}

	// #6b proses: validasi refresh token dan ambil claims
	claims, err := utilspostgre.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token tidak valid atau sudah expired")
	}

	// #6c proses: cari user berdasarkan user ID dari claims
	user, err := s.userRepo.FindUserByID(ctx, claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data user tidak ditemukan di database")
		}
		return nil, err
	}

	// #6d proses: cek apakah user masih aktif
	if !user.IsActive {
		return nil, errors.New("akun Anda tidak aktif. Silakan hubungi administrator")
	}

	// #6e proses: generate access token dan refresh token baru
	token, err := utilspostgre.GenerateToken(*user)
	if err != nil {
		return nil, errors.New("error generating token: " + err.Error())
	}

	newRefreshToken, err := utilspostgre.GenerateRefreshToken(*user)
	if err != nil {
		return nil, errors.New("error generating refresh token: " + err.Error())
	}

	// #6f proses: build response dengan token baru
	response := &model.RefreshTokenResponse{
		Status: "success",
		Data: struct {
			Token        string `json:"token"`
			RefreshToken string `json:"refreshToken"`
		}{
			Token:        token,
			RefreshToken: newRefreshToken,
		},
	}

	return response, nil
}

// #7 proses: proses logout user, saat ini hanya return success
func (s *AuthService) Logout(ctx context.Context, userID string) error {
	return nil
}

// #8 proses: ambil profil user lengkap dengan role dan permissions
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.GetProfileResponse, error) {
	// #8a proses: cari user berdasarkan user ID
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data user tidak ditemukan di database")
		}
		return nil, err
	}

	// #8b proses: ambil role name dan permissions user
	roleName, err := s.userRepo.GetRoleName(ctx, user.RoleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	permissions, err := s.userRepo.GetUserPermissions(ctx, user.ID)
	if err != nil {
		return nil, errors.New("error mengambil permissions: " + err.Error())
	}

	// #8c proses: build response dengan data profil lengkap
	response := &model.GetProfileResponse{
		Status: "success",
		Data: struct {
			UserID      string   `json:"user_id"`
			Username    string   `json:"username"`
			Email       string   `json:"email"`
			FullName    string   `json:"full_name"`
			RoleID      string   `json:"role_id"`
			Role        string   `json:"role"`
			Permissions []string `json:"permissions"`
		}{
			UserID:      user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FullName:    user.FullName,
			RoleID:      user.RoleID,
			Role:        roleName,
			Permissions: permissions,
		},
	}

	return response, nil
}
