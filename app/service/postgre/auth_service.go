package service

import (
	"context"
	"database/sql"
	"errors"
	model "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repository "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"
	utilspostgre "sistem-pelaporan-prestasi-mahasiswa/utils/postgre"
	"time"
)

type IAuthService interface {
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.RefreshTokenResponse, error)
	Logout(ctx context.Context, userID string) error
	GetProfile(ctx context.Context, userID string) (*model.GetProfileResponse, error)
}

type AuthService struct {
	userRepo repository.IUserRepository
}

func NewAuthService(userRepo repository.IUserRepository) IAuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username dan password wajib diisi")
	}

	user, err := s.userRepo.FindUserByUsernameOrEmail(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("username atau password tidak valid")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("akun Anda tidak aktif. Silakan hubungi administrator")
	}

	if !utilspostgre.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("username atau password tidak valid")
	}

	token, err := utilspostgre.GenerateToken(*user)
	if err != nil {
		return nil, errors.New("error generating token: " + err.Error())
	}

	refreshToken, err := utilspostgre.GenerateRefreshToken(*user)
	if err != nil {
		return nil, errors.New("error generating refresh token: " + err.Error())
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
	if err := s.userRepo.SaveRefreshToken(ctx, user.ID, refreshToken, expiresAt); err != nil {
		return nil, errors.New("error menyimpan refresh token: " + err.Error())
	}

	permissions, err := s.userRepo.GetUserPermissions(ctx, user.ID)
	if err != nil {
		return nil, errors.New("error mengambil permissions: " + err.Error())
	}

	roleName, err := s.userRepo.GetRoleName(ctx, user.RoleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	response := &model.LoginResponse{
		Status: "success",
		Data: struct {
			Token        string                   `json:"token"`
			RefreshToken string                   `json:"refreshToken"`
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

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.RefreshTokenResponse, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token wajib diisi")
	}

	_, err := s.userRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("refresh token tidak valid atau sudah expired")
		}
		return nil, err
	}

	claims, err := utilspostgre.ValidateRefreshToken(refreshToken)
	if err != nil {
		s.userRepo.DeleteRefreshToken(ctx, refreshToken)
		return nil, errors.New("refresh token tidak valid atau sudah expired")
	}

	user, err := s.userRepo.FindUserByID(ctx, claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data user tidak ditemukan di database")
		}
		return nil, err
	}

	if !user.IsActive {
		s.userRepo.DeleteRefreshToken(ctx, refreshToken)
		return nil, errors.New("akun Anda tidak aktif. Silakan hubungi administrator")
	}

	token, err := utilspostgre.GenerateToken(*user)
	if err != nil {
		return nil, errors.New("error generating token: " + err.Error())
	}

	newRefreshToken, err := utilspostgre.GenerateRefreshToken(*user)
	if err != nil {
		return nil, errors.New("error generating refresh token: " + err.Error())
	}

	s.userRepo.DeleteRefreshToken(ctx, refreshToken)

	expiresAt := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
	if err := s.userRepo.SaveRefreshToken(ctx, user.ID, newRefreshToken, expiresAt); err != nil {
		return nil, errors.New("error menyimpan refresh token: " + err.Error())
	}

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

func (s *AuthService) Logout(ctx context.Context, userID string) error {
	return s.userRepo.DeleteUserRefreshTokens(ctx, userID)
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.GetProfileResponse, error) {
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data user tidak ditemukan di database")
		}
		return nil, err
	}

	roleName, err := s.userRepo.GetRoleName(ctx, user.RoleID)
	if err != nil {
		return nil, errors.New("error mengambil role name: " + err.Error())
	}

	permissions, err := s.userRepo.GetUserPermissions(ctx, user.ID)
	if err != nil {
		return nil, errors.New("error mengambil permissions: " + err.Error())
	}

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

