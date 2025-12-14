package model_test

import (
	"encoding/json"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

func TestRole_StructCreation(t *testing.T) {
	now := time.Now()
	role := modelpostgre.Role{
		ID:          "role-id-1",
		Name:        "Mahasiswa",
		Description: "Role for students",
		CreatedAt:   now,
	}

	if role.ID != "role-id-1" {
		t.Errorf("Expected ID 'role-id-1', got '%s'", role.ID)
	}
	if role.Name != "Mahasiswa" {
		t.Errorf("Expected Name 'Mahasiswa', got '%s'", role.Name)
	}
	if role.Description != "Role for students" {
		t.Errorf("Expected Description 'Role for students', got '%s'", role.Description)
	}
}

func TestRole_JSONMarshalling(t *testing.T) {
	now := time.Now()
	role := modelpostgre.Role{
		ID:          "role-id-1",
		Name:        "Mahasiswa",
		Description: "Role for students",
		CreatedAt:   now,
	}

	jsonData, err := json.Marshal(role)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != "role-id-1" {
		t.Errorf("Expected id 'role-id-1', got '%v'", result["id"])
	}
	if result["name"] != "Mahasiswa" {
		t.Errorf("Expected name 'Mahasiswa', got '%v'", result["name"])
	}
	if result["description"] != "Role for students" {
		t.Errorf("Expected description 'Role for students', got '%v'", result["description"])
	}
}

func TestRole_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"id": "role-id-1",
		"name": "Mahasiswa",
		"description": "Role for students",
		"created_at": "2024-01-01T00:00:00Z"
	}`

	var role modelpostgre.Role
	if err := json.Unmarshal([]byte(jsonStr), &role); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if role.ID != "role-id-1" {
		t.Errorf("Expected ID 'role-id-1', got '%s'", role.ID)
	}
	if role.Name != "Mahasiswa" {
		t.Errorf("Expected Name 'Mahasiswa', got '%s'", role.Name)
	}
}

func TestRole_ZeroValues(t *testing.T) {
	var role modelpostgre.Role

	if role.ID != "" {
		t.Errorf("Expected empty ID, got '%s'", role.ID)
	}
	if role.Name != "" {
		t.Errorf("Expected empty Name, got '%s'", role.Name)
	}
	if role.Description != "" {
		t.Errorf("Expected empty Description, got '%s'", role.Description)
	}
	if !role.CreatedAt.IsZero() {
		t.Errorf("Expected zero CreatedAt, got %v", role.CreatedAt)
	}
}

func TestCreateRoleRequest_StructCreation(t *testing.T) {
	req := modelpostgre.CreateRoleRequest{
		Name:        "Dosen Wali",
		Description: "Role for academic advisors",
	}

	if req.Name != "Dosen Wali" {
		t.Errorf("Expected Name 'Dosen Wali', got '%s'", req.Name)
	}
	if req.Description != "Role for academic advisors" {
		t.Errorf("Expected Description 'Role for academic advisors', got '%s'", req.Description)
	}
}

func TestCreateRoleRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.CreateRoleRequest{
		Name:        "Dosen Wali",
		Description: "Role for academic advisors",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["name"] != "Dosen Wali" {
		t.Errorf("Expected name 'Dosen Wali', got '%v'", result["name"])
	}
	if result["description"] != "Role for academic advisors" {
		t.Errorf("Expected description 'Role for academic advisors', got '%v'", result["description"])
	}
}

func TestCreateRoleRequest_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"name": "Dosen Wali",
		"description": "Role for academic advisors"
	}`

	var req modelpostgre.CreateRoleRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.Name != "Dosen Wali" {
		t.Errorf("Expected Name 'Dosen Wali', got '%s'", req.Name)
	}
	if req.Description != "Role for academic advisors" {
		t.Errorf("Expected Description 'Role for academic advisors', got '%s'", req.Description)
	}
}

func TestUpdateRoleRequest_StructCreation(t *testing.T) {
	req := modelpostgre.UpdateRoleRequest{
		Name:        "Admin",
		Description: "Administrator role",
	}

	if req.Name != "Admin" {
		t.Errorf("Expected Name 'Admin', got '%s'", req.Name)
	}
	if req.Description != "Administrator role" {
		t.Errorf("Expected Description 'Administrator role', got '%s'", req.Description)
	}
}

func TestUpdateRoleRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.UpdateRoleRequest{
		Name:        "Admin",
		Description: "Administrator role",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["name"] != "Admin" {
		t.Errorf("Expected name 'Admin', got '%v'", result["name"])
	}
	if result["description"] != "Administrator role" {
		t.Errorf("Expected description 'Administrator role', got '%v'", result["description"])
	}
}

func TestGetAllRolesResponse_StructCreation(t *testing.T) {
	now := time.Now()
	roles := []modelpostgre.Role{
		{
			ID:          "role-id-1",
			Name:        "Mahasiswa",
			Description: "Role for students",
			CreatedAt:   now,
		},
		{
			ID:          "role-id-2",
			Name:        "Dosen Wali",
			Description: "Role for academic advisors",
			CreatedAt:   now,
		},
	}

	resp := modelpostgre.GetAllRolesResponse{
		Status: "success",
		Data:   roles,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(resp.Data))
	}
	if resp.Data[0].Name != "Mahasiswa" {
		t.Errorf("Expected first role Name 'Mahasiswa', got '%s'", resp.Data[0].Name)
	}
}

func TestGetAllRolesResponse_JSONMarshalling(t *testing.T) {
	now := time.Now()
	roles := []modelpostgre.Role{
		{
			ID:          "role-id-1",
			Name:        "Mahasiswa",
			Description: "Role for students",
			CreatedAt:   now,
		},
	}

	resp := modelpostgre.GetAllRolesResponse{
		Status: "success",
		Data:   roles,
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
			t.Errorf("Expected 1 role in data, got %d", len(data))
		}
	} else {
		t.Error("Expected data to be an array")
	}
}

func TestGetRoleByIDResponse_StructCreation(t *testing.T) {
	now := time.Now()
	role := modelpostgre.Role{
		ID:          "role-id-1",
		Name:        "Mahasiswa",
		Description: "Role for students",
		CreatedAt:   now,
	}

	resp := modelpostgre.GetRoleByIDResponse{
		Status: "success",
		Data:   role,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.ID != "role-id-1" {
		t.Errorf("Expected Data.ID 'role-id-1', got '%s'", resp.Data.ID)
	}
}

func TestCreateRoleResponse_StructCreation(t *testing.T) {
	now := time.Now()
	role := modelpostgre.Role{
		ID:          "role-id-1",
		Name:        "Dosen Wali",
		Description: "Role for academic advisors",
		CreatedAt:   now,
	}

	resp := modelpostgre.CreateRoleResponse{
		Status: "success",
		Data:   role,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Name != "Dosen Wali" {
		t.Errorf("Expected Data.Name 'Dosen Wali', got '%s'", resp.Data.Name)
	}
}

func TestUpdateRoleResponse_StructCreation(t *testing.T) {
	now := time.Now()
	role := modelpostgre.Role{
		ID:          "role-id-1",
		Name:        "Admin",
		Description: "Administrator role",
		CreatedAt:   now,
	}

	resp := modelpostgre.UpdateRoleResponse{
		Status: "success",
		Data:   role,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Name != "Admin" {
		t.Errorf("Expected Data.Name 'Admin', got '%s'", resp.Data.Name)
	}
}

func TestDeleteRoleResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.DeleteRoleResponse{
		Status: "success",
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
}

func TestDeleteRoleResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.DeleteRoleResponse{
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
