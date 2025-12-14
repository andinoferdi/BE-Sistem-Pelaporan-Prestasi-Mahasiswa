package route_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
	routepostgre "sistem-pelaporan-prestasi-mahasiswa/route/postgre"
)

type mockNotificationService struct {
	getNotificationsResp *modelpostgre.GetNotificationsResponse
	getNotificationsErr  error
	getUnreadCountResp   *modelpostgre.GetUnreadCountResponse
	getUnreadCountErr    error
	markAsReadResp       *modelpostgre.MarkAsReadResponse
	markAsReadErr        error
	markAllAsReadResp    *modelpostgre.MarkAllAsReadResponse
	markAllAsReadErr     error
}

func (m *mockNotificationService) GetNotifications(ctx context.Context, userID string, page, limit int) (*modelpostgre.GetNotificationsResponse, error) {
	if m.getNotificationsErr != nil {
		return nil, m.getNotificationsErr
	}
	return m.getNotificationsResp, nil
}

func (m *mockNotificationService) GetUnreadCount(ctx context.Context, userID string) (*modelpostgre.GetUnreadCountResponse, error) {
	if m.getUnreadCountErr != nil {
		return nil, m.getUnreadCountErr
	}
	return m.getUnreadCountResp, nil
}

func (m *mockNotificationService) MarkAsRead(ctx context.Context, notificationID string, userID string) (*modelpostgre.MarkAsReadResponse, error) {
	if m.markAsReadErr != nil {
		return nil, m.markAsReadErr
	}
	return m.markAsReadResp, nil
}

func (m *mockNotificationService) MarkAllAsRead(ctx context.Context, userID string) (*modelpostgre.MarkAllAsReadResponse, error) {
	if m.markAllAsReadErr != nil {
		return nil, m.markAllAsReadErr
	}
	return m.markAllAsReadResp, nil
}

func (m *mockNotificationService) CreateAchievementNotification(ctx context.Context, studentUserID string, mongoAchievementID string, achievementRefID string, rejectionNote string) error {
	return nil
}

func (m *mockNotificationService) CreateSubmissionNotification(ctx context.Context, studentID string, mongoAchievementID string, achievementRefID string) error {
	return nil
}

func TestGetNotificationsRoute_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockNotificationService{
		getNotificationsResp: &modelpostgre.GetNotificationsResponse{
			Status: "success",
			Data: []modelpostgre.Notification{
				{
					ID:        "notif-id-1",
					UserID:    userID,
					Title:     "Test Notification",
					Message:   "Test message",
					IsRead:    false,
					CreatedAt: time.Now(),
				},
			},
			Pagination: struct {
				Page       int `json:"page"`
				Limit      int `json:"limit"`
				Total      int `json:"total"`
				TotalPages int `json:"total_pages"`
			}{
				Page:       1,
				Limit:      10,
				Total:      1,
				TotalPages: 1,
			},
		},
	}

	app := setupTestApp()
	routepostgre.NotificationRoutes(app, mockService)

	req := createRequestWithToken("GET", "/api/v1/notifications", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestGetNotificationsRoute_MissingUserID(t *testing.T) {
	app := setupTestApp()
	routepostgre.NotificationRoutes(app, &mockNotificationService{})

	req := httptest.NewRequest("GET", "/api/v1/notifications", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusUnauthorized)
}

func TestGetUnreadCountRoute_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockNotificationService{
		getUnreadCountResp: &modelpostgre.GetUnreadCountResponse{
			Status: "success",
			Data: struct {
				Count int `json:"count"`
			}{
				Count: 5,
			},
		},
	}

	app := setupTestApp()
	routepostgre.NotificationRoutes(app, mockService)

	req := createRequestWithToken("GET", "/api/v1/notifications/unread-count", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestMarkAsReadRoute_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockNotificationService{
		markAsReadResp: &modelpostgre.MarkAsReadResponse{
			Status: "success",
			Data: modelpostgre.Notification{
				ID:        "notif-id-1",
				UserID:    userID,
				Title:     "Test Notification",
				Message:   "Test message",
				IsRead:    true,
				ReadAt:    timePtr(time.Now()),
				CreatedAt: time.Now(),
			},
		},
	}

	app := setupTestApp()
	routepostgre.NotificationRoutes(app, mockService)

	req := createRequestWithToken("PUT", "/api/v1/notifications/notif-id-1/read", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func TestMarkAsReadRoute_NotFound(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockNotificationService{
		markAsReadErr: errors.New("notification tidak ditemukan"),
	}

	app := setupTestApp()
	routepostgre.NotificationRoutes(app, mockService)

	req := createRequestWithToken("PUT", "/api/v1/notifications/nonexistent-id/read", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusNotFound)
}

func TestMarkAllAsReadRoute_Success(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email := "test@example.com"
	roleID := "550e8400-e29b-41d4-a716-446655440001"

	token, err := createTestToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	mockService := &mockNotificationService{
		markAllAsReadResp: &modelpostgre.MarkAllAsReadResponse{
			Status: "success",
			Data: struct {
				Message string `json:"message"`
			}{
				Message: "Semua notification berhasil ditandai sebagai read",
			},
		},
	}

	app := setupTestApp()
	routepostgre.NotificationRoutes(app, mockService)

	req := createRequestWithToken("PUT", "/api/v1/notifications/read-all", nil, token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	assertStatusCode(t, resp, http.StatusOK)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
