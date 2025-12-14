package model_test

import (
	"encoding/json"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

func TestUser_StructCreation(t *testing.T) {
	now := time.Now()
	user := modelpostgre.User{
		ID:           "user-id-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		RoleID:       "role-id-1",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if user.ID != "user-id-1" {
		t.Errorf("Expected ID 'user-id-1', got '%s'", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got '%s'", user.Email)
	}
	if user.FullName != "Test User" {
		t.Errorf("Expected FullName 'Test User', got '%s'", user.FullName)
	}
	if user.RoleID != "role-id-1" {
		t.Errorf("Expected RoleID 'role-id-1', got '%s'", user.RoleID)
	}
	if user.IsActive != true {
		t.Errorf("Expected IsActive true, got %v", user.IsActive)
	}
}

func TestUser_JSONMarshalling(t *testing.T) {
	now := time.Now()
	user := modelpostgre.User{
		ID:           "user-id-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		RoleID:       "role-id-1",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	jsonData, err := json.Marshal(user)
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
	if result["email"] != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%v'", result["email"])
	}
	if result["full_name"] != "Test User" {
		t.Errorf("Expected full_name 'Test User', got '%v'", result["full_name"])
	}
	if result["is_active"] != true {
		t.Errorf("Expected is_active true, got '%v'", result["is_active"])
	}
}

func TestUser_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"id": "user-id-1",
		"username": "testuser",
		"email": "test@example.com",
		"password_hash": "hashed_password",
		"full_name": "Test User",
		"role_id": "role-id-1",
		"is_active": true,
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-01-01T00:00:00Z"
	}`

	var user modelpostgre.User
	if err := json.Unmarshal([]byte(jsonStr), &user); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.ID != "user-id-1" {
		t.Errorf("Expected ID 'user-id-1', got '%s'", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got '%s'", user.Email)
	}
}

func TestUser_ZeroValues(t *testing.T) {
	var user modelpostgre.User

	if user.ID != "" {
		t.Errorf("Expected empty ID, got '%s'", user.ID)
	}
	if user.Username != "" {
		t.Errorf("Expected empty Username, got '%s'", user.Username)
	}
	if user.IsActive != false {
		t.Errorf("Expected IsActive false, got %v", user.IsActive)
	}
	if !user.CreatedAt.IsZero() {
		t.Errorf("Expected zero CreatedAt, got %v", user.CreatedAt)
	}
}

func TestCreateUserRequest_StructCreation(t *testing.T) {
	isActive := true
	req := modelpostgre.CreateUserRequest{
		Username:     "testuser",
		Email:        "test@example.com",
		Password:     "password123",
		FullName:     "Test User",
		RoleID:       "role-id-1",
		IsActive:     &isActive,
		StudentID:    "student-id-1",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		AdvisorID:    "advisor-id-1",
		LecturerID:   "lecturer-id-1",
		Department:   "IT",
	}

	if req.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", req.Username)
	}
	if req.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got '%s'", req.Email)
	}
	if req.Password != "password123" {
		t.Errorf("Expected Password 'password123', got '%s'", req.Password)
	}
	if req.IsActive == nil || *req.IsActive != true {
		t.Errorf("Expected IsActive true, got %v", req.IsActive)
	}
}

func TestCreateUserRequest_JSONMarshalling(t *testing.T) {
	isActive := true
	req := modelpostgre.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
		RoleID:   "role-id-1",
		IsActive: &isActive,
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
	if result["email"] != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%v'", result["email"])
	}
}

func TestCreateUserRequest_JSONMarshallingOmitempty(t *testing.T) {
	req := modelpostgre.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
		RoleID:   "role-id-1",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, exists := result["is_active"]; exists {
		t.Error("Expected is_active to be omitted when nil")
	}
}

func TestUpdateUserRequest_StructCreation(t *testing.T) {
	isActive := false
	req := modelpostgre.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
		FullName: "Updated User",
		RoleID:   "role-id-2",
		IsActive: &isActive,
	}

	if req.Username != "updateduser" {
		t.Errorf("Expected Username 'updateduser', got '%s'", req.Username)
	}
	if req.Email != "updated@example.com" {
		t.Errorf("Expected Email 'updated@example.com', got '%s'", req.Email)
	}
	if req.IsActive == nil || *req.IsActive != false {
		t.Errorf("Expected IsActive false, got %v", req.IsActive)
	}
}

func TestUpdateUserRequest_JSONMarshalling(t *testing.T) {
	isActive := false
	req := modelpostgre.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
		FullName: "Updated User",
		RoleID:   "role-id-2",
		IsActive: &isActive,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["username"] != "updateduser" {
		t.Errorf("Expected username 'updateduser', got '%v'", result["username"])
	}
	if result["is_active"] != false {
		t.Errorf("Expected is_active false, got '%v'", result["is_active"])
	}
}

func TestGetAllUsersResponse_StructCreation(t *testing.T) {
	now := time.Now()
	users := []modelpostgre.User{
		{
			ID:        "user-id-1",
			Username:  "user1",
			Email:     "user1@example.com",
			FullName:  "User One",
			RoleID:    "role-id-1",
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "user-id-2",
			Username:  "user2",
			Email:     "user2@example.com",
			FullName:  "User Two",
			RoleID:    "role-id-2",
			IsActive:  false,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	resp := modelpostgre.GetAllUsersResponse{
		Status: "success",
		Data:   users,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 users, got %d", len(resp.Data))
	}
	if resp.Data[0].Username != "user1" {
		t.Errorf("Expected first user username 'user1', got '%s'", resp.Data[0].Username)
	}
}

func TestGetAllUsersResponse_JSONMarshalling(t *testing.T) {
	now := time.Now()
	users := []modelpostgre.User{
		{
			ID:        "user-id-1",
			Username:  "user1",
			Email:     "user1@example.com",
			FullName:  "User One",
			RoleID:    "role-id-1",
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	resp := modelpostgre.GetAllUsersResponse{
		Status: "success",
		Data:   users,
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
	if data, ok := result["data"].([]interface{}); ok {
		if len(data) != 1 {
			t.Errorf("Expected 1 user in data, got %d", len(data))
		}
	} else {
		t.Error("Expected data to be an array")
	}
}

func TestGetUserByIDResponse_StructCreation(t *testing.T) {
	now := time.Now()
	user := modelpostgre.User{
		ID:        "user-id-1",
		Username:  "testuser",
		Email:     "test@example.com",
		FullName:  "Test User",
		RoleID:    "role-id-1",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := modelpostgre.GetUserByIDResponse{
		Status: "success",
		Data:   user,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.ID != "user-id-1" {
		t.Errorf("Expected Data.ID 'user-id-1', got '%s'", resp.Data.ID)
	}
}

func TestCreateUserResponse_StructCreation(t *testing.T) {
	now := time.Now()
	user := modelpostgre.User{
		ID:        "user-id-1",
		Username:  "testuser",
		Email:     "test@example.com",
		FullName:  "Test User",
		RoleID:    "role-id-1",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := modelpostgre.CreateUserResponse{
		Status: "success",
		Data:   user,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Username != "testuser" {
		t.Errorf("Expected Data.Username 'testuser', got '%s'", resp.Data.Username)
	}
}

func TestUpdateUserResponse_StructCreation(t *testing.T) {
	now := time.Now()
	user := modelpostgre.User{
		ID:        "user-id-1",
		Username:  "updateduser",
		Email:     "updated@example.com",
		FullName:  "Updated User",
		RoleID:    "role-id-2",
		IsActive:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := modelpostgre.UpdateUserResponse{
		Status: "success",
		Data:   user,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Username != "updateduser" {
		t.Errorf("Expected Data.Username 'updateduser', got '%s'", resp.Data.Username)
	}
}

func TestDeleteUserResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.DeleteUserResponse{
		Status: "success",
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
}

func TestDeleteUserResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.DeleteUserResponse{
		Status: "success",
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
}
