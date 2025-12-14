package model_test

import (
	"encoding/json"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

func TestStudent_StructCreation(t *testing.T) {
	now := time.Now()
	student := modelpostgre.Student{
		ID:           "student-id-1",
		UserID:       "user-id-1",
		StudentID:    "2024001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		AdvisorID:    "advisor-id-1",
		FullName:     "Student Name",
		CreatedAt:    now,
	}

	if student.ID != "student-id-1" {
		t.Errorf("Expected ID 'student-id-1', got '%s'", student.ID)
	}
	if student.UserID != "user-id-1" {
		t.Errorf("Expected UserID 'user-id-1', got '%s'", student.UserID)
	}
	if student.StudentID != "2024001" {
		t.Errorf("Expected StudentID '2024001', got '%s'", student.StudentID)
	}
	if student.ProgramStudy != "Computer Science" {
		t.Errorf("Expected ProgramStudy 'Computer Science', got '%s'", student.ProgramStudy)
	}
	if student.AcademicYear != "2024" {
		t.Errorf("Expected AcademicYear '2024', got '%s'", student.AcademicYear)
	}
	if student.AdvisorID != "advisor-id-1" {
		t.Errorf("Expected AdvisorID 'advisor-id-1', got '%s'", student.AdvisorID)
	}
}

func TestStudent_JSONMarshalling(t *testing.T) {
	now := time.Now()
	student := modelpostgre.Student{
		ID:           "student-id-1",
		UserID:       "user-id-1",
		StudentID:    "2024001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		AdvisorID:    "advisor-id-1",
		FullName:     "Student Name",
		CreatedAt:    now,
	}

	jsonData, err := json.Marshal(student)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != "student-id-1" {
		t.Errorf("Expected id 'student-id-1', got '%v'", result["id"])
	}
	if result["user_id"] != "user-id-1" {
		t.Errorf("Expected user_id 'user-id-1', got '%v'", result["user_id"])
	}
	if result["student_id"] != "2024001" {
		t.Errorf("Expected student_id '2024001', got '%v'", result["student_id"])
	}
	if result["program_study"] != "Computer Science" {
		t.Errorf("Expected program_study 'Computer Science', got '%v'", result["program_study"])
	}
}

func TestStudent_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"id": "student-id-1",
		"user_id": "user-id-1",
		"student_id": "2024001",
		"program_study": "Computer Science",
		"academic_year": "2024",
		"advisor_id": "advisor-id-1",
		"full_name": "Student Name",
		"created_at": "2024-01-01T00:00:00Z"
	}`

	var student modelpostgre.Student
	if err := json.Unmarshal([]byte(jsonStr), &student); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if student.ID != "student-id-1" {
		t.Errorf("Expected ID 'student-id-1', got '%s'", student.ID)
	}
	if student.StudentID != "2024001" {
		t.Errorf("Expected StudentID '2024001', got '%s'", student.StudentID)
	}
}

func TestStudent_ZeroValues(t *testing.T) {
	var student modelpostgre.Student

	if student.ID != "" {
		t.Errorf("Expected empty ID, got '%s'", student.ID)
	}
	if student.UserID != "" {
		t.Errorf("Expected empty UserID, got '%s'", student.UserID)
	}
	if student.StudentID != "" {
		t.Errorf("Expected empty StudentID, got '%s'", student.StudentID)
	}
	if !student.CreatedAt.IsZero() {
		t.Errorf("Expected zero CreatedAt, got %v", student.CreatedAt)
	}
}

func TestStudent_JSONMarshallingOmitempty(t *testing.T) {
	now := time.Now()
	student := modelpostgre.Student{
		ID:           "student-id-1",
		UserID:       "user-id-1",
		StudentID:    "2024001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		AdvisorID:    "advisor-id-1",
		CreatedAt:    now,
	}

	jsonData, err := json.Marshal(student)
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

func TestCreateStudentRequest_StructCreation(t *testing.T) {
	req := modelpostgre.CreateStudentRequest{
		UserID:       "user-id-1",
		StudentID:    "2024001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		AdvisorID:    "advisor-id-1",
	}

	if req.UserID != "user-id-1" {
		t.Errorf("Expected UserID 'user-id-1', got '%s'", req.UserID)
	}
	if req.StudentID != "2024001" {
		t.Errorf("Expected StudentID '2024001', got '%s'", req.StudentID)
	}
	if req.ProgramStudy != "Computer Science" {
		t.Errorf("Expected ProgramStudy 'Computer Science', got '%s'", req.ProgramStudy)
	}
	if req.AcademicYear != "2024" {
		t.Errorf("Expected AcademicYear '2024', got '%s'", req.AcademicYear)
	}
}

func TestCreateStudentRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.CreateStudentRequest{
		UserID:       "user-id-1",
		StudentID:    "2024001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		AdvisorID:    "advisor-id-1",
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
	if result["student_id"] != "2024001" {
		t.Errorf("Expected student_id '2024001', got '%v'", result["student_id"])
	}
}

func TestCreateStudentRequest_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"user_id": "user-id-1",
		"student_id": "2024001",
		"program_study": "Computer Science",
		"academic_year": "2024",
		"advisor_id": "advisor-id-1"
	}`

	var req modelpostgre.CreateStudentRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.UserID != "user-id-1" {
		t.Errorf("Expected UserID 'user-id-1', got '%s'", req.UserID)
	}
	if req.StudentID != "2024001" {
		t.Errorf("Expected StudentID '2024001', got '%s'", req.StudentID)
	}
}

func TestUpdateStudentRequest_StructCreation(t *testing.T) {
	req := modelpostgre.UpdateStudentRequest{
		StudentID:    "2024001",
		ProgramStudy: "Information Technology",
		AcademicYear: "2025",
		AdvisorID:    "advisor-id-2",
	}

	if req.StudentID != "2024001" {
		t.Errorf("Expected StudentID '2024001', got '%s'", req.StudentID)
	}
	if req.ProgramStudy != "Information Technology" {
		t.Errorf("Expected ProgramStudy 'Information Technology', got '%s'", req.ProgramStudy)
	}
	if req.AcademicYear != "2025" {
		t.Errorf("Expected AcademicYear '2025', got '%s'", req.AcademicYear)
	}
}

func TestUpdateStudentRequest_JSONMarshalling(t *testing.T) {
	req := modelpostgre.UpdateStudentRequest{
		StudentID:    "2024001",
		ProgramStudy: "Information Technology",
		AcademicYear: "2025",
		AdvisorID:    "advisor-id-2",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["student_id"] != "2024001" {
		t.Errorf("Expected student_id '2024001', got '%v'", result["student_id"])
	}
	if result["program_study"] != "Information Technology" {
		t.Errorf("Expected program_study 'Information Technology', got '%v'", result["program_study"])
	}
}

func TestGetAllStudentsResponse_StructCreation(t *testing.T) {
	now := time.Now()
	students := []modelpostgre.Student{
		{
			ID:           "student-id-1",
			UserID:       "user-id-1",
			StudentID:    "2024001",
			ProgramStudy: "Computer Science",
			AcademicYear: "2024",
			CreatedAt:    now,
		},
		{
			ID:           "student-id-2",
			UserID:       "user-id-2",
			StudentID:    "2024002",
			ProgramStudy: "Information Technology",
			AcademicYear: "2024",
			CreatedAt:    now,
		},
	}

	resp := modelpostgre.GetAllStudentsResponse{
		Status: "success",
		Data:   students,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 students, got %d", len(resp.Data))
	}
	if resp.Data[0].StudentID != "2024001" {
		t.Errorf("Expected first student StudentID '2024001', got '%s'", resp.Data[0].StudentID)
	}
}

func TestGetAllStudentsResponse_JSONMarshalling(t *testing.T) {
	now := time.Now()
	students := []modelpostgre.Student{
		{
			ID:           "student-id-1",
			UserID:       "user-id-1",
			StudentID:    "2024001",
			ProgramStudy: "Computer Science",
			AcademicYear: "2024",
			CreatedAt:    now,
		},
	}

	resp := modelpostgre.GetAllStudentsResponse{
		Status: "success",
		Data:   students,
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
			t.Errorf("Expected 1 student in data, got %d", len(data))
		}
	} else {
		t.Error("Expected data to be an array")
	}
}

func TestGetStudentByIDResponse_StructCreation(t *testing.T) {
	now := time.Now()
	student := modelpostgre.Student{
		ID:           "student-id-1",
		UserID:       "user-id-1",
		StudentID:    "2024001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		CreatedAt:    now,
	}

	resp := modelpostgre.GetStudentByIDResponse{
		Status: "success",
		Data:   student,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.ID != "student-id-1" {
		t.Errorf("Expected Data.ID 'student-id-1', got '%s'", resp.Data.ID)
	}
}

func TestCreateStudentResponse_StructCreation(t *testing.T) {
	now := time.Now()
	student := modelpostgre.Student{
		ID:           "student-id-1",
		UserID:       "user-id-1",
		StudentID:    "2024001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		CreatedAt:    now,
	}

	resp := modelpostgre.CreateStudentResponse{
		Status: "success",
		Data:   student,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.StudentID != "2024001" {
		t.Errorf("Expected Data.StudentID '2024001', got '%s'", resp.Data.StudentID)
	}
}

func TestUpdateStudentResponse_StructCreation(t *testing.T) {
	now := time.Now()
	student := modelpostgre.Student{
		ID:           "student-id-1",
		UserID:       "user-id-1",
		StudentID:    "2024001",
		ProgramStudy: "Information Technology",
		AcademicYear: "2025",
		CreatedAt:    now,
	}

	resp := modelpostgre.UpdateStudentResponse{
		Status: "success",
		Data:   student,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.ProgramStudy != "Information Technology" {
		t.Errorf("Expected Data.ProgramStudy 'Information Technology', got '%s'", resp.Data.ProgramStudy)
	}
}

func TestDeleteStudentResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.DeleteStudentResponse{
		Status: "success",
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
}

func TestDeleteStudentResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.DeleteStudentResponse{
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
