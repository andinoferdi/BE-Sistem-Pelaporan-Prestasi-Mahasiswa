package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepository_FindUserByID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedUser := &modelpostgre.User{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		RoleID:       "role-id-1",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "full_name", "role_id", "is_active", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.PasswordHash,
			expectedUser.FullName, expectedUser.RoleID, expectedUser.IsActive, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(`SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.FindUserByID(ctx, userID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("Expected user, got nil")
	}

	if user.ID != expectedUser.ID {
		t.Errorf("Expected ID %s, got %s", expectedUser.ID, user.ID)
	}

	if user.Username != expectedUser.Username {
		t.Errorf("Expected Username %s, got %s", expectedUser.Username, user.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_FindUserByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"

	mock.ExpectQuery(`SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.id = \$1`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.FindUserByID(ctx, userID)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}

	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_FindUserByEmail_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	email := "test@example.com"
	expectedUser := &modelpostgre.User{
		ID:           "550e8400-e29b-41d4-a716-446655440000",
		Username:     "testuser",
		Email:        email,
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		RoleID:       "role-id-1",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "full_name", "role_id", "is_active", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.PasswordHash,
			expectedUser.FullName, expectedUser.RoleID, expectedUser.IsActive, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(`SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.email = \$1`).
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.FindUserByEmail(ctx, email)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("Expected user, got nil")
	}

	if user.Email != email {
		t.Errorf("Expected Email %s, got %s", email, user.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_FindUserByUsernameOrEmail_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	usernameOrEmail := "testuser"
	expectedUser := &modelpostgre.User{
		ID:           "550e8400-e29b-41d4-a716-446655440000",
		Username:     usernameOrEmail,
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		RoleID:       "role-id-1",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "full_name", "role_id", "is_active", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.PasswordHash,
			expectedUser.FullName, expectedUser.RoleID, expectedUser.IsActive, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(`SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		WHERE u.username = \$1 OR u.email = \$1`).
		WithArgs(usernameOrEmail).
		WillReturnRows(rows)

	user, err := repo.FindUserByUsernameOrEmail(ctx, usernameOrEmail)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("Expected user, got nil")
	}

	if user.Username != usernameOrEmail && user.Email != usernameOrEmail {
		t.Errorf("Expected Username or Email %s, got Username: %s, Email: %s", usernameOrEmail, user.Username, user.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_CreateUser_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	user := modelpostgre.User{
		Username:     "newuser",
		Email:        "newuser@example.com",
		PasswordHash: "hashed_password",
		FullName:     "New User",
		RoleID:       "role-id-1",
		IsActive:     true,
	}

	expectedID := "550e8400-e29b-41d4-a716-446655440000"
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "full_name", "role_id", "is_active", "created_at", "updated_at"}).
		AddRow(expectedID, user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive, expectedCreatedAt, expectedUpdatedAt)

	mock.ExpectQuery(`INSERT INTO users \(username, email, password_hash, full_name, role_id, is_active, created_at, updated_at\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, NOW\(\), NOW\(\)\)
		RETURNING id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at`).
		WithArgs(user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive).
		WillReturnRows(rows)

	createdUser, err := repo.CreateUser(ctx, user)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createdUser == nil {
		t.Fatal("Expected created user, got nil")
	}

	if createdUser.ID != expectedID {
		t.Errorf("Expected ID %s, got %s", expectedID, createdUser.ID)
	}

	if createdUser.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, createdUser.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_UpdateUser_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	user := modelpostgre.User{
		Username:     "updateduser",
		Email:        "updated@example.com",
		FullName:     "Updated User",
		RoleID:       "role-id-2",
		IsActive:     false,
		PasswordHash: "old_hash",
	}

	expectedUpdatedAt := time.Now()
	expectedCreatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "full_name", "role_id", "is_active", "created_at", "updated_at"}).
		AddRow(userID, user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive, expectedCreatedAt, expectedUpdatedAt)

	mock.ExpectQuery(`UPDATE users
		SET username = \$1, email = \$2, full_name = \$3, role_id = \$4, is_active = \$5, updated_at = NOW\(\)
		WHERE id = \$6
		RETURNING id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at`).
		WithArgs(user.Username, user.Email, user.FullName, user.RoleID, user.IsActive, userID).
		WillReturnRows(rows)

	updatedUser, err := repo.UpdateUser(ctx, userID, user)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedUser == nil {
		t.Fatal("Expected updated user, got nil")
	}

	if updatedUser.ID != userID {
		t.Errorf("Expected ID %s, got %s", userID, updatedUser.ID)
	}

	if updatedUser.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, updatedUser.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_DeleteUser_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUser(ctx, userID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_DeleteUser_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteUser(ctx, userID)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetUserPermissions_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedPermissions := []string{"achievement:create", "achievement:read", "achievement:update"}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("achievement:create").
		AddRow("achievement:read").
		AddRow("achievement:update")

	mock.ExpectQuery(`SELECT p.name
		FROM role_permissions rp
		INNER JOIN permissions p ON rp.permission_id = p.id
		INNER JOIN users u ON u.role_id = rp.role_id
		WHERE u.id = \$1
		ORDER BY p.name`).
		WithArgs(userID).
		WillReturnRows(rows)

	permissions, err := repo.GetUserPermissions(ctx, userID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(permissions) != len(expectedPermissions) {
		t.Errorf("Expected %d permissions, got %d", len(expectedPermissions), len(permissions))
	}

	for i, perm := range expectedPermissions {
		if permissions[i] != perm {
			t.Errorf("Expected permission %s at index %d, got %s", perm, i, permissions[i])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetRoleName_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	roleID := "550e8400-e29b-41d4-a716-446655440000"
	expectedRoleName := "Mahasiswa"

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(expectedRoleName)

	mock.ExpectQuery(`SELECT name FROM roles WHERE id = \$1`).
		WithArgs(roleID).
		WillReturnRows(rows)

	roleName, err := repo.GetRoleName(ctx, roleID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if roleName != expectedRoleName {
		t.Errorf("Expected role name %s, got %s", expectedRoleName, roleName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetRoleName_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	roleID := "550e8400-e29b-41d4-a716-446655440000"

	mock.ExpectQuery(`SELECT name FROM roles WHERE id = \$1`).
		WithArgs(roleID).
		WillReturnError(sql.ErrNoRows)

	roleName, err := repo.GetRoleName(ctx, roleID)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}

	if roleName != "" {
		t.Errorf("Expected empty role name, got %s", roleName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetAllRoles_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at"}).
		AddRow("role-id-1", "Admin", "Administrator role", time.Now()).
		AddRow("role-id-2", "Mahasiswa", "Student role", time.Now()).
		AddRow("role-id-3", "Dosen Wali", "Advisor role", time.Now())

	mock.ExpectQuery(`SELECT id, name, description, created_at
		FROM roles
		ORDER BY name`).
		WillReturnRows(rows)

	roles, err := repo.GetAllRoles(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(roles) != 3 {
		t.Errorf("Expected 3 roles, got %d", len(roles))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetAllUsers_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "full_name", "role_id", "is_active", "created_at", "updated_at"}).
		AddRow("user-id-1", "user1", "user1@example.com", "hash1", "User One", "role-id-1", true, time.Now(), time.Now()).
		AddRow("user-id-2", "user2", "user2@example.com", "hash2", "User Two", "role-id-2", true, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.role_id, u.is_active, u.created_at, u.updated_at
		FROM users u
		ORDER BY u.created_at DESC`).
		WillReturnRows(rows)

	users, err := repo.GetAllUsers(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_UpdateUserRole_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	mock.ExpectExec(`UPDATE users
		SET role_id = \$1, updated_at = NOW\(\)
		WHERE id = \$2`).
		WithArgs(roleID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateUserRole(ctx, userID, roleID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_UpdateUserRole_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	mock.ExpectExec(`UPDATE users
		SET role_id = \$1, updated_at = NOW\(\)
		WHERE id = \$2`).
		WithArgs(roleID, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateUserRole(ctx, userID, roleID)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetLecturerByUserID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedLecturer := &modelpostgre.Lecturer{
		ID:         "lecturer-id-1",
		UserID:     userID,
		LecturerID: "LEC001",
		Department: "Computer Science",
		CreatedAt:  time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "lecturer_id", "department", "created_at"}).
		AddRow(expectedLecturer.ID, expectedLecturer.UserID, expectedLecturer.LecturerID, expectedLecturer.Department, expectedLecturer.CreatedAt)

	mock.ExpectQuery(`SELECT l.id, l.user_id, l.lecturer_id, l.department, l.created_at
		FROM lecturers l
		WHERE l.user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	lecturer, err := repo.GetLecturerByUserID(ctx, userID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if lecturer == nil {
		t.Fatal("Expected lecturer, got nil")
	}

	if lecturer.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, lecturer.UserID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserRepository_GetLecturerByID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewUserRepository(db)
	ctx := context.Background()

	lecturerID := "550e8400-e29b-41d4-a716-446655440000"
	expectedLecturer := &modelpostgre.Lecturer{
		ID:         lecturerID,
		UserID:     "user-id-1",
		LecturerID: "LEC001",
		Department: "Computer Science",
		CreatedAt:  time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "lecturer_id", "department", "created_at"}).
		AddRow(expectedLecturer.ID, expectedLecturer.UserID, expectedLecturer.LecturerID, expectedLecturer.Department, expectedLecturer.CreatedAt)

	mock.ExpectQuery(`SELECT l.id, l.user_id, l.lecturer_id, l.department, l.created_at
		FROM lecturers l
		WHERE l.id = \$1`).
		WithArgs(lecturerID).
		WillReturnRows(rows)

	lecturer, err := repo.GetLecturerByID(ctx, lecturerID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if lecturer == nil {
		t.Fatal("Expected lecturer, got nil")
	}

	if lecturer.ID != lecturerID {
		t.Errorf("Expected ID %s, got %s", lecturerID, lecturer.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
