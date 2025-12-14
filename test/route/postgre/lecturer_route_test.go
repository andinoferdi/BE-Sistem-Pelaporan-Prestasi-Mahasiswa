package route_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	routepostgre "sistem-pelaporan-prestasi-mahasiswa/route/postgre"

	"github.com/DATA-DOG/go-sqlmock"
)

type mockLecturerServiceForRoute struct {
	allLecturers     []modelpostgre.Lecturer
	lecturerByID     *modelpostgre.Lecturer
	lecturerByUserID *modelpostgre.Lecturer
	err              error
}

func (m *mockLecturerServiceForRoute) GetAllLecturers(ctx context.Context) ([]modelpostgre.Lecturer, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allLecturers, nil
}

func (m *mockLecturerServiceForRoute) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.lecturerByUserID, nil
}

func (m *mockLecturerServiceForRoute) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, errors.New("lecturer dengan ID tersebut tidak ditemukan")
		}
		return nil, m.err
	}
	return m.lecturerByID, nil
}

func (m *mockLecturerServiceForRoute) CreateLecturer(ctx context.Context, req modelpostgre.CreateLecturerRequest) (*modelpostgre.Lecturer, error) {
	return nil, errors.New("not implemented in route test")
}

type mockStudentServiceForLecturerRoute struct {
	allStudents         []modelpostgre.Student
	studentByID         *modelpostgre.Student
	studentByUserID     *modelpostgre.Student
	studentsByAdvisorID []modelpostgre.Student
	studentIDByUserID   string
	err                 error
}

func (m *mockStudentServiceForLecturerRoute) GetAllStudents(ctx context.Context) ([]modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allStudents, nil
}

func (m *mockStudentServiceForLecturerRoute) GetStudentByID(ctx context.Context, id string) (*modelpostgre.Student, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.studentByID, nil
}

func (m *mockStudentServiceForLecturerRoute) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.studentIDByUserID, nil
}

func (m *mockStudentServiceForLecturerRoute) GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.studentByUserID, nil
}

func (m *mockStudentServiceForLecturerRoute) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.studentsByAdvisorID, nil
}

func (m *mockStudentServiceForLecturerRoute) CreateStudent(ctx context.Context, req modelpostgre.CreateStudentRequest) (*modelpostgre.Student, error) {
	return nil, errors.New("not implemented")
}

func (m *mockStudentServiceForLecturerRoute) UpdateStudentAdvisor(ctx context.Context, studentID string, advisorID string) error {
	return errors.New("not implemented")
}

func TestGetAllLecturersRoute_Success(t *testing.T) {
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

	mockLecturerService := &mockLecturerServiceForRoute{
		allLecturers: []modelpostgre.Lecturer{
			{
				ID:         "lecturer-id-1",
				UserID:     "user-id-1",
				LecturerID: "D001",
				Department: "Teknik Informatika",
				CreatedAt:  time.Now(),
			},
		},
	}

	app := setupTestApp()
	routepostgre.LecturerRoutes(app, mockLecturerService, &mockStudentServiceForLecturerRoute{}, db)

	req := createRequestWithToken("GET", "/api/v1/lecturers", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetLecturerAdviseesRoute_Success(t *testing.T) {
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

	mockLecturerService := &mockLecturerServiceForRoute{
		lecturerByID: &modelpostgre.Lecturer{
			ID:         "lecturer-id-1",
			UserID:     "user-id-1",
			LecturerID: "D001",
			Department: "Teknik Informatika",
			CreatedAt:  time.Now(),
		},
	}

	mockStudentService := &mockStudentServiceForLecturerRoute{
		studentsByAdvisorID: []modelpostgre.Student{
			{
				ID:           "student-id-1",
				UserID:       "user-id-2",
				StudentID:    "M001",
				ProgramStudy: "Teknik Informatika",
				AcademicYear: "2024",
				AdvisorID:    "lecturer-id-1",
				CreatedAt:    time.Now(),
			},
		},
	}

	app := setupTestApp()
	routepostgre.LecturerRoutes(app, mockLecturerService, mockStudentService, db)

	req := createRequestWithToken("GET", "/api/v1/lecturers/lecturer-id-1/advisees", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetLecturerAdviseesRoute_LecturerNotFound(t *testing.T) {
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

	mockLecturerService := &mockLecturerServiceForRoute{
		err: sql.ErrNoRows,
	}

	app := setupTestApp()
	routepostgre.LecturerRoutes(app, mockLecturerService, &mockStudentServiceForLecturerRoute{}, db)

	req := createRequestWithToken("GET", "/api/v1/lecturers/nonexistent-id/advisees", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusNotFound)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
