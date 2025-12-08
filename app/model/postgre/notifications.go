package model

import "time"

const (
	NotificationTypeAchievementRejected = "achievement_rejected"
	NotificationTypeAchievementSubmitted = "achievement_submitted"
)

type Notification struct {
	ID                string     `json:"id"`
	UserID            string     `json:"user_id"`
	Type              string     `json:"type"`
	Title             string     `json:"title"`
	Message           string     `json:"message"`
	AchievementID     *string    `json:"achievement_id"`
	MongoAchievementID *string   `json:"mongo_achievement_id"`
	IsRead            bool       `json:"is_read"`
	ReadAt            *time.Time `json:"read_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type CreateNotificationRequest struct {
	UserID            string  `json:"user_id"`
	Type              string  `json:"type"`
	Title             string  `json:"title"`
	Message           string  `json:"message"`
	AchievementID     *string `json:"achievement_id"`
	MongoAchievementID *string `json:"mongo_achievement_id"`
}

type GetNotificationsResponse struct {
	Status string        `json:"status"`
	Data   []Notification `json:"data"`
	Pagination struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	} `json:"pagination"`
}

type GetUnreadCountResponse struct {
	Status string `json:"status"`
	Data   struct {
		Count int `json:"count"`
	} `json:"data"`
}

type MarkAsReadResponse struct {
	Status string      `json:"status"`
	Data   Notification `json:"data"`
}

type MarkAllAsReadResponse struct {
	Status string `json:"status"`
	Data   struct {
		Message string `json:"message"`
	} `json:"data"`
}

