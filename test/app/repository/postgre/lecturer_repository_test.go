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

func TestLecturerRepository_GetLecturerByUserID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewLecturerRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedLecturer := &modelpostgre.Lecturer{
		ID:         "550e8400-e29b-41d4-a716-446655440001",
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

func TestLecturerRepository_GetLecturerByUserID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewLecturerRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"

	mock.ExpectQuery(`SELECT l.id, l.user_id, l.lecturer_id, l.department, l.created_at
		FROM lecturers l
		WHERE l.user_id = \$1`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	lecturer, err := repo.GetLecturerByUserID(ctx, userID)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}

	if lecturer != nil {
		t.Errorf("Expected nil lecturer, got %v", lecturer)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestLecturerRepository_GetLecturerByID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewLecturerRepository(db)
	ctx := context.Background()

	lecturerID := "550e8400-e29b-41d4-a716-446655440001"
	expectedLecturer := &modelpostgre.Lecturer{
		ID:         lecturerID,
		UserID:     "550e8400-e29b-41d4-a716-446655440000",
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

func TestLecturerRepository_GetAllLecturers_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewLecturerRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "user_id", "lecturer_id", "department", "full_name", "created_at"}).
		AddRow("lecturer-id-1", "user-id-1", "LEC001", "Computer Science", "Lecturer One", time.Now()).
		AddRow("lecturer-id-2", "user-id-2", "LEC002", "Information Technology", "Lecturer Two", time.Now())

	mock.ExpectQuery(`SELECT l.id, l.user_id, l.lecturer_id, l.department, COALESCE\(u.full_name, ''\) as full_name, l.created_at
		FROM lecturers l
		LEFT JOIN users u ON l.user_id = u.id
		ORDER BY l.created_at DESC`).
		WillReturnRows(rows)

	lecturers, err := repo.GetAllLecturers(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(lecturers) != 2 {
		t.Errorf("Expected 2 lecturers, got %d", len(lecturers))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestLecturerRepository_CreateLecturer_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewLecturerRepository(db)
	ctx := context.Background()

	req := modelpostgre.CreateLecturerRequest{
		UserID:     "550e8400-e29b-41d4-a716-446655440000",
		LecturerID: "LEC001",
		Department: "Computer Science",
	}

	expectedID := "550e8400-e29b-41d4-a716-446655440001"
	expectedCreatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "lecturer_id", "department", "created_at"}).
		AddRow(expectedID, req.UserID, req.LecturerID, req.Department, expectedCreatedAt)

	mock.ExpectQuery(`INSERT INTO lecturers \(user_id, lecturer_id, department, created_at\)
		VALUES \(\$1, \$2, \$3, NOW\(\)\)
		RETURNING id, user_id, lecturer_id, department, created_at`).
		WithArgs(req.UserID, req.LecturerID, req.Department).
		WillReturnRows(rows)

	lecturer, err := repo.CreateLecturer(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if lecturer == nil {
		t.Fatal("Expected lecturer, got nil")
	}

	if lecturer.ID != expectedID {
		t.Errorf("Expected ID %s, got %s", expectedID, lecturer.ID)
	}

	if lecturer.LecturerID != req.LecturerID {
		t.Errorf("Expected LecturerID %s, got %s", req.LecturerID, lecturer.LecturerID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
