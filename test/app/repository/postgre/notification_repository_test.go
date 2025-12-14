package repository_test

import (
	"context"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorypostgre "sistem-pelaporan-prestasi-mahasiswa/app/repository/postgre"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNotificationRepository_CreateNotification_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewNotificationRepository(db)
	ctx := context.Background()

	achievementID := "550e8400-e29b-41d4-a716-446655440001"
	mongoAchievementID := "507f1f77bcf86cd799439011"
	req := modelpostgre.CreateNotificationRequest{
		UserID:             "550e8400-e29b-41d4-a716-446655440000",
		Type:               "achievement_submitted",
		Title:              "Prestasi Dikirim",
		Message:            "Prestasi Anda telah dikirim untuk verifikasi",
		AchievementID:      &achievementID,
		MongoAchievementID: &mongoAchievementID,
	}

	expectedID := "550e8400-e29b-41d4-a716-446655440002"
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "type", "title", "message", "achievement_id", "mongo_achievement_id", "is_read", "read_at", "created_at", "updated_at"}).
		AddRow(expectedID, req.UserID, req.Type, req.Title, req.Message, req.AchievementID, req.MongoAchievementID, false, nil, expectedCreatedAt, expectedUpdatedAt)

	mock.ExpectQuery(`INSERT INTO notifications \(user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, created_at, updated_at\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, false, NOW\(\), NOW\(\)\)
		RETURNING id, user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, read_at, created_at, updated_at`).
		WithArgs(req.UserID, req.Type, req.Title, req.Message, req.AchievementID, req.MongoAchievementID).
		WillReturnRows(rows)

	notif, err := repo.CreateNotification(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if notif == nil {
		t.Fatal("Expected notification, got nil")
	}

	if notif.ID != expectedID {
		t.Errorf("Expected ID %s, got %s", expectedID, notif.ID)
	}

	if notif.Type != req.Type {
		t.Errorf("Expected Type %s, got %s", req.Type, notif.Type)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestNotificationRepository_GetNotificationsByUserIDPaginated_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewNotificationRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	page := 1
	limit := 10

	countRows := sqlmock.NewRows([]string{"count"}).
		AddRow(2)

	mock.ExpectQuery(`SELECT COUNT\(\*\)
		FROM notifications
		WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(countRows)

	dataRows := sqlmock.NewRows([]string{"id", "user_id", "type", "title", "message", "achievement_id", "mongo_achievement_id", "is_read", "read_at", "created_at", "updated_at"}).
		AddRow("notif-id-1", userID, "achievement_submitted", "Title 1", "Message 1", nil, nil, false, nil, time.Now(), time.Now()).
		AddRow("notif-id-2", userID, "achievement_verified", "Title 2", "Message 2", nil, nil, true, timePtr(time.Now()), time.Now(), time.Now())

	mock.ExpectQuery(`SELECT id, user_id, type, title, message, achievement_id, mongo_achievement_id, is_read, read_at, created_at, updated_at
		FROM notifications
		WHERE user_id = \$1
		ORDER BY created_at DESC
		LIMIT \$2 OFFSET \$3`).
		WithArgs(userID, limit, 0).
		WillReturnRows(dataRows)

	notifs, total, err := repo.GetNotificationsByUserIDPaginated(ctx, userID, page, limit)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if total != 2 {
		t.Errorf("Expected total 2, got %d", total)
	}

	if len(notifs) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(notifs))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestNotificationRepository_GetUnreadCountByUserID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewNotificationRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	expectedCount := 5

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(expectedCount)

	mock.ExpectQuery(`SELECT COUNT\(\*\)
		FROM notifications
		WHERE user_id = \$1 AND is_read = false`).
		WithArgs(userID).
		WillReturnRows(rows)

	count, err := repo.GetUnreadCountByUserID(ctx, userID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestNotificationRepository_MarkAsRead_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewNotificationRepository(db)
	ctx := context.Background()

	notificationID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440000"

	mock.ExpectExec(`UPDATE notifications
		SET is_read = true, read_at = NOW\(\), updated_at = NOW\(\)
		WHERE id = \$1 AND user_id = \$2`).
		WithArgs(notificationID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.MarkAsRead(ctx, notificationID, userID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestNotificationRepository_MarkAllAsRead_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := repositorypostgre.NewNotificationRepository(db)
	ctx := context.Background()

	userID := "550e8400-e29b-41d4-a716-446655440000"

	mock.ExpectExec(`UPDATE notifications
		SET is_read = true, read_at = NOW\(\), updated_at = NOW\(\)
		WHERE user_id = \$1 AND is_read = false`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 3))

	err := repo.MarkAllAsRead(ctx, userID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
