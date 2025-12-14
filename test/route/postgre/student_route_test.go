package route_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	modelmongo "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	routepostgre "sistem-pelaporan-prestasi-mahasiswa/route/postgre"

	"github.com/DATA-DOG/go-sqlmock"
)

type mockStudentServiceForRoute struct {
	allStudents         []modelpostgre.Student
	studentByID         *modelpostgre.Student
	studentByUserID     *modelpostgre.Student
	studentsByAdvisorID []modelpostgre.Student
	studentIDByUserID   string
	err                 error
	updateAdvisorErr    error
}

func (m *mockStudentServiceForRoute) GetAllStudents(ctx context.Context) ([]modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allStudents, nil
}

func (m *mockStudentServiceForRoute) GetStudentByID(ctx context.Context, id string) (*modelpostgre.Student, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, errors.New("student dengan ID tersebut tidak ditemukan")
		}
		return nil, m.err
	}
	return m.studentByID, nil
}

func (m *mockStudentServiceForRoute) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.studentIDByUserID, nil
}

func (m *mockStudentServiceForRoute) GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.studentByUserID, nil
}

func (m *mockStudentServiceForRoute) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.studentsByAdvisorID, nil
}

func (m *mockStudentServiceForRoute) CreateStudent(ctx context.Context, req modelpostgre.CreateStudentRequest) (*modelpostgre.Student, error) {
	return nil, errors.New("not implemented in route test")
}

func (m *mockStudentServiceForRoute) UpdateStudentAdvisor(ctx context.Context, studentID string, advisorID string) error {
	if m.updateAdvisorErr != nil {
		if m.updateAdvisorErr == sql.ErrNoRows {
			return errors.New("student dengan ID tersebut tidak ditemukan")
		}
		return m.updateAdvisorErr
	}
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return errors.New("student dengan ID tersebut tidak ditemukan")
		}
		return m.err
	}
	return nil
}

type mockAchievementServiceForRoute struct {
	getByStudentIDResp map[string]interface{}
	getByStudentIDErr  error
}

func (m *mockAchievementServiceForRoute) CreateAchievement(ctx context.Context, userID string, roleID string, req modelmongo.CreateAchievementRequest) (*modelmongo.CreateAchievementResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAchievementServiceForRoute) SubmitAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.UpdateAchievementReferenceResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAchievementServiceForRoute) VerifyAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelpostgre.VerifyAchievementResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAchievementServiceForRoute) RejectAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelpostgre.RejectAchievementRequest) (*modelpostgre.RejectAchievementResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAchievementServiceForRoute) DeleteAchievement(ctx context.Context, userID string, roleID string, mongoID string) (*modelmongo.DeleteAchievementResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAchievementServiceForRoute) GetAchievements(ctx context.Context, userID string, roleID string, page, limit int, statusFilter string, achievementTypeFilter string, sortBy string, sortOrder string) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAchievementServiceForRoute) GetAchievementsByStudentID(ctx context.Context, studentID string, page, limit int) (map[string]interface{}, error) {
	if m.getByStudentIDErr != nil {
		return nil, m.getByStudentIDErr
	}
	return m.getByStudentIDResp, nil
}

func (m *mockAchievementServiceForRoute) GetAchievementByID(ctx context.Context, userID string, roleID string, mongoID string) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAchievementServiceForRoute) UpdateAchievement(ctx context.Context, userID string, roleID string, mongoID string, req modelmongo.UpdateAchievementRequest) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAchievementServiceForRoute) GetAchievementStats(ctx context.Context) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAchievementServiceForRoute) UploadFile(ctx context.Context, userID string, roleID string, mongoID string, fileName string, fileURL string, fileType string) (*modelmongo.Attachment, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAchievementServiceForRoute) GetAchievementHistory(ctx context.Context, userID string, roleID string, mongoID string) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func TestGetAllStudentsRoute_Success(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "user:manage").
		WillReturnRows(rows)

	mockStudentService := &mockStudentServiceForRoute{
		allStudents: []modelpostgre.Student{
			{
				ID:           "student-id-1",
				UserID:       "user-id-1",
				StudentID:    "M001",
				ProgramStudy: "Teknik Informatika",
				AcademicYear: "2024",
				CreatedAt:    time.Now(),
			},
		},
	}

	app := setupTestApp()
	routepostgre.StudentRoutes(app, mockStudentService, &mockAchievementServiceForRoute{}, db)

	req := createRequestWithToken("GET", "/api/v1/students", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetStudentByIDRoute_Success(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "user:manage").
		WillReturnRows(rows)

	mockStudentService := &mockStudentServiceForRoute{
		studentByID: &modelpostgre.Student{
			ID:           "student-id-1",
			UserID:       "user-id-1",
			StudentID:    "M001",
			ProgramStudy: "Teknik Informatika",
			AcademicYear: "2024",
			CreatedAt:    time.Now(),
		},
	}

	app := setupTestApp()
	routepostgre.StudentRoutes(app, mockStudentService, &mockAchievementServiceForRoute{}, db)

	req := createRequestWithToken("GET", "/api/v1/students/student-id-1", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetStudentByIDRoute_NotFound(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "user:manage").
		WillReturnRows(rows)

	mockStudentService := &mockStudentServiceForRoute{
		err: sql.ErrNoRows,
	}

	app := setupTestApp()
	routepostgre.StudentRoutes(app, mockStudentService, &mockAchievementServiceForRoute{}, db)

	req := createRequestWithToken("GET", "/api/v1/students/nonexistent-id", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusNotFound)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetStudentAchievementsRoute_Success(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "achievement:read").
		WillReturnRows(rows)

	mockAchievementService := &mockAchievementServiceForRoute{
		getByStudentIDResp: map[string]interface{}{
			"status": "success",
			"data":   []interface{}{},
			"pagination": map[string]interface{}{
				"page":        1,
				"limit":       10,
				"total":       0,
				"total_pages": 0,
			},
		},
	}

	app := setupTestApp()
	routepostgre.StudentRoutes(app, &mockStudentServiceForRoute{}, mockAchievementService, db)

	req := createRequestWithToken("GET", "/api/v1/students/student-id-1/achievements", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateStudentAdvisorRoute_Success(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "user:manage").
		WillReturnRows(rows)

	app := setupTestApp()
	routepostgre.StudentRoutes(app, &mockStudentServiceForRoute{}, &mockAchievementServiceForRoute{}, db)

	reqBody := map[string]string{
		"advisor_id": "lecturer-id-1",
	}

	req := createRequestWithToken("PUT", "/api/v1/students/student-id-1/advisor", reqBody, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateStudentAdvisorRoute_NotFound(t *testing.T) {
	db, mock := setupTestDBForRoute(t)
	defer db.Close()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(true)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "user:manage").
		WillReturnRows(rows)

	mockStudentService := &mockStudentServiceForRoute{
		updateAdvisorErr: sql.ErrNoRows,
	}

	app := setupTestApp()
	routepostgre.StudentRoutes(app, mockStudentService, &mockAchievementServiceForRoute{}, db)

	reqBody := map[string]string{
		"advisor_id": "lecturer-id-1",
	}

	req := createRequestWithToken("PUT", "/api/v1/students/nonexistent-id/advisor", reqBody, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusNotFound)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
