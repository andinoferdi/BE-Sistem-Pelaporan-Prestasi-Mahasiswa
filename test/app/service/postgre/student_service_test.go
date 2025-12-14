package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"
)

type mockStudentServiceStudentRepo struct {
	byID              *modelpostgre.Student
	byUserID          *modelpostgre.Student
	studentIDByUserID string
	byAdvisorID       []modelpostgre.Student
	allStudents       []modelpostgre.Student
	err               error
}

func (m *mockStudentServiceStudentRepo) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", m.err
	}
	return m.studentIDByUserID, nil
}

func (m *mockStudentServiceStudentRepo) GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byUserID, nil
}

func (m *mockStudentServiceStudentRepo) GetStudentByID(ctx context.Context, id string) (*modelpostgre.Student, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockStudentServiceStudentRepo) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byAdvisorID, nil
}

func (m *mockStudentServiceStudentRepo) GetAllStudents(ctx context.Context) ([]modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.allStudents, nil
}

func (m *mockStudentServiceStudentRepo) CreateStudent(ctx context.Context, req modelpostgre.CreateStudentRequest) (*modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	student := &modelpostgre.Student{
		ID:           "student-id-1",
		UserID:       req.UserID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    req.AdvisorID,
		CreatedAt:    time.Now(),
	}
	return student, nil
}

func (m *mockStudentServiceStudentRepo) UpdateStudent(ctx context.Context, id string, req modelpostgre.UpdateStudentRequest) (*modelpostgre.Student, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockStudentServiceStudentRepo) UpdateStudentAdvisor(ctx context.Context, id string, advisorID string) error {
	return m.err
}

type mockStudentServiceUserRepo struct {
	byID             *modelpostgre.User
	roleName         string
	lecturerByUserID *modelpostgre.Lecturer
	lecturerByID     *modelpostgre.Lecturer
	err              error
}

func (m *mockStudentServiceUserRepo) FindUserByID(ctx context.Context, id string) (*modelpostgre.User, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockStudentServiceUserRepo) FindUserByEmail(ctx context.Context, email string) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockStudentServiceUserRepo) FindUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockStudentServiceUserRepo) GetAllUsers(ctx context.Context) ([]modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockStudentServiceUserRepo) CreateUser(ctx context.Context, user modelpostgre.User) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockStudentServiceUserRepo) UpdateUser(ctx context.Context, id string, user modelpostgre.User) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockStudentServiceUserRepo) DeleteUser(ctx context.Context, id string) error {
	return m.err
}

func (m *mockStudentServiceUserRepo) UpdateUserRole(ctx context.Context, id string, roleID string) error {
	return m.err
}

func (m *mockStudentServiceUserRepo) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return nil, m.err
}

func (m *mockStudentServiceUserRepo) GetRoleName(ctx context.Context, roleID string) (string, error) {
	return m.roleName, m.err
}

func (m *mockStudentServiceUserRepo) GetAllRoles(ctx context.Context) ([]modelpostgre.Role, error) {
	return nil, m.err
}

func (m *mockStudentServiceUserRepo) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.lecturerByUserID, nil
}

func (m *mockStudentServiceUserRepo) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.lecturerByID, nil
}

type mockStudentServiceLecturerRepo struct {
	byID *modelpostgre.Lecturer
	err  error
}

func (m *mockStudentServiceLecturerRepo) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	return nil, m.err
}

func (m *mockStudentServiceLecturerRepo) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockStudentServiceLecturerRepo) GetAllLecturers(ctx context.Context) ([]modelpostgre.Lecturer, error) {
	return nil, m.err
}

func (m *mockStudentServiceLecturerRepo) CreateLecturer(ctx context.Context, req modelpostgre.CreateLecturerRequest) (*modelpostgre.Lecturer, error) {
	return nil, m.err
}

func (m *mockStudentServiceLecturerRepo) UpdateLecturer(ctx context.Context, id string, req modelpostgre.UpdateLecturerRequest) (*modelpostgre.Lecturer, error) {
	return nil, m.err
}

func TestGetAllStudents_Success(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockStudentServiceStudentRepo{
		allStudents: []modelpostgre.Student{
			{
				ID:           "student-id-1",
				UserID:       "user-id-1",
				StudentID:    "STU001",
				ProgramStudy: "Teknik Informatika",
				AcademicYear: "2023",
				CreatedAt:    time.Now(),
			},
			{
				ID:           "student-id-2",
				UserID:       "user-id-2",
				StudentID:    "STU002",
				ProgramStudy: "Sistem Informasi",
				AcademicYear: "2023",
				CreatedAt:    time.Now(),
			},
		},
	}

	service := servicepostgre.NewStudentService(
		mockStudentRepo,
		&mockStudentServiceUserRepo{},
		&mockStudentServiceLecturerRepo{},
	)

	result, err := service.GetAllStudents(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 students, got %d", len(result))
	}
}

func TestGetStudentByID_Success(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockStudentServiceStudentRepo{
		byID: &modelpostgre.Student{
			ID:           "student-id-1",
			UserID:       "user-id-1",
			StudentID:    "STU001",
			ProgramStudy: "Teknik Informatika",
			AcademicYear: "2023",
			CreatedAt:    time.Now(),
		},
	}

	service := servicepostgre.NewStudentService(
		mockStudentRepo,
		&mockStudentServiceUserRepo{},
		&mockStudentServiceLecturerRepo{},
	)

	result, err := service.GetStudentByID(ctx, "student-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != "student-id-1" {
		t.Errorf("Expected ID 'student-id-1', got '%s'", result.ID)
	}
}

func TestGetStudentByID_EmptyID(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewStudentService(
		&mockStudentServiceStudentRepo{},
		&mockStudentServiceUserRepo{},
		&mockStudentServiceLecturerRepo{},
	)

	_, err := service.GetStudentByID(ctx, "")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "student ID wajib diisi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestGetStudentByID_NotFound(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockStudentServiceStudentRepo{
		err: sql.ErrNoRows,
	}

	service := servicepostgre.NewStudentService(
		mockStudentRepo,
		&mockStudentServiceUserRepo{},
		&mockStudentServiceLecturerRepo{},
	)

	_, err := service.GetStudentByID(ctx, "nonexistent-id")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != sql.ErrNoRows && err.Error() != "student tidak ditemukan" {
		t.Errorf("Expected sql.ErrNoRows or 'student tidak ditemukan', got: %v", err)
	}
}

func TestGetStudentIDByUserID_Success(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockStudentServiceStudentRepo{
		studentIDByUserID: "student-id-1",
	}

	service := servicepostgre.NewStudentService(
		mockStudentRepo,
		&mockStudentServiceUserRepo{},
		&mockStudentServiceLecturerRepo{},
	)

	result, err := service.GetStudentIDByUserID(ctx, "user-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != "student-id-1" {
		t.Errorf("Expected 'student-id-1', got '%s'", result)
	}
}

func TestGetStudentByUserID_Success(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockStudentServiceStudentRepo{
		byUserID: &modelpostgre.Student{
			ID:           "student-id-1",
			UserID:       "user-id-1",
			StudentID:    "STU001",
			ProgramStudy: "Teknik Informatika",
			AcademicYear: "2023",
			CreatedAt:    time.Now(),
		},
	}

	service := servicepostgre.NewStudentService(
		mockStudentRepo,
		&mockStudentServiceUserRepo{},
		&mockStudentServiceLecturerRepo{},
	)

	result, err := service.GetStudentByUserID(ctx, "user-id-1")

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

func TestGetStudentsByAdvisorID_Success(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockStudentServiceStudentRepo{
		byAdvisorID: []modelpostgre.Student{
			{
				ID:           "student-id-1",
				UserID:       "user-id-1",
				StudentID:    "STU001",
				AdvisorID:    "lecturer-id-1",
				ProgramStudy: "Teknik Informatika",
				AcademicYear: "2023",
				CreatedAt:    time.Now(),
			},
			{
				ID:           "student-id-2",
				UserID:       "user-id-2",
				StudentID:    "STU002",
				AdvisorID:    "lecturer-id-1",
				ProgramStudy: "Sistem Informasi",
				AcademicYear: "2023",
				CreatedAt:    time.Now(),
			},
		},
	}

	service := servicepostgre.NewStudentService(
		mockStudentRepo,
		&mockStudentServiceUserRepo{},
		&mockStudentServiceLecturerRepo{},
	)

	result, err := service.GetStudentsByAdvisorID(ctx, "lecturer-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 students, got %d", len(result))
	}
}

func TestGetStudentsByAdvisorID_EmptyAdvisorID(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewStudentService(
		&mockStudentServiceStudentRepo{},
		&mockStudentServiceUserRepo{},
		&mockStudentServiceLecturerRepo{},
	)

	_, err := service.GetStudentsByAdvisorID(ctx, "")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "advisor ID wajib diisi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestCreateStudent_Success(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockStudentServiceUserRepo{
		byID: &modelpostgre.User{
			ID:     "user-id-1",
			RoleID: "role-id-1",
		},
		roleName: "Mahasiswa",
	}

	mockStudentRepo := &mockStudentServiceStudentRepo{}

	service := servicepostgre.NewStudentService(
		mockStudentRepo,
		mockUserRepo,
		&mockStudentServiceLecturerRepo{},
	)

	req := modelpostgre.CreateStudentRequest{
		UserID:       "user-id-1",
		StudentID:    "STU001",
		ProgramStudy: "Teknik Informatika",
		AcademicYear: "2023",
		AdvisorID:    "lecturer-id-1",
	}

	result, err := service.CreateStudent(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.StudentID != "STU001" {
		t.Errorf("Expected StudentID 'STU001', got '%s'", result.StudentID)
	}
}

func TestCreateStudent_ValidationErrors(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockStudentServiceUserRepo{
		byID: &modelpostgre.User{
			ID:     "user-id-1",
			RoleID: "role-id-1",
		},
		roleName: "Mahasiswa",
	}

	service := servicepostgre.NewStudentService(
		&mockStudentServiceStudentRepo{},
		mockUserRepo,
		&mockStudentServiceLecturerRepo{},
	)

	testCases := []struct {
		name string
		req  modelpostgre.CreateStudentRequest
		want string
	}{
		{
			name: "empty user ID",
			req: modelpostgre.CreateStudentRequest{
				StudentID:    "STU001",
				ProgramStudy: "Teknik Informatika",
				AcademicYear: "2023",
			},
			want: "user ID wajib diisi",
		},
		{
			name: "empty student ID",
			req: modelpostgre.CreateStudentRequest{
				UserID:       "user-id-1",
				ProgramStudy: "Teknik Informatika",
				AcademicYear: "2023",
			},
			want: "student ID wajib diisi",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateStudent(ctx, tc.req)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != tc.want {
				t.Errorf("Expected error '%s', got: %v", tc.want, err)
			}
		})
	}
}

func TestCreateStudent_UserNotFound(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockStudentServiceUserRepo{
		err: sql.ErrNoRows,
	}

	service := servicepostgre.NewStudentService(
		&mockStudentServiceStudentRepo{},
		mockUserRepo,
		&mockStudentServiceLecturerRepo{},
	)

	req := modelpostgre.CreateStudentRequest{
		UserID:       "nonexistent-id",
		StudentID:    "STU001",
		ProgramStudy: "Teknik Informatika",
		AcademicYear: "2023",
	}

	_, err := service.CreateStudent(ctx, req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "user tidak ditemukan" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestCreateStudent_WrongRole(t *testing.T) {
	ctx := setupTestContext()

	mockUserRepo := &mockStudentServiceUserRepo{
		byID: &modelpostgre.User{
			ID:     "user-id-1",
			RoleID: "role-id-1",
		},
		roleName: "Admin",
	}

	service := servicepostgre.NewStudentService(
		&mockStudentServiceStudentRepo{},
		mockUserRepo,
		&mockStudentServiceLecturerRepo{},
	)

	req := modelpostgre.CreateStudentRequest{
		UserID:       "user-id-1",
		StudentID:    "STU001",
		ProgramStudy: "Teknik Informatika",
		AcademicYear: "2023",
	}

	_, err := service.CreateStudent(ctx, req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "user harus memiliki role Mahasiswa untuk membuat student profile" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestUpdateStudentAdvisor_Success(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockStudentServiceStudentRepo{
		byID: &modelpostgre.Student{
			ID:        "student-id-1",
			AdvisorID: "lecturer-id-1",
		},
	}

	mockLecturerRepo := &mockStudentServiceLecturerRepo{
		byID: &modelpostgre.Lecturer{
			ID:     "lecturer-id-2",
			UserID: "lecturer-user-id-2",
		},
	}

	service := servicepostgre.NewStudentService(
		mockStudentRepo,
		&mockStudentServiceUserRepo{},
		mockLecturerRepo,
	)

	err := service.UpdateStudentAdvisor(ctx, "student-id-1", "lecturer-id-2")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestUpdateStudentAdvisor_EmptyStudentID(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewStudentService(
		&mockStudentServiceStudentRepo{},
		&mockStudentServiceUserRepo{},
		&mockStudentServiceLecturerRepo{},
	)

	err := service.UpdateStudentAdvisor(ctx, "", "lecturer-id-1")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "student ID wajib diisi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestUpdateStudentAdvisor_StudentNotFound(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockStudentServiceStudentRepo{
		err: sql.ErrNoRows,
	}

	service := servicepostgre.NewStudentService(
		mockStudentRepo,
		&mockStudentServiceUserRepo{},
		&mockStudentServiceLecturerRepo{},
	)

	err := service.UpdateStudentAdvisor(ctx, "nonexistent-id", "lecturer-id-1")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "student tidak ditemukan" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestUpdateStudentAdvisor_InvalidAdvisor(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockStudentServiceStudentRepo{
		byID: &modelpostgre.Student{
			ID: "student-id-1",
		},
	}

	mockLecturerRepo := &mockStudentServiceLecturerRepo{
		err: sql.ErrNoRows,
	}

	service := servicepostgre.NewStudentService(
		mockStudentRepo,
		&mockStudentServiceUserRepo{},
		mockLecturerRepo,
	)

	err := service.UpdateStudentAdvisor(ctx, "student-id-1", "nonexistent-lecturer-id")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "advisor ID tidak ditemukan" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}
