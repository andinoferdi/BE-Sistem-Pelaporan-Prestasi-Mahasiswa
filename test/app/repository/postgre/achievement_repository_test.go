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

func TestAchievementReferenceRepository_CreateAchievementReference_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewAchievementReferenceRepository(db)
	ctx := context.Background()

	req := modelpostgre.CreateAchievementReferenceRequest{
		StudentID:          "550e8400-e29b-41d4-a716-446655440000",
		MongoAchievementID: "507f1f77bcf86cd799439011",
		Status:             modelpostgre.AchievementStatusDraft,
	}

	expectedID := "550e8400-e29b-41d4-a716-446655440001"
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "student_id", "mongo_achievement_id", "status", "submitted_at", "verified_at", "verified_by", "rejection_note", "created_at", "updated_at"}).
		AddRow(expectedID, req.StudentID, req.MongoAchievementID, req.Status, nil, nil, nil, nil, expectedCreatedAt, expectedUpdatedAt)

	mock.ExpectQuery(`INSERT INTO achievement_references \(student_id, mongo_achievement_id, status, created_at, updated_at\)
		VALUES \(\$1, \$2, \$3, NOW\(\), NOW\(\)\)
		RETURNING id, student_id, mongo_achievement_id, status, submitted_at, 
		          verified_at, verified_by, rejection_note, created_at, updated_at`).
		WithArgs(req.StudentID, req.MongoAchievementID, req.Status).
		WillReturnRows(rows)

	ref, err := repo.CreateAchievementReference(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ref == nil {
		t.Fatal("Expected achievement reference, got nil")
	}

	if ref.ID != expectedID {
		t.Errorf("Expected ID %s, got %s", expectedID, ref.ID)
	}

	if ref.Status != req.Status {
		t.Errorf("Expected Status %s, got %s", req.Status, ref.Status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAchievementReferenceRepository_GetAchievementReferenceByMongoID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewAchievementReferenceRepository(db)
	ctx := context.Background()

	mongoID := "507f1f77bcf86cd799439011"
	expectedRef := &modelpostgre.AchievementReference{
		ID:                 "550e8400-e29b-41d4-a716-446655440001",
		StudentID:          "550e8400-e29b-41d4-a716-446655440000",
		MongoAchievementID: mongoID,
		Status:             modelpostgre.AchievementStatusDraft,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "student_id", "mongo_achievement_id", "status", "submitted_at", "verified_at", "verified_by", "rejection_note", "created_at", "updated_at"}).
		AddRow(expectedRef.ID, expectedRef.StudentID, expectedRef.MongoAchievementID, expectedRef.Status,
			expectedRef.SubmittedAt, expectedRef.VerifiedAt, expectedRef.VerifiedBy, expectedRef.RejectionNote,
			expectedRef.CreatedAt, expectedRef.UpdatedAt)

	mock.ExpectQuery(`SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = \$1 AND status != 'deleted'`).
		WithArgs(mongoID).
		WillReturnRows(rows)

	ref, err := repo.GetAchievementReferenceByMongoID(ctx, mongoID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ref == nil {
		t.Fatal("Expected achievement reference, got nil")
	}

	if ref.MongoAchievementID != mongoID {
		t.Errorf("Expected MongoAchievementID %s, got %s", mongoID, ref.MongoAchievementID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAchievementReferenceRepository_GetAchievementReferenceByMongoID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewAchievementReferenceRepository(db)
	ctx := context.Background()

	mongoID := "507f1f77bcf86cd799439011"

	mock.ExpectQuery(`SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = \$1 AND status != 'deleted'`).
		WithArgs(mongoID).
		WillReturnError(sql.ErrNoRows)

	ref, err := repo.GetAchievementReferenceByMongoID(ctx, mongoID)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}

	if ref != nil {
		t.Errorf("Expected nil reference, got %v", ref)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAchievementReferenceRepository_UpdateAchievementReferenceStatus_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewAchievementReferenceRepository(db)
	ctx := context.Background()

	refID := "550e8400-e29b-41d4-a716-446655440001"
	status := modelpostgre.AchievementStatusSubmitted
	submittedAt := time.Now()

	mock.ExpectExec(`UPDATE achievement_references
		SET status = \$1, submitted_at = \$2, updated_at = NOW\(\)
		WHERE id = \$3`).
		WithArgs(status, submittedAt, refID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateAchievementReferenceStatus(ctx, refID, status, &submittedAt)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAchievementReferenceRepository_DeleteAchievementReference_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewAchievementReferenceRepository(db)
	ctx := context.Background()

	refID := "550e8400-e29b-41d4-a716-446655440001"

	mock.ExpectExec(`DELETE FROM achievement_references WHERE id = \$1`).
		WithArgs(refID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteAchievementReference(ctx, refID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAchievementReferenceRepository_GetAchievementReferenceByStudentID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewAchievementReferenceRepository(db)
	ctx := context.Background()

	studentID := "550e8400-e29b-41d4-a716-446655440000"

	rows := sqlmock.NewRows([]string{"id", "student_id", "mongo_achievement_id", "status", "submitted_at", "verified_at", "verified_by", "rejection_note", "created_at", "updated_at"}).
		AddRow("ref-id-1", studentID, "mongo-id-1", modelpostgre.AchievementStatusDraft, nil, nil, nil, nil, time.Now(), time.Now()).
		AddRow("ref-id-2", studentID, "mongo-id-2", modelpostgre.AchievementStatusSubmitted, timePtr(time.Now()), nil, nil, nil, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT id, student_id, mongo_achievement_id, status, submitted_at,
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = \$1 AND status != 'deleted'
		ORDER BY created_at DESC`).
		WithArgs(studentID).
		WillReturnRows(rows)

	refs, err := repo.GetAchievementReferenceByStudentID(ctx, studentID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(refs) != 2 {
		t.Errorf("Expected 2 references, got %d", len(refs))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
