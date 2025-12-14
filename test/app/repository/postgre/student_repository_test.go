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

func TestStudentRepository_GetStudentIDByUserID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedStudentID := "550e8400-e29b-41d4-a716-446655440001"

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(expectedStudentID)

	mock.ExpectQuery(`SELECT s.id
		FROM students s
		WHERE s.user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	studentID, err := repo.GetStudentIDByUserID(ctx, userID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if studentID != expectedStudentID {
		t.Errorf("Expected StudentID %s, got %s", expectedStudentID, studentID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestStudentRepository_GetStudentIDByUserID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"

	mock.ExpectQuery(`SELECT s.id
		FROM students s
		WHERE s.user_id = \$1`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	studentID, err := repo.GetStudentIDByUserID(ctx, userID)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}

	if studentID != "" {
		t.Errorf("Expected empty studentID, got %s", studentID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestStudentRepository_GetStudentByUserID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedStudent := &modelpostgre.Student{
		ID:           "550e8400-e29b-41d4-a716-446655440001",
		UserID:       userID,
		StudentID:    "STU001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		AdvisorID:    "550e8400-e29b-41d4-a716-446655440002",
		FullName:     "Test Student",
		CreatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "student_id", "program_study", "academic_year", "advisor_id", "full_name", "created_at"}).
		AddRow(expectedStudent.ID, expectedStudent.UserID, expectedStudent.StudentID,
			expectedStudent.ProgramStudy, expectedStudent.AcademicYear, expectedStudent.AdvisorID,
			expectedStudent.FullName, expectedStudent.CreatedAt)

	mock.ExpectQuery(`SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE\(u.full_name, ''\) as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	student, err := repo.GetStudentByUserID(ctx, userID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if student == nil {
		t.Fatal("Expected student, got nil")
	}

	if student.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, student.UserID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestStudentRepository_GetStudentByID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	studentID := "550e8400-e29b-41d4-a716-446655440001"
	expectedStudent := &modelpostgre.Student{
		ID:           studentID,
		UserID:       "550e8400-e29b-41d4-a716-446655440000",
		StudentID:    "STU001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		AdvisorID:    "550e8400-e29b-41d4-a716-446655440002",
		FullName:     "Test Student",
		CreatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "student_id", "program_study", "academic_year", "advisor_id", "full_name", "created_at"}).
		AddRow(expectedStudent.ID, expectedStudent.UserID, expectedStudent.StudentID,
			expectedStudent.ProgramStudy, expectedStudent.AcademicYear, expectedStudent.AdvisorID,
			expectedStudent.FullName, expectedStudent.CreatedAt)

	mock.ExpectQuery(`SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE\(u.full_name, ''\) as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.id = \$1`).
		WithArgs(studentID).
		WillReturnRows(rows)

	student, err := repo.GetStudentByID(ctx, studentID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if student == nil {
		t.Fatal("Expected student, got nil")
	}

	if student.ID != studentID {
		t.Errorf("Expected ID %s, got %s", studentID, student.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestStudentRepository_GetStudentByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	studentID := "550e8400-e29b-41d4-a716-446655440001"

	mock.ExpectQuery(`SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE\(u.full_name, ''\) as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.id = \$1`).
		WithArgs(studentID).
		WillReturnError(sql.ErrNoRows)

	student, err := repo.GetStudentByID(ctx, studentID)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}

	if student != nil {
		t.Errorf("Expected nil student, got %v", student)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestStudentRepository_GetStudentsByAdvisorID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	advisorID := "550e8400-e29b-41d4-a716-446655440002"

	rows := sqlmock.NewRows([]string{"id", "user_id", "student_id", "program_study", "academic_year", "advisor_id", "full_name", "created_at"}).
		AddRow("student-id-1", "user-id-1", "STU001", "Computer Science", "2024", advisorID, "Student One", time.Now()).
		AddRow("student-id-2", "user-id-2", "STU002", "Computer Science", "2024", advisorID, "Student Two", time.Now())

	mock.ExpectQuery(`SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE\(u.full_name, ''\) as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.advisor_id = \$1
		ORDER BY s.created_at DESC`).
		WithArgs(advisorID).
		WillReturnRows(rows)

	students, err := repo.GetStudentsByAdvisorID(ctx, advisorID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(students) != 2 {
		t.Errorf("Expected 2 students, got %d", len(students))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestStudentRepository_GetAllStudents_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "user_id", "student_id", "program_study", "academic_year", "advisor_id", "full_name", "created_at"}).
		AddRow("student-id-1", "user-id-1", "STU001", "Computer Science", "2024", "advisor-id-1", "Student One", time.Now()).
		AddRow("student-id-2", "user-id-2", "STU002", "Computer Science", "2024", "advisor-id-2", "Student Two", time.Now())

	mock.ExpectQuery(`SELECT s.id, s.user_id, s.student_id, s.program_study, 
		       s.academic_year, s.advisor_id, COALESCE\(u.full_name, ''\) as full_name, s.created_at
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		ORDER BY s.created_at DESC`).
		WillReturnRows(rows)

	students, err := repo.GetAllStudents(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(students) != 2 {
		t.Errorf("Expected 2 students, got %d", len(students))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestStudentRepository_CreateStudent_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	req := modelpostgre.CreateStudentRequest{
		UserID:       "550e8400-e29b-41d4-a716-446655440000",
		StudentID:    "STU001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		AdvisorID:    "550e8400-e29b-41d4-a716-446655440002",
	}

	expectedID := "550e8400-e29b-41d4-a716-446655440001"
	expectedCreatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "student_id", "program_study", "academic_year", "advisor_id", "created_at"}).
		AddRow(expectedID, req.UserID, req.StudentID, req.ProgramStudy, req.AcademicYear, req.AdvisorID, expectedCreatedAt)

	mock.ExpectQuery(`INSERT INTO students \(user_id, student_id, program_study, academic_year, advisor_id, created_at\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, NOW\(\)\)
		RETURNING id, user_id, student_id, program_study, academic_year, advisor_id, created_at`).
		WithArgs(req.UserID, req.StudentID, req.ProgramStudy, req.AcademicYear, req.AdvisorID).
		WillReturnRows(rows)

	student, err := repo.CreateStudent(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if student == nil {
		t.Fatal("Expected student, got nil")
	}

	if student.ID != expectedID {
		t.Errorf("Expected ID %s, got %s", expectedID, student.ID)
	}

	if student.StudentID != req.StudentID {
		t.Errorf("Expected StudentID %s, got %s", req.StudentID, student.StudentID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestStudentRepository_CreateStudent_WithNullAdvisorID(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	req := modelpostgre.CreateStudentRequest{
		UserID:       "550e8400-e29b-41d4-a716-446655440000",
		StudentID:    "STU001",
		ProgramStudy: "Computer Science",
		AcademicYear: "2024",
		AdvisorID:    "",
	}

	expectedID := "550e8400-e29b-41d4-a716-446655440001"
	expectedCreatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "student_id", "program_study", "academic_year", "advisor_id", "created_at"}).
		AddRow(expectedID, req.UserID, req.StudentID, req.ProgramStudy, req.AcademicYear, nil, expectedCreatedAt)

	mock.ExpectQuery(`INSERT INTO students \(user_id, student_id, program_study, academic_year, advisor_id, created_at\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, NOW\(\)\)
		RETURNING id, user_id, student_id, program_study, academic_year, advisor_id, created_at`).
		WithArgs(req.UserID, req.StudentID, req.ProgramStudy, req.AcademicYear, req.AdvisorID).
		WillReturnRows(rows)

	student, err := repo.CreateStudent(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if student == nil {
		t.Fatal("Expected student, got nil")
	}

	if student.AdvisorID != "" {
		t.Errorf("Expected empty AdvisorID, got %s", student.AdvisorID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestStudentRepository_UpdateStudent_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	studentID := "550e8400-e29b-41d4-a716-446655440001"
	req := modelpostgre.UpdateStudentRequest{
		StudentID:    "STU002",
		ProgramStudy: "Information Technology",
		AcademicYear: "2025",
		AdvisorID:    "550e8400-e29b-41d4-a716-446655440003",
	}

	expectedCreatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "student_id", "program_study", "academic_year", "advisor_id", "created_at"}).
		AddRow(studentID, "user-id-1", req.StudentID, req.ProgramStudy, req.AcademicYear, req.AdvisorID, expectedCreatedAt)

	mock.ExpectQuery(`UPDATE students
		SET student_id = \$1, program_study = \$2, academic_year = \$3, advisor_id = \$4
		WHERE id = \$5
		RETURNING id, user_id, student_id, program_study, academic_year, advisor_id, created_at`).
		WithArgs(req.StudentID, req.ProgramStudy, req.AcademicYear, req.AdvisorID, studentID).
		WillReturnRows(rows)

	student, err := repo.UpdateStudent(ctx, studentID, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if student == nil {
		t.Fatal("Expected student, got nil")
	}

	if student.StudentID != req.StudentID {
		t.Errorf("Expected StudentID %s, got %s", req.StudentID, student.StudentID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestStudentRepository_UpdateStudentAdvisor_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	studentID := "550e8400-e29b-41d4-a716-446655440001"
	advisorID := "550e8400-e29b-41d4-a716-446655440002"

	mock.ExpectExec(`UPDATE students
		SET advisor_id = \$1
		WHERE id = \$2`).
		WithArgs(advisorID, studentID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateStudentAdvisor(ctx, studentID, advisorID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestStudentRepository_UpdateStudentAdvisor_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	studentID := "550e8400-e29b-41d4-a716-446655440001"
	advisorID := "550e8400-e29b-41d4-a716-446655440002"

	mock.ExpectExec(`UPDATE students
		SET advisor_id = \$1
		WHERE id = \$2`).
		WithArgs(advisorID, studentID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateStudentAdvisor(ctx, studentID, advisorID)

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

func TestStudentRepository_UpdateStudentAdvisor_WithEmptyAdvisorID(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewStudentRepository(db)
	ctx := context.Background()

	studentID := "550e8400-e29b-41d4-a716-446655440001"
	advisorID := ""

	mock.ExpectExec(`UPDATE students
		SET advisor_id = \$1
		WHERE id = \$2`).
		WithArgs(nil, studentID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateStudentAdvisor(ctx, studentID, advisorID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
