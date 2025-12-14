package model_test

import (
	"encoding/json"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

func TestLecturer_StructCreation(t *testing.T) {
	now := time.Now()
	lecturer := modelpostgre.Lecturer{
		ID:         "lecturer-id-1",
		UserID:     "user-id-1",
		LecturerID: "LEC001",
		Department: "Computer Science",
		FullName:   "Lecturer Name",
		CreatedAt:  now,
	}

	if lecturer.ID != "lecturer-id-1" {
		t.Errorf("Expected ID 'lecturer-id-1', got '%s'", lecturer.ID)
	}
	if lecturer.UserID != "user-id-1" {
		t.Errorf("Expected UserID 'user-id-1', got '%s'", lecturer.UserID)
	}
	if lecturer.LecturerID != "LEC001" {
		t.Errorf("Expected LecturerID 'LEC001', got '%s'", lecturer.LecturerID)
	}
	if lecturer.Department != "Computer Science" {
		t.Errorf("Expected Department 'Computer Science', got '%s'", lecturer.Department)
	}
}

func TestLecturer_JSONMarshalling(t *testing.T) {
	now := time.Now()
	lecturer := modelpostgre.Lecturer{
		ID:         "lecturer-id-1",
		UserID:     "user-id-1",
		LecturerID: "LEC001",
		Department: "Computer Science",
		FullName:   "Lecturer Name",
		CreatedAt:  now,
	}

	jsonData, err := json.Marshal(lecturer)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != "lecturer-id-1" {
		t.Errorf("Expected id 'lecturer-id-1', got '%v'", result["id"])
	}
	if result["user_id"] != "user-id-1" {
		t.Errorf("Expected user_id 'user-id-1', got '%v'", result["user_id"])
	}
	if result["lecturer_id"] != "LEC001" {
		t.Errorf("Expected lecturer_id 'LEC001', got '%v'", result["lecturer_id"])
	}
	if result["department"] != "Computer Science" {
		t.Errorf("Expected department 'Computer Science', got '%v'", result["department"])
	}
}

func TestLecturer_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"id": "lecturer-id-1",
		"user_id": "user-id-1",
		"lecturer_id": "LEC001",
		"department": "Computer Science",
		"full_name": "Lecturer Name",
		"created_at": "2024-01-01T00:00:00Z"
	}`

	var lecturer modelpostgre.Lecturer
	if err := json.Unmarshal([]byte(jsonStr), &lecturer); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if lecturer.ID != "lecturer-id-1" {
		t.Errorf("Expected ID 'lecturer-id-1', got '%s'", lecturer.ID)
	}
	if lecturer.LecturerID != "LEC001" {
		t.Errorf("Expected LecturerID 'LEC001', got '%s'", lecturer.LecturerID)
	}
}

func TestLecturer_ZeroValues(t *testing.T) {
	var lecturer modelpostgre.Lecturer

	if lecturer.ID != "" {
		t.Errorf("Expected empty ID, got '%s'", lecturer.ID)
	}
	if lecturer.UserID != "" {
		t.Errorf("Expected empty UserID, got '%s'", lecturer.UserID)
	}
	if lecturer.LecturerID != "" {
		t.Errorf("Expected empty LecturerID, got '%s'", lecturer.LecturerID)
	}
	if !lecturer.CreatedAt.IsZero() {
		t.Errorf("Expected zero CreatedAt, got %v", lecturer.CreatedAt)
	}
}

func TestLecturer_JSONMarshallingOmitempty(t *testing.T) {
	now := time.Now()
	lecturer := modelpostgre.Lecturer{
		ID:         "lecturer-id-1",
		UserID:     "user-id-1",
		LecturerID: "LEC001",
		Department: "Computer Science",
		CreatedAt:  now,
	}

	jsonData, err := json.Marshal(lecturer)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, exists := result["full_name"]; exists {
		t.Error("Expected full_name to be omitted when empty")
	}
}

func TestCreateLecturerRequest_StructCreation(t *testing.T) {
	req := modelpostgre.CreateLecturerRequest{
		UserID:     "user-id-1",
		LecturerID: "LEC001",
		Department: "Computer Science",
	}

	if req.UserID != "user-id-1" {
		t.Errorf("Expected UserID 'user-id-1', got '%s'", req.UserID)
	}
	if req.LecturerID != "LEC001" {
		t.Errorf("Expected LecturerID 'LEC001', got '%s'", req.LecturerID)
	}
	if req.Department != "Computer Science" {
		t.Errorf("Expected Department 'Computer Science', got '%s'", req.Department)
	}
}

func TestCreateLecturerRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.CreateLecturerRequest{
		UserID:     "user-id-1",
		LecturerID: "LEC001",
		Department: "Computer Science",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["user_id"] != "user-id-1" {
		t.Errorf("Expected user_id 'user-id-1', got '%v'", result["user_id"])
	}
	if result["lecturer_id"] != "LEC001" {
		t.Errorf("Expected lecturer_id 'LEC001', got '%v'", result["lecturer_id"])
	}
}

func TestCreateLecturerRequest_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"user_id": "user-id-1",
		"lecturer_id": "LEC001",
		"department": "Computer Science"
	}`

	var req modelpostgre.CreateLecturerRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.UserID != "user-id-1" {
		t.Errorf("Expected UserID 'user-id-1', got '%s'", req.UserID)
	}
	if req.LecturerID != "LEC001" {
		t.Errorf("Expected LecturerID 'LEC001', got '%s'", req.LecturerID)
	}
}

func TestUpdateLecturerRequest_StructCreation(t *testing.T) {
	req := modelpostgre.UpdateLecturerRequest{
		LecturerID: "LEC002",
		Department: "Information Technology",
	}

	if req.LecturerID != "LEC002" {
		t.Errorf("Expected LecturerID 'LEC002', got '%s'", req.LecturerID)
	}
	if req.Department != "Information Technology" {
		t.Errorf("Expected Department 'Information Technology', got '%s'", req.Department)
	}
}

func TestUpdateLecturerRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.UpdateLecturerRequest{
		LecturerID: "LEC002",
		Department: "Information Technology",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["lecturer_id"] != "LEC002" {
		t.Errorf("Expected lecturer_id 'LEC002', got '%v'", result["lecturer_id"])
	}
	if result["department"] != "Information Technology" {
		t.Errorf("Expected department 'Information Technology', got '%v'", result["department"])
	}
}

func TestGetAllLecturersResponse_StructCreation(t *testing.T) {
	now := time.Now()
	lecturers := []modelpostgre.Lecturer{
		{
			ID:         "lecturer-id-1",
			UserID:     "user-id-1",
			LecturerID: "LEC001",
			Department: "Computer Science",
			CreatedAt:  now,
		},
		{
			ID:         "lecturer-id-2",
			UserID:     "user-id-2",
			LecturerID: "LEC002",
			Department: "Information Technology",
			CreatedAt:  now,
		},
	}

	resp := modelpostgre.GetAllLecturersResponse{
		Status: "success",
		Data:   lecturers,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 lecturers, got %d", len(resp.Data))
	}
	if resp.Data[0].LecturerID != "LEC001" {
		t.Errorf("Expected first lecturer LecturerID 'LEC001', got '%s'", resp.Data[0].LecturerID)
	}
}

func TestGetAllLecturersResponse_JSONMarshalling(t *testing.T) {
	now := time.Now()
	lecturers := []modelpostgre.Lecturer{
		{
			ID:         "lecturer-id-1",
			UserID:     "user-id-1",
			LecturerID: "LEC001",
			Department: "Computer Science",
			CreatedAt:  now,
		},
	}

	resp := modelpostgre.GetAllLecturersResponse{
		Status: "success",
		Data:   lecturers,
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
			t.Errorf("Expected 1 lecturer in data, got %d", len(data))
		}
	} else {
		t.Error("Expected data to be an array")
	}
}

func TestGetLecturerByIDResponse_StructCreation(t *testing.T) {
	now := time.Now()
	lecturer := modelpostgre.Lecturer{
		ID:         "lecturer-id-1",
		UserID:     "user-id-1",
		LecturerID: "LEC001",
		Department: "Computer Science",
		CreatedAt:  now,
	}

	resp := modelpostgre.GetLecturerByIDResponse{
		Status: "success",
		Data:   lecturer,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.ID != "lecturer-id-1" {
		t.Errorf("Expected Data.ID 'lecturer-id-1', got '%s'", resp.Data.ID)
	}
}

func TestCreateLecturerResponse_StructCreation(t *testing.T) {
	now := time.Now()
	lecturer := modelpostgre.Lecturer{
		ID:         "lecturer-id-1",
		UserID:     "user-id-1",
		LecturerID: "LEC001",
		Department: "Computer Science",
		CreatedAt:  now,
	}

	resp := modelpostgre.CreateLecturerResponse{
		Status: "success",
		Data:   lecturer,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.LecturerID != "LEC001" {
		t.Errorf("Expected Data.LecturerID 'LEC001', got '%s'", resp.Data.LecturerID)
	}
}

func TestUpdateLecturerResponse_StructCreation(t *testing.T) {
	now := time.Now()
	lecturer := modelpostgre.Lecturer{
		ID:         "lecturer-id-1",
		UserID:     "user-id-1",
		LecturerID: "LEC002",
		Department: "Information Technology",
		CreatedAt:  now,
	}

	resp := modelpostgre.UpdateLecturerResponse{
		Status: "success",
		Data:   lecturer,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Department != "Information Technology" {
		t.Errorf("Expected Data.Department 'Information Technology', got '%s'", resp.Data.Department)
	}
}

func TestDeleteLecturerResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.DeleteLecturerResponse{
		Status: "success",
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
}

func TestDeleteLecturerResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.DeleteLecturerResponse{
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
