package model

// #1 proses: import library time untuk handle timestamp
import "time"

// #2 proses: definisikan konstanta tipe notifikasi yang tersedia
const (
	NotificationTypeAchievementRejected  = "achievement_rejected"
	NotificationTypeAchievementSubmitted = "achievement_submitted"
)

// #3 proses: struct utama untuk menyimpan data notifikasi di database
type Notification struct {
	ID                 string     `json:"id"`
	UserID             string     `json:"user_id"`
	Type               string     `json:"type"`
	Title              string     `json:"title"`
	Message            string     `json:"message"`
	AchievementID      *string    `json:"achievement_id"`
	MongoAchievementID *string    `json:"mongo_achievement_id"`
	IsRead             bool       `json:"is_read"`
	ReadAt             *time.Time `json:"read_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// #4 proses: struct untuk request create notifikasi baru
type CreateNotificationRequest struct {
	UserID             string  `json:"user_id"`
	Type               string  `json:"type"`
	Title              string  `json:"title"`
	Message            string  `json:"message"`
	AchievementID      *string `json:"achievement_id"`
	MongoAchievementID *string `json:"mongo_achievement_id"`
}

// #5 proses: struct response untuk get notifications dengan pagination
type GetNotificationsResponse struct {
	Status     string         `json:"status"`
	Data       []Notification `json:"data"`
	Pagination struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	} `json:"pagination"`
}

// #6 proses: struct response untuk get unread count, return jumlah notifikasi belum dibaca
type GetUnreadCountResponse struct {
	Status string `json:"status"`
	Data   struct {
		Count int `json:"count"`
	} `json:"data"`
}

// #7 proses: struct response untuk mark as read, return notifikasi yang sudah ditandai
type MarkAsReadResponse struct {
	Status string       `json:"status"`
	Data   Notification `json:"data"`
}

// #8 proses: struct response untuk mark all as read, return pesan konfirmasi
type MarkAllAsReadResponse struct {
	Status string `json:"status"`
	Data   struct {
		Message string `json:"message"`
	} `json:"data"`
}
