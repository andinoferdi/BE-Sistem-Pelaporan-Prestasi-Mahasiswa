package model_test

import (
	"encoding/json"
	"testing"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

func TestRolePermission_StructCreation(t *testing.T) {
	rolePermission := modelpostgre.RolePermission{
		RoleID:       "role-id-1",
		PermissionID: "permission-id-1",
	}

	if rolePermission.RoleID != "role-id-1" {
		t.Errorf("Expected RoleID 'role-id-1', got '%s'", rolePermission.RoleID)
	}
	if rolePermission.PermissionID != "permission-id-1" {
		t.Errorf("Expected PermissionID 'permission-id-1', got '%s'", rolePermission.PermissionID)
	}
}

func TestRolePermission_JSONMarshalling(t *testing.T) {
	rolePermission := modelpostgre.RolePermission{
		RoleID:       "role-id-1",
		PermissionID: "permission-id-1",
	}

	jsonData, err := json.Marshal(rolePermission)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["role_id"] != "role-id-1" {
		t.Errorf("Expected role_id 'role-id-1', got '%v'", result["role_id"])
	}
	if result["permission_id"] != "permission-id-1" {
		t.Errorf("Expected permission_id 'permission-id-1', got '%v'", result["permission_id"])
	}
}

func TestRolePermission_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"role_id": "role-id-1",
		"permission_id": "permission-id-1"
	}`

	var rolePermission modelpostgre.RolePermission
	if err := json.Unmarshal([]byte(jsonStr), &rolePermission); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rolePermission.RoleID != "role-id-1" {
		t.Errorf("Expected RoleID 'role-id-1', got '%s'", rolePermission.RoleID)
	}
	if rolePermission.PermissionID != "permission-id-1" {
		t.Errorf("Expected PermissionID 'permission-id-1', got '%s'", rolePermission.PermissionID)
	}
}

func TestRolePermission_ZeroValues(t *testing.T) {
	var rolePermission modelpostgre.RolePermission

	if rolePermission.RoleID != "" {
		t.Errorf("Expected empty RoleID, got '%s'", rolePermission.RoleID)
	}
	if rolePermission.PermissionID != "" {
		t.Errorf("Expected empty PermissionID, got '%s'", rolePermission.PermissionID)
	}
}

func TestCreateRolePermissionRequest_StructCreation(t *testing.T) {
	req := modelpostgre.CreateRolePermissionRequest{
		RoleID:       "role-id-1",
		PermissionID: "permission-id-1",
	}

	if req.RoleID != "role-id-1" {
		t.Errorf("Expected RoleID 'role-id-1', got '%s'", req.RoleID)
	}
	if req.PermissionID != "permission-id-1" {
		t.Errorf("Expected PermissionID 'permission-id-1', got '%s'", req.PermissionID)
	}
}

func TestCreateRolePermissionRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.CreateRolePermissionRequest{
		RoleID:       "role-id-1",
		PermissionID: "permission-id-1",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["role_id"] != "role-id-1" {
		t.Errorf("Expected role_id 'role-id-1', got '%v'", result["role_id"])
	}
	if result["permission_id"] != "permission-id-1" {
		t.Errorf("Expected permission_id 'permission-id-1', got '%v'", result["permission_id"])
	}
}

func TestCreateRolePermissionRequest_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"role_id": "role-id-1",
		"permission_id": "permission-id-1"
	}`

	var req modelpostgre.CreateRolePermissionRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.RoleID != "role-id-1" {
		t.Errorf("Expected RoleID 'role-id-1', got '%s'", req.RoleID)
	}
	if req.PermissionID != "permission-id-1" {
		t.Errorf("Expected PermissionID 'permission-id-1', got '%s'", req.PermissionID)
	}
}

func TestGetAllRolePermissionsResponse_StructCreation(t *testing.T) {
	rolePermissions := []modelpostgre.RolePermission{
		{
			RoleID:       "role-id-1",
			PermissionID: "permission-id-1",
		},
		{
			RoleID:       "role-id-1",
			PermissionID: "permission-id-2",
		},
	}

	resp := modelpostgre.GetAllRolePermissionsResponse{
		Status: "success",
		Data:   rolePermissions,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 role permissions, got %d", len(resp.Data))
	}
	if resp.Data[0].PermissionID != "permission-id-1" {
		t.Errorf("Expected first role permission PermissionID 'permission-id-1', got '%s'", resp.Data[0].PermissionID)
	}
}

func TestGetAllRolePermissionsResponse_JSONMarshalling(t *testing.T) {
	rolePermissions := []modelpostgre.RolePermission{
		{
			RoleID:       "role-id-1",
			PermissionID: "permission-id-1",
		},
	}

	resp := modelpostgre.GetAllRolePermissionsResponse{
		Status: "success",
		Data:   rolePermissions,
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
			t.Errorf("Expected 1 role permission in data, got %d", len(data))
		}
	} else {
		t.Error("Expected data to be an array")
	}
}

func TestCreateRolePermissionResponse_StructCreation(t *testing.T) {
	rolePermission := modelpostgre.RolePermission{
		RoleID:       "role-id-1",
		PermissionID: "permission-id-1",
	}

	resp := modelpostgre.CreateRolePermissionResponse{
		Status: "success",
		Data:   rolePermission,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.RoleID != "role-id-1" {
		t.Errorf("Expected Data.RoleID 'role-id-1', got '%s'", resp.Data.RoleID)
	}
}

func TestDeleteRolePermissionResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.DeleteRolePermissionResponse{
		Status: "success",
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
}

func TestDeleteRolePermissionResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.DeleteRolePermissionResponse{
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
