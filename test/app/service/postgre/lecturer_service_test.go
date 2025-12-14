package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
)

type mockLecturerServiceLecturerRepo struct {
	byID         *modelpostgre.Lecturer
	byUserID     *modelpostgre.Lecturer
	allLecturers []modelpostgre.Lecturer
	err          error
}

func (m *mockLecturerServiceLecturerRepo) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byUserID, nil
}

func (m *mockLecturerServiceLecturerRepo) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockLecturerServiceLecturerRepo) GetAllLecturers(ctx context.Context) ([]modelpostgre.Lecturer, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allLecturers, nil
}

func (m *mockLecturerServiceLecturerRepo) CreateLecturer(ctx context.Context, req modelpostgre.CreateLecturerRequest) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		return nil, m.err
	}
	lecturer := &modelpostgre.Lecturer{
		ID:         "lecturer-id-1",
		UserID:     req.UserID,
		LecturerID: req.LecturerID,
		Department: req.Department,
		CreatedAt:  time.Now(),
	}
	return lecturer, nil
}

func (m *mockLecturerServiceLecturerRepo) UpdateLecturer(ctx context.Context, id string, req modelpostgre.UpdateLecturerRequest) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byID, nil
}

type mockLecturerServiceUserRepo struct {
	byID             *modelpostgre.User
	roleName         string
	lecturerByUserID *modelpostgre.Lecturer
	lecturerByID     *modelpostgre.Lecturer
	err              error
}

func (m *mockLecturerServiceUserRepo) FindUserByID(ctx context.Context, id string) (*modelpostgre.User, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockLecturerServiceUserRepo) FindUserByEmail(ctx context.Context, email string) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockLecturerServiceUserRepo) FindUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockLecturerServiceUserRepo) GetAllUsers(ctx context.Context) ([]modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockLecturerServiceUserRepo) CreateUser(ctx context.Context, user modelpostgre.User) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockLecturerServiceUserRepo) UpdateUser(ctx context.Context, id string, user modelpostgre.User) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockLecturerServiceUserRepo) DeleteUser(ctx context.Context, id string) error {
	return m.err
}

func (m *mockLecturerServiceUserRepo) UpdateUserRole(ctx context.Context, id string, roleID string) error {
	return m.err
}

func (m *mockLecturerServiceUserRepo) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return nil, m.err
}

func (m *mockLecturerServiceUserRepo) GetRoleName(ctx context.Context, roleID string) (string, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", m.err
	}
	return m.roleName, nil
}

func (m *mockLecturerServiceUserRepo) GetAllRoles(ctx context.Context) ([]modelpostgre.Role, error) {
	return nil, m.err
}

func (m *mockLecturerServiceUserRepo) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.lecturerByUserID, nil
}

func (m *mockLecturerServiceUserRepo) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.lecturerByID, nil
}

func TestGetAllLecturers_Success(t *testing.T) {
	ctx := setupTestContext()

	mockLecturerRepo := &mockLecturerServiceLecturerRepo{
		allLecturers: []modelpostgre.Lecturer{
			{
				ID:         "lecturer-id-1",
				UserID:     "user-id-1",
				LecturerID: "LEC001",
				Department: "Teknik Informatika",
				CreatedAt:  time.Now(),
			},
			{
				ID:         "lecturer-id-2",
				UserID:     "user-id-2",
				LecturerID: "LEC002",
				Department: "Sistem Informasi",
				CreatedAt:  time.Now(),
			},
		},
	}

	service := servicepostgre.NewLecturerService(
		&mockLecturerServiceUserRepo{},
		mockLecturerRepo,
	)

	result, err := service.GetAllLecturers(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 lecturers, got %d", len(result))
	}
}

func TestGetAllLecturers_RepositoryNil(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewLecturerService(
		&mockLecturerServiceUserRepo{},
		nil,
	)

	_, err := service.GetAllLecturers(ctx)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "lecturer repository tidak tersedia" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGetLecturerByUserID_Success(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockLecturerServiceUserRepo{
		lecturerByUserID: &modelpostgre.Lecturer{
			ID:         "lecturer-id-1",
			UserID:     "user-id-1",
			LecturerID: "LEC001",
			Department: "Teknik Informatika",
			CreatedAt:  time.Now(),
		},
	}

	service := servicepostgre.NewLecturerService(
		mockUserRepo,
		&mockLecturerServiceLecturerRepo{},
	)

	result, err := service.GetLecturerByUserID(ctx, "user-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.UserID != "user-id-1" {
		t.Errorf("Expected UserID 'user-id-1', got '%s'", result.UserID)
	}
}

func TestGetLecturerByUserID_EmptyUserID(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewLecturerService(
		&mockLecturerServiceUserRepo{},
		&mockLecturerServiceLecturerRepo{},
	)

	_, err := service.GetLecturerByUserID(ctx, "")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "user ID wajib diisi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGetLecturerByID_Success(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockLecturerServiceUserRepo{
		lecturerByID: &modelpostgre.Lecturer{
			ID:         "lecturer-id-1",
			UserID:     "user-id-1",
			LecturerID: "LEC001",
			Department: "Teknik Informatika",
			CreatedAt:  time.Now(),
		},
	}

	service := servicepostgre.NewLecturerService(
		mockUserRepo,
		&mockLecturerServiceLecturerRepo{},
	)

	result, err := service.GetLecturerByID(ctx, "lecturer-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != "lecturer-id-1" {
		t.Errorf("Expected ID 'lecturer-id-1', got '%s'", result.ID)
	}
}

func TestGetLecturerByID_EmptyID(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewLecturerService(
		&mockLecturerServiceUserRepo{},
		&mockLecturerServiceLecturerRepo{},
	)

	_, err := service.GetLecturerByID(ctx, "")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "lecturer ID wajib diisi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestCreateLecturer_Success(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockLecturerServiceUserRepo{
		byID: &modelpostgre.User{
			ID:     "user-id-1",
			RoleID: "role-id-1",
		},
		roleName: "Dosen Wali",
	}

	mockLecturerRepo := &mockLecturerServiceLecturerRepo{}

	service := servicepostgre.NewLecturerService(
		mockUserRepo,
		mockLecturerRepo,
	)

	req := modelpostgre.CreateLecturerRequest{
		UserID:     "user-id-1",
		LecturerID: "LEC001",
		Department: "Teknik Informatika",
	}

	result, err := service.CreateLecturer(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.LecturerID != "LEC001" {
		t.Errorf("Expected LecturerID 'LEC001', got '%s'", result.LecturerID)
	}
}

func TestCreateLecturer_ValidationErrors(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewLecturerService(
		&mockLecturerServiceUserRepo{},
		&mockLecturerServiceLecturerRepo{},
	)

	testCases := []struct {
		name string
		req  modelpostgre.CreateLecturerRequest
		want string
	}{
		{
			name: "empty user ID",
			req: modelpostgre.CreateLecturerRequest{
				LecturerID: "LEC001",
				Department: "Teknik Informatika",
			},
			want: "user ID wajib diisi",
		},
		{
			name: "empty lecturer ID",
			req: modelpostgre.CreateLecturerRequest{
				UserID:     "user-id-1",
				Department: "Teknik Informatika",
			},
			want: "lecturer ID wajib diisi",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateLecturer(ctx, tc.req)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != tc.want {
				t.Errorf("Expected error '%s', got: %v", tc.want, err)
			}
		})
	}
}

func TestCreateLecturer_UserNotFound(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockLecturerServiceUserRepo{
		err: sql.ErrNoRows,
	}

	service := servicepostgre.NewLecturerService(
		mockUserRepo,
		&mockLecturerServiceLecturerRepo{},
	)

	req := modelpostgre.CreateLecturerRequest{
		UserID:     "nonexistent-id",
		LecturerID: "LEC001",
		Department: "Teknik Informatika",
	}

	_, err := service.CreateLecturer(ctx, req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "user tidak ditemukan" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestCreateLecturer_WrongRole(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockLecturerServiceUserRepo{
		byID: &modelpostgre.User{
			ID:     "user-id-1",
			RoleID: "role-id-1",
		},
		roleName: "Mahasiswa",
	}

	service := servicepostgre.NewLecturerService(
		mockUserRepo,
		&mockLecturerServiceLecturerRepo{},
	)

	req := modelpostgre.CreateLecturerRequest{
		UserID:     "user-id-1",
		LecturerID: "LEC001",
		Department: "Teknik Informatika",
	}

	_, err := service.CreateLecturer(ctx, req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "user harus memiliki role Dosen Wali untuk membuat lecturer profile" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestCreateLecturer_ExistingProfile(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockLecturerServiceUserRepo{
		byID: &modelpostgre.User{
			ID:     "user-id-1",
			RoleID: "role-id-1",
		},
		roleName: "Dosen Wali",
		lecturerByUserID: &modelpostgre.Lecturer{
			ID:     "lecturer-id-1",
			UserID: "user-id-1",
		},
	}

	service := servicepostgre.NewLecturerService(
		mockUserRepo,
		&mockLecturerServiceLecturerRepo{},
	)

	req := modelpostgre.CreateLecturerRequest{
		UserID:     "user-id-1",
		LecturerID: "LEC001",
		Department: "Teknik Informatika",
	}

	_, err := service.CreateLecturer(ctx, req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "user sudah memiliki lecturer profile" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}
