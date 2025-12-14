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

type mockUserService struct {
	allUsers      []modelpostgre.User
	userByID      *modelpostgre.User
	createdUser   *modelpostgre.User
	updatedUser   *modelpostgre.User
	allRoles      []modelpostgre.Role
	err           error
	createErr     error
	updateErr     error
	deleteErr     error
	updateRoleErr error
}

func (m *mockUserService) GetAllUsers(ctx context.Context) ([]modelpostgre.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allUsers, nil
}

func (m *mockUserService) GetUserByID(ctx context.Context, id string) (*modelpostgre.User, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, errors.New("user dengan ID tersebut tidak ditemukan")
		}
		return nil, m.err
	}
	return m.userByID, nil
}

func (m *mockUserService) CreateUser(ctx context.Context, req modelpostgre.CreateUserRequest) (*modelpostgre.User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.createdUser, nil
}

func (m *mockUserService) UpdateUser(ctx context.Context, id string, req modelpostgre.UpdateUserRequest) (*modelpostgre.User, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.updatedUser, nil
}

func (m *mockUserService) DeleteUser(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	return m.err
}

func (m *mockUserService) UpdateUserRole(ctx context.Context, id string, roleID string) error {
	if m.updateRoleErr != nil {
		return m.updateRoleErr
	}
	return m.err
}

func (m *mockUserService) GetAllRoles(ctx context.Context) ([]modelpostgre.Role, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allRoles, nil
}

type mockStudentService struct {
	allStudents         []modelpostgre.Student
	studentByID         *modelpostgre.Student
	studentByUserID     *modelpostgre.Student
	studentsByAdvisorID []modelpostgre.Student
	studentIDByUserID   string
	createdStudent      *modelpostgre.Student
	err                 error
	createErr           error
	updateAdvisorErr    error
}

func (m *mockStudentService) GetAllStudents(ctx context.Context) ([]modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allStudents, nil
}

func (m *mockStudentService) GetStudentByID(ctx context.Context, id string) (*modelpostgre.Student, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.studentByID, nil
}

func (m *mockStudentService) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.studentIDByUserID, nil
}

func (m *mockStudentService) GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.studentByUserID, nil
}

func (m *mockStudentService) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.studentsByAdvisorID, nil
}

func (m *mockStudentService) CreateStudent(ctx context.Context, req modelpostgre.CreateStudentRequest) (*modelpostgre.Student, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.createdStudent, nil
}

func (m *mockStudentService) UpdateStudentAdvisor(ctx context.Context, studentID string, advisorID string) error {
	if m.updateAdvisorErr != nil {
		return m.updateAdvisorErr
	}
	return m.err
}

type mockLecturerService struct {
	allLecturers     []modelpostgre.Lecturer
	lecturerByID     *modelpostgre.Lecturer
	lecturerByUserID *modelpostgre.Lecturer
	createdLecturer  *modelpostgre.Lecturer
	err              error
	createErr        error
}

func (m *mockLecturerService) GetAllLecturers(ctx context.Context) ([]modelpostgre.Lecturer, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allLecturers, nil
}

func (m *mockLecturerService) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.lecturerByUserID, nil
}

func (m *mockLecturerService) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.lecturerByID, nil
}

func (m *mockLecturerService) CreateLecturer(ctx context.Context, req modelpostgre.CreateLecturerRequest) (*modelpostgre.Lecturer, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.createdLecturer, nil
}

func TestGetAllUsersRoute_Success(t *testing.T) {
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

	mockUserService := &mockUserService{
		allUsers: []modelpostgre.User{
			{
				ID:        "user-id-1",
				Username:  "user1",
				Email:     "user1@example.com",
				FullName:  "User One",
				RoleID:    "role-id-1",
				IsActive:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	app := setupTestApp()
	routepostgre.UserRoutes(app, mockUserService, &mockStudentService{}, &mockLecturerService{}, db)

	req := createRequestWithToken("GET", "/api/v1/users", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllUsersRoute_NoPermission(t *testing.T) {
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
		AddRow(false)

	mock.ExpectQuery(getPermissionQuery()).
		WithArgs(userID, "user:manage").
		WillReturnRows(rows)

	app := setupTestApp()
	routepostgre.UserRoutes(app, &mockUserService{}, &mockStudentService{}, &mockLecturerService{}, db)

	req := createRequestWithToken("GET", "/api/v1/users", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusForbidden)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetUserByIDRoute_Success(t *testing.T) {
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

	mockUserService := &mockUserService{
		userByID: &modelpostgre.User{
			ID:        "user-id-1",
			Username:  "testuser",
			Email:     "test@example.com",
			FullName:  "Test User",
			RoleID:    "role-id-1",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	app := setupTestApp()
	routepostgre.UserRoutes(app, mockUserService, &mockStudentService{}, &mockLecturerService{}, db)

	req := createRequestWithToken("GET", "/api/v1/users/user-id-1", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetUserByIDRoute_NotFound(t *testing.T) {
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

	mockUserService := &mockUserService{
		err: sql.ErrNoRows,
	}

	app := setupTestApp()
	routepostgre.UserRoutes(app, mockUserService, &mockStudentService{}, &mockLecturerService{}, db)

	req := createRequestWithToken("GET", "/api/v1/users/nonexistent-id", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusNotFound)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateUserRoute_Success(t *testing.T) {
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

	mockUserService := &mockUserService{
		createdUser: &modelpostgre.User{
			ID:        "new-user-id",
			Username:  "newuser",
			Email:     "newuser@example.com",
			FullName:  "New User",
			RoleID:    "role-id-1",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	app := setupTestApp()
	routepostgre.UserRoutes(app, mockUserService, &mockStudentService{}, &mockLecturerService{}, db)

	reqBody := modelpostgre.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
		RoleID:   "role-id-1",
	}

	req := createRequestWithToken("POST", "/api/v1/users", reqBody, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateUserRoute_InvalidBody(t *testing.T) {
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
	routepostgre.UserRoutes(app, &mockUserService{}, &mockStudentService{}, &mockLecturerService{}, db)

	req := createRequestWithToken("POST", "/api/v1/users", "invalid json", token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusBadRequest)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateUserRoute_DuplicateEmail(t *testing.T) {
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

	mockUserService := &mockUserService{
		createErr: errors.New("email sudah digunakan"),
	}

	app := setupTestApp()
	routepostgre.UserRoutes(app, mockUserService, &mockStudentService{}, &mockLecturerService{}, db)

	reqBody := modelpostgre.CreateUserRequest{
		Username: "newuser",
		Email:    "existing@example.com",
		Password: "password123",
		FullName: "New User",
		RoleID:   "role-id-1",
	}

	req := createRequestWithToken("POST", "/api/v1/users", reqBody, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusConflict)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateUserRoute_Success(t *testing.T) {
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

	mockUserService := &mockUserService{
		updatedUser: &modelpostgre.User{
			ID:        "user-id-1",
			Username:  "updateduser",
			Email:     "updated@example.com",
			FullName:  "Updated User",
			RoleID:    "role-id-1",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	app := setupTestApp()
	routepostgre.UserRoutes(app, mockUserService, &mockStudentService{}, &mockLecturerService{}, db)

	reqBody := modelpostgre.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
		FullName: "Updated User",
		RoleID:   "role-id-1",
	}

	req := createRequestWithToken("PUT", "/api/v1/users/user-id-1", reqBody, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteUserRoute_Success(t *testing.T) {
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
	routepostgre.UserRoutes(app, &mockUserService{}, &mockStudentService{}, &mockLecturerService{}, db)

	req := createRequestWithToken("DELETE", "/api/v1/users/user-id-1", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateUserRoleRoute_Success(t *testing.T) {
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
	routepostgre.UserRoutes(app, &mockUserService{}, &mockStudentService{}, &mockLecturerService{}, db)

	reqBody := map[string]string{
		"role_id": "role-id-2",
	}

	req := createRequestWithToken("PUT", "/api/v1/users/user-id-1/role", reqBody, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateStudentProfileRoute_Success(t *testing.T) {
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

	mockStudentService := &mockStudentService{
		createdStudent: &modelpostgre.Student{
			ID:           "student-id-1",
			UserID:       "user-id-1",
			StudentID:    "M001",
			ProgramStudy: "Teknik Informatika",
			AcademicYear: "2024",
			CreatedAt:    time.Now(),
		},
	}

	app := setupTestApp()
	routepostgre.UserRoutes(app, &mockUserService{}, mockStudentService, &mockLecturerService{}, db)

	reqBody := modelpostgre.CreateStudentRequest{
		StudentID:    "M001",
		ProgramStudy: "Teknik Informatika",
		AcademicYear: "2024",
	}

	req := createRequestWithToken("POST", "/api/v1/users/user-id-1/student-profile", reqBody, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateLecturerProfileRoute_Success(t *testing.T) {
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

	mockLecturerService := &mockLecturerService{
		createdLecturer: &modelpostgre.Lecturer{
			ID:         "lecturer-id-1",
			UserID:     "user-id-1",
			LecturerID: "D001",
			Department: "Teknik Informatika",
			CreatedAt:  time.Now(),
		},
	}

	app := setupTestApp()
	routepostgre.UserRoutes(app, &mockUserService{}, &mockStudentService{}, mockLecturerService, db)

	reqBody := modelpostgre.CreateLecturerRequest{
		LecturerID: "D001",
		Department: "Teknik Informatika",
	}

	req := createRequestWithToken("POST", "/api/v1/users/user-id-1/lecturer-profile", reqBody, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllRolesRoute_Success(t *testing.T) {
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

	mockUserService := &mockUserService{
		allRoles: []modelpostgre.Role{
			{
				ID:          "role-id-1",
				Name:        "Mahasiswa",
				Description: "Role untuk mahasiswa",
				CreatedAt:   time.Now(),
			},
			{
				ID:          "role-id-2",
				Name:        "Dosen Wali",
				Description: "Role untuk dosen wali",
				CreatedAt:   time.Now(),
			},
		},
	}

	app := setupTestApp()
	routepostgre.UserRoutes(app, mockUserService, &mockStudentService{}, &mockLecturerService{}, db)

	req := createRequestWithToken("GET", "/api/v1/roles", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
