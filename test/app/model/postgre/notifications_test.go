package model_test

import (
	"encoding/json"
	"testing"
	"time"

	modelpostgre "sistem-pelaporan-prestasi-mahasiswa/app/model/postgre"
)

func TestNotificationTypeConstants(t *testing.T) {
	if modelpostgre.NotificationTypeAchievementRejected != "achievement_rejected" {
		t.Errorf("Expected NotificationTypeAchievementRejected 'achievement_rejected', got '%s'", modelpostgre.NotificationTypeAchievementRejected)
	}
	if modelpostgre.NotificationTypeAchievementSubmitted != "achievement_submitted" {
		t.Errorf("Expected NotificationTypeAchievementSubmitted 'achievement_submitted', got '%s'", modelpostgre.NotificationTypeAchievementSubmitted)
	}
}

func TestNotification_StructCreation(t *testing.T) {
	now := time.Now()
	readAt := now.Add(1 * time.Hour)
	achievementID := "achievement-id-1"
	mongoAchievementID := "mongo-id-1"
	notification := modelpostgre.Notification{
		ID:                 "notification-id-1",
		UserID:             "user-id-1",
		Type:               modelpostgre.NotificationTypeAchievementRejected,
		Title:              "Prestasi Ditolak",
		Message:            "Prestasi Anda ditolak",
		AchievementID:      &achievementID,
		MongoAchievementID: &mongoAchievementID,
		IsRead:             true,
		ReadAt:             &readAt,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if notification.ID != "notification-id-1" {
		t.Errorf("Expected ID 'notification-id-1', got '%s'", notification.ID)
	}
	if notification.UserID != "user-id-1" {
		t.Errorf("Expected UserID 'user-id-1', got '%s'", notification.UserID)
	}
	if notification.Type != modelpostgre.NotificationTypeAchievementRejected {
		t.Errorf("Expected Type '%s', got '%s'", modelpostgre.NotificationTypeAchievementRejected, notification.Type)
	}
	if notification.Title != "Prestasi Ditolak" {
		t.Errorf("Expected Title 'Prestasi Ditolak', got '%s'", notification.Title)
	}
	if notification.Message != "Prestasi Anda ditolak" {
		t.Errorf("Expected Message 'Prestasi Anda ditolak', got '%s'", notification.Message)
	}
	if notification.IsRead != true {
		t.Errorf("Expected IsRead true, got %v", notification.IsRead)
	}
	if notification.AchievementID == nil || *notification.AchievementID != "achievement-id-1" {
		t.Errorf("Expected AchievementID 'achievement-id-1', got '%v'", notification.AchievementID)
	}
}

func TestNotification_JSONMarshalling(t *testing.T) {
	now := time.Now()
	achievementID := "achievement-id-1"
	mongoAchievementID := "mongo-id-1"
	notification := modelpostgre.Notification{
		ID:                 "notification-id-1",
		UserID:             "user-id-1",
		Type:               modelpostgre.NotificationTypeAchievementRejected,
		Title:              "Prestasi Ditolak",
		Message:            "Prestasi Anda ditolak",
		AchievementID:      &achievementID,
		MongoAchievementID: &mongoAchievementID,
		IsRead:             false,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != "notification-id-1" {
		t.Errorf("Expected id 'notification-id-1', got '%v'", result["id"])
	}
	if result["user_id"] != "user-id-1" {
		t.Errorf("Expected user_id 'user-id-1', got '%v'", result["user_id"])
	}
	if result["type"] != "achievement_rejected" {
		t.Errorf("Expected type 'achievement_rejected', got '%v'", result["type"])
	}
	if result["title"] != "Prestasi Ditolak" {
		t.Errorf("Expected title 'Prestasi Ditolak', got '%v'", result["title"])
	}
	if result["is_read"] != false {
		t.Errorf("Expected is_read false, got '%v'", result["is_read"])
	}
}

func TestNotification_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"id": "notification-id-1",
		"user_id": "user-id-1",
		"type": "achievement_rejected",
		"title": "Prestasi Ditolak",
		"message": "Prestasi Anda ditolak",
		"achievement_id": "achievement-id-1",
		"mongo_achievement_id": "mongo-id-1",
		"is_read": true,
		"read_at": "2024-01-01T01:00:00Z",
		"created_at": "2024-01-01T00:00:00Z",
		"updated_at": "2024-01-01T00:00:00Z"
	}`

	var notification modelpostgre.Notification
	if err := json.Unmarshal([]byte(jsonStr), &notification); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if notification.ID != "notification-id-1" {
		t.Errorf("Expected ID 'notification-id-1', got '%s'", notification.ID)
	}
	if notification.Type != "achievement_rejected" {
		t.Errorf("Expected Type 'achievement_rejected', got '%s'", notification.Type)
	}
	if notification.IsRead != true {
		t.Errorf("Expected IsRead true, got %v", notification.IsRead)
	}
}

func TestNotification_ZeroValues(t *testing.T) {
	var notification modelpostgre.Notification

	if notification.ID != "" {
		t.Errorf("Expected empty ID, got '%s'", notification.ID)
	}
	if notification.UserID != "" {
		t.Errorf("Expected empty UserID, got '%s'", notification.UserID)
	}
	if notification.Type != "" {
		t.Errorf("Expected empty Type, got '%s'", notification.Type)
	}
	if notification.IsRead != false {
		t.Errorf("Expected IsRead false, got %v", notification.IsRead)
	}
	if notification.AchievementID != nil {
		t.Error("Expected AchievementID to be nil")
	}
	if notification.MongoAchievementID != nil {
		t.Error("Expected MongoAchievementID to be nil")
	}
	if notification.ReadAt != nil {
		t.Error("Expected ReadAt to be nil")
	}
	if !notification.CreatedAt.IsZero() {
		t.Errorf("Expected zero CreatedAt, got %v", notification.CreatedAt)
	}
}

func TestNotification_PointerFieldsNil(t *testing.T) {
	now := time.Now()
	notification := modelpostgre.Notification{
		ID:        "notification-id-1",
		UserID:    "user-id-1",
		Type:      modelpostgre.NotificationTypeAchievementSubmitted,
		Title:     "Prestasi Dikirim",
		Message:   "Prestasi Anda telah dikirim",
		IsRead:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if notification.AchievementID != nil {
		t.Error("Expected AchievementID to be nil")
	}
	if notification.MongoAchievementID != nil {
		t.Error("Expected MongoAchievementID to be nil")
	}
	if notification.ReadAt != nil {
		t.Error("Expected ReadAt to be nil")
	}
}

func TestNotification_JSONMarshallingNilPointers(t *testing.T) {
	now := time.Now()
	notification := modelpostgre.Notification{
		ID:        "notification-id-1",
		UserID:    "user-id-1",
		Type:      modelpostgre.NotificationTypeAchievementSubmitted,
		Title:     "Prestasi Dikirim",
		Message:   "Prestasi Anda telah dikirim",
		IsRead:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if achievementID, exists := result["achievement_id"]; exists {
		if achievementID != nil {
			t.Errorf("Expected achievement_id to be null when nil, got '%v'", achievementID)
		}
	} else {
		t.Error("Expected achievement_id field to exist in JSON")
	}
	if mongoAchievementID, exists := result["mongo_achievement_id"]; exists {
		if mongoAchievementID != nil {
			t.Errorf("Expected mongo_achievement_id to be null when nil, got '%v'", mongoAchievementID)
		}
	} else {
		t.Error("Expected mongo_achievement_id field to exist in JSON")
	}
	if readAt, exists := result["read_at"]; exists {
		if readAt != nil {
			t.Errorf("Expected read_at to be null when nil, got '%v'", readAt)
		}
	} else {
		t.Error("Expected read_at field to exist in JSON")
	}
}

func TestCreateNotificationRequest_StructCreation(t *testing.T) {
	achievementID := "achievement-id-1"
	mongoAchievementID := "mongo-id-1"
	req := modelpostgre.CreateNotificationRequest{
		UserID:             "user-id-1",
		Type:               modelpostgre.NotificationTypeAchievementRejected,
		Title:              "Prestasi Ditolak",
		Message:            "Prestasi Anda ditolak",
		AchievementID:      &achievementID,
		MongoAchievementID: &mongoAchievementID,
	}

	if req.UserID != "user-id-1" {
		t.Errorf("Expected UserID 'user-id-1', got '%s'", req.UserID)
	}
	if req.Type != modelpostgre.NotificationTypeAchievementRejected {
		t.Errorf("Expected Type '%s', got '%s'", modelpostgre.NotificationTypeAchievementRejected, req.Type)
	}
	if req.Title != "Prestasi Ditolak" {
		t.Errorf("Expected Title 'Prestasi Ditolak', got '%s'", req.Title)
	}
	if req.AchievementID == nil || *req.AchievementID != "achievement-id-1" {
		t.Errorf("Expected AchievementID 'achievement-id-1', got '%v'", req.AchievementID)
	}
}

func TestCreateNotificationRequest_JSONMarshalling(t *testing.T) {
	achievementID := "achievement-id-1"
	mongoAchievementID := "mongo-id-1"
	req := modelpostgre.CreateNotificationRequest{
		UserID:             "user-id-1",
		Type:               modelpostgre.NotificationTypeAchievementRejected,
		Title:              "Prestasi Ditolak",
		Message:            "Prestasi Anda ditolak",
		AchievementID:      &achievementID,
		MongoAchievementID: &mongoAchievementID,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["user_id"] != "user-id-1" {
		t.Errorf("Expected user_id 'user-id-1', got '%v'", result["user_id"])
	}
	if result["type"] != "achievement_rejected" {
		t.Errorf("Expected type 'achievement_rejected', got '%v'", result["type"])
	}
	if result["title"] != "Prestasi Ditolak" {
		t.Errorf("Expected title 'Prestasi Ditolak', got '%v'", result["title"])
	}
}

func TestCreateNotificationRequest_JSONUnmarshalling(t *testing.T) {
	jsonStr := `{
		"user_id": "user-id-1",
		"type": "achievement_rejected",
		"title": "Prestasi Ditolak",
		"message": "Prestasi Anda ditolak",
		"achievement_id": "achievement-id-1",
		"mongo_achievement_id": "mongo-id-1"
	}`

	var req modelpostgre.CreateNotificationRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if req.UserID != "user-id-1" {
		t.Errorf("Expected UserID 'user-id-1', got '%s'", req.UserID)
	}
	if req.Type != "achievement_rejected" {
		t.Errorf("Expected Type 'achievement_rejected', got '%s'", req.Type)
	}
}

func TestGetNotificationsResponse_StructCreation(t *testing.T) {
	now := time.Now()
	notifications := []modelpostgre.Notification{
		{
			ID:        "notification-id-1",
			UserID:    "user-id-1",
			Type:      modelpostgre.NotificationTypeAchievementRejected,
			Title:     "Prestasi Ditolak",
			Message:   "Prestasi Anda ditolak",
			IsRead:    false,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "notification-id-2",
			UserID:    "user-id-1",
			Type:      modelpostgre.NotificationTypeAchievementSubmitted,
			Title:     "Prestasi Dikirim",
			Message:   "Prestasi Anda telah dikirim",
			IsRead:    true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	resp := modelpostgre.GetNotificationsResponse{
		Status: "success",
		Data:   notifications,
	}
	resp.Pagination.Page = 1
	resp.Pagination.Limit = 10
	resp.Pagination.Total = 2
	resp.Pagination.TotalPages = 1

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(resp.Data))
	}
	if resp.Pagination.Page != 1 {
		t.Errorf("Expected Pagination.Page 1, got %d", resp.Pagination.Page)
	}
	if resp.Pagination.Total != 2 {
		t.Errorf("Expected Pagination.Total 2, got %d", resp.Pagination.Total)
	}
}

func TestGetNotificationsResponse_JSONMarshalling(t *testing.T) {
	now := time.Now()
	notifications := []modelpostgre.Notification{
		{
			ID:        "notification-id-1",
			UserID:    "user-id-1",
			Type:      modelpostgre.NotificationTypeAchievementRejected,
			Title:     "Prestasi Ditolak",
			Message:   "Prestasi Anda ditolak",
			IsRead:    false,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	resp := modelpostgre.GetNotificationsResponse{
		Status: "success",
		Data:   notifications,
	}
	resp.Pagination.Page = 1
	resp.Pagination.Limit = 10
	resp.Pagination.Total = 1
	resp.Pagination.TotalPages = 1

	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", result["status"])
	}
	if pagination, ok := result["pagination"].(map[string]interface{}); ok {
		if pagination["page"] != float64(1) {
			t.Errorf("Expected pagination.page 1, got '%v'", pagination["page"])
		}
		if pagination["total"] != float64(1) {
			t.Errorf("Expected pagination.total 1, got '%v'", pagination["total"])
		}
	} else {
		t.Error("Expected pagination to be a map")
	}
}

func TestGetUnreadCountResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.GetUnreadCountResponse{
		Status: "success",
	}
	resp.Data.Count = 5

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Count != 5 {
		t.Errorf("Expected Data.Count 5, got %d", resp.Data.Count)
	}
}

func TestGetUnreadCountResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.GetUnreadCountResponse{
		Status: "success",
	}
	resp.Data.Count = 5

	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", result["status"])
	}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if data["count"] != float64(5) {
			t.Errorf("Expected data.count 5, got '%v'", data["count"])
		}
	} else {
		t.Error("Expected data to be a map")
	}
}

func TestMarkAsReadResponse_StructCreation(t *testing.T) {
	now := time.Now()
	readAt := now.Add(1 * time.Hour)
	notification := modelpostgre.Notification{
		ID:        "notification-id-1",
		UserID:    "user-id-1",
		Type:      modelpostgre.NotificationTypeAchievementRejected,
		Title:     "Prestasi Ditolak",
		Message:   "Prestasi Anda ditolak",
		IsRead:    true,
		ReadAt:    &readAt,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := modelpostgre.MarkAsReadResponse{
		Status: "success",
		Data:   notification,
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.ID != "notification-id-1" {
		t.Errorf("Expected Data.ID 'notification-id-1', got '%s'", resp.Data.ID)
	}
	if resp.Data.IsRead != true {
		t.Errorf("Expected Data.IsRead true, got %v", resp.Data.IsRead)
	}
}

func TestMarkAllAsReadResponse_StructCreation(t *testing.T) {
	resp := modelpostgre.MarkAllAsReadResponse{
		Status: "success",
	}
	resp.Data.Message = "Semua notifikasi telah ditandai sebagai dibaca"

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got '%s'", resp.Status)
	}
	if resp.Data.Message != "Semua notifikasi telah ditandai sebagai dibaca" {
		t.Errorf("Expected Data.Message 'Semua notifikasi telah ditandai sebagai dibaca', got '%s'", resp.Data.Message)
	}
}

func TestMarkAllAsReadResponse_JSONMarshalling(t *testing.T) {
	resp := modelpostgre.MarkAllAsReadResponse{
		Status: "success",
	}
	resp.Data.Message = "Semua notifikasi telah ditandai sebagai dibaca"

	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", result["status"])
	}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if data["message"] != "Semua notifikasi telah ditandai sebagai dibaca" {
			t.Errorf("Expected data.message 'Semua notifikasi telah ditandai sebagai dibaca', got '%v'", data["message"])
		}
	} else {
		t.Error("Expected data to be a map")
	}
}
