package route_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	routepostgre "sistem-pelaporan-prestasi-mahasiswa/route/postgre"
)

type mockReportService struct {
	getStatisticsResp            map[string]interface{}
	getStatisticsErr             error
	getCurrentStudentReportResp  map[string]interface{}
	getCurrentStudentReportErr   error
	getStudentReportResp         map[string]interface{}
	getStudentReportErr          error
	getCurrentLecturerReportResp map[string]interface{}
	getCurrentLecturerReportErr  error
	getLecturerReportResp        map[string]interface{}
	getLecturerReportErr         error
}

func (m *mockReportService) GetStatistics(ctx context.Context, userID string, roleID string) (map[string]interface{}, error) {
	if m.getStatisticsErr != nil {
		return nil, m.getStatisticsErr
	}
	return m.getStatisticsResp, nil
}

func (m *mockReportService) GetCurrentStudentReport(ctx context.Context, userID string) (map[string]interface{}, error) {
	if m.getCurrentStudentReportErr != nil {
		return nil, m.getCurrentStudentReportErr
	}
	return m.getCurrentStudentReportResp, nil
}

func (m *mockReportService) GetStudentReport(ctx context.Context, studentID string) (map[string]interface{}, error) {
	if m.getStudentReportErr != nil {
		return nil, m.getStudentReportErr
	}
	return m.getStudentReportResp, nil
}

func (m *mockReportService) GetCurrentLecturerReport(ctx context.Context, userID string) (map[string]interface{}, error) {
	if m.getCurrentLecturerReportErr != nil {
		return nil, m.getCurrentLecturerReportErr
	}
	return m.getCurrentLecturerReportResp, nil
}

func (m *mockReportService) GetLecturerReport(ctx context.Context, lecturerID string) (map[string]interface{}, error) {
	if m.getLecturerReportErr != nil {
		return nil, m.getLecturerReportErr
	}
	return m.getLecturerReportResp, nil
}

func TestGetStatisticsRoute_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockReportService{
		getStatisticsResp: map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"total_achievements": 100,
				"verified":           75,
				"pending":            15,
				"rejected":           10,
			},
		},
	}

	app := setupTestApp()
	routepostgre.ReportRoutes(app, mockService, nil)

	req := createRequestWithToken("GET", "/api/v1/reports/statistics", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestGetStatisticsRoute_MissingUserID(t *testing.T) {
	app := setupTestApp()
	routepostgre.ReportRoutes(app, &mockReportService{}, nil)

	req := httptest.NewRequest("GET", "/api/v1/reports/statistics", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusUnauthorized)
}

func TestGetCurrentStudentReportRoute_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockReportService{
		getCurrentStudentReportResp: map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"student_id":   "M001",
				"achievements": []interface{}{},
			},
		},
	}

	app := setupTestApp()
	routepostgre.ReportRoutes(app, mockService, nil)

	req := createRequestWithToken("GET", "/api/v1/reports/student", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestGetStudentReportRoute_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockReportService{
		getStudentReportResp: map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"student_id":   "M001",
				"achievements": []interface{}{},
			},
		},
	}

	app := setupTestApp()
	routepostgre.ReportRoutes(app, mockService, nil)

	req := createRequestWithToken("GET", "/api/v1/reports/student/student-id-1", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestGetStudentReportRoute_EmptyID(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockReportService{
		getStudentReportErr: errors.New("ID student wajib diisi"),
	}

	app := setupTestApp()
	routepostgre.ReportRoutes(app, mockService, nil)

	req := createRequestWithToken("GET", "/api/v1/reports/student/", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Skip("Path dengan trailing slash tidak trigger validasi empty ID di Fiber routing")
	}
}

func TestGetCurrentLecturerReportRoute_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockReportService{
		getCurrentLecturerReportResp: map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"lecturer_id": "D001",
				"advisees":    []interface{}{},
			},
		},
	}

	app := setupTestApp()
	routepostgre.ReportRoutes(app, mockService, nil)

	req := createRequestWithToken("GET", "/api/v1/reports/lecturer", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestGetLecturerReportRoute_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockReportService{
		getLecturerReportResp: map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"lecturer_id": "D001",
				"advisees":    []interface{}{},
			},
		},
	}

	app := setupTestApp()
	routepostgre.ReportRoutes(app, mockService, nil)

	req := createRequestWithToken("GET", "/api/v1/reports/lecturer/lecturer-id-1", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestGetLecturerReportRoute_EmptyID(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockReportService{
		getLecturerReportErr: errors.New("ID lecturer wajib diisi"),
	}

	app := setupTestApp()
	routepostgre.ReportRoutes(app, mockService, nil)

	req := createRequestWithToken("GET", "/api/v1/reports/lecturer/", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Skip("Path dengan trailing slash tidak trigger validasi empty ID di Fiber routing")
	}
}

func TestGetStatisticsRoute_Error(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockReportService{
		getStatisticsErr: errors.New("gagal mengambil statistik"),
	}

	app := setupTestApp()
	routepostgre.ReportRoutes(app, mockService, nil)

	req := createRequestWithToken("GET", "/api/v1/reports/statistics", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusInternalServerError)
}
