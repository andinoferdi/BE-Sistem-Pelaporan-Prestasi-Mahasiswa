package model_test

import (
	"encoding/json"
	"testing"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

func TestLoginRequest_StructCreation(t *testing.T) {
	req := modelpostgre.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	if req.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", req.Username)
	}
	if req.Password != "password123" {
		t.Errorf("Expected Password 'password123', got '%s'", req.Password)
	}
}

func TestLoginRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["username"] != "testuser" {
		t.Errorf("Expected username 'testuser', got '%v'", result["username"])
	}
	if result["password"] != "password123" {
		t.Errorf("Expected password 'password123', got '%v'", result["password"])
	}
}

func TestLoginRequest_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"username": "testuser",
		"password": "password123"
	}`

	var req modelpostgre.LoginRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", req.Username)
	}
	if req.Password != "password123" {
		t.Errorf("Expected Password 'password123', got '%s'", req.Password)
	}
}

func TestLoginRequest_ZeroValues(t *testing.T) {
	var req modelpostgre.LoginRequest

	if req.Username != "" {
		t.Errorf("Expected empty Username, got '%s'", req.Username)
	}
	if req.Password != "" {
		t.Errorf("Expected empty Password, got '%s'", req.Password)
	}
}

func TestLoginUserResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.LoginUserResponse{
		ID:          "user-id-1",
		Username:    "testuser",
		FullName:    "Test User",
		Role:        "Mahasiswa",
		Permissions: []string{"read:achievements", "write:achievements"},
	}

	if resp.ID != "user-id-1" {
		t.Errorf("Expected ID 'user-id-1', got '%s'", resp.ID)
	}
	if resp.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", resp.Username)
	}
	if resp.FullName != "Test User" {
		t.Errorf("Expected FullName 'Test User', got '%s'", resp.FullName)
	}
	if resp.Role != "Mahasiswa" {
		t.Errorf("Expected Role 'Mahasiswa', got '%s'", resp.Role)
	}
	if len(resp.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(resp.Permissions))
	}
}

func TestLoginUserResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.LoginUserResponse{
		ID:          "user-id-1",
		Username:    "testuser",
		FullName:    "Test User",
		Role:        "Mahasiswa",
		Permissions: []string{"read:achievements", "write:achievements"},
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != "user-id-1" {
		t.Errorf("Expected id 'user-id-1', got '%v'", result["id"])
	}
	if result["username"] != "testuser" {
		t.Errorf("Expected username 'testuser', got '%v'", result["username"])
	}
	if result["fullName"] != "Test User" {
		t.Errorf("Expected fullName 'Test User', got '%v'", result["fullName"])
	}
	if result["role"] != "Mahasiswa" {
		t.Errorf("Expected role 'Mahasiswa', got '%v'", result["role"])
	}
}

func TestLoginResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.LoginResponse{
		Status: "success",
	}
	resp.Data.Token = "access_token_123"
	resp.Data.RefreshToken = "refresh_token_123"
	resp.Data.User = modelpostgre.LoginUserResponse{
		ID:       "user-id-1",
		Username: "testuser",
		FullName: "Test User",
		Role:     "Mahasiswa",
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Token != "access_token_123" {
		t.Errorf("Expected Data.Token 'access_token_123', got '%s'", resp.Data.Token)
	}
	if resp.Data.RefreshToken != "refresh_token_123" {
		t.Errorf("Expected Data.RefreshToken 'refresh_token_123', got '%s'", resp.Data.RefreshToken)
	}
	if resp.Data.User.Username != "testuser" {
		t.Errorf("Expected Data.User.Username 'testuser', got '%s'", resp.Data.User.Username)
	}
}

func TestLoginResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.LoginResponse{
		Status: "success",
	}
	resp.Data.Token = "access_token_123"
	resp.Data.RefreshToken = "refresh_token_123"
	resp.Data.User = modelpostgre.LoginUserResponse{
		ID:       "user-id-1",
		Username: "testuser",
		FullName: "Test User",
		Role:     "Mahasiswa",
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", result["status"])
	}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if data["token"] != "access_token_123" {
			t.Errorf("Expected token 'access_token_123', got '%v'", data["token"])
		}
		if data["refreshToken"] != "refresh_token_123" {
			t.Errorf("Expected refreshToken 'refresh_token_123', got '%v'", data["refreshToken"])
		}
	} else {
		t.Error("Expected data to be a map")
	}
}

func TestRefreshTokenRequest_StructCreation(t *testing.T) {
	req := modelpostgre.RefreshTokenRequest{
		RefreshToken: "refresh_token_123",
	}

	if req.RefreshToken != "refresh_token_123" {
		t.Errorf("Expected RefreshToken 'refresh_token_123', got '%s'", req.RefreshToken)
	}
}

func TestRefreshTokenRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.RefreshTokenRequest{
		RefreshToken: "refresh_token_123",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["refreshToken"] != "refresh_token_123" {
		t.Errorf("Expected refreshToken 'refresh_token_123', got '%v'", result["refreshToken"])
	}
}

func TestRefreshTokenRequest_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"refreshToken": "refresh_token_123"
	}`

	var req modelpostgre.RefreshTokenRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.RefreshToken != "refresh_token_123" {
		t.Errorf("Expected RefreshToken 'refresh_token_123', got '%s'", req.RefreshToken)
	}
}

func TestRefreshTokenResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.RefreshTokenResponse{
		Status: "success",
	}
	resp.Data.Token = "new_access_token_123"
	resp.Data.RefreshToken = "new_refresh_token_123"

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Token != "new_access_token_123" {
		t.Errorf("Expected Data.Token 'new_access_token_123', got '%s'", resp.Data.Token)
	}
	if resp.Data.RefreshToken != "new_refresh_token_123" {
		t.Errorf("Expected Data.RefreshToken 'new_refresh_token_123', got '%s'", resp.Data.RefreshToken)
	}
}

func TestRefreshTokenResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.RefreshTokenResponse{
		Status: "success",
	}
	resp.Data.Token = "new_access_token_123"
	resp.Data.RefreshToken = "new_refresh_token_123"

	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", result["status"])
	}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if data["token"] != "new_access_token_123" {
			t.Errorf("Expected token 'new_access_token_123', got '%v'", data["token"])
		}
	} else {
		t.Error("Expected data to be a map")
	}
}

func TestGetProfileResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.GetProfileResponse{
		Status: "success",
	}
	resp.Data.UserID = "user-id-1"
	resp.Data.Username = "testuser"
	resp.Data.Email = "test@example.com"
	resp.Data.FullName = "Test User"
	resp.Data.RoleID = "role-id-1"
	resp.Data.Role = "Mahasiswa"
	resp.Data.Permissions = []string{"read:achievements", "write:achievements"}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.UserID != "user-id-1" {
		t.Errorf("Expected Data.UserID 'user-id-1', got '%s'", resp.Data.UserID)
	}
	if resp.Data.Username != "testuser" {
		t.Errorf("Expected Data.Username 'testuser', got '%s'", resp.Data.Username)
	}
	if resp.Data.Email != "test@example.com" {
		t.Errorf("Expected Data.Email 'test@example.com', got '%s'", resp.Data.Email)
	}
	if resp.Data.Role != "Mahasiswa" {
		t.Errorf("Expected Data.Role 'Mahasiswa', got '%s'", resp.Data.Role)
	}
	if len(resp.Data.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(resp.Data.Permissions))
	}
}

func TestGetProfileResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.GetProfileResponse{
		Status: "success",
	}
	resp.Data.UserID = "user-id-1"
	resp.Data.Username = "testuser"
	resp.Data.Email = "test@example.com"
	resp.Data.FullName = "Test User"
	resp.Data.RoleID = "role-id-1"
	resp.Data.Role = "Mahasiswa"
	resp.Data.Permissions = []string{"read:achievements"}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", result["status"])
	}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if data["user_id"] != "user-id-1" {
			t.Errorf("Expected user_id 'user-id-1', got '%v'", data["user_id"])
		}
		if data["username"] != "testuser" {
			t.Errorf("Expected username 'testuser', got '%v'", data["username"])
		}
		if data["role"] != "Mahasiswa" {
			t.Errorf("Expected role 'Mahasiswa', got '%v'", data["role"])
		}
	} else {
		t.Error("Expected data to be a map")
	}
}
