package route_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	routepostgre "sistem-pelaporan-prestasi-mahasiswa/route/postgre"

	"github.com/gofiber/fiber/v2"
)

type mockAuthService struct {
	loginResponse   *modelpostgre.LoginResponse
	loginErr        error
	refreshResponse *modelpostgre.RefreshTokenResponse
	refreshErr      error
	profileResponse *modelpostgre.GetProfileResponse
	profileErr      error
	logoutErr       error
}

func (m *mockAuthService) Login(ctx context.Context, req modelpostgre.LoginRequest) (*modelpostgre.LoginResponse, error) {
	if m.loginErr != nil {
		return nil, m.loginErr
	}
	return m.loginResponse, nil
}

func (m *mockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*modelpostgre.RefreshTokenResponse, error) {
	if m.refreshErr != nil {
		return nil, m.refreshErr
	}
	return m.refreshResponse, nil
}

func (m *mockAuthService) GetProfile(ctx context.Context, userID string) (*modelpostgre.GetProfileResponse, error) {
	if m.profileErr != nil {
		return nil, m.profileErr
	}
	return m.profileResponse, nil
}

func (m *mockAuthService) Logout(ctx context.Context, userID string) error {
	return m.logoutErr
}

func TestLoginRoute_Success(t *testing.T) {
	app := fiber.New()
	mockAuthService := &mockAuthService{
		loginResponse: &modelpostgre.LoginResponse{
			Status: "success",
			Data: struct {
				Token        string                         `json:"token"`
				RefreshToken string                         `json:"refreshToken"`
				User         modelpostgre.LoginUserResponse `json:"user"`
			}{
				Token:        "access_token",
				RefreshToken: "refresh_token",
				User: modelpostgre.LoginUserResponse{
					ID:          "user-id-1",
					Username:    "testuser",
					FullName:    "Test User",
					Role:        "Mahasiswa",
					Permissions: []string{"achievement:create", "achievement:read"},
				},
			},
		},
	}

	routepostgre.AuthRoutes(app, mockAuthService, "instance-id")

	reqBody := modelpostgre.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestLoginRoute_InvalidBody(t *testing.T) {
	app := fiber.New()
	mockAuthService := &mockAuthService{}
	routepostgre.AuthRoutes(app, mockAuthService, "instance-id")

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestLoginRoute_InvalidCredentials(t *testing.T) {
	app := fiber.New()
	mockAuthService := &mockAuthService{
		loginErr: errors.New("username atau password tidak valid"),
	}
	routepostgre.AuthRoutes(app, mockAuthService, "instance-id")

	reqBody := modelpostgre.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestRefreshTokenRoute_Success(t *testing.T) {
	app := fiber.New()
	mockAuthService := &mockAuthService{
		refreshResponse: &modelpostgre.RefreshTokenResponse{
			Status: "success",
			Data: struct {
				Token        string `json:"token"`
				RefreshToken string `json:"refreshToken"`
			}{
				Token:        "new_access_token",
				RefreshToken: "new_refresh_token",
			},
		},
	}
	routepostgre.AuthRoutes(app, mockAuthService, "instance-id")

	reqBody := modelpostgre.RefreshTokenRequest{
		RefreshToken: "valid_refresh_token",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestRefreshTokenRoute_InvalidToken(t *testing.T) {
	app := fiber.New()
	mockAuthService := &mockAuthService{
		refreshErr: errors.New("refresh token tidak valid atau sudah expired"),
	}
	routepostgre.AuthRoutes(app, mockAuthService, "instance-id")

	reqBody := modelpostgre.RefreshTokenRequest{
		RefreshToken: "invalid_token",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestGetProfileRoute_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	app := fiber.New()
	mockAuthService := &mockAuthService{
		profileResponse: &modelpostgre.GetProfileResponse{
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
				UserID:      userID,
				Username:    "testuser",
				Email:       email,
				FullName:    "Test User",
				RoleID:      roleID,
				Role:        "Mahasiswa",
				Permissions: []string{"achievement:create", "achievement:read"},
			},
		},
	}
	routepostgre.AuthRoutes(app, mockAuthService, "instance-id")

	req := createRequestWithToken("GET", "/api/v1/auth/profile", nil, token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestHealthRoute_Success(t *testing.T) {
	app := fiber.New()
	mockAuthService := &mockAuthService{}
	routepostgre.AuthRoutes(app, mockAuthService, "test-instance-id")

	req := httptest.NewRequest("GET", "/api/v1/health", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestLogoutRoute_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	app := fiber.New()
	mockAuthService := &mockAuthService{}

	routepostgre.AuthRoutes(app, mockAuthService, "instance-id")

	req := createRequestWithToken("POST", "/api/v1/auth/logout", nil, token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestLogoutRoute_MissingToken(t *testing.T) {
	app := fiber.New()
	mockAuthService := &mockAuthService{}

	routepostgre.AuthRoutes(app, mockAuthService, "instance-id")

	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestGetProfileRoute_MissingToken(t *testing.T) {
	app := fiber.New()
	mockAuthService := &mockAuthService{
		profileResponse: &modelpostgre.GetProfileResponse{
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
				UserID:      "user-id-1",
				Username:    "testuser",
				Email:       "test@example.com",
				FullName:    "Test User",
				RoleID:      "role-id-1",
				Role:        "Mahasiswa",
				Permissions: []string{"achievement:create", "achievement:read"},
			},
		},
	}
	routepostgre.AuthRoutes(app, mockAuthService, "instance-id")

	req := httptest.NewRequest("GET", "/api/v1/auth/profile", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestGetProfileRoute_NotFound(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	app := fiber.New()
	mockAuthService := &mockAuthService{
		profileErr: errors.New("data user tidak ditemukan di database"),
	}
	routepostgre.AuthRoutes(app, mockAuthService, "instance-id")

	req := createRequestWithToken("GET", "/api/v1/auth/profile", nil, token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusNotFound)
}
