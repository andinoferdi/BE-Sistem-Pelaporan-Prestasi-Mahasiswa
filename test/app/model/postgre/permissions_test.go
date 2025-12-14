package model_test

import (
	"encoding/json"
	"testing"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

func TestPermission_StructCreation(t *testing.T) {
	permission := modelpostgre.Permission{
		ID:          "permission-id-1",
		Name:        "read:achievements",
		Resource:    "achievements",
		Action:      "read",
		Description: "Permission to read achievements",
	}

	if permission.ID != "permission-id-1" {
		t.Errorf("Expected ID 'permission-id-1', got '%s'", permission.ID)
	}
	if permission.Name != "read:achievements" {
		t.Errorf("Expected Name 'read:achievements', got '%s'", permission.Name)
	}
	if permission.Resource != "achievements" {
		t.Errorf("Expected Resource 'achievements', got '%s'", permission.Resource)
	}
	if permission.Action != "read" {
		t.Errorf("Expected Action 'read', got '%s'", permission.Action)
	}
	if permission.Description != "Permission to read achievements" {
		t.Errorf("Expected Description 'Permission to read achievements', got '%s'", permission.Description)
	}
}

func TestPermission_JSONMarshalling(t *testing.T) {
	permission := modelpostgre.Permission{
		ID:          "permission-id-1",
		Name:        "read:achievements",
		Resource:    "achievements",
		Action:      "read",
		Description: "Permission to read achievements",
	}

	jsonData, err := json.Marshal(permission)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != "permission-id-1" {
		t.Errorf("Expected id 'permission-id-1', got '%v'", result["id"])
	}
	if result["name"] != "read:achievements" {
		t.Errorf("Expected name 'read:achievements', got '%v'", result["name"])
	}
	if result["resource"] != "achievements" {
		t.Errorf("Expected resource 'achievements', got '%v'", result["resource"])
	}
	if result["action"] != "read" {
		t.Errorf("Expected action 'read', got '%v'", result["action"])
	}
}

func TestPermission_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"id": "permission-id-1",
		"name": "read:achievements",
		"resource": "achievements",
		"action": "read",
		"description": "Permission to read achievements"
	}`

	var permission modelpostgre.Permission
	if err := json.Unmarshal([]byte(jsonStr), &permission); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if permission.ID != "permission-id-1" {
		t.Errorf("Expected ID 'permission-id-1', got '%s'", permission.ID)
	}
	if permission.Name != "read:achievements" {
		t.Errorf("Expected Name 'read:achievements', got '%s'", permission.Name)
	}
}

func TestPermission_ZeroValues(t *testing.T) {
	var permission modelpostgre.Permission

	if permission.ID != "" {
		t.Errorf("Expected empty ID, got '%s'", permission.ID)
	}
	if permission.Name != "" {
		t.Errorf("Expected empty Name, got '%s'", permission.Name)
	}
	if permission.Resource != "" {
		t.Errorf("Expected empty Resource, got '%s'", permission.Resource)
	}
	if permission.Action != "" {
		t.Errorf("Expected empty Action, got '%s'", permission.Action)
	}
}

func TestCreatePermissionRequest_StructCreation(t *testing.T) {
	req := modelpostgre.CreatePermissionRequest{
		Name:        "write:achievements",
		Resource:    "achievements",
		Action:      "write",
		Description: "Permission to write achievements",
	}

	if req.Name != "write:achievements" {
		t.Errorf("Expected Name 'write:achievements', got '%s'", req.Name)
	}
	if req.Resource != "achievements" {
		t.Errorf("Expected Resource 'achievements', got '%s'", req.Resource)
	}
	if req.Action != "write" {
		t.Errorf("Expected Action 'write', got '%s'", req.Action)
	}
}

func TestCreatePermissionRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.CreatePermissionRequest{
		Name:        "write:achievements",
		Resource:    "achievements",
		Action:      "write",
		Description: "Permission to write achievements",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["name"] != "write:achievements" {
		t.Errorf("Expected name 'write:achievements', got '%v'", result["name"])
	}
	if result["resource"] != "achievements" {
		t.Errorf("Expected resource 'achievements', got '%v'", result["resource"])
	}
	if result["action"] != "write" {
		t.Errorf("Expected action 'write', got '%v'", result["action"])
	}
}

func TestCreatePermissionRequest_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"name": "write:achievements",
		"resource": "achievements",
		"action": "write",
		"description": "Permission to write achievements"
	}`

	var req modelpostgre.CreatePermissionRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.Name != "write:achievements" {
		t.Errorf("Expected Name 'write:achievements', got '%s'", req.Name)
	}
	if req.Resource != "achievements" {
		t.Errorf("Expected Resource 'achievements', got '%s'", req.Resource)
	}
}

func TestUpdatePermissionRequest_StructCreation(t *testing.T) {
	req := modelpostgre.UpdatePermissionRequest{
		Name:        "delete:achievements",
		Resource:    "achievements",
		Action:      "delete",
		Description: "Permission to delete achievements",
	}

	if req.Name != "delete:achievements" {
		t.Errorf("Expected Name 'delete:achievements', got '%s'", req.Name)
	}
	if req.Resource != "achievements" {
		t.Errorf("Expected Resource 'achievements', got '%s'", req.Resource)
	}
	if req.Action != "delete" {
		t.Errorf("Expected Action 'delete', got '%s'", req.Action)
	}
}

func TestUpdatePermissionRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.UpdatePermissionRequest{
		Name:        "delete:achievements",
		Resource:    "achievements",
		Action:      "delete",
		Description: "Permission to delete achievements",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["name"] != "delete:achievements" {
		t.Errorf("Expected name 'delete:achievements', got '%v'", result["name"])
	}
	if result["action"] != "delete" {
		t.Errorf("Expected action 'delete', got '%v'", result["action"])
	}
}

func TestGetAllPermissionsResponse_StructCreation(t *testing.T) {
	permissions := []modelpostgre.Permission{
		{
			ID:          "permission-id-1",
			Name:        "read:achievements",
			Resource:    "achievements",
			Action:      "read",
			Description: "Permission to read achievements",
		},
		{
			ID:          "permission-id-2",
			Name:        "write:achievements",
			Resource:    "achievements",
			Action:      "write",
			Description: "Permission to write achievements",
		},
	}

	resp := modelpostgre.GetAllPermissionsResponse{
		Status: "success",
		Data:   permissions,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(resp.Data))
	}
	if resp.Data[0].Name != "read:achievements" {
		t.Errorf("Expected first permission Name 'read:achievements', got '%s'", resp.Data[0].Name)
	}
}

func TestGetAllPermissionsResponse_JSONMarshalling(t *testing.T) {
	permissions := []modelpostgre.Permission{
		{
			ID:          "permission-id-1",
			Name:        "read:achievements",
			Resource:    "achievements",
			Action:      "read",
			Description: "Permission to read achievements",
		},
	}

	resp := modelpostgre.GetAllPermissionsResponse{
		Status: "success",
		Data:   permissions,
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
			t.Errorf("Expected 1 permission in data, got %d", len(data))
		}
	} else {
		t.Error("Expected data to be an array")
	}
}

func TestGetPermissionByIDResponse_StructCreation(t *testing.T) {
	permission := modelpostgre.Permission{
		ID:          "permission-id-1",
		Name:        "read:achievements",
		Resource:    "achievements",
		Action:      "read",
		Description: "Permission to read achievements",
	}

	resp := modelpostgre.GetPermissionByIDResponse{
		Status: "success",
		Data:   permission,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.ID != "permission-id-1" {
		t.Errorf("Expected Data.ID 'permission-id-1', got '%s'", resp.Data.ID)
	}
}

func TestCreatePermissionResponse_StructCreation(t *testing.T) {
	permission := modelpostgre.Permission{
		ID:          "permission-id-1",
		Name:        "write:achievements",
		Resource:    "achievements",
		Action:      "write",
		Description: "Permission to write achievements",
	}

	resp := modelpostgre.CreatePermissionResponse{
		Status: "success",
		Data:   permission,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Name != "write:achievements" {
		t.Errorf("Expected Data.Name 'write:achievements', got '%s'", resp.Data.Name)
	}
}

func TestUpdatePermissionResponse_StructCreation(t *testing.T) {
	permission := modelpostgre.Permission{
		ID:          "permission-id-1",
		Name:        "delete:achievements",
		Resource:    "achievements",
		Action:      "delete",
		Description: "Permission to delete achievements",
	}

	resp := modelpostgre.UpdatePermissionResponse{
		Status: "success",
		Data:   permission,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Action != "delete" {
		t.Errorf("Expected Data.Action 'delete', got '%s'", resp.Data.Action)
	}
}

func TestDeletePermissionResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.DeletePermissionResponse{
		Status: "success",
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
}

func TestDeletePermissionResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.DeletePermissionResponse{
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
