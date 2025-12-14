package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	modelmongo "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"
	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
	servicepostgre "sistem-pelaporan-prestasi-mahasiswa/app/service/postgre"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockNotificationServiceNotificationRepo struct {
	notifications    []modelpostgre.Notification
	unreadCount      int
	err              error
	createErr        error
	markAsReadErr    error
	markAllAsReadErr error
}

func (m *mockNotificationServiceNotificationRepo) CreateNotification(ctx context.Context, req modelpostgre.CreateNotificationRequest) (*modelpostgre.Notification, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.err != nil {
		return nil, m.err
	}
	notif := &modelpostgre.Notification{
		ID:        "notif-id-1",
		UserID:    req.UserID,
		Type:      req.Type,
		Title:     req.Title,
		Message:   req.Message,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if req.AchievementID != nil {
		notif.AchievementID = req.AchievementID
	}
	if req.MongoAchievementID != nil {
		notif.MongoAchievementID = req.MongoAchievementID
	}
	return notif, nil
}

func (m *mockNotificationServiceNotificationRepo) GetNotificationsByUserIDPaginated(ctx context.Context, userID string, page, limit int) ([]modelpostgre.Notification, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.notifications, len(m.notifications), nil
}

func (m *mockNotificationServiceNotificationRepo) GetUnreadCountByUserID(ctx context.Context, userID string) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.unreadCount, nil
}

func (m *mockNotificationServiceNotificationRepo) MarkAsRead(ctx context.Context, notificationID string, userID string) error {
	if m.markAsReadErr != nil {
		return m.markAsReadErr
	}
	return m.err
}

func (m *mockNotificationServiceNotificationRepo) MarkAllAsRead(ctx context.Context, userID string) error {
	if m.markAllAsReadErr != nil {
		return m.markAllAsReadErr
	}
	return m.err
}

type mockNotificationServiceStudentRepo struct {
	byID              *modelpostgre.Student
	studentIDByUserID string
	err               error
}

func (m *mockNotificationServiceStudentRepo) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.studentIDByUserID, nil
}

func (m *mockNotificationServiceStudentRepo) GetStudentByUserID(ctx context.Context, userID string) (*modelpostgre.Student, error) {
	return nil, m.err
}

func (m *mockNotificationServiceStudentRepo) GetStudentByID(ctx context.Context, id string) (*modelpostgre.Student, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockNotificationServiceStudentRepo) GetStudentsByAdvisorID(ctx context.Context, advisorID string) ([]modelpostgre.Student, error) {
	return nil, m.err
}

func (m *mockNotificationServiceStudentRepo) GetAllStudents(ctx context.Context) ([]modelpostgre.Student, error) {
	return nil, m.err
}

func (m *mockNotificationServiceStudentRepo) CreateStudent(ctx context.Context, req modelpostgre.CreateStudentRequest) (*modelpostgre.Student, error) {
	return nil, m.err
}

func (m *mockNotificationServiceStudentRepo) UpdateStudent(ctx context.Context, id string, req modelpostgre.UpdateStudentRequest) (*modelpostgre.Student, error) {
	return nil, m.err
}

func (m *mockNotificationServiceStudentRepo) UpdateStudentAdvisor(ctx context.Context, id string, advisorID string) error {
	return m.err
}

type mockNotificationServiceAchievementRepo struct {
	byID *modelmongo.Achievement
	err  error
}

func (m *mockNotificationServiceAchievementRepo) CreateAchievement(ctx context.Context, achievement *modelmongo.Achievement) (*modelmongo.Achievement, error) {
	return nil, m.err
}

func (m *mockNotificationServiceAchievementRepo) GetAchievementByID(ctx context.Context, id string) (*modelmongo.Achievement, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.byID, nil
}

func (m *mockNotificationServiceAchievementRepo) UpdateAchievement(ctx context.Context, id string, req modelmongo.UpdateAchievementRequest) (*modelmongo.Achievement, error) {
	return nil, m.err
}

func (m *mockNotificationServiceAchievementRepo) DeleteAchievement(ctx context.Context, id string) error {
	return m.err
}

func (m *mockNotificationServiceAchievementRepo) GetAchievementsByStudentID(ctx context.Context, studentID string) ([]modelmongo.Achievement, error) {
	return nil, m.err
}

func (m *mockNotificationServiceAchievementRepo) GetAchievementsByIDs(ctx context.Context, ids []string) ([]modelmongo.Achievement, error) {
	return nil, m.err
}

func (m *mockNotificationServiceAchievementRepo) AddAttachmentToAchievement(ctx context.Context, id string, attachment modelmongo.Attachment) (*modelmongo.Achievement, error) {
	return nil, m.err
}

func (m *mockNotificationServiceAchievementRepo) GetAchievementsByType(ctx context.Context) (map[string]int, error) {
	return nil, m.err
}

func (m *mockNotificationServiceAchievementRepo) GetCompetitionLevelDistribution(ctx context.Context) (map[string]int, error) {
	return nil, m.err
}

func (m *mockNotificationServiceAchievementRepo) GetTopStudentsByPoints(ctx context.Context, limit int) ([]repositorymongo.TopStudentResult, error) {
	return nil, m.err
}

type mockNotificationServiceUserRepo struct {
	roleName         string
	lecturerByID     *modelpostgre.Lecturer
	lecturerByUserID *modelpostgre.Lecturer
	err              error
}

func (m *mockNotificationServiceUserRepo) FindUserByID(ctx context.Context, id string) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockNotificationServiceUserRepo) FindUserByEmail(ctx context.Context, email string) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockNotificationServiceUserRepo) FindUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockNotificationServiceUserRepo) GetAllUsers(ctx context.Context) ([]modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockNotificationServiceUserRepo) CreateUser(ctx context.Context, user modelpostgre.User) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockNotificationServiceUserRepo) UpdateUser(ctx context.Context, id string, user modelpostgre.User) (*modelpostgre.User, error) {
	return nil, m.err
}

func (m *mockNotificationServiceUserRepo) DeleteUser(ctx context.Context, id string) error {
	return m.err
}

func (m *mockNotificationServiceUserRepo) UpdateUserRole(ctx context.Context, id string, roleID string) error {
	return m.err
}

func (m *mockNotificationServiceUserRepo) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return nil, m.err
}

func (m *mockNotificationServiceUserRepo) GetRoleName(ctx context.Context, roleID string) (string, error) {
	return m.roleName, m.err
}

func (m *mockNotificationServiceUserRepo) GetAllRoles(ctx context.Context) ([]modelpostgre.Role, error) {
	return nil, m.err
}

func (m *mockNotificationServiceUserRepo) GetLecturerByUserID(ctx context.Context, userID string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.lecturerByUserID, nil
}

func (m *mockNotificationServiceUserRepo) GetLecturerByID(ctx context.Context, id string) (*modelpostgre.Lecturer, error) {
	if m.err != nil {
		if m.err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, m.err
	}
	return m.lecturerByID, nil
}

func TestGetNotifications_Success(t *testing.T) {
	ctx := setupTestContext()

	mockNotificationRepo := &mockNotificationServiceNotificationRepo{
		notifications: []modelpostgre.Notification{
			{
				ID:        "notif-id-1",
				UserID:    "user-id-1",
				Type:      "achievement_verified",
				Title:     "Prestasi Diverifikasi",
				Message:   "Prestasi Anda telah diverifikasi",
				IsRead:    false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "notif-id-2",
				UserID:    "user-id-1",
				Type:      "achievement_rejected",
				Title:     "Prestasi Ditolak",
				Message:   "Prestasi Anda ditolak",
				IsRead:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	service := servicepostgre.NewNotificationService(
		mockNotificationRepo,
		&mockNotificationServiceStudentRepo{},
		&mockNotificationServiceUserRepo{},
		&mockNotificationServiceAchievementRepo{},
	)

	result, err := service.GetNotifications(ctx, "user-id-1", 1, 10)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if len(result.Data) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(result.Data))
	}
}

func TestGetNotifications_PaginationNormalization(t *testing.T) {
	ctx := setupTestContext()

	mockNotificationRepo := &mockNotificationServiceNotificationRepo{
		notifications: []modelpostgre.Notification{},
	}

	service := servicepostgre.NewNotificationService(
		mockNotificationRepo,
		&mockNotificationServiceStudentRepo{},
		&mockNotificationServiceUserRepo{},
		&mockNotificationServiceAchievementRepo{},
	)

	result, err := service.GetNotifications(ctx, "user-id-1", 0, 200)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Pagination.Page < 1 {
		t.Errorf("Expected page >= 1, got %d", result.Pagination.Page)
	}

	if result.Pagination.Limit > 100 {
		t.Errorf("Expected limit <= 100, got %d", result.Pagination.Limit)
	}
}

func TestGetUnreadCount_Success(t *testing.T) {
	ctx := setupTestContext()

	mockNotificationRepo := &mockNotificationServiceNotificationRepo{
		unreadCount: 5,
	}

	service := servicepostgre.NewNotificationService(
		mockNotificationRepo,
		&mockNotificationServiceStudentRepo{},
		&mockNotificationServiceUserRepo{},
		&mockNotificationServiceAchievementRepo{},
	)

	result, err := service.GetUnreadCount(ctx, "user-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data.Count != 5 {
		t.Errorf("Expected count 5, got %d", result.Data.Count)
	}
}

func TestMarkAsRead_Success(t *testing.T) {
	ctx := setupTestContext()

	mockNotificationRepo := &mockNotificationServiceNotificationRepo{
		notifications: []modelpostgre.Notification{
			{
				ID:        "notif-id-1",
				UserID:    "user-id-1",
				Type:      modelpostgre.NotificationTypeAchievementRejected,
				Title:     "Test Notification",
				Message:   "Test message",
				IsRead:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	service := servicepostgre.NewNotificationService(
		mockNotificationRepo,
		&mockNotificationServiceStudentRepo{},
		&mockNotificationServiceUserRepo{},
		&mockNotificationServiceAchievementRepo{},
	)

	result, err := service.MarkAsRead(ctx, "notif-id-1", "user-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}
}

func TestMarkAsRead_EmptyNotificationID(t *testing.T) {
	ctx := setupTestContext()

	service := servicepostgre.NewNotificationService(
		&mockNotificationServiceNotificationRepo{},
		&mockNotificationServiceStudentRepo{},
		&mockNotificationServiceUserRepo{},
		&mockNotificationServiceAchievementRepo{},
	)

	_, err := service.MarkAsRead(ctx, "", "user-id-1")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "ID notification wajib diisi" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestMarkAsRead_NotFound(t *testing.T) {
	ctx := setupTestContext()

	mockNotificationRepo := &mockNotificationServiceNotificationRepo{
		markAsReadErr: sql.ErrNoRows,
	}

	service := servicepostgre.NewNotificationService(
		mockNotificationRepo,
		&mockNotificationServiceStudentRepo{},
		&mockNotificationServiceUserRepo{},
		&mockNotificationServiceAchievementRepo{},
	)

	_, err := service.MarkAsRead(ctx, "nonexistent-id", "user-id-1")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "notifikasi tidak ditemukan atau bukan milik Anda" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestMarkAllAsRead_Success(t *testing.T) {
	ctx := setupTestContext()

	mockNotificationRepo := &mockNotificationServiceNotificationRepo{}

	service := servicepostgre.NewNotificationService(
		mockNotificationRepo,
		&mockNotificationServiceStudentRepo{},
		&mockNotificationServiceUserRepo{},
		&mockNotificationServiceAchievementRepo{},
	)

	result, err := service.MarkAllAsRead(ctx, "user-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}
}

func TestCreateAchievementNotification_Success(t *testing.T) {
	ctx := setupTestContext()

	mockNotificationRepo := &mockNotificationServiceNotificationRepo{}

	mockStudentRepo := &mockNotificationServiceStudentRepo{
		byID: &modelpostgre.Student{
			ID:     "student-id-1",
			UserID: "user-id-1",
		},
	}

	mockAchievementRepo := &mockNotificationServiceAchievementRepo{
		byID: &modelmongo.Achievement{
			ID:              primitive.NewObjectID(),
			StudentID:       "student-id-1",
			AchievementType: "competition",
			Title:           "Test Achievement",
			Description:     "Test description",
			Points:          100,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	service := servicepostgre.NewNotificationService(
		mockNotificationRepo,
		mockStudentRepo,
		&mockNotificationServiceUserRepo{},
		mockAchievementRepo,
	)

	err := service.CreateAchievementNotification(ctx, "user-id-1", "mongo-id-1", "ref-id-1", "Data tidak lengkap")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCreateAchievementNotification_AchievementNotFound(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockNotificationServiceStudentRepo{
		byID: &modelpostgre.Student{
			ID:     "student-id-1",
			UserID: "user-id-1",
		},
	}

	mockAchievementRepo := &mockNotificationServiceAchievementRepo{
		byID: nil,
		err:  nil,
	}

	service := servicepostgre.NewNotificationService(
		&mockNotificationServiceNotificationRepo{},
		mockStudentRepo,
		&mockNotificationServiceUserRepo{},
		mockAchievementRepo,
	)

	err := service.CreateAchievementNotification(ctx, "user-id-1", "mongo-id-1", "ref-id-1", "Data tidak lengkap")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "prestasi tidak ditemukan" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestCreateSubmissionNotification_Success(t *testing.T) {
	ctx := setupTestContext()

	mockNotificationRepo := &mockNotificationServiceNotificationRepo{}

	mockStudentRepo := &mockNotificationServiceStudentRepo{
		byID: &modelpostgre.Student{
			ID:        "student-id-1",
			UserID:    "user-id-1",
			AdvisorID: "lecturer-id-1",
		},
	}

	mockUserRepo := &mockNotificationServiceUserRepo{
		lecturerByID: &modelpostgre.Lecturer{
			ID:     "lecturer-id-1",
			UserID: "lecturer-user-id-1",
		},
	}

	mockAchievementRepo := &mockNotificationServiceAchievementRepo{
		byID: &modelmongo.Achievement{
			ID:              primitive.NewObjectID(),
			StudentID:       "student-id-1",
			AchievementType: "competition",
			Title:           "Test Achievement",
			Description:     "Test description",
			Points:          100,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	service := servicepostgre.NewNotificationService(
		mockNotificationRepo,
		mockStudentRepo,
		mockUserRepo,
		mockAchievementRepo,
	)

	err := service.CreateSubmissionNotification(ctx, "student-id-1", "mongo-id-1", "ref-id-1")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCreateSubmissionNotification_NoAdvisor(t *testing.T) {
	ctx := setupTestContext()

	mockStudentRepo := &mockNotificationServiceStudentRepo{
		byID: &modelpostgre.Student{
			ID:        "student-id-1",
			UserID:    "user-id-1",
			AdvisorID: "",
		},
	}

	service := servicepostgre.NewNotificationService(
		&mockNotificationServiceNotificationRepo{},
		mockStudentRepo,
		&mockNotificationServiceUserRepo{},
		&mockNotificationServiceAchievementRepo{},
	)

	err := service.CreateSubmissionNotification(ctx, "student-id-1", "mongo-id-1", "ref-id-1")

	if err != nil {
		t.Fatalf("Expected no error when no advisor, got %v", err)
	}
}
